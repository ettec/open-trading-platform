package marketdata

import (
	"fmt"

	"github.com/ettec/open-trading-platform/go/model"
	"log"
	"os"
)

type ConflatedQuoteConnection interface {
	GetId() string
	Subscribe(listingId int32) error
	Close()
}

type conflatedQuoteConnection struct {
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

func (c *conflatedQuoteConnection) Close() {
	c.quoteDistributor.RemoveOutQuoteChan(c.distToConflatorChan)
	c.quoteConflator.Close()
}

func (c *conflatedQuoteConnection) Subscribe(listingId int32) error {

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

func (c *conflatedQuoteConnection) GetId() string {
	return c.id
}

func NewConflatedQuoteConnection(id string, out chan<- *model.ClobQuote, quoteDistributor QuoteDistributor,
	maxSubscriptions int) ConflatedQuoteConnection {

	distToConflatorChan := make(chan *model.ClobQuote, 200)

	quoteDistributor.AddOutQuoteChan(distToConflatorChan)

	conflator := NewQuoteConflator(distToConflatorChan, out, maxSubscriptions)

	c := &conflatedQuoteConnection{id: id,
		maxSubscriptions:    maxSubscriptions,
		quoteDistributor:    quoteDistributor,
		quoteConflator:      conflator,
		out:                 out,
		subscriptions:       map[int32]bool{},
		distToConflatorChan: distToConflatorChan,
		log:                 log.New(os.Stdout, "conflatedQuoteConnection:"+id, log.LstdFlags),
		errLog:              log.New(os.Stderr, "conflatedQuoteConnection:"+id, log.LstdFlags)}

	return c
}
