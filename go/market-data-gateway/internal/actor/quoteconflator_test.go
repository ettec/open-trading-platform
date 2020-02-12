package actor

import (
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/model"
	"reflect"
	"testing"
)

func Test_quotesAreConflated(t *testing.T) {

	in := make(chan *model.ClobQuote)
	out := make(chan *model.ClobQuote)
	close := make(chan bool)

	NewQuoteConflator(in, out, close)

	in <- &model.ClobQuote{ListingId: 1, XXX_sizecache:        1,}
	in <- &model.ClobQuote{ListingId: 1, XXX_sizecache:        2,}
	in <- &model.ClobQuote{ListingId: 1, XXX_sizecache:        3,}

	in <- &model.ClobQuote{ListingId: 2, XXX_sizecache:        6,}
	in <- &model.ClobQuote{ListingId: 2, XXX_sizecache:        7,}

	q := <-out
	if q.XXX_sizecache != 3 {
		t.Fatalf("expected last sent quote")
	}

}


func Test_circularBuffer(t *testing.T) {

	b := newBoundedCircularIntBuffer(4)

	in := []int32{1,2,3,4}

	allAdded := true
	for _, val := range in {
		allAdded = b.addHead(val) && allAdded
	}

	if !allAdded {
		t.Errorf("expected all values to be added")
	}

	var out []int32

	i, ok := b.removeTail()
	for ok {
		out = append(out,i)
		i, ok = b.removeTail()
	}

	if ! reflect.DeepEqual(in, out) {
		t.Errorf("expected in to equal out")
	}


}

func Test_circularBufferReadOverCapacity(t *testing.T) {

	b := newBoundedCircularIntBuffer(4)

	in := []int32{1,2,3,4}

	for _, val := range in {
		b.addHead(val)
	}

	i, ok := b.removeTail()
	if i != 1 ||!ok {
		t.FailNow()
	}

	i, ok = b.removeTail()
	if i != 2 ||!ok {
		t.FailNow()
	}

	if !b.addHead(5) {
		t.FailNow()
	}

	if !b.addHead(6) {
		t.FailNow()
	}

	var out []int32

	i, ok = b.removeTail()
	for ok {
		out = append(out,i)
		i, ok = b.removeTail()
	}

	expected := []int32{3,4,5,6}
	if ! reflect.DeepEqual(expected, out) {
		t.Fatalf("expected out %v to equal %v",out,expected)
	}

}