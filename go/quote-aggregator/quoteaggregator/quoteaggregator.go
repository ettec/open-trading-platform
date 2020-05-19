package quoteaggregator

import (
	"github.com/ettec/open-trading-platform/go/common/marketdata"
	"github.com/ettec/open-trading-platform/go/common/marketdata/quotestream"
	"github.com/ettec/open-trading-platform/go/model"
	"log"
)

const QuoteAggregatorMic = "XOSR"

type quoteAggregator struct {
	getListings     getListingsWithSameInstrument
	listingGroupsIn chan []*model.Listing
	stream          chan *model.ClobQuote
	closeChan       chan bool
}

func (q quoteAggregator) GetStream() <-chan *model.ClobQuote {
	return q.stream
}

func (q quoteAggregator) Close() {
	q.closeChan <- true
}

type getListingsWithSameInstrument = func(listingId int32, listingGroupsIn chan<- []*model.Listing)

func New(id string, getListingsWithSameInstrument getListingsWithSameInstrument, micToMdsAddress map[string]string,
	bufferSize int, mdsClientFn quotestream.GetMdsClientFn) *quoteAggregator {

	qa := &quoteAggregator{
		getListings:     getListingsWithSameInstrument,
		listingGroupsIn: make(chan []*model.Listing, 1000),
		stream:          make(chan *model.ClobQuote),
		closeChan:       make(chan bool),
	}

	micToStream := map[string]marketdata.MdsQuoteStream{}

	quoteStreamsOut := make(chan *model.ClobQuote, bufferSize)

	for mic, targetAddress := range micToMdsAddress {
		stream, err := quotestream.NewMdsQuoteStreamFromFn(id, targetAddress, quoteStreamsOut, mdsClientFn)
		if err != nil {
			log.Panicf("failed to created quote stream for mic %v, targetAddress %v. Error:%v", mic, targetAddress, err)
		}
		micToStream[mic] = stream
	}

	listingIdToQuoteChan := map[int32]chan<- *model.ClobQuote{}

	go func() {
		for {
			select {
			case <-qa.closeChan:
				break
			case listings := <-qa.listingGroupsIn:
				quoteChan := make(chan *model.ClobQuote)
				numStreams := 0
				var quoteAggListingId int32 = -1
				for _, listing := range listings {
					if listing.Market.Mic != QuoteAggregatorMic {
						listingIdToQuoteChan[listing.Id] = quoteChan
						if stream, ok := micToStream[listing.Market.Mic]; ok {
							stream.Subscribe(listing.Id)
							numStreams++
						} else {
							log.Printf("no quote stream available for mic %v, instrument %v ", listing.Market.Mic, listing.Instrument.DisplaySymbol)
						}
					} else {
						quoteAggListingId = listing.Id
					}
				}

				go func() {
					listingIdToLastQuote := map[int32]*model.ClobQuote{}
					quotes := make([]*model.ClobQuote, 0, numStreams)
					for {
						select {
						case q := <-quoteChan:
							listingIdToLastQuote[q.ListingId] = q
							quotes = quotes[:0]
							for _, q := range listingIdToLastQuote {
								quotes = append(quotes, q)
							}
							qa.stream <- combineQuotes(quoteAggListingId, quotes)
						}
					}
				}()

			case q := <-quoteStreamsOut:
				listingIdToQuoteChan[q.ListingId] <- q
			}
		}
	}()

	return qa
}

func combineQuotes(combinedListingId int32, quotes []*model.ClobQuote) *model.ClobQuote {

	bids := getCombinedLines(quotes,
		func(quote *model.ClobQuote) []*model.ClobLine {
			return quote.Bids
		}, func(a *model.Decimal64, b *model.Decimal64) bool {
			return a.GreaterThan(b)
		})

	offers := getCombinedLines(quotes,
		func(quote *model.ClobQuote) []*model.ClobLine {
			return quote.Offers
		}, func(a *model.Decimal64, b *model.Decimal64) bool {
			return a.LessThan(b)
		})

	streamInterrupted := false
	streamStatusMsg := ""

	for _, quote := range quotes {
		if !streamInterrupted {
			streamInterrupted = quote.StreamInterrupted
		}
		if quote.StreamStatusMsg != "" {
			streamStatusMsg = quote.StreamStatusMsg
		}
	}

	quote := &model.ClobQuote{
		ListingId:         combinedListingId,
		Bids:              bids,
		Offers:            offers,
		StreamInterrupted: streamInterrupted,
		StreamStatusMsg:   streamStatusMsg,
	}

	return quote
}

func getCombinedLines(quotes []*model.ClobQuote, getQuoteLines func(quote *model.ClobQuote) []*model.ClobLine, compare func(a *model.Decimal64, b *model.Decimal64) bool) []*model.ClobLine {
	result := []*model.ClobLine{}
	levelIdxs := make([]int, len(quotes), len(quotes))

	for {
		var bestPrice *model.Decimal64 = nil
		var bestSize *model.Decimal64 = nil
		var bestListingId int32
		bestQuoteIdx := 0

		for quoteIdx, quote := range quotes {
			lines := getQuoteLines(quote)
			if levelIdxs[quoteIdx] < len(lines) {
				price := lines[levelIdxs[quoteIdx]].Price
				size := lines[levelIdxs[quoteIdx]].Size
				listingId := quote.ListingId

				if bestPrice == nil || compare(price, bestPrice) {
					bestPrice = price
					bestSize = size
					bestQuoteIdx = quoteIdx
					bestListingId = listingId
				}
			}
		}

		if bestPrice != nil {
			levelIdxs[bestQuoteIdx] = levelIdxs[bestQuoteIdx] + 1
			result = append(result, &model.ClobLine{Price: bestPrice, Size: bestSize, ListingId: bestListingId})
		} else {
			break
		}

	}
	return result
}

func (q quoteAggregator) getListingsWithSameInstrument(listingId int32) {
	q.getListings(listingId, q.listingGroupsIn)
}

func (q quoteAggregator) Subscribe(listingId int32) {
	q.getListingsWithSameInstrument(listingId)
}
