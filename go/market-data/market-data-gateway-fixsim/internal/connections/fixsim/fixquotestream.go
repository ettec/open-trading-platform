package fixsim

import (
	"context"
	"fmt"
	"github.com/ettec/open-trading-platform/go/market-data/market-data-gateway-fixsim/internal/fix/marketdata"
	"github.com/ettec/otp-common/model"
	"github.com/ettec/otp-common/staticdata"
	"github.com/golang/protobuf/proto"
	"log/slog"
)

type GetListingFn func(ctx context.Context, listingId int32, resultChan chan<- staticdata.ListingResult)

type FixQuoteStream struct {
	ctx                  context.Context
	connectionName       string
	getListingResultChan chan staticdata.ListingResult
	out                  chan *model.ClobQuote
	fixMarketDataClient  fixMarketDataClient
	cancelCtx            func()
	getListing           GetListingFn
}

func (n *FixQuoteStream) Chan() <-chan *model.ClobQuote {
	return n.out
}

func (n *FixQuoteStream) Close() {
	n.cancelCtx()
}

func (n *FixQuoteStream) Subscribe(listingId int32) error {
	n.getListing(n.ctx, listingId, n.getListingResultChan)
	return nil
}

type fixMarketDataClient interface {
	Subscribe(symbol string) error
	Chan() <-chan *marketdata.MarketDataIncrementalRefresh
}

func NewQuoteStreamFromFixClient(parentCtx context.Context,
	fixMarketDataClient fixMarketDataClient, connectionName string, symbolLookup GetListingFn,
	sendBufferSize int) (*FixQuoteStream, error) {

	ctx, cancel := context.WithCancel(parentCtx)
	out := make(chan *model.ClobQuote, sendBufferSize)

	n := &FixQuoteStream{
		ctx:                  ctx,
		connectionName:       connectionName,
		out:                  out,
		getListingResultChan: make(chan staticdata.ListingResult, 1000),
		fixMarketDataClient:  fixMarketDataClient,
		cancelCtx:            cancel,
		getListing:           symbolLookup,
	}

	log := slog.With(slog.Default(), "connectionName", connectionName)
	symbolToListingId := make(map[string]int32)
	idToQuote := map[int32]*model.ClobQuote{}

	go func() {
		defer close(out)

		for {
			select {
			case <-ctx.Done():
				break
			case lr := <-n.getListingResultChan:

				if lr.Err != nil {
					log.Error("failed to get listing", "error", lr.Err)
					continue
				}
				symbolToListingId[lr.Listing.MarketSymbol] = lr.Listing.Id
				go func() {
					if err := n.fixMarketDataClient.Subscribe(lr.Listing.MarketSymbol); err != nil {
						log.Error("failed to subscribe", "error", err)
					}
				}()
			case r, ok := <-n.fixMarketDataClient.Chan():
				if !ok {
					log.Warn("fix sim client closed")
					return
				}

				if r != nil {
					for _, incGrp := range r.MdIncGrp {
						symbol := incGrp.GetInstrument().GetSymbol()
						if listingId, ok := symbolToListingId[symbol]; ok {

							var originalQuote *model.ClobQuote
							if originalQuote, ok = idToQuote[listingId]; !ok {
								originalQuote = newClobQuote(listingId)
							}

							newQuote, err := copyQuote(originalQuote)
							newQuote.StreamInterrupted = false
							newQuote.StreamStatusMsg = ""
							if err != nil {
								log.Error("failed to copy originalQuote", "error", err)
							}

							linesUpdate := incGrp.MdEntryType == marketdata.MDEntryTypeEnum_MD_ENTRY_TYPE_BID ||
								incGrp.MdEntryType == marketdata.MDEntryTypeEnum_MD_ENTRY_TYPE_OFFER
							if linesUpdate {

								bids := incGrp.MdEntryType == marketdata.MDEntryTypeEnum_MD_ENTRY_TYPE_BID
								updateQuoteDepth(originalQuote, newQuote, incGrp, bids)

							} else if incGrp.MdEntryType == marketdata.MDEntryTypeEnum_MD_ENTRY_TYPE_TRADE {
								newQuote.LastPrice = &model.Decimal64{Mantissa: incGrp.MdEntryPx.Mantissa, Exponent: incGrp.MdEntryPx.Exponent}
								newQuote.LastQuantity = &model.Decimal64{Mantissa: incGrp.MdEntrySize.Mantissa, Exponent: incGrp.MdEntrySize.Exponent}

							} else if incGrp.MdEntryType == marketdata.MDEntryTypeEnum_MD_ENTRY_TYPE_TRADE_VOLUME {
								newQuote.TradedVolume = &model.Decimal64{Mantissa: incGrp.MdEntrySize.Mantissa, Exponent: incGrp.MdEntrySize.Exponent}
							} else {
								continue
							}

							idToQuote[listingId] = newQuote
							n.out <- newQuote

						} else {
							log.Warn("received refresh for unknown symbol", "symbol", symbol)
						}

					}
				} else {
					for id := range idToQuote {
						emptyQuote := newClobQuote(id)
						emptyQuote.StreamInterrupted = true
						emptyQuote.StreamStatusMsg = "fix sim adapter stream interrupted"
						idToQuote[id] = emptyQuote
						n.out <- emptyQuote
					}
				}
			}
		}
	}()

	return n, nil
}

