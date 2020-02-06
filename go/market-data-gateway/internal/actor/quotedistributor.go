package actor

import (
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/model"
)

type QuoteDistributor interface {
	Actor
	Send(quote *model.ClobQuote)
	AddConnection(sink IdentifiableQuoteSink)
	RemoveConnection(sink IdentifiableQuoteSink)
}

type IdentifiableQuoteSink interface {
	ClobQuoteSink
	GetId() string
}

type quoteDistributor struct {
	actorImpl
	sinks          []IdentifiableQuoteSink
	addConnChan    chan IdentifiableQuoteSink
	removeConnChan chan IdentifiableQuoteSink
	inQuoteChan    chan *model.ClobQuote
}

func NewQuoteDistributor() *quoteDistributor {
	q := &quoteDistributor{sinks: make([]IdentifiableQuoteSink, 0),
		addConnChan:    make(chan IdentifiableQuoteSink, 100),
		removeConnChan: make(chan IdentifiableQuoteSink, 100),
		inQuoteChan:    make(chan *model.ClobQuote, 1000),
	}

	q.actorImpl = newActorImpl("clobDistributor", q.readInputChannels)
	return q
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

func (q *quoteDistributor) readInputChannels() (chan<- bool, error) {
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
		return d, nil
	}

	return nil, nil
}
