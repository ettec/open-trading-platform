package quoteaggregator

import (
	"context"
	common "github.com/ettec/otp-common"
	"github.com/ettec/otp-common/marketdata"
	"github.com/ettec/otp-common/model"
	"github.com/ettec/otp-common/staticdata"
	"log/slog"
)

type getListingsWithSameInstrument = func(ctx context.Context, listingId int32, resultChan chan<- staticdata.ListingsResult)

type quoteAggregator struct {
	ctx                           context.Context
	cancel                        context.CancelFunc
	getListingsWithSameInstrument getListingsWithSameInstrument
	listingGroupsIn               chan staticdata.ListingsResult
	outChan                       chan *model.ClobQuote
}

func (q *quoteAggregator) Chan() <-chan *model.ClobQuote {
	return q.outChan
}

func (q *quoteAggregator) Subscribe(listingId int32) error {
	q.getListingsWithSameInstrument(q.ctx, listingId, q.listingGroupsIn)
	return nil
}

func (q *quoteAggregator) Close() {
	q.cancel()
}

func New(ctx context.Context, getListingsWithSameInstrument getListingsWithSameInstrument, stream marketdata.QuoteStream,
	inboundListingsBufferSize int) *quoteAggregator {

	ctx, cancel := context.WithCancel(ctx)

	qa := &quoteAggregator{
		ctx:                           ctx,
		cancel:                        cancel,
		getListingsWithSameInstrument: getListingsWithSameInstrument,
		listingGroupsIn:               make(chan staticdata.ListingsResult, inboundListingsBufferSize),
		outChan:                       make(chan *model.ClobQuote),
	}

	go func() {
		listingIdToQuoteChan := map[int32]chan<- *model.ClobQuote{}
		subscribedListings := map[int32]bool{}

		for {
			select {
			case <-ctx.Done():
				return
			case q := <-stream.Chan():
				listingIdToQuoteChan[q.ListingId] <- q
			case listingsResult := <-qa.listingGroupsIn:
				if listingsResult.Err != nil {
					slog.Error("failed to get listings", "error", listingsResult.Err)
					continue
				}

				var quoteAggListingId int32 = -1
				for _, listing := range listingsResult.Listings {
					if listing.Market.Mic == common.SR_MIC {
						quoteAggListingId = listing.Id
						if _, ok := subscribedListings[listing.Id]; ok {
							slog.Warn("already subscribed to quote stream", "listingId", listing.Id)
							continue
						} else {
							subscribedListings[listing.Id] = true
						}
					}
				}

				quoteChan := make(chan *model.ClobQuote)
				numStreams := 0
				for _, listing := range listingsResult.Listings {
					if listing.Market.Mic != common.SR_MIC {
						listingIdToQuoteChan[listing.Id] = quoteChan
						if err := stream.Subscribe(listing.Id); err != nil {
							slog.Error("failed to subscribe to quote stream", "listingId", listing.Id, "error", err)
						}
						numStreams++
					}
				}

				go func() {
					listingIdToLastQuote := map[int32]*model.ClobQuote{}
					quotes := make([]*model.ClobQuote, 0, numStreams)
					for {
						select {
						case <-ctx.Done():
							return
						case q := <-quoteChan:
							listingIdToLastQuote[q.ListingId] = q
							quotes = quotes[:0]
							for _, q := range listingIdToLastQuote {
								quotes = append(quotes, q)
							}
							qa.outChan <- combineQuotes(quoteAggListingId, quotes, q)
						}
					}
				}()
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
		if !streamInterrupted && quote.StreamInterrupted {
			streamInterrupted = true
		}
		if quote.StreamStatusMsg != "" {
			streamStatusMsg += quote.StreamStatusMsg
		}
	}

	quote := &model.ClobQuote{
		ListingId:         combinedListingId,
		Bids:              bids,
		Offers:            offers,
		StreamInterrupted: streamInterrupted,
		StreamStatusMsg:   streamStatusMsg,
		LastPrice:         lastQuote.LastPrice,
		LastQuantity:      lastQuote.LastQuantity,
		TradedVolume:      tradedVolume,
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
