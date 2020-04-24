package marketdata

import (
	"github.com/ettec/open-trading-platform/go/model"
	"testing"
	"time"
)

func Test_quoteDistributor_Send(t *testing.T) {

	in := make(chan *model.ClobQuote, 10)

	d := NewQuoteDistributor(testMdsQuoteStream{
		func(listingId int32) {
		}, in}, 100)

	s1 := d.GetNewQuoteStream()

	s2 := d.GetNewQuoteStream()

	s1.Subscribe(1)
	s2.Subscribe(1)

	in <- &model.ClobQuote{ListingId: 1}

	if q := <-s1.GetStream(); q.ListingId != 1 {
		t.Errorf("expected quote not received")
	}

	if q := <-s2.GetStream(); q.ListingId != 1 {
		t.Errorf("expected quote not received")
	}

}

func Test_subscriptionReceivesLastSentQuote(t *testing.T) {

	in := make(chan *model.ClobQuote, 10)

	d := NewQuoteDistributor(testMdsQuoteStream{
		func(listingId int32) {
		}, in}, 100)

	s1 := d.GetNewQuoteStream()
	s2 := d.GetNewQuoteStream()

	s1.Subscribe(1)

	in <- &model.ClobQuote{ListingId: 1}

	if q := <-s1.GetStream(); q.ListingId != 1 {
		t.Errorf("expected quote note received")
	}

	s2.Subscribe(1)

	if q := <-s2.GetStream(); q.ListingId != 1 {
		t.Errorf("expected quote note received")
	}

}

func Test_subscribeOnlyCalledOnceForAGivenListing(t *testing.T) {

	in := make(chan *model.ClobQuote)

	subscribeCalls := make(chan int32, 10)
	d := NewQuoteDistributor(testMdsQuoteStream{
		func(listingId int32) {
			subscribeCalls <- listingId
		}, in}, 100)

	s1 := d.GetNewQuoteStream()
	s2 := d.GetNewQuoteStream()

	s1.Subscribe(1)
	s2.Subscribe(1)

	time.Sleep(2 * time.Second)

	if len(subscribeCalls) != 1 {
		t.FailNow()
	}

}

func Test_onlySubscribedQuotesReceived(t *testing.T) {

	in := make(chan *model.ClobQuote)

	d := NewQuoteDistributor(testMdsQuoteStream{
		func(listingId int32) {
		}, in}, 100)

	s1 := d.GetNewQuoteStream()

	s1.Subscribe(1)
	s1.Subscribe(2)

	in <- &model.ClobQuote{ListingId: 3}
	in <- &model.ClobQuote{ListingId: 1}
	in <- &model.ClobQuote{ListingId: 4}
	in <- &model.ClobQuote{ListingId: 2}

	q := <-s1.GetStream()
	if q.ListingId != 1 {
		t.Errorf("unexpected quote")
	}

	q = <-s1.GetStream()
	if q.ListingId != 2 {
		t.Errorf("unexpected quote")
	}

}
