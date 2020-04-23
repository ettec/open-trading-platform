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
	managedOrderId  string
	lastStoredOrder model.Order
	store           func(model.Order) error
}

func (om *orderManager) GetManagedOrderId() string {
	return om.managedOrderId
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
	execVenueId string, underlyingListings map[int32]*model.Listing, doneChan chan<- string, store func(model.Order) error, orderRouter api.ExecutionVenueClient,
	quoteChan <-chan *model.ClobQuote, childOrderUpdates <-chan *model.Order) (*orderManager, error) {

	uniqueId, err := uuid.NewUUID()

	if err != nil {
		return nil, err
	}

	om := &orderManager{
		managedOrderId:  uniqueId.String(),
		lastStoredOrder: model.Order{},
		store:           store,
	}

	orderState := model.NewOrder(uniqueId.String(), params.OrderSide, params.Quantity, params.Price, params.Listing.Id,
		params.OriginatorId, params.OriginatorRef)

	po := newParentOrder(*orderState)
	po.SetTargetStatus(model.OrderStatus_LIVE)

	go func() {

		ordersSubmitted := false

		for {
			om.persist(po.Order)
			if po.IsTerminalState() {
				doneChan <- po.GetId()
				break
			}

			select {
			case co := <-childOrderUpdates:
				po.onChildOrderUpdate(co)
			case q := <-quoteChan:

				if !ordersSubmitted && !q.StreamInterrupted {

					if po.GetSide() == model.Side_BUY {
						submitBuyOrders(q, po, underlyingListings, execVenueId, orderRouter)
					} else {
						submitSellOrders(q, po, underlyingListings, execVenueId, orderRouter)
					}

					ordersSubmitted = true

					po.SetStatus(model.OrderStatus_LIVE)
				}
			}

		}

	}()

	return om, nil
}

func submitBuyOrders(q *model.ClobQuote, po *parentOrder, underlyingListings map[int32]*model.Listing, execVenueId string, orderRouter api.ExecutionVenueClient) {
	submitOrders(q.Offers, func(line *model.ClobLine) bool {
		return line.Price.LessThanOrEqual(po.GetPrice())
	}, po, model.Side_BUY, underlyingListings, execVenueId, orderRouter)
}

func submitSellOrders(q *model.ClobQuote, po *parentOrder, underlyingListings map[int32]*model.Listing, execVenueId string, orderRouter api.ExecutionVenueClient) {
	submitOrders(q.Bids, func(line *model.ClobLine) bool {
		return line.Price.GreaterThanOrEqual(po.GetPrice())
	}, po, model.Side_SELL, underlyingListings, execVenueId, orderRouter)
}

func submitOrders(oppositeClobLines []*model.ClobLine, willTrade func(line *model.ClobLine) bool, po *parentOrder,
	side model.Side, underlyingListings map[int32]*model.Listing, execVenueId string, orderRouter api.ExecutionVenueClient) {
	listingIdToQnt := map[int32]*model.Decimal64{}
	for _, line := range oppositeClobLines {
		if po.GetAvailableQty().GreaterThan(zero) && willTrade(line) {
			quantity := line.Size

			if line.Size.GreaterThanOrEqual(po.GetAvailableQty()) {
				quantity = po.GetAvailableQty()
			}

			sendChildOrder(side, quantity, line.Price, underlyingListings[line.ListingId], execVenueId, po, orderRouter)

		} else {
			if qnt, ok := listingIdToQnt[line.ListingId]; ok {
				qnt.Add(line.Size)
			} else {
				qntCopy := *line.Size
				listingIdToQnt[line.ListingId] = &qntCopy
			}
		}

	}

	if po.GetAvailableQty().GreaterThan(zero) {
		var l int32
		greatestQty := zero
		for listingId, qnt := range listingIdToQnt {
			if qnt.GreaterThan(greatestQty) {
				l = listingId
				greatestQty = qnt
			}
		}

		if greatestQty.GreaterThan(zero) {
			sendChildOrder(side, po.GetAvailableQty(), po.Price, underlyingListings[l], execVenueId, po, orderRouter)
		}

	}

}

func sendChildOrder(side model.Side, quantity *model.Decimal64, price *model.Decimal64, listing *model.Listing,
	execVenueId string, mo *parentOrder, orderRouter api.ExecutionVenueClient) {
	params := &api.CreateAndRouteOrderParams{
		OrderSide:     side,
		Quantity:      quantity,
		Price:         price,
		Listing:       listing,
		OriginatorId:  execVenueId,
		OriginatorRef: mo.GetId(),
	}

	id, err := orderRouter.CreateAndRouteOrder(context.Background(), params)

	if err != nil {
		msg := fmt.Sprintf("failed to submit child order:%v", err)
		cancelOrderWithErrorMsg(mo, msg)
	}

	childOrder := model.NewOrder(id.OrderId, params.OrderSide, params.Quantity, params.Price, params.Listing.GetId(),
		execVenueId, mo.GetId())

	mo.onChildOrderUpdate(childOrder)

}

func cancelOrderWithErrorMsg(po *parentOrder, msg string) {
	po.ErrorMessage = msg
	//here - cancel the order - self cancel
}
