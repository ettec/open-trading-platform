package actor

import (
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/model"
	"log"
)

type QuoteDistributor interface {
	Start() Actor
	Send(quote *model.ClobQuote)
	AddConnection(sink IdentifiableQuoteSink)
	RemoveConnection(sink IdentifiableQuoteSink)
	Close(chan<- bool)
}

type IdentifiableQuoteSink interface {
	ClobQuoteSink
	GetId() string
}

type quoteDistributor struct {
	sinks          []IdentifiableQuoteSink
	addConnChan    chan IdentifiableQuoteSink
	removeConnChan chan IdentifiableQuoteSink
	inQuoteChan    chan *model.ClobQuote
	closeChan      chan chan<- bool
}

func NewQuoteDistributor() *quoteDistributor  {
	return &quoteDistributor{sinks: make([]IdentifiableQuoteSink, 0),
		addConnChan:    make(chan IdentifiableQuoteSink, 100),
		removeConnChan: make(chan IdentifiableQuoteSink, 100),
		inQuoteChan:    make(chan *model.ClobQuote, 1000),
		closeChan:      make(chan chan<- bool, 1)}
}


func (q *quoteDistributor) Start() Actor {

	go func() {
		for {
			if d := q.readInputChannels(); d != nil {
				log.Println("closing quote distributor")
				d<-true
				return
			}
		}
	}()

	return q
}

func (q *quoteDistributor) Close(d chan<- bool) {
	q.closeChan <- d
}

func (q *quoteDistributor) Send(quote *model.ClobQuote) {
	q.inQuoteChan <- quote
}

func (q *quoteDistributor) AddConnection(sink IdentifiableQuoteSink) {
	q.addConnChan <- sink
}

func (q *quoteDistributor) RemoveConnection(sink IdentifiableQuoteSink) {
	q.removeConnChan <- sink
}

func (q *quoteDistributor) readInputChannels() chan<- bool {
	select {
	case s := <-q.addConnChan:
		q.sinks = append(q.sinks, s)
	case s := <-q.removeConnChan:
		for idx, o := range q.sinks {
			if o.GetId() == s.GetId() {
				q.sinks = append(q.sinks[:idx], q.sinks[idx+1:]...)
				break
			}
		}
	case cq := <-q.inQuoteChan:
		for _, sink := range q.sinks {
			sink.Send(cq)
		}
	case d := <-q.closeChan:
		return d
	}

	return nil
}

type SubscriptionHandler interface {
	Close()
	Subscribe(listingId int)
	Start()
}

