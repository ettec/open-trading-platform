package main

import (
	"log"
	"os"
)

type quoteNormaliser struct {
	symbolToListingId map[string]int
	idToQuote         map[int]*fullQuote
	inboundChan       <-chan mdupdate
	outboundChan      chan<- *snapshot
	closeChan         chan bool
	log               *log.Logger
}

func newQuoteNormaliser(inboundChan <-chan mdupdate,
	outboundChan chan<- *snapshot) *quoteNormaliser {

	return &quoteNormaliser{
		symbolToListingId: make(map[string]int),
		idToQuote:         make(map[int]*fullQuote),
		inboundChan:       inboundChan,
		outboundChan:      outboundChan,
		closeChan:         make(chan bool, 1),
		log:               log.New(os.Stdout, "quoteNormalise:", log.LstdFlags),
	}
}

func (n *quoteNormaliser) close() {
	n.closeChan <- true
}

func (n *quoteNormaliser) processUpdates() {

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
			break
		}
	}

}
