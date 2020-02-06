package actor

import (
	"fmt"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/model"
	"log"
	"os"
	"sync/atomic"
)

type ClientConnection interface {
	Actor
	GetId() string
	Send(q *model.ClobQuote)
	Subscribe(listingId int)
}

type clientConnection struct {
	actorImpl
	id              string
	subscriptions   SubscriptionHandler
	stream          model.MarketDataGateway_ConnectServer
	quotesInChan    chan *model.ClobQuote
	droppedQuoteCnt int32
	closeChan       chan chan<- bool
	log             *log.Logger
}

func (c *clientConnection) Subscribe(listingId int) {
	c.subscriptions.Subscribe(listingId)
	c.log.Println("subscribed to listing id:", listingId)
}

func (c *clientConnection) GetId() string {
	return c.id
}

func (c *clientConnection) Send(q *model.ClobQuote) {
	select {
	case c.quotesInChan <- q:
	default: // Must not block downstream actors under any circumstances, add conflation here in addition to this safety valve
		atomic.AddInt32(&c.droppedQuoteCnt, 1)
	}

}

func NewClientConnection(id string, subClient SubscriptionClient, stream model.MarketDataGateway_ConnectServer,
	clientConnBufferSize int) *clientConnection {
	ss := &testSymbolSource{
		mappings: map[int]string{1: "A", 2: "B", 3: "C", 4: "D"},
	}

	sh := NewSubscriptionHandler(id, ss, subClient)

	cc := &clientConnection{id:id, subscriptions: sh, stream: stream,
		quotesInChan: make(chan *model.ClobQuote, clientConnBufferSize), droppedQuoteCnt:0,
		log: log.New(os.Stdout, "clientConnection:"+id, log.LstdFlags)}

	cc.actorImpl = newActorImpl("connection:" + id, cc.readInputChannels)

	return cc
}

func (c *clientConnection) readInputChannels() (chan<- bool, error) {

	select {
	case q := <-c.quotesInChan:
		if err := c.stream.Send(q); err != nil {
			return nil, fmt.Errorf("error occurred whilst sending to grpc output stream:%w", err)
		}
	case d := <-c.closeChan:
		return d, nil
	}

	return nil, nil
}

type testSymbolSource struct {
	mappings map[int]string
}

func (t *testSymbolSource) FetchSymbol(listingId int, onSymbol chan<- ListingIdSymbol) {

	if sym, ok := t.mappings[listingId]; ok {
		onSymbol <- ListingIdSymbol{
			ListingId: listingId,
			Symbol:    sym,
		}
	} else {
		log.Println("no symbol mapping found for listing id ", listingId)
	}

}
