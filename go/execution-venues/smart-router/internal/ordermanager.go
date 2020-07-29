package internal

import (
	"github.com/ettec/otp-common"
	api "github.com/ettec/otp-common/api/executionvenue"
	"github.com/ettec/otp-common/executionvenue"
	"github.com/ettec/otp-common/marketdata"
	"github.com/ettec/otp-common/model"
)

var zero *model.Decimal64

func init() {
	zero = &model.Decimal64{}
}

type orderManager struct {
	commonOrderManager
	underlyingListings map[int32]*model.Listing
}

type GetListingsWithSameInstrument = func(listingId int32, listingGroupsIn chan<- []*model.Listing)

func NewOrderManagerFromParams(id string, params *api.CreateAndRouteOrderParams,
	orderManagerId string, getListingsWithSameInstrument GetListingsWithSameInstrument, doneChan chan<- string, store func(*model.Order) error, orderRouter api.ExecutionVenueClient,
	quoteStream marketdata.MdsQuoteStream, childOrderStream executionvenue.ChildOrderStream) (*commonOrderManager, error) {

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
	getListingsWithSameInstrument GetListingsWithSameInstrument, quoteStream marketdata.MdsQuoteStream, childOrderStream executionvenue.ChildOrderStream,
	doneChan chan<- string) *commonOrderManager {
	po := executionvenue.NewParentOrder(*initialState)

	om := newCommonOrderManager(store, po, orderManagerId, orderRouter, childOrderStream, doneChan)

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

		quoteStream.Subscribe(po.ListingId)

		om.log.Println("order initialised")

		for {
			close := om.persistChanges()

			if close {
				quoteStream.Close()
				break
			}

			select {
			case <-om.cancelChan:
				om.cancelOrder( orderRouter, func(listingId int32) *model.Listing {
					return underlyingListings[listingId]
				})
			case co, ok := <-om.childOrderStream.GetStream():
				om.onChildOrderUpdate(ok, co)
			case q, ok := <-quoteStream.GetStream():

				if ok {

					if !q.StreamInterrupted {

						if po.GetAvailableQty().GreaterThan(zero) {
							if po.GetSide() == model.Side_BUY {
								om.submitBuyOrders(q, underlyingListings)
							} else {
								om.submitSellOrders(q, underlyingListings)
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



func (om *commonOrderManager) submitBuyOrders(q *model.ClobQuote, underlyingListings map[int32]*model.Listing) {
	om.submitOrders(q.Offers, func(line *model.ClobLine) bool {
		return line.Price.LessThanOrEqual(om.managedOrder.GetPrice())
	}, model.Side_BUY, underlyingListings)
}

func (om *commonOrderManager) submitSellOrders(q *model.ClobQuote, underlyingListings map[int32]*model.Listing) {
	om.submitOrders(q.Bids, func(line *model.ClobLine) bool {
		return line.Price.GreaterThanOrEqual(om.managedOrder.GetPrice())
	}, model.Side_SELL, underlyingListings)
}

func (om *commonOrderManager) submitOrders(oppositeClobLines []*model.ClobLine, willTrade func(line *model.ClobLine) bool,
	side model.Side, underlyingListings map[int32]*model.Listing) {
	listingIdToQnt := map[int32]*model.Decimal64{}
	for _, line := range oppositeClobLines {
		if om.managedOrder.GetAvailableQty().GreaterThan(zero) && willTrade(line) {
			quantity := line.Size

			if line.Size.GreaterThanOrEqual(om.managedOrder.GetAvailableQty()) {
				quantity = om.managedOrder.GetAvailableQty()
			}

			om.sendChildOrder(side, quantity, line.Price, underlyingListings[line.ListingId])

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

