package actor

import (
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/model"
	"log"
	"os"
)

type ClientConnection interface {
	GetId() string
	Subscribe(listingId int)
	Close()
}

type sendQuoteFn = func(quote *model.ClobQuote) error

type clientConnection struct {
	id            string
	subscribeChan chan int
	closeFilterChan     chan  bool
	closeToClientChan     chan  bool
	log           *log.Logger
	errLog        *log.Logger
}

func (c *clientConnection) Close() {
	c.closeFilterChan <- true
	c.closeToClientChan <-true
}

func (c *clientConnection) Subscribe(listingId int) {
	c.subscribeChan <- listingId
	c.log.Println("subscribed to listing id:", listingId)
}

func (c *clientConnection) GetId() string {
	return c.id
}

func NewClientConnection(id string, sendQuoteFn sendQuoteFn, subscribe subscribeToListing, in  <-chan *model.ClobQuote,
	maxSubscriptions int) *clientConnection {

	c := &clientConnection{id: id,
		closeFilterChan: make(chan  bool, 1),
		closeToClientChan: make(chan  bool, 1),
		subscribeChan: make(chan int),
		log:       log.New(os.Stdout, "clientConnection:"+id, log.LstdFlags),
		errLog:    log.New(os.Stderr, "clientConnection:"+id, log.LstdFlags)}


	toConflator := make(chan *model.ClobQuote)
	toClient := make(chan *model.ClobQuote)

	conflator := NewQuoteConflator(toConflator, toClient, maxSubscriptions)


	subscribedListings := map[int32]bool{}

	go func() {
		for {
			select {
			case q := <-in:
				if subscribedListings[q.ListingId] {
					toConflator <- q
				}
			case l := <-c.subscribeChan:
				subscribedListings[int32(l)] = true
				subscribe(l)
			case <-c.closeFilterChan:
				conflator.Close()
				return
			}
		}

	}()


	go func() {
		for {
			select {
			case q := <-toClient:
					if err := sendQuoteFn(q); err != nil {
						c.errLog.Printf(" closing as error occurred whilst sending quote:%v", err)
						c.Close()
					}
			case <-c.closeToClientChan:
				return
			}
		}

	}()

	return c
}
