package actor

import (
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/model"
	"testing"
)

func Test_clientConnection_Send(t *testing.T) {

}

type testClobQuoteSink struct {
	send  func(quote *model.ClobQuote) error
	close func() error
}

func (s *testClobQuoteSink) Send(quote *model.ClobQuote) error {
	return s.send(quote)
}

func (s *testClobQuoteSink) Close() error {
	return s.close()
}

func Test_clientConnection_Subscribe(t *testing.T) {

	out := make(chan *model.ClobQuote, 100)

	quoteSink := &testClobQuoteSink{send: func(quote *model.ClobQuote) error {
		out <- quote
		return nil
	}}

	c := NewClientConnection("testId", func(listingId int) {

	},
		quoteSink, 100)

	c.Subscribe(1)
	c.Subscribe(2)

	c.quotesInChan <- &model.ClobQuote{ListingId: 4}
	c.quotesInChan <- &model.ClobQuote{ListingId: 1}
	c.quotesInChan <- &model.ClobQuote{ListingId: 3}
	c.quotesInChan <- &model.ClobQuote{ListingId: 2}

	if q := <-out; q.ListingId != 1 {
		t.Errorf("expected quote with listing id 1")
	}
	if q := <-out; q.ListingId != 2 {
		t.Errorf("expected quote with listing id 2")
	}

	select {
	case <-out:
		t.Errorf("no more quotes expected")
	default:
	}

}

func Test_quoteSinkClosedWhenInboundChannelClosed(t *testing.T) {



	closed := false
	quoteSink := &testClobQuoteSink{close: func() error {
		closed = true
		return nil
	}}


	c := NewClientConnection("testId", func(listingId int) {

	},
		quoteSink, 100)

	close(c.quotesInChan)

	c.readInputChannels()

	if ! closed {
		t.Error("expected quote sink to be closed")
	}

}
