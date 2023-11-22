package main

import (
	"context"
	"fmt"
	"github.com/ettec/otp-common"
	"github.com/ettec/otp-common/marketdata"
	"github.com/ettec/otp-common/model"
	"github.com/ettec/otp-common/staticdata"
	"github.com/ettec/otp-common/strategy"
)

var zero *model.Decimal64

func init() {
	zero = &model.Decimal64{}
}

type GetListingsWithSameInstrument = func(ctx context.Context, listingId int32, listingGroupsIn chan<- staticdata.ListingsResult)

func ExecuteAsSmartRouterStrategy(ctx context.Context, om *strategy.Strategy,
	getListingsWithSameInstrument GetListingsWithSameInstrument, stream marketdata.QuoteStream) {

	go func() {

		defer stream.Close()

		om.Log.Info("initialising order")

		listingsIn := make(chan staticdata.ListingsResult)

		getListingsWithSameInstrument(ctx, om.ParentOrder.ListingId, listingsIn)

		instrumentListings := map[int32]*model.Listing{}
		select {
		case <-ctx.Done():
			return
		case ls := <-listingsIn:
			if ls.Err != nil {
				om.Log.Error("failed to get listings with same instrument", "listingId", om.ParentOrder.ListingId, "error", ls.Err)
				om.CancelChan <- fmt.Sprintf("failed to get listings for same instrument for listingId:%v, err:%v", om.ParentOrder.ListingId, ls.Err)
			}

			for _, listing := range ls.Listings {
				if listing.Market.Mic != common.SR_MIC {
					instrumentListings[listing.Id] = listing
				}
			}
		}

		err := stream.Subscribe(om.ParentOrder.ListingId)
		if err != nil {
			om.Log.Error("failed to subscribe to listing", "error", err)
			om.CancelChan <- fmt.Sprintf("failed to subscribe to listing:%v", err)
		}

		om.Log.Info("order initialised", "status", om.ParentOrder.GetStatus(),
			"targetStatus", om.ParentOrder.GetTargetStatus())

		for {
			done, err := om.CheckIfDone(ctx)
			if err != nil {
				msg := fmt.Sprintf("failed to check if done, cancelling order:%v", err)
				om.Log.Error(msg)
				om.CancelChan <- msg
			}

			if done {
				return
			}

			select {
			case <-ctx.Done():
				return
			case errMsg := <-om.CancelChan:
				if errMsg != "" {
					om.ParentOrder.ErrorMessage = errMsg
				}
				err := om.CancelChildOrdersAndStrategyOrder()
				if err != nil {
					om.Log.Error("failed to cancel order", "error", err)
				}

			case co, ok := <-om.ChildOrderUpdateChan:
				err = om.OnChildOrderUpdate(ok, co)
				if err != nil {
					om.Log.Error("error processing child order update", "error", err)
				}

			case quote, ok := <-stream.Chan():
				if ok {
					if !quote.StreamInterrupted {

						if om.ParentOrder.GetAvailableQty().GreaterThan(zero) {
							if om.ParentOrder.GetSide() == model.Side_BUY {
								submitBuyOrders(om, quote, instrumentListings)
							} else {
								submitSellOrders(om, quote, instrumentListings)
							}

						}

						if om.ParentOrder.GetTargetStatus() == model.OrderStatus_LIVE {
							err := om.ParentOrder.SetStatus(model.OrderStatus_LIVE)
							om.Log.Info("order status changed to LIVE")
							if err != nil {
								msg := fmt.Sprintf("failed to set managed order status, cancelling order:%v", err)
								om.Log.Error(msg)
								om.CancelChan <- msg
							}

						}
					}

				} else {
					msg := "quote chan unexpectedly closed, cancelling order"
					om.Log.Error(msg)
					om.CancelChan <- msg
				}
			}

		}

	}()
}

func submitBuyOrders(om *strategy.Strategy, q *model.ClobQuote, instrumentListings map[int32]*model.Listing) {
	submitOrders(om, q.Offers, func(line *model.ClobLine) bool {
		return line.Price.LessThanOrEqual(om.ParentOrder.GetPrice())
	}, model.Side_BUY, instrumentListings)
}

func submitSellOrders(om *strategy.Strategy, q *model.ClobQuote, instrumentListings map[int32]*model.Listing) {
	submitOrders(om, q.Bids, func(line *model.ClobLine) bool {
		return line.Price.GreaterThanOrEqual(om.ParentOrder.GetPrice())
	}, model.Side_SELL, instrumentListings)
}

func submitOrders(om *strategy.Strategy, oppositeClobLines []*model.ClobLine, willTrade func(line *model.ClobLine) bool,
	side model.Side, instrumentListings map[int32]*model.Listing) {
	listingIdToQnt := map[int32]*model.Decimal64{}
	for _, line := range oppositeClobLines {
		if om.ParentOrder.GetAvailableQty().GreaterThan(zero) && willTrade(line) {
			quantity := line.Size

			if line.Size.GreaterThanOrEqual(om.ParentOrder.GetAvailableQty()) {
				quantity = om.ParentOrder.GetAvailableQty()
			}

			listing := instrumentListings[line.ListingId]
			err := om.SendChildOrder(side, quantity, line.Price, listing.Id, listing.Market.Mic, "")
			if err != nil {
				om.CancelChan <- fmt.Sprintf("failed to send child order:%v", err)
			}

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
