package actor

import (
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/model"
	"testing"
)

type testQuoteDistributor struct {

}

func (t *testQuoteDistributor) Subscribe(listingId int) {

}

func (t *testQuoteDistributor) AddOutQuoteChan(sink chan<- *model.ClobQuote) {

}

func (t *testQuoteDistributor) RemoveOutQuoteChan(sink chan<- *model.ClobQuote) {

}



func Test_clientConnection_Send(t *testing.T) {

}



func Test_clientConnection_Subscribe(t *testing.T) {

	out := make(chan *model.ClobQuote, 100)



	c := NewClientConnection("testId", func(quote *model.ClobQuote) error {
		out <- quote
		return nil
	},
		&testQuoteDistributor{}, 100)

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

