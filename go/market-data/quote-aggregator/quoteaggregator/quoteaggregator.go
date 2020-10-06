package quoteaggregator

import (
	common "github.com/ettec/otp-common"
	"github.com/ettec/otp-common/marketdata"
	"github.com/ettec/otp-common/model"
)

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

func New(getListingsWithSameInstrument getListingsWithSameInstrument, stream marketdata.MdsQuoteStream,
	inboundListingsBufferSize int) *quoteAggregator {

	qa := &quoteAggregator{
		getListings:     getListingsWithSameInstrument,
		listingGroupsIn: make(chan []*model.Listing, inboundListingsBufferSize),
		stream:          make(chan *model.ClobQuote),
		closeChan:       make(chan bool),
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
					if listing.Market.Mic != common.SR_MIC {
						listingIdToQuoteChan[listing.Id] = quoteChan
						stream.Subscribe(listing.Id)
						numStreams++
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
							qa.stream <- combineQuotes(quoteAggListingId, quotes,q)
						}
					}
				}()

			case q := <-stream.GetStream():
				listingIdToQuoteChan[q.ListingId] <- q
			}
		}
	}()

	return qa
}

func combineQuotes(combinedListingId int32, quotes []*model.ClobQuote, lastQuote *model.ClobQuote) *model.ClobQuote {

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

	tradedVolume := &model.Decimal64{}
	for _, quote := range quotes {
		tradedVolume.Add(quote.TradedVolume)
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
		LastPrice: lastQuote.LastPrice,
		LastQuantity: lastQuote.LastQuantity,
		TradedVolume: tradedVolume,
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
