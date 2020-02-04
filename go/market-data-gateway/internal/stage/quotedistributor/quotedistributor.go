package quotedistributor

import (
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/model"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/stage"
)

type IdentifiableQuoteSink interface {
	stage.ClobQuoteSink
	GetId() string
}

type quoteDistributor struct {
	sinks          []IdentifiableQuoteSink
	addConnChan    chan IdentifiableQuoteSink
	removeConnChan chan IdentifiableQuoteSink
	inQuoteChan    chan *model.ClobQuote
}

func newQuoteDistributor() *quoteDistributor {
	return &quoteDistributor{sinks: make([]IdentifiableQuoteSink, 0),
		addConnChan:    make(chan IdentifiableQuoteSink, 100),
		removeConnChan: make(chan IdentifiableQuoteSink, 100),
		inQuoteChan:    make(chan *model.ClobQuote, 1000)}
}

func (q *quoteDistributor) start() {

	go func() {
		for {
			q.readInputChannels()
		}
	}()

}

func (q *quoteDistributor) addConnection(sink IdentifiableQuoteSink) {
	q.addConnChan <- sink
}

func (q *quoteDistributor) removeConnection(sink IdentifiableQuoteSink) {
	q.removeConnChan <- sink
}

func (q *quoteDistributor) Send(quote *model.ClobQuote) {
	q.inQuoteChan<-quote
}

func (q *quoteDistributor) readInputChannels() {
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

	}
}
