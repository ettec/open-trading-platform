package main

import (
	"bytes"
	"context"
	"fmt"
	api "github.com/ettec/open-trading-platform/go/common/api/executionvenue"
	"github.com/ettec/open-trading-platform/go/model"
	"github.com/gogo/protobuf/proto"
	"github.com/google/uuid"
)

var zero *model.Decimal64

func init() {
	zero = &model.Decimal64{}
}

type orderManager struct {
	lastStoredOrder    model.Order
	cancelChan         chan bool
	store              func(model.Order) error
	managedOrder       *parentOrder
	underlyingListings map[int32]*model.Listing
	Id                 string
	orderRouter        api.ExecutionVenueClient
}

func (om *orderManager) GetManagedOrderId() string {
	return om.managedOrder.GetId()
}

func (om *orderManager) Cancel() {
	om.cancelChan <- true
}

func (om *orderManager) persist(order model.Order) error {

	last, err := proto.Marshal(&om.lastStoredOrder)
	if err != nil {
		return err
	}

	new, err := proto.Marshal(&order)

	if bytes.Compare(last, new) != 0 {
		newVersionNum := om.lastStoredOrder.Version + 1
		om.lastStoredOrder = order
		om.lastStoredOrder.Version = newVersionNum
		om.store(order)
	}

	return err
}

func NewOrderManager(params *api.CreateAndRouteOrderParams,
	orderManagerId string, underlyingListings map[int32]*model.Listing, doneChan chan<- string, store func(model.Order) error, orderRouter api.ExecutionVenueClient,
	quoteChan <-chan *model.ClobQuote, childOrderUpdates <-chan *model.Order) (*orderManager, error) {

	uniqueId, err := uuid.NewUUID()

	if err != nil {
		return nil, err
	}

	initialState := model.NewOrder(uniqueId.String(), params.OrderSide, params.Quantity, params.Price, params.Listing.Id,
		params.OriginatorId, params.OriginatorRef)
	initialState.SetTargetStatus(model.OrderStatus_LIVE)

	po := newParentOrder(*initialState)

	om := newOrderManager(store, po, underlyingListings, orderManagerId, orderRouter)

	go func() {

		ordersSubmitted := false

		for {
			om.persist(po.Order)
			if po.IsTerminalState() {
				doneChan <- po.GetId()
				break
			}

			select {
			case <-om.cancelChan:

				if !po.IsTerminalState() {
					po.SetTargetStatus(model.OrderStatus_CANCELLED)
					for _, co := range po.childOrders {
						orderRouter.CancelOrder(context.Background(), &api.CancelOrderParams{
							OrderId: co.Id,
							Listing: underlyingListings[co.ListingId],
						})
					}
				}
			case co := <-childOrderUpdates:
				po.onChildOrderUpdate(co)
			case q := <-quoteChan:

				if !ordersSubmitted && !q.StreamInterrupted {

					if po.GetSide() == model.Side_BUY {
						om.submitBuyOrders(q)
					} else {
						om.submitSellOrders(q)
					}

					ordersSubmitted = true

					po.SetStatus(model.OrderStatus_LIVE)
				}
			}

		}

	}()

	return om, nil
}

func newOrderManager(store func(model.Order) error, po *parentOrder, underlyingListings map[int32]*model.Listing, execVenueId string, orderRouter api.ExecutionVenueClient) *orderManager {
	om := &orderManager{
		lastStoredOrder:    model.Order{},
		cancelChan:         make(chan bool, 1),
		store:              store,
		managedOrder:       po,
		underlyingListings: underlyingListings,
		Id:                 execVenueId,
		orderRouter:        orderRouter,
	}
	return om
}

func (om *orderManager) submitBuyOrders(q *model.ClobQuote) {
	om.submitOrders(q.Offers, func(line *model.ClobLine) bool {
		return line.Price.LessThanOrEqual(om.managedOrder.GetPrice())
	}, model.Side_BUY)
}

func (om *orderManager) submitSellOrders(q *model.ClobQuote) {
	om.submitOrders(q.Bids, func(line *model.ClobLine) bool {
		return line.Price.GreaterThanOrEqual(om.managedOrder.GetPrice())
	}, model.Side_SELL)
}

func (om *orderManager) submitOrders(oppositeClobLines []*model.ClobLine, willTrade func(line *model.ClobLine) bool,
	side model.Side) {
	listingIdToQnt := map[int32]*model.Decimal64{}
	for _, line := range oppositeClobLines {
		if om.managedOrder.GetAvailableQty().GreaterThan(zero) && willTrade(line) {
			quantity := line.Size

			if line.Size.GreaterThanOrEqual(om.managedOrder.GetAvailableQty()) {
				quantity = om.managedOrder.GetAvailableQty()
			}

			om.sendChildOrder(side, quantity, line.Price, line.ListingId)

		} else {
			if qnt, ok := listingIdToQnt[line.ListingId]; ok {
				qnt.Add(line.Size)
			} else {
				qntCopy := *line.Size
				listingIdToQnt[line.ListingId] = &qntCopy
			}
		}

	}

	if om.managedOrder.GetAvailableQty().GreaterThan(zero) {
		var l int32
		greatestQty := zero
		for listingId, qnt := range listingIdToQnt {
			if qnt.GreaterThan(greatestQty) {
				l = listingId
				greatestQty = qnt
			}
		}

		if greatestQty.GreaterThan(zero) {
			om.sendChildOrder(side, om.managedOrder.GetAvailableQty(), om.managedOrder.Price, l)
		}

	}

}

func (om *orderManager) sendChildOrder(side model.Side, quantity *model.Decimal64, price *model.Decimal64, listingId int32) {
	params := &api.CreateAndRouteOrderParams{
		OrderSide:     side,
		Quantity:      quantity,
		Price:         price,
		Listing:       om.underlyingListings[listingId],
		OriginatorId:  om.Id,
		OriginatorRef: om.GetManagedOrderId(),
	}

	id, err := om.orderRouter.CreateAndRouteOrder(context.Background(), params)

	if err != nil {
		msg := fmt.Sprintf("failed to submit child order:%v", err)
		om.cancelOrderWithErrorMsg(msg)
		return
	}

	childOrder := model.NewOrder(id.OrderId, params.OrderSide, params.Quantity, params.Price, params.Listing.GetId(),
		om.Id, om.GetManagedOrderId())

	om.managedOrder.onChildOrderUpdate(childOrder)

}

func (om *orderManager)  cancelOrderWithErrorMsg(msg string) {
	om.managedOrder.ErrorMessage = msg
	om.cancelChan <-true
}
