package main

import (
	"fmt"
	"github.com/ettec/otp-common"
	"github.com/ettec/otp-common/marketdata"
	"github.com/ettec/otp-common/model"
	"github.com/ettec/otp-common/strategy"
)

var zero *model.Decimal64

func init() {
	zero = &model.Decimal64{}
}

type GetListingsWithSameInstrument = func(listingId int32, listingGroupsIn chan<- []*model.Listing)

func ExecuteAsSmartRouterStrategy(om *strategy.Strategy,
	getListingsWithSameInstrument GetListingsWithSameInstrument, quoteStream marketdata.MdsQuoteStream) {

	go func() {

		om.Log.Println("initialising order")

		listingsIn := make(chan []*model.Listing)

		getListingsWithSameInstrument(om.ParentOrder.ListingId, listingsIn)

		instrumentListings := map[int32]*model.Listing{}
		select {
		case ls := <-listingsIn:
			for _, listing := range ls {
				if listing.Market.Mic != common.SR_MIC {
					instrumentListings[listing.Id] = listing
				}
			}
		}

		quoteStream.Subscribe(om.ParentOrder.ListingId)

		om.Log.Println("order initialised")

		for {
			done, err := om.CheckIfDone()
			if err != nil {
				msg := fmt.Sprintf("failed to check if done, cancelling order:%v", err)
				om.ErrLog.Print(msg)
				om.CancelChan <- msg
			}

			if done {
				quoteStream.Close()
				break
			}

			select {
			case errMsg := <-om.CancelChan:
				if errMsg != "" {
					om.ParentOrder.ErrorMessage = errMsg
				}
				err := om.CancelChildOrdersAndStrategyOrder()
				if err != nil {
					om.ErrLog.Printf("failed to cancel order:%v", err)
				}
			case co, ok := <-om.ChildOrderUpdateChan:
				err = om.OnChildOrderUpdate(ok, co)
				if err != nil {
					om.ErrLog.Printf("error whilst applying child order update:%v", err)
				}
			case q, ok := <-quoteStream.GetStream():

				if ok {

					if !q.StreamInterrupted {

						if om.ParentOrder.GetAvailableQty().GreaterThan(zero) {
							if om.ParentOrder.GetSide() == model.Side_BUY {
								submitBuyOrders(om, q, instrumentListings)
							} else {
								submitSellOrders(om, q, instrumentListings)
							}

						}

						if om.ParentOrder.GetTargetStatus() == model.OrderStatus_LIVE {
							err := om.ParentOrder.SetStatus(model.OrderStatus_LIVE)
							if err != nil {
								msg := fmt.Sprintf("failed to set managed order status, cancelling order:%v", err)
								om.ErrLog.Print(msg)
								om.CancelChan <- msg
							}

						}
					}

				} else {
					msg := "quote chan unexpectedly closed, cancelling order"
					om.ErrLog.Print(msg)
					om.CancelChan <- msg
				}
			}

		}

	}()
}

func  submitBuyOrders(om *strategy.Strategy, q *model.ClobQuote, instrumentListings map[int32]*model.Listing) {
	submitOrders(om, q.Offers, func(line *model.ClobLine) bool {
		return line.Price.LessThanOrEqual(om.ParentOrder.GetPrice())
	}, model.Side_BUY, instrumentListings)
}

func submitSellOrders(om *strategy.Strategy, q *model.ClobQuote, instrumentListings map[int32]*model.Listing) {
	submitOrders(om, q.Bids, func(line *model.ClobLine) bool {
		return line.Price.GreaterThanOrEqual(om.ParentOrder.GetPrice())
	}, model.Side_SELL, instrumentListings)
}

func  submitOrders(om *strategy.Strategy, oppositeClobLines []*model.ClobLine, willTrade func(line *model.ClobLine) bool,
	side model.Side, instrumentListings map[int32]*model.Listing) {
	listingIdToQnt := map[int32]*model.Decimal64{}
	for _, line := range oppositeClobLines {
		if om.ParentOrder.GetAvailableQty().GreaterThan(zero) && willTrade(line) {
			quantity := line.Size

			if line.Size.GreaterThanOrEqual(om.ParentOrder.GetAvailableQty()) {
				quantity = om.ParentOrder.GetAvailableQty()
			}

			listing :=  instrumentListings[line.ListingId]
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
