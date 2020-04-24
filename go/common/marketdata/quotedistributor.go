package marketdata

import (
	"github.com/ettec/open-trading-platform/go/model"
	"log"
	"os"
)

type QuoteDistributor interface {
  GetNewQuoteStream() MdsQuoteStream
}

type QuoteSource interface {
	Subscribe(listingId int)
}

type subscribeRequest struct {
	listingId int32
	out       chan<- *model.ClobQuote
}

type distConnection struct {
	out           chan<- *model.ClobQuote
	subscriptions map[int32]bool
}

type subscribeToListing = func(listingId int32)

type quoteDistributor struct {
	connections         []*distConnection
	addOutChan          chan chan<- *model.ClobQuote
	removeOutChan       chan chan<- *model.ClobQuote
	subscriptionChan    chan subscribeRequest
	lastQuote           map[int32]*model.ClobQuote
	subscribedFn        subscribeToListing
	subscribedToListing map[int32]bool
	sendBufferSize int
	log                 *log.Logger
	errLog              *log.Logger
}

func NewQuoteDistributor(stream MdsQuoteStream, sendBufferSize int)  *quoteDistributor {
	q := &quoteDistributor{connections: make([]*distConnection, 0),
		addOutChan:          make(chan chan<- *model.ClobQuote),
		removeOutChan:       make(chan chan<- *model.ClobQuote),
		subscriptionChan:    make(chan subscribeRequest),
		lastQuote:           map[int32]*model.ClobQuote{},
		subscribedFn:        stream.Subscribe,
		subscribedToListing: map[int32]bool{},
		sendBufferSize: sendBufferSize,
		log:                 log.New(os.Stdout, "", log.Ltime|log.Lshortfile),
		errLog:              log.New(os.Stderr, "", log.Ltime|log.Lshortfile),
	}

	streamChan := stream.GetStream()

	go func() {

		for {
			select {
			case s := <-q.subscriptionChan:

				if conn, exists := q.getConnection(s.out); exists {
					conn.subscriptions[s.listingId] = true
					if quote, ok := q.lastQuote[s.listingId]; ok {
						conn.out <- quote
					}

				} else {
					q.errLog.Printf("failed to subscribe to listing id %v as no connection exists", s.listingId)
				}

				if !q.subscribedToListing[s.listingId] {
					q.subscribedToListing[s.listingId] = true
					go q.subscribedFn(s.listingId)
				}
			case cq := <-streamChan:
				q.lastQuote[cq.ListingId] = cq
				for _, subscription := range q.connections {

					if subscription.subscriptions[cq.ListingId] {
						subscription.out <- cq
					}
				}
			case s := <-q.addOutChan:
				q.connections = append(q.connections, &distConnection{
					out:           s,
					subscriptions: map[int32]bool{},
				})
			case s := <-q.removeOutChan:
				if conn, exists := q.getConnection(s); exists {
					q.removeConnection(conn)
				} else {
					q.errLog.Println("no matching connection exists, connection not removed")
				}
			}
		}

	}()

	return q
}

func (q *quoteDistributor) Subscribe(listingId int32, out chan<- *model.ClobQuote) {

	q.subscriptionChan <- subscribeRequest{
		listingId: listingId,
		out:       out,
	}
}

type quoteDistributorQuoteStream struct {
	out chan *model.ClobQuote
	distributor *quoteDistributor
}

func newQuoteDistributorQuoteStream(distributor *quoteDistributor) *quoteDistributorQuoteStream {
	result := &quoteDistributorQuoteStream{make(chan *model.ClobQuote, distributor.sendBufferSize),
		distributor}
	distributor.addOutQuoteChan(result.out)
	return result
}

func (q *quoteDistributorQuoteStream) Subscribe(listingId int32) {
	q.distributor.Subscribe(listingId, q.out)
}

func (q *quoteDistributorQuoteStream) GetStream() <-chan *model.ClobQuote {
	return q.out
}

func (q *quoteDistributorQuoteStream) Close() {
	q.distributor.removeOutQuoteChan(q.out)
}

func (q *quoteDistributor) GetNewQuoteStream() MdsQuoteStream {
	return newQuoteDistributorQuoteStream(q)
}


func (q *quoteDistributor) addOutQuoteChan(out chan<- *model.ClobQuote) {
	q.addOutChan <- out
}

func (q *quoteDistributor) removeOutQuoteChan(out chan<- *model.ClobQuote) {
	q.removeOutChan <- out
}

func (q *quoteDistributor) getConnection(out chan<- *model.ClobQuote) (*distConnection, bool) {
	for _, o := range q.connections {
		if o.out == out {
			return o, true
		}
	}

	return nil, false
}

func (q *quoteDistributor) removeConnection(s *distConnection) {
	for idx, o := range q.connections {
		if o.out == s.out {
			q.connections = append(q.connections[:idx], q.connections[idx+1:]...)
			break
		}
	}
}
