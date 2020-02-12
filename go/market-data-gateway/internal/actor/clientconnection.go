package actor

import (
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/model"
	"log"
	"os"
)

type ClientConnection interface {
	Actor
	GetId() string
	Subscribe(listingId int)
}

type sendQuoteFn = func(quote *model.ClobQuote) error


type clientConnection struct {
	actorImpl
	id                string
	sendQuoteFn       sendQuoteFn
	quotesInChan      chan *model.ClobQuote
	quoteDistributor  QuoteDistributor
	subscribeChan     chan int
	closeChan         chan chan<- bool
	subscribeListings map[int32]bool
	maxPendingQuotes  int
	log               *log.Logger
}

func (c *clientConnection) Subscribe(listingId int) {
	c.subscribeChan <- listingId
	c.log.Println("subscribed to listing id:", listingId)
}

func (c *clientConnection) GetId() string {
	return c.id
}



func NewClientConnection(id string,  sendQuoteFn sendQuoteFn,
	quoteDistributor QuoteDistributor,
	maxPendingQuotes int) *clientConnection {

	cc := &clientConnection{id: id, sendQuoteFn: sendQuoteFn,
		quotesInChan: make(chan *model.ClobQuote, maxPendingQuotes), subscribeChan: make(chan int),
		quoteDistributor: quoteDistributor,
		subscribeListings: map[int32]bool{},
		maxPendingQuotes:   maxPendingQuotes,
		log:               log.New(os.Stdout, "clientConnection:"+id, log.LstdFlags)}

	cc.actorImpl = newActorImpl("connection:"+id, cc.readInputChannels)
	cc.quoteDistributor.AddOutQuoteChan(cc.quotesInChan)

	log.Println("the id of the inbound quotes channel for the connection is:", cc.quotesInChan)

	return cc
}

func (c *clientConnection) readInputChannels() (chan<- bool, error) {

	select {
	case q := <-c.quotesInChan:

			if c.subscribeListings[q.ListingId] {
				if err := c.sendQuoteFn(q); err != nil {
					log.Printf(" closing as error occurred whilst sending quote:%w", err)
					c.closeChan <- make(chan bool, 1)
				}
			}

	case l := <-c.subscribeChan:
		c.subscribeListings[int32(l)] = true
		c.quoteDistributor.Subscribe(l)
	case d := <-c.closeChan:
		defer c.quoteDistributor.RemoveOutQuoteChan(c.quotesInChan)
		return d, nil
	}

	return nil, nil
}


