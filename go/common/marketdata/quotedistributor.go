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

type quoteDistributorQuoteStream struct {
	out         chan *model.ClobQuote
	distributor *quoteDistributor
}

func newQuoteDistributorQuoteStream(distributor *quoteDistributor) *quoteDistributorQuoteStream {
	result := &quoteDistributorQuoteStream{make(chan *model.ClobQuote, distributor.sendBufferSize),
		distributor}
	return result
}

func (q *quoteDistributorQuoteStream) Subscribe(listingId int32) {

	q.distributor.subscriptionChan <- subscribeRequest{
		listingId: listingId,
		out:       q.out,
	}

}

func (q *quoteDistributorQuoteStream) GetStream() <-chan *model.ClobQuote {
	return q.out
}

func (q *quoteDistributorQuoteStream) Close() {
	q.distributor.removeOutChan <- q.out
}

type subscribeToListing = func(listingId int32)

type quoteDistributor struct {
	listingToStreams    map[int32][]chan<- *model.ClobQuote
	streamToListings    map[chan<- *model.ClobQuote][]int32
	removeOutChan       chan chan<- *model.ClobQuote
	subscriptionChan    chan subscribeRequest
	lastQuote           map[int32]*model.ClobQuote
	subscribedFn        subscribeToListing
	subscribedToListing map[int32]bool
	sendBufferSize      int
	log                 *log.Logger
	errLog              *log.Logger
}

func NewQuoteDistributor(stream MdsQuoteStream, sendBufferSize int) *quoteDistributor {
	q := &quoteDistributor{
		listingToStreams:    map[int32][]chan<- *model.ClobQuote{},
		streamToListings:    map[chan<- *model.ClobQuote][]int32{},
		removeOutChan:       make(chan chan<- *model.ClobQuote),
		subscriptionChan:    make(chan subscribeRequest),
		lastQuote:           map[int32]*model.ClobQuote{},
		subscribedFn:        stream.Subscribe,
		subscribedToListing: map[int32]bool{},
		sendBufferSize:      sendBufferSize,
		log:                 log.New(os.Stdout, "", log.Ltime|log.Lshortfile),
		errLog:              log.New(os.Stderr, "", log.Ltime|log.Lshortfile),
	}

	streamChan := stream.GetStream()

	go func() {

		for {
			select {
			case s := <-q.subscriptionChan:

				subscribedToQuotes := q.listingToStreams[s.listingId] == nil
				if !subscribedToQuotes {
					go q.subscribedFn(s.listingId)
				}

				streamSubscribed := false
				for _, stream := range q.listingToStreams[s.listingId] {
					if stream == s.out {
						streamSubscribed = true
						break
					}
				}

				if !streamSubscribed {
					q.listingToStreams[s.listingId] = append(q.listingToStreams[s.listingId], s.out)
					q.streamToListings[s.out] = append(q.streamToListings[s.out], s.listingId)
					if lastQuote, exists := q.lastQuote[s.listingId]; exists {
						s.out <- lastQuote
					}
				}

			case cq := <-streamChan:
				q.lastQuote[cq.ListingId] = cq

				for _, stream := range q.listingToStreams[cq.ListingId] {
					stream <- cq
				}
			case s := <-q.removeOutChan:
				subscribedListings := q.streamToListings[s]
				for _, listingId := range subscribedListings {
					for idx, o := range q.listingToStreams[listingId] {
						if o == s {
							q.listingToStreams[listingId] = append(q.listingToStreams[listingId][:idx], q.listingToStreams[listingId][idx+1:]...)
							break
						}
					}
				}

				delete(q.streamToListings, s)
				close(s)
			}
		}

	}()

	return q
}

func (q *quoteDistributor) GetNewQuoteStream() MdsQuoteStream {
	return newQuoteDistributorQuoteStream(q)
}
