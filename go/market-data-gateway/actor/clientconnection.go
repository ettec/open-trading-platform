package actor

import (
	"fmt"
	"github.com/ettec/open-trading-platform/go/model"
	"log"
	"os"
)

type ClientConnection interface {
	GetId() string
	Subscribe(listingId int32) error
	Close()
}

type clientConnection struct {
	id                  string
	maxSubscriptions    int
	quoteDistributor    QuoteDistributor
	quoteConflator      *quoteConflator
	distToConflatorChan chan *model.ClobQuote
	out                 chan<- *model.ClobQuote
	subscriptions       map[int32]bool
	log                 *log.Logger
	errLog              *log.Logger
}

func (c *clientConnection) Close() {

	c.quoteDistributor.RemoveOutQuoteChan(c.distToConflatorChan)
	c.quoteConflator.Close()
	close(c.out)

}

func (c *clientConnection) Subscribe(listingId int32) error {

	if len(c.subscriptions) == c.maxSubscriptions {
		return fmt.Errorf("max number of subscriptions for this connection has been reached: %v", c.maxSubscriptions)
	}

	if c.subscriptions[listingId] {
		return fmt.Errorf("already subscribed to listing id: %v", listingId)
	}

	c.quoteDistributor.Subscribe(listingId, c.distToConflatorChan)
	c.log.Println("subscribed to listing id:", listingId)

	return nil
}

func (c *clientConnection) GetId() string {
	return c.id
}

type subscribeToListing = func(listingId int32)

func NewClientConnection(id string, out chan<- *model.ClobQuote, quoteDistributor QuoteDistributor,
	maxSubscriptions int) *clientConnection {

	distToConflatorChan := make(chan *model.ClobQuote, 200)

	quoteDistributor.AddOutQuoteChan(distToConflatorChan)

	conflator := NewQuoteConflator(distToConflatorChan, out, maxSubscriptions)

	c := &clientConnection{id: id,
		maxSubscriptions:    maxSubscriptions,
		quoteDistributor:    quoteDistributor,
		quoteConflator:      conflator,
		out:				 out,
		subscriptions:       map[int32]bool{},
		distToConflatorChan: distToConflatorChan,
		log:                 log.New(os.Stdout, "clientConnection:"+id, log.LstdFlags),
		errLog:              log.New(os.Stderr, "clientConnection:"+id, log.LstdFlags)}

	return c
}
