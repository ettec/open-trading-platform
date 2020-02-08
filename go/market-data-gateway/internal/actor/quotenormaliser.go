package actor

import (
	"fmt"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/fix/marketdata"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/model"
	"log"
	"os"
)

type clobQuoteNormaliser struct {
	actorImpl
	symbolToListingId map[string]int
	idToQuote         map[int]*model.ClobQuote
	refreshChan       chan *marketdata.MarketDataIncrementalRefresh
	mappingChan       chan ListingIdSymbol
	out               clobQuoteSink
	log               *log.Logger
}

func NewClobQuoteNormaliser(
	out clobQuoteSink) *clobQuoteNormaliser {

	q := &clobQuoteNormaliser{
		symbolToListingId: make(map[string]int),
		idToQuote:         make(map[int]*model.ClobQuote),
		refreshChan:       make(chan *marketdata.MarketDataIncrementalRefresh, 10000),
		mappingChan:       make(chan ListingIdSymbol),
		out:               out,
		log:               log.New(os.Stdout, "clobQuoteNormaliser:", log.LstdFlags),
	}

	q.actorImpl = newActorImpl("quoteNormaliser", q.readInputChannel)

	return q
}

func (n *clobQuoteNormaliser) SendRefresh(refresh *marketdata.MarketDataIncrementalRefresh) {
	n.refreshChan <- refresh
}

func (n *clobQuoteNormaliser) registerMapping(lts ListingIdSymbol) {
	n.mappingChan <- lts
}

func (n *clobQuoteNormaliser) readInputChannel() (chan<- bool, error) {
	select {
	case m := <-n.mappingChan:
		n.symbolToListingId[m.Symbol] = m.ListingId
	case r := <-n.refreshChan:

		for _, incGrp := range r.MdIncGrp {
			symbol := incGrp.GetInstrument().GetSymbol()
			bids := incGrp.MdEntryType == marketdata.MDEntryTypeEnum_MD_ENTRY_TYPE_BID
			if listingId, ok := n.symbolToListingId[symbol]; ok {
				if quote, ok := n.idToQuote[listingId]; ok {
					updatedQuote := updateQuote(quote, incGrp, bids)
					n.idToQuote[listingId] = updatedQuote
					n.out.Send(updatedQuote)
				} else {
					quote := newClobQuote(listingId)
					updatedQuote := updateQuote(quote, incGrp, bids)
					n.idToQuote[listingId] = updatedQuote
					n.out.Send(updatedQuote)
				}
			} else {
				return nil, fmt.Errorf("no listing found for symbol: %v", symbol)
			}
		}
	case d := <-n.closeChan:
		return d, nil
	}

	return nil, nil
}

func newClobQuote(listingId int) *model.ClobQuote {
	bids := make([]*model.ClobLine, 0)
	offers := make([]*model.ClobLine, 0)

	return &model.ClobQuote{
		ListingId: int32(listingId),
		Bids:      bids,
		Offers:    offers,
	}
}

func updateQuote(quote *model.ClobQuote, update *marketdata.MDIncGrp, bidSide bool) *model.ClobQuote {

	newQuote := model.ClobQuote{
		ListingId: quote.ListingId,
	}

	if bidSide {
		newQuote.Offers = quote.Offers
		newQuote.Bids = updateClobLines(quote.Bids, update, bidSide)
	} else {
		newQuote.Bids = quote.Bids
		newQuote.Offers = updateClobLines(quote.Offers, update, bidSide)
	}

	return &newQuote
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
