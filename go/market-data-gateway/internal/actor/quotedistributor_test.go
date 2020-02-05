package actor

import (
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/model"
	"testing"
)


type testClobSink struct {
	received []*model.ClobQuote
	id string
}


func (s *testClobSink) Send(quote *model.ClobQuote) {
	s.received = append(s.received, quote)
}

func (s *testClobSink) GetId() string {
	return s.id
}

func Test_quoteDistributor_Send(t *testing.T) {

	d := NewQuoteDistributor()

	s1 := &testClobSink{
		id:       "s1",
	}

	s2 := &testClobSink{
		id:       "s2",
	}

	d.AddConnection(s1)
	d.AddConnection(s2)

	d.readInputChannels()
	d.readInputChannels()

	d.Send(&model.ClobQuote{ListingId:1})

	d.readInputChannels()

	if s1.received[0].ListingId != 1 {
		t.Errorf("expected quote note received")
	}

	if s2.received[0].ListingId != 1 {
		t.Errorf("expected quote note received")
	}

}


