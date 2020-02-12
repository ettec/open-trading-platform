package actor

import (
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/model"
	"testing"
)

type testQuoteSource struct {
	out chan<- *model.ClobQuote
}

func (t *testQuoteSource) Connect(out chan<- *model.ClobQuote) error {
	t.out = out
	return nil
}

func (t *testQuoteSource)  Subscribe(listingId int) {

}


func Test_quoteDistributor_Send(t *testing.T) {

	tqs := &testQuoteSource{}
	d := NewQuoteDistributor(tqs)

	s1 := make( chan *model.ClobQuote, 100)

	s2 := make( chan *model.ClobQuote, 100)

	d.AddOutQuoteChan(s1)
	d.AddOutQuoteChan(s2)

	d.readInputChannels()
	d.readInputChannels()

	tqs.out <- &model.ClobQuote{ListingId:1}


	d.readInputChannels()

	if q:=<-s1; q.ListingId != 1 {
		t.Errorf("expected quote note received")
	}

	if q:=<-s2; q.ListingId != 1 {
		t.Errorf("expected quote note received")
	}

}


func Test_quoteDistributorRemovesFullChan(t *testing.T) {

	tqs := &testQuoteSource{}
	d := NewQuoteDistributor(tqs)

	s1 := make( chan *model.ClobQuote, 2)

	s2 := make( chan *model.ClobQuote, 100)

	d.AddOutQuoteChan(s1)
	d.AddOutQuoteChan(s2)

	d.readInputChannels()
	d.readInputChannels()

	tqs.out <- &model.ClobQuote{ListingId:1}
	tqs.out <- &model.ClobQuote{ListingId:2}
	tqs.out <- &model.ClobQuote{ListingId:3}

	d.readInputChannels()
	d.readInputChannels()
	d.readInputChannels()

	<-s1
	<-s1

	if _, ok :=<-s1; ok {
		t.Errorf("expected 3rd message to be empty and the channel to be closed")
	}

	// check all messages still sent to second channel
	<-s2
	<-s2
	<-s2

}