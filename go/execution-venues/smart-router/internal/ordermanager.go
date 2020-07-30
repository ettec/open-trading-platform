package internal

import (
	"github.com/ettec/otp-common"
	"github.com/ettec/otp-common/marketdata"
	"github.com/ettec/otp-common/model"
)

var zero *model.Decimal64

func init() {
	zero = &model.Decimal64{}
}

type GetListingsWithSameInstrument = func(listingId int32, listingGroupsIn chan<- []*model.Listing)

func ExecuteAsSmartRouterStrategy(om *commonOrderManager,
	getListingsWithSameInstrument GetListingsWithSameInstrument, quoteStream marketdata.MdsQuoteStream) {

	go func() {

		om.log.Println("initialising order")

		listingsIn := make(chan []*model.Listing)

		getListingsWithSameInstrument(om.managedOrder.ListingId, listingsIn)

		underlyingListings := map[int32]*model.Listing{}
		select {
		case ls := <-listingsIn:
			for _, listing := range ls {
				if listing.Market.Mic != common.SR_MIC {
					underlyingListings[listing.Id] = listing
				}
			}
		}

		quoteStream.Subscribe(om.managedOrder.ListingId)

		om.log.Println("order initialised")

		for {
			close := om.persistChanges()

			if close {
				quoteStream.Close()
				break
			}

			select {
			case <-om.cancelChan:
				om.cancelOrder(func(listingId int32) *model.Listing {
					return underlyingListings[listingId]
				})
			case co, ok := <-om.childOrderStream.GetStream():
				om.onChildOrderUpdate(ok, co)
			case q, ok := <-quoteStream.GetStream():

				if ok {

					if !q.StreamInterrupted {

						if om.managedOrder.GetAvailableQty().GreaterThan(zero) {
							if om.managedOrder.GetSide() == model.Side_BUY {
								om.submitBuyOrders(q, underlyingListings)
							} else {
								om.submitSellOrders(q, underlyingListings)
							}

						}

						if om.managedOrder.GetTargetStatus() == model.OrderStatus_LIVE {
							om.managedOrder.SetStatus(model.OrderStatus_LIVE)
						}
					}

				} else {
					om.errLog.Printf("quote chan unexpectedly closed, cancelling order")
					om.Cancel()
				}
			}

		}

	}()
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
