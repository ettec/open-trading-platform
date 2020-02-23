package actor

import (
	"github.com/ettec/open-trading-platform/go/model"
	"log"
	"os"
)

type QuoteDistributor interface {
	Subscribe(listingId int32)

	AddOutQuoteChan(sink chan<- *model.ClobQuote)

	RemoveOutQuoteChan(sink chan<- *model.ClobQuote)
}

type QuoteSource interface {
	Subscribe(listingId int)
}

type quoteDistributor struct {
	outQuoteChans []chan<- *model.ClobQuote
	addOutChan    chan chan<- *model.ClobQuote
	removeOutChan chan chan<- *model.ClobQuote
	subscribedFn  subscribeToListing
	log           *log.Logger
	errLog        *log.Logger
}

func NewQuoteDistributor(subscribedFn subscribeToListing, in <-chan *model.ClobQuote) *quoteDistributor {
	q := &quoteDistributor{outQuoteChans: make([]chan<- *model.ClobQuote, 0),
		addOutChan:    make(chan chan<- *model.ClobQuote),
		removeOutChan: make(chan chan<- *model.ClobQuote),
		subscribedFn:   subscribedFn,
		log:           log.New(os.Stdout, "", log.Ltime | log.Lshortfile),
		errLog:           log.New(os.Stderr, "", log.Ltime | log.Lshortfile),
	}

	go func() {

		for {
			select {
			case s := <-q.addOutChan:
				q.outQuoteChans = append(q.outQuoteChans, s)
			case s := <-q.removeOutChan:
				q.removeOutChannel(s)
			case cq := <-in:
				var toRemove []chan<- *model.ClobQuote
				for _, quoteChan := range q.outQuoteChans {

					select {
					case quoteChan <- cq:
					default:
						q.errLog.Printf("removing outbound quote channel %v as it is full", quoteChan)
						close(quoteChan)
						toRemove = append(toRemove, quoteChan)
					}
				}

				if toRemove != nil {
					for _, sink := range toRemove {
						q.removeOutChannel(sink)
					}
				}
			}
		}

	}()

	return q
}

func (q *quoteDistributor) Subscribe(listingId int32) {
	q.subscribedFn(listingId)
}

func (q *quoteDistributor) AddOutQuoteChan(out chan<- *model.ClobQuote) {
	q.addOutChan <- out
}

func (q *quoteDistributor) RemoveOutQuoteChan(out chan<- *model.ClobQuote) {
	q.removeOutChan <- out
}

func (q *quoteDistributor) removeOutChannel(s chan<- *model.ClobQuote) {
	for idx, o := range q.outQuoteChans {
		if o == s {
			q.outQuoteChans = append(q.outQuoteChans[:idx], q.outQuoteChans[idx+1:]...)
			break
		}
	}
}
