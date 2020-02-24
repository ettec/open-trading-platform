package actor

import (
	"github.com/ettec/open-trading-platform/go/model"
	"testing"
)


func Test_quoteDistributor_Send(t *testing.T) {

	in := make(chan *model.ClobQuote, 10)

	d := NewQuoteDistributor(func(listingId int32) {}, in)

	s1 := make(chan *model.ClobQuote, 100)

	s2 := make(chan *model.ClobQuote, 100)

	d.AddOutQuoteChan(s1)
	d.AddOutQuoteChan(s2)

	d.Subscribe(1, s1)
	d.Subscribe(1, s2)

	in <- &model.ClobQuote{ListingId: 1}

	if q := <-s1; q.ListingId != 1 {
		t.Errorf("expected quote not received")
	}

	if q := <-s2; q.ListingId != 1 {
		t.Errorf("expected quote not received")
	}

}


func Test_subscriptionReceivesLastSentQuote(t *testing.T) {


	in := make(chan *model.ClobQuote, 10)

	d := NewQuoteDistributor(func(listingId int32) {}, in)

	s1 := make(chan *model.ClobQuote, 100)

	s2 := make(chan *model.ClobQuote, 100)

	d.AddOutQuoteChan(s1)
	d.Subscribe(1,s1)

	in <- &model.ClobQuote{ListingId: 1}

	if q := <-s1; q.ListingId != 1 {
		t.Errorf("expected quote note received")
	}

	d.AddOutQuoteChan(s2)

	d.Subscribe(1, s2)

	if q := <-s2; q.ListingId != 1 {
		t.Errorf("expected quote note received")
	}

}

func Test_onlySubscribedQuotesReceived(t *testing.T) {

	in := make(chan *model.ClobQuote)

	d := NewQuoteDistributor(func(listingId int32) {}, in)

	s1 := make(chan *model.ClobQuote, 50)

	d.AddOutQuoteChan(s1)

	d.Subscribe(1,s1)
	d.Subscribe(2,s1)


	in <- &model.ClobQuote{ListingId: 3}
	in <- &model.ClobQuote{ListingId: 1}
	in <- &model.ClobQuote{ListingId: 4}
	in <- &model.ClobQuote{ListingId: 2}


	q := <-s1
	if q.ListingId != 1 {
		t.Errorf("unexpected quote")
	}

	q = <-s1
	if q.ListingId != 2 {
		t.Errorf("unexpected quote")
	}


}
