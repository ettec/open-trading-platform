package actor

import (
	"fmt"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/model"
	"log"
	"os"
)

type ClientConnection interface {
	Actor
	GetId() string
	Send(q *model.ClobQuote)
	Subscribe(listingId int)
}

type clobQuoteSink interface {
	Send (quote *model.ClobQuote) error
	Close() error
}

type clientConnection struct {
	actorImpl
	id                string
	subscribeFn       subscribeToListing
	clobQuoteSink       clobQuoteSink
	quotesInChan      chan *model.ClobQuote
	subscribeChan     chan int
	closeChan         chan chan<- bool
	subscribeListings map[int32]bool
	maxPendingQuotes  int
	log               *log.Logger
}

func (c *clientConnection) Subscribe(listingId int) {
	c.subscribeChan <- listingId
	c.subscribeFn(listingId)
	c.log.Println("subscribed to listing id:", listingId)
}

func (c *clientConnection) GetId() string {
	return c.id
}



func NewClientConnection(id string, subscribeFn subscribeToListing, clobQuoteSink clobQuoteSink,
	maxPendingQuotes int) *clientConnection {

	cc := &clientConnection{id: id, subscribeFn: subscribeFn, clobQuoteSink: clobQuoteSink,
		quotesInChan: make(chan *model.ClobQuote, maxPendingQuotes), subscribeChan: make(chan int),
		subscribeListings: map[int32]bool{},
		maxPendingQuotes:   maxPendingQuotes,
		log:               log.New(os.Stdout, "clientConnection:"+id, log.LstdFlags)}

	cc.actorImpl = newActorImpl("connection:"+id, cc.readInputChannels)

	log.Println("the id of the inbound quotes channel for the connection is:", cc.quotesInChan)

	return cc
}

func (c *clientConnection) readInputChannels() (chan<- bool, error) {

	select {
	case q, ok := <-c.quotesInChan:
		if ok {
			if c.subscribeListings[q.ListingId] {
				if err := c.clobQuoteSink.Send(q); err != nil {
					return nil, fmt.Errorf("error occurred whilst sending quote:%w", err)
				}
			}
		} else {
			log.Printf("closing as inbound quote channel %v is closed", c.id)
			c.closeChan <- make(chan bool, 1)
		}

	case l := <-c.subscribeChan:
		c.subscribeListings[int32(l)] = true
	case d := <-c.closeChan:
		defer c.clobQuoteSink.Close()
		return d, nil
	}

	return nil, nil
}


