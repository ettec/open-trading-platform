package main

import (
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/api"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/fix/marketdata"
)

type quote struct {
	refresh   *marketdata.MarketDataSnapshotFullRefresh
	listingId int
}



func newConnection(stream api.MarketDataGateway_ConnectServer) *connection {

	c := &connection{make(chan *quote), stream, make(map[int]bool),
		make(chan bool, 1)}

	go func() {
		select {
		case q := <-c.QuoteChan:
			if c.subscriptions[q.listingId] {
				stream.Send(q.refresh)
			}
		case <-c.closeChan:
			return
		}
	}()

	return c
}

func (c *connection) subscribe(listingId int) {
	c.subscriptions[listingId] = true

}

func (c *connection) close() {
	c.closeChan <- true
}
