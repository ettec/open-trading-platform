package quoteaggregator

import (
	"github.com/ettec/open-trading-platform/go/common/marketdata"
	"github.com/ettec/open-trading-platform/go/model"
	"log"
)

const QuoteAggregatorMic  = "XOSR"

type quoteAggregator struct {
	getListings     getListingsWithSameInstrument
	listingGroupsIn chan []model.Listing
}

type getListingsWithSameInstrument = func(listingId int32, listingGroupsIn chan<- []model.Listing)

func New(id string, getListingsWithSameInstrument getListingsWithSameInstrument, micToMdsAddress map[string]string,
	out chan<- *model.ClobQuote, mdsClientFn marketdata.GetMdsClientFn) *quoteAggregator {

	q := &quoteAggregator{
		getListings:     getListingsWithSameInstrument,
		listingGroupsIn: make(chan []model.Listing, 1000),
	}

	quoteStreamsOut := make(chan *model.ClobQuote, 1000)
	micToStream := map[string]marketdata.MdsQuoteStream{}

	for mic, targetAddress := range micToMdsAddress {
		stream, err := marketdata.NewMdsQuoteStream(id, targetAddress, quoteStreamsOut, mdsClientFn)
		if err != nil {
			log.Panicf("failed to created quote stream for mic %v, targetAddress %v. Error:%v", mic, targetAddress, err)
		}
		micToStream[mic] = stream
	}

	listingIdToQuoteChan := map[int32]chan<- *model.ClobQuote{}

	go func() {
		select {
		case listings := <-q.listingGroupsIn:
			quoteChan := make(chan *model.ClobQuote)
			numStreams := 0
			for _, listing := range listings {
				if listing.Market.Mic != QuoteAggregatorMic {
					listingIdToQuoteChan[listing.Id] = quoteChan
					if stream, ok := micToStream[listing.Market.Mic]; ok {
						stream.Subscribe(listing.Id)
						numStreams++
					} else {
						log.Printf("no quote stream available for mic %v, instrument %v ", listing.Market.Mic, listing.Instrument.DisplaySymbol)
					}
				}
			}

			go func() {
				listingIdToLastQuote := map[int32]*model.ClobQuote{}
				quotes := make([]*model.ClobQuote, 0, numStreams)
				select {
				case q := <-quoteChan:
					listingIdToLastQuote[q.ListingId] = q

					quotes = quotes[:0]
					idx := 0
					for _, q := range listingIdToLastQuote {
						quotes[idx] = q
						idx++
					}
					out <- combineQuotes(quotes)
				}

			}()

		case q := <-quoteStreamsOut:
			listingIdToQuoteChan[q.ListingId] <- q
		}

	}()

	return q
}

func combineQuotes(quotes []*model.ClobQuote) *model.ClobQuote {
  here impl this and test
}

func (q quoteAggregator) getListingsWithSameInstrument(listingId int32) {
	q.getListings(listingId, q.listingGroupsIn)
}

func (q quoteAggregator) Subscribe(listingId int32) {
	q.getListingsWithSameInstrument(listingId)
}
