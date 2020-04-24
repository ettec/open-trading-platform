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
	id               string
	maxSubscriptions int
	stream           MdsQuoteStream
	subscriptions    map[int32]bool
	log              *log.Logger
	errLog           *log.Logger
}

func (c *conflatedQuoteConnection) GetStream() <-chan *model.ClobQuote {
	return c.stream.GetStream()
}

func (c *conflatedQuoteConnection) Close() {
	c.stream.Close()
}

func (c *conflatedQuoteConnection) Subscribe(listingId int32) error {

	if len(c.subscriptions) == c.maxSubscriptions {
		return fmt.Errorf("max number of subscriptions for this connection has been reached: %v", c.maxSubscriptions)
	}

	if c.subscriptions[listingId] {
		return fmt.Errorf("already subscribed to listing id: %v", listingId)
	}

	c.stream.Subscribe(listingId)
	c.log.Println("subscribed to listing id:", listingId)

	return nil
}

func (c *conflatedQuoteConnection) GetId() string {
	return c.id
}

func NewConflatedQuoteConnection(id string, stream MdsQuoteStream,
	maxSubscriptions int) *conflatedQuoteConnection {

	conflatedStream := NewConflatedQuoteStream(stream, maxSubscriptions)

	c := &conflatedQuoteConnection{id: id,
		maxSubscriptions: maxSubscriptions,
		subscriptions:    map[int32]bool{},
		stream:           conflatedStream,
		log:              log.New(os.Stdout, "conflatedQuoteConnection:"+id, log.LstdFlags),
		errLog:           log.New(os.Stderr, "conflatedQuoteConnection:"+id, log.LstdFlags)}

	return c
}
