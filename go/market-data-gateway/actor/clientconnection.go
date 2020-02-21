package actor

import (
	"github.com/ettec/open-trading-platform/go/model"
	"log"
	"os"
)

type ClientConnection interface {
	GetId() string
	Subscribe(listingId int)
	Close()
}

type clientConnection struct {
	id            string
	subscribeChan chan int
	closeChan     chan  bool
	log           *log.Logger
	errLog        *log.Logger
}

func (c *clientConnection) Close() {
	c.closeChan <- true
}

func (c *clientConnection) Subscribe(listingId int) {
	c.subscribeChan <- listingId
	c.log.Println("subscribed to listing id:", listingId)
}

func (c *clientConnection) GetId() string {
	return c.id
}

func NewClientConnection(id string, out chan<- *model.ClobQuote , subscribe subscribeToListing, in  <-chan *model.ClobQuote,
	maxSubscriptions int) *clientConnection {

	c := &clientConnection{id: id,
		closeChan:     make(chan  bool, 1),
		subscribeChan: make(chan int),
		log:           log.New(os.Stdout, "clientConnection:"+id, log.LstdFlags),
		errLog:        log.New(os.Stderr, "clientConnection:"+id, log.LstdFlags)}


	toConflator := make(chan *model.ClobQuote)

	conflator := NewQuoteConflator(toConflator, out, maxSubscriptions)


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
			case <-c.closeChan:
				conflator.Close()
				return
			}
		}

	}()

	return c
}
