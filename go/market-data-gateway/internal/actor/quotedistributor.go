package actor

import (
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/model"
	"log"
	"os"
)

type QuoteDistributor interface {
	Subscribe(listingId int)

	AddOutQuoteChan(sink chan<- *model.ClobQuote)

	RemoveOutQuoteChan(sink chan<- *model.ClobQuote)
}

type QuoteSource interface {
	Connect(out chan<- *model.ClobQuote) error
	Subscribe(listingId int)
}

type quoteDistributor struct {
	actorImpl
	outQuoteChans []chan<- *model.ClobQuote
	addOutChan    chan chan<- *model.ClobQuote
	removeOutChan chan chan<- *model.ClobQuote
	inQuoteChan   chan *model.ClobQuote
	quoteSource   QuoteSource
	log           *log.Logger
}

func NewQuoteDistributor(quoteSource QuoteSource) *quoteDistributor {
	q := &quoteDistributor{outQuoteChans: make([]chan<- *model.ClobQuote, 0),
		addOutChan:    make(chan chan<- *model.ClobQuote, 100),
		removeOutChan: make(chan chan<- *model.ClobQuote, 100),
		quoteSource:   quoteSource,
		log:           log.New(os.Stdout, "quoteDistributor", log.LstdFlags),
	}

	q.inQuoteChan = make(chan *model.ClobQuote, 1000)
	q.quoteSource.Connect(q.inQuoteChan)

	q.actorImpl = newActorImpl("quoteDistributor", q.readInputChannels)
	return q
}

func (q *quoteDistributor) Subscribe(listingId int) {
	q.quoteSource.Subscribe(listingId)
}

func (q *quoteDistributor) AddOutQuoteChan(sink chan<- *model.ClobQuote) {
	q.addOutChan <- sink
}

func (q *quoteDistributor) RemoveOutQuoteChan(sink chan<- *model.ClobQuote) {
	q.removeOutChan <- sink
}

func (q *quoteDistributor) readInputChannels() (chan<- bool, error) {
	select {
	case s := <-q.addOutChan:
		q.outQuoteChans = append(q.outQuoteChans, s)
	case s := <-q.removeOutChan:
		q.removeOutChannel(s)
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
				q.removeOutChannel(sink)
			}
		}

	case d := <-q.closeChan:
		for _, out := range q.outQuoteChans {
			close(out)
		}
		return d, nil
	}

	return nil, nil
}

func (q *quoteDistributor) removeOutChannel(s chan<- *model.ClobQuote) {
	for idx, o := range q.outQuoteChans {
		if o == s {
			q.outQuoteChans = append(q.outQuoteChans[:idx], q.outQuoteChans[idx+1:]...)
			break
		}
	}
}
