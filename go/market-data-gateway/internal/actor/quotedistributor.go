package actor

import (
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/model"
	"log"
	"os"
)

type quoteDistributor struct {
	actorImpl
	outQuoteChans  []chan<- *model.ClobQuote
	addConnChan    chan chan<- *model.ClobQuote
	removeConnChan chan chan<- *model.ClobQuote
	inQuoteChan    chan *model.ClobQuote
	log            *log.Logger
}

func NewQuoteDistributor() *quoteDistributor {
	q := &quoteDistributor{outQuoteChans: make([]chan<- *model.ClobQuote, 0),
		addConnChan:    make(chan chan<- *model.ClobQuote, 100),
		removeConnChan: make(chan chan<- *model.ClobQuote, 100),
		inQuoteChan:    make(chan *model.ClobQuote, 1000),
		log:            log.New(os.Stdout, "quoteDistributor", log.LstdFlags),
	}

	q.actorImpl = newActorImpl("quoteDistributor", q.readInputChannels)
	return q
}

func (q *quoteDistributor) Send(quote *model.ClobQuote) {
	q.inQuoteChan <- quote
}

func (q *quoteDistributor) AddConnection(sink chan<- *model.ClobQuote) {
	q.addConnChan <- sink
}

func (q *quoteDistributor) RemoveConnection(sink chan<- *model.ClobQuote) {
	q.removeConnChan <- sink
}

func (q *quoteDistributor) readInputChannels() (chan<- bool, error) {
	select {
	case s := <-q.addConnChan:
		q.outQuoteChans = append(q.outQuoteChans, s)
	case s := <-q.removeConnChan:
		q.removeSink(s)
	case cq := <-q.inQuoteChan:
		var toRemove []chan<- *model.ClobQuote
		for _, quoteChan := range q.outQuoteChans {

			select {
			case quoteChan <- cq:
			default:
				log.Printf("removing outbound quote channel %v as it is full", quoteChan)
				close(quoteChan)
				toRemove = append(toRemove, quoteChan)
			}
		}

		if toRemove != nil {
			for _, sink := range toRemove {
				q.removeSink(sink)
			}
		}

	case d := <-q.closeChan:
		return d, nil
	}

	return nil, nil
}

func (q *quoteDistributor) removeSink(s chan<- *model.ClobQuote) {
	for idx, o := range q.outQuoteChans {
		if o == s {
			q.outQuoteChans = append(q.outQuoteChans[:idx], q.outQuoteChans[idx+1:]...)
			break
		}
	}
}
