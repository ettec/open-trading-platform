package actor

import (
	"github.com/ettec/open-trading-platform/go/model"
	"testing"
)

func Test_clientConnection_Subscribe(t *testing.T) {

	out := make(chan *model.ClobQuote, 100)
	in := make(chan *model.ClobQuote, 100)



	c := NewClientConnection("testId", func(quote *model.ClobQuote) error {
		out <- quote
		return nil
	},
		func(listingId int) {

		}, in, 100)

	c.Subscribe(1)
	c.Subscribe(2)

	in <- &model.ClobQuote{ListingId: 4}
	in <- &model.ClobQuote{ListingId: 1}
	in <- &model.ClobQuote{ListingId: 3}
	in <- &model.ClobQuote{ListingId: 2}

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

func Test_slowConnectionDoesNotBlockDownstreamSender(t *testing.T) {

	out := make(chan *model.ClobQuote)
	in := make(chan *model.ClobQuote)


	c := NewClientConnection("testId", func(quote *model.ClobQuote) error {
		out <- quote
		return nil
	},
		func(listingId int) {

		}, in, 100)

	c.Subscribe(1)
	c.Subscribe(2)

	for i:=0; i<2000; i++ {
		in <- &model.ClobQuote{ListingId: 1, XXX_sizecache:int32(i)}
		in <- &model.ClobQuote{ListingId: 2,  XXX_sizecache:int32(i)}
	}


	if q := <-out; q.ListingId != 1 && q.XXX_sizecache != 1999{
		t.Errorf("expected quote with listing id 1")
	}
	if q := <-out; q.ListingId != 2 && q.XXX_sizecache != 1999 {
		t.Errorf("expected quote with listing id 2")
	}

}

