package main

import (
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/fix/marketdata"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/model"
	"log"
	"os"
)

type quoteNormaliser struct {
	symbolToListingId map[string]int
	idToQuote         map[int]*model.ClobQuote
	inboundChan       <-chan mdupdate
	outboundChan      chan<- *model.ClobQuote
	closeChan         chan bool
	log               *log.Logger
}

func newQuoteNormaliser(inboundChan <-chan mdupdate,
	outboundChan chan<- *model.ClobQuote) *quoteNormaliser {

	q := &quoteNormaliser{
		symbolToListingId: make(map[string]int),
		idToQuote:         make(map[int]*model.ClobQuote),
		inboundChan:       inboundChan,
		outboundChan:      outboundChan,
		closeChan:         make(chan bool, 1),
		log:               log.New(os.Stdout, "quoteNormaliser:", log.LstdFlags),
	}
	go q.processUpdates()

	return q
}

func (n *quoteNormaliser) close() {
	n.closeChan <- true
}

func (n *quoteNormaliser) processUpdates() {

Loop:
	for {
		select {
		case u := <-n.inboundChan:
			if u.listingIdToSymbol != nil {
				n.symbolToListingId[u.listingIdToSymbol.symbol] = u.listingIdToSymbol.listingId
			}

			if u.refresh != nil {
				for _, incGrp := range u.refresh.MdIncGrp {
					symbol := incGrp.GetInstrument().GetSymbol()
					if listingId, ok := n.symbolToListingId[symbol]; ok {

						if fullQuote, ok := n.idToQuote[listingId]; ok {
							n.outboundChan <- fullQuote.onIncRefresh(incGrp)
						} else {
							newQuote := newFullQuote(listingId)
							n.idToQuote[listingId] = newQuote
							n.outboundChan <- newQuote.onIncRefresh(incGrp)
						}
					} else {
						n.log.Println("no listing found for symbol:", symbol)
					}
				}
			}
		case <-n.closeChan:
			break Loop
		}
	}
}

func updateQuote(quote model.ClobQuote, update marketdata.MDIncGrp) model.ClobQuote {

	id := update.GetMdEntryId()
	updateAction := update.GetMdUpdateAction()


	bids := make([]*model.ClobLine, len(quote.Bids) )

	switch updateAction {
	case marketdata.MDUpdateActionEnum_MD_UPDATE_ACTION_NEW:
		for idx, line := range quote.Bids {
			 compareResult := model.Compare(*line.Price, model.Decimal64(*update.GetMdEntryPx()))
		}


	case marketdata.MDUpdateActionEnum_MD_UPDATE_ACTION_CHANGE:
		fullGrp := marketdata.MDFullGrp{
			MdEntryPx:   inc.GetMdEntryPx(),
			MdEntrySize: inc.GetMdEntrySize(),
			MdEntryId:   inc.GetMdEntryId(),
			MdEntryType: inc.GetMdEntryType(),
		}
		q.entryIdToEntry[id] = &fullGrp
	case marketdata.MDUpdateActionEnum_MD_UPDATE_ACTION_DELETE:
		delete(q.entryIdToEntry, id)
	}

}





func (q *fullQuote) onIncRefresh(inc *marketdata.MDIncGrp) *snapshot {

	id := inc.GetMdEntryId()
	updateAction := inc.GetMdUpdateAction()

	switch updateAction {
	case marketdata.MDUpdateActionEnum_MD_UPDATE_ACTION_NEW:
		fallthrough
	case marketdata.MDUpdateActionEnum_MD_UPDATE_ACTION_CHANGE:
		fullGrp := marketdata.MDFullGrp{
			MdEntryPx:   inc.GetMdEntryPx(),
			MdEntrySize: inc.GetMdEntrySize(),
			MdEntryId:   inc.GetMdEntryId(),
			MdEntryType: inc.GetMdEntryType(),
		}
		q.entryIdToEntry[id] = &fullGrp
	case marketdata.MDUpdateActionEnum_MD_UPDATE_ACTION_DELETE:
		delete(q.entryIdToEntry, id)
	}

	entries := make([]*marketdata.MDFullGrp, len(q.entryIdToEntry))
	idx := 0
	for _, value := range q.entryIdToEntry {
		entries[idx] = value
		idx++
	}

	return &snapshot{
		Instrument: q.instrument,
		MdFullGrp:  entries,
	}
}


