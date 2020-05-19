package internal

import (
	"bytes"
	"context"
	"fmt"
	"github.com/ettec/open-trading-platform/go/common"
	api "github.com/ettec/open-trading-platform/go/common/api/executionvenue"
	"github.com/ettec/open-trading-platform/go/common/marketdata"
	"github.com/ettec/open-trading-platform/go/model"
	"github.com/golang/protobuf/proto"
	logger "log"
	"os"
)

var zero *model.Decimal64

func init() {
	zero = &model.Decimal64{}
}

type orderManager struct {
	lastStoredOrder    []byte
	cancelChan         chan bool
	store              func(*model.Order) error
	managedOrder       *parentOrder
	underlyingListings map[int32]*model.Listing
	Id                 string
	orderRouter        api.ExecutionVenueClient
	log                *logger.Logger
	errLog             *logger.Logger
}

func (om *orderManager) GetManagedOrderId() string {
	return om.managedOrder.GetId()
}

func (om *orderManager) Cancel() {
	om.cancelChan <- true
}

func (om *orderManager) persistManagedOrderChanges() error {

	orderAsBytes, err := proto.Marshal(&om.managedOrder.Order)

	if bytes.Compare(om.lastStoredOrder, orderAsBytes) != 0 {

		if om.lastStoredOrder != nil {
			om.managedOrder.Version = om.managedOrder.Version + 1
		}

		toStore, err := proto.Marshal(&om.managedOrder.Order)

		om.lastStoredOrder = toStore

		orderCopy := &model.Order{}
		proto.Unmarshal(toStore, orderCopy)
		om.store(orderCopy)

		if err != nil {
			return err
		}

	}

	return err
}

type GetListingsWithSameInstrument = func(listingId int32, listingGroupsIn chan<- []*model.Listing)

func NewOrderManagerFromParams(id string, params *api.CreateAndRouteOrderParams,
	orderManagerId string, getListingsWithSameInstrument GetListingsWithSameInstrument, doneChan chan<- string, store func(*model.Order) error, orderRouter api.ExecutionVenueClient,
	quoteStream marketdata.MdsQuoteStream, childOrderStream ChildOrderStream) (*orderManager, error) {

	initialState := model.NewOrder(id, params.OrderSide, params.Quantity, params.Price, params.Listing.Id,
		params.OriginatorId, params.OriginatorRef, params.RootOriginatorId, params.RootOriginatorRef)
	err := initialState.SetTargetStatus(model.OrderStatus_LIVE)
	if err != nil {
		return nil, err
	}

	om := NewOrderManager(initialState, store, orderManagerId, orderRouter, getListingsWithSameInstrument, quoteStream, childOrderStream, doneChan)

	return om, nil
}

func NewOrderManager(initialState *model.Order, store func(*model.Order) error, orderManagerId string, orderRouter api.ExecutionVenueClient,
	getListingsWithSameInstrument GetListingsWithSameInstrument, quoteStream marketdata.MdsQuoteStream, childOrderStream ChildOrderStream,
	doneChan chan<- string) *orderManager {
	po := newParentOrder(*initialState)

	om := newOrderManager(store, po, orderManagerId, orderRouter)

	go func() {

		om.log.Println("initialising order")

		listingsIn := make(chan []*model.Listing)

		getListingsWithSameInstrument(po.ListingId, listingsIn)

		underlyingListings := map[int32]*model.Listing{}
		select {
		case ls := <-listingsIn:
			for _, listing := range ls {
				if listing.Market.Mic != common.SR_MIC {
					underlyingListings[listing.Id] = listing
				}
			}
		}

		om.underlyingListings = underlyingListings

		quoteStream.Subscribe(po.ListingId)

		om.log.Println("order initialised")

		for {
			om.persistManagedOrderChanges()
			if po.IsTerminalState() {
				quoteStream.Close()
				childOrderStream.Close()
				doneChan <- po.GetId()
				break
			}

			select {
			case <-om.cancelChan:

				if !po.IsTerminalState() {
					om.log.Print("cancelling order")
					po.SetTargetStatus(model.OrderStatus_CANCELLED)

					pendingChildOrderCancels := false
					for _, co := range po.childOrders {
						if !co.IsTerminalState() {
							pendingChildOrderCancels = true
							orderRouter.CancelOrder(context.Background(), &api.CancelOrderParams{
								OrderId: co.Id,
								Listing: underlyingListings[co.ListingId],
							})
						}

					}

					if !pendingChildOrderCancels {
						po.SetStatus(model.OrderStatus_CANCELLED)
					}

				}
			case co, ok := <-childOrderStream.GetStream():
				if ok {
					po.onChildOrderUpdate(co)
				} else {
					om.errLog.Printf("child order update chan unexpectedly closed, cancelling order")
					om.Cancel()
				}

			case q, ok := <-quoteStream.GetStream():

				if ok {

					if !q.StreamInterrupted {

						if po.GetAvailableQty().GreaterThan(zero) {
							if po.GetSide() == model.Side_BUY {
								om.submitBuyOrders(q)
							} else {
								om.submitSellOrders(q)
							}

						}

						if po.GetTargetStatus() == model.OrderStatus_LIVE {
							po.SetStatus(model.OrderStatus_LIVE)
						}
					}

				} else {
					om.errLog.Printf("quote chan unexpectedly closed, cancelling order")
					om.Cancel()
				}
			}

		}

	}()
	return om
}

func newOrderManager(store func(*model.Order) error, po *parentOrder, execVenueId string, orderRouter api.ExecutionVenueClient) *orderManager {
	om := &orderManager{
		lastStoredOrder: nil,
		cancelChan:      make(chan bool, 1),
		store:           store,
		managedOrder:    po,
		Id:              execVenueId,
		orderRouter:     orderRouter,
		log:             logger.New(os.Stdout, "order:"+po.Id, logger.Lshortfile|logger.Ltime),
		errLog:          logger.New(os.Stderr, "order:"+po.Id, logger.Lshortfile|logger.Ltime),
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

}

func (om *orderManager) sendChildOrder(side model.Side, quantity *model.Decimal64, price *model.Decimal64, listingId int32) {
	params := &api.CreateAndRouteOrderParams{
		OrderSide:         side,
		Quantity:          quantity,
		Price:             price,
		Listing:           om.underlyingListings[listingId],
		OriginatorId:      om.Id,
		OriginatorRef:     om.GetManagedOrderId(),
		RootOriginatorId:  om.managedOrder.RootOriginatorId,
		RootOriginatorRef: om.managedOrder.RootOriginatorRef,
	}

	id, err := om.orderRouter.CreateAndRouteOrder(context.Background(), params)

	if err != nil {
		msg := fmt.Sprintf("failed to submit child order:%v", err)
		om.cancelOrderWithErrorMsg(msg)
		return
	}

	pendingOrder := model.NewOrder(id.OrderId, params.OrderSide, params.Quantity, params.Price, params.Listing.GetId(),
		om.Id, om.GetManagedOrderId(), om.managedOrder.RootOriginatorId, om.managedOrder.RootOriginatorRef)

	// First persisted orders start at version 0, this is a placeholder until the first child order update is received
	pendingOrder.Version = -1

	om.managedOrder.onChildOrderUpdate(pendingOrder)

}

func (om *orderManager) cancelOrderWithErrorMsg(msg string) {
	om.managedOrder.ErrorMessage = msg
	om.cancelChan <- true
}