func copyQuote(quote *model.ClobQuote) (*model.ClobQuote, error) {

	bytes, err := proto.Marshal(quote)
	if err != nil {
		return nil, fmt.Errorf("failed to copy quote:%w", err)
	}

	quoteCopy := &model.ClobQuote{}
	err = proto.Unmarshal(bytes, quoteCopy)
	if err != nil {
		return nil, fmt.Errorf("failed to copy quote:%w", err)
	}

	return quoteCopy, nil
}

func newClobQuote(listingId int32) *model.ClobQuote {
	bids := make([]*model.ClobLine, 0)
	offers := make([]*model.ClobLine, 0)

	return &model.ClobQuote{
		ListingId: listingId,
		Bids:      bids,
		Offers:    offers,
	}
}

func updateQuoteDepth(originalQuote *model.ClobQuote, newQuote *model.ClobQuote, update *marketdata.MDIncGrp, bidSide bool) {
	if bidSide {
		newQuote.Bids = updateClobLines(originalQuote.Bids, update, bidSide)
	} else {
		newQuote.Offers = updateClobLines(originalQuote.Offers, update, bidSide)
	}
}

func updateClobLines(lines []*model.ClobLine, update *marketdata.MDIncGrp, bids bool) []*model.ClobLine {

	updateAction := update.GetMdUpdateAction()
	newClobLines := make([]*model.ClobLine, 0, len(lines)+1)

	compareTest := 1
	if bids {
		compareTest = -1
	}

	switch updateAction {
	case marketdata.MDUpdateActionEnum_MD_UPDATE_ACTION_NEW:
		inserted := false

		newLine := &model.ClobLine{
			Size:    &model.Decimal64{Mantissa: update.MdEntrySize.Mantissa, Exponent: update.MdEntrySize.Exponent},
			Price:   &model.Decimal64{Mantissa: update.MdEntryPx.Mantissa, Exponent: update.MdEntryPx.Exponent},
			EntryId: update.MdEntryId,
		}

		for _, line := range lines {
			compareResult := model.Compare(*line.Price, model.Decimal64(*update.GetMdEntryPx()))
			if !inserted && compareResult == compareTest {
				newClobLines = append(newClobLines, newLine)
				inserted = true
			}
			newClobLines = append(newClobLines, line)
		}

		if !inserted {
			newClobLines = append(newClobLines, newLine)
		}

	case marketdata.MDUpdateActionEnum_MD_UPDATE_ACTION_CHANGE:
		inserted := false

		newLine := &model.ClobLine{
			Size:    &model.Decimal64{Mantissa: update.MdEntrySize.Mantissa, Exponent: update.MdEntrySize.Exponent},
			Price:   &model.Decimal64{Mantissa: update.MdEntryPx.Mantissa, Exponent: update.MdEntryPx.Exponent},
			EntryId: update.MdEntryId,
		}

		for _, line := range lines {
			compareResult := model.Compare(*line.Price, model.Decimal64(*update.GetMdEntryPx()))
			if !inserted && compareResult == compareTest {
				newClobLines = append(newClobLines, newLine)
				inserted = true
			}
			if line.EntryId != newLine.EntryId {
				newClobLines = append(newClobLines, line)
			}

		}

		if !inserted {
			newClobLines = append(newClobLines, newLine)
		}

	case marketdata.MDUpdateActionEnum_MD_UPDATE_ACTION_DELETE:
		for _, line := range lines {
			if line.EntryId != update.MdEntryId {
				newClobLines = append(newClobLines, line)
			}
		}
	}

	return newClobLines

}
