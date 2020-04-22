package main

import (
	"context"
	"fmt"
	api "github.com/ettec/open-trading-platform/go/common/api/executionvenue"
	"github.com/ettec/open-trading-platform/go/model"
	"github.com/google/uuid"
)

type execution struct {
	id    string
	price *model.Decimal64
	qty   *model.Decimal64
}

type managedOrder struct {
	modified    bool
	order       *model.Order
	childOrders map[string]*model.Order
	executions  map[string]*execution
}

func newManagedOrder(order *model.Order) *managedOrder {

	return &managedOrder{
		modified:    true,
		order:       order,
		childOrders: map[string]*model.Order{},
		executions:  map[string]*execution{},
	}
}

func (m *managedOrder) GetID() string {
	return m.order.GetId()
}

func (m *managedOrder) isTerminalState() bool {
	return m.order.IsTerminalState()
}

func (m *managedOrder) getListingId() int32 {
	return m.order.GetListingId()
}

func (m *managedOrder) setTargetStatus(status model.OrderStatus) error {
	err := m.order.SetTargetStatus(status)
	if err != nil {
		return err
	}

	m.modified = true
	return nil
}

func (m *managedOrder) setErrorMsg(msg string) {
	m.order.ErrorMessage = msg
	m.modified = true
}

func (m *managedOrder) setStatus(status model.OrderStatus) error {
	err := m.order.SetStatus(status)
	if err != nil {
		return err
	}

	m.modified = true
	return nil
}

func (m *managedOrder) getAvailableQuantity() *model.Decimal64 {
	d := m.order.GetAvailableQty()
	return &d
}

func (m *managedOrder) onChildOrderUpdate(order *model.Order) {

	var lastExecSeqNo int32

	if previous, ok := m.childOrders[order.Id]; ok {
		lastExecSeqNo = previous.LastExecSeqNo
	}

	if order.LastExecSeqNo > lastExecSeqNo {
		execId := order.Id + ":" + order.LastExecId

		execution := execution{
			id:    execId,
			price: order.LastExecPrice,
			qty:   order.LastExecQuantity,
		}

		m.executions[execId] = &execution

		order.AddExecution(*execution.price, *execution.qty, execution.id)
		m.modified = true

	}

	m.childOrders[order.Id] = order

	exposedQnt := model.IasD(0)
	for _, order := range m.childOrders {
		if !order.IsTerminalState() {
			exposedQnt.Add(order.RemainingQuantity)
		}

		if !m.order.ExposedQuantity.Equal(exposedQnt) {
			m.order.ExposedQuantity = exposedQnt
			m.modified = true
		}
	}

}

func NewOrderManager(params *api.CreateAndRouteOrderParams,
	execVenueId string, underlyingListings map[int32]*model.Listing, doneChan chan<- *managedOrder, store func(model.Order) error, orderRouter api.ExecutionVenueClient,
	quoteChan <-chan *model.ClobQuote, childOrderUpdates <-chan *model.Order) (*managedOrder, error) {

	uniqueId, err := uuid.NewUUID()

	if err != nil {
		return nil, err
	}

	orderState := model.NewOrder(uniqueId.String(), params.OrderSide, params.Quantity, params.Price, params.Listing.Id,
		params.OriginatorId, params.OriginatorRef)

	mo := newManagedOrder(orderState)

	go func() {

		ordersSubmitted := false

		for {
			persist(mo, store)
			if mo.isTerminalState() {
				doneChan <- mo
				break
			}

			select {
			case co := <-childOrderUpdates:
				mo.onChildOrderUpdate(co)
			case q := <-quoteChan:

				if !ordersSubmitted && !q.StreamInterrupted {

					if mo.order.GetSide() == model.Side_BUY {
						submitBuyOrders(q, mo, underlyingListings, execVenueId, orderRouter)
					} else {
						submitSellOrders(q, mo, underlyingListings, execVenueId, orderRouter)
					}

					mo.setStatus(model.OrderStatus_LIVE)
				}
			}

		}

	}()

	return mo, nil
}

func submitBuyOrders(q *model.ClobQuote, mo *managedOrder, underlyingListings map[int32]*model.Listing, execVenueId string, orderRouter api.ExecutionVenueClient) {
	submitOrders(q.Offers, func(line *model.ClobLine) bool {
		return line.Price.LessThanOrEqual(mo.order.GetPrice())
	}, mo, model.Side_BUY, underlyingListings, execVenueId, orderRouter)
}

func submitSellOrders(q *model.ClobQuote, mo *managedOrder, underlyingListings map[int32]*model.Listing, execVenueId string, orderRouter api.ExecutionVenueClient) {
	submitOrders(q.Bids, func(line *model.ClobLine) bool {
		return line.Price.GreaterThanOrEqual(mo.order.GetPrice())
	}, mo, model.Side_SELL, underlyingListings, execVenueId, orderRouter)
}

func submitOrders(oppositeClobLines []*model.ClobLine, willTrade func(line *model.ClobLine) bool, mo *managedOrder, side model.Side, underlyingListings map[int32]*model.Listing, execVenueId string, orderRouter api.ExecutionVenueClient) {
	listingIdToQnt := map[int32]*model.Decimal64{}
	for _, line := range oppositeClobLines {
		if willTrade(line) {
			quantity := line.Size

			if line.Size.GreaterThanOrEqual(mo.getAvailableQuantity()) {
				quantity = mo.getAvailableQuantity()
			}

			sendChildOrder(side, quantity, line.Price, underlyingListings[line.ListingId], execVenueId, mo, orderRouter)

		} else {
			if qnt, ok := listingIdToQnt[line.ListingId]; ok {
				qnt.Add(line.Size)
			} else {
				qntCopy := *line.Size
				listingIdToQnt[line.ListingId] = &qntCopy
			}
		}

	}

	if mo.getAvailableQuantity().GreaterThan(model.IasD(0)) {
		var l int32
		greatestQty := model.IasD(0)
		for listingId, qnt := range listingIdToQnt {
			if qnt.GreaterThan(greatestQty) {
				l = listingId
				greatestQty = qnt
			}
		}

		if greatestQty.GreaterThan(model.IasD(0)) {
			sendChildOrder(side, mo.getAvailableQuantity(), mo.order.Price, underlyingListings[l], execVenueId, mo, orderRouter)
		}

	}

}

func sendChildOrder(side model.Side, quantity *model.Decimal64, price *model.Decimal64, listing *model.Listing,
	execVenueId string, mo *managedOrder, orderRouter api.ExecutionVenueClient) {
	params := &api.CreateAndRouteOrderParams{
		OrderSide:     side,
		Quantity:      quantity,
		Price:         price,
		Listing:       listing,
		OriginatorId:  execVenueId,
		OriginatorRef: mo.GetID(),
	}

	id, err := orderRouter.CreateAndRouteOrder(context.Background(), params)

	if err != nil {
		msg := fmt.Sprintf("failed to submit child order:%v", err)
		cancelOrderWithErrorMsg(mo, msg)
	}

	childOrder := model.NewOrder(id.OrderId, params.OrderSide, params.Quantity, params.Price, params.Listing.GetId(),
		execVenueId, mo.GetID())

	mo.onChildOrderUpdate(childOrder)

}

func cancelOrderWithErrorMsg(mo *managedOrder, msg string) {
	mo.setErrorMsg(msg)
	//here - cancel the order - self cancel
}

func persist(order *managedOrder, store func(model.Order) error) {
	if order.modified {
		order.order.Version = order.order.Version + 1
		store(*order.order)
		order.modified = false
	}
}
