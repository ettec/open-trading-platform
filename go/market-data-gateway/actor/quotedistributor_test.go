package actor

import (
	"github.com/ettec/open-trading-platform/go/model"
	"testing"
)


func Test_quoteDistributor_Send(t *testing.T) {

	in := make(chan *model.ClobQuote, 10)

	d := NewQuoteDistributor(func(listingId int) {}, in)

	s1 := make(chan *model.ClobQuote, 100)

	s2 := make(chan *model.ClobQuote, 100)

	d.AddOutQuoteChan(s1)
	d.AddOutQuoteChan(s2)

	in <- &model.ClobQuote{ListingId: 1}

	if q := <-s1; q.ListingId != 1 {
		t.Errorf("expected quote note received")
	}

	if q := <-s2; q.ListingId != 1 {
		t.Errorf("expected quote note received")
	}

}

func Test_quoteDistributorRemovesFullChan(t *testing.T) {

	in := make(chan *model.ClobQuote)

	d := NewQuoteDistributor(func(listingId int) {}, in)

	s1 := make(chan *model.ClobQuote, 2)

	s2 := make(chan *model.ClobQuote, 100)

	d.AddOutQuoteChan(s1)
	d.AddOutQuoteChan(s2)



	in <- &model.ClobQuote{ListingId: 1}
	in <- &model.ClobQuote{ListingId: 2}
	in <- &model.ClobQuote{ListingId: 3}

	<-s1
	<-s1

	if _, ok := <-s1; ok {
		t.Errorf("expected 3rd message to be empty and the channel to be closed")
	}

	// check all messages still sent to second channel
	<-s2
	<-s2
	<-s2

}
