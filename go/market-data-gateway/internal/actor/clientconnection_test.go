package actor

import (
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/model"
	"testing"
)

func Test_clientConnection_Send(t *testing.T) {

}


type testClobQuoteSink struct {
	send func(quote *model.ClobQuote) error
	close func() error
}

func (s *testClobQuoteSink) Send (quote *model.ClobQuote) error {
	return s.send(quote)
}

func (s *testClobQuoteSink) Close() error {
	return s.close()
}


func Test_clientConnection_Subscribe(t *testing.T) {



	subscribes := make(chan string, 100)

	quoteSink := &testClobQuoteSink{}


	c := NewClientConnection("testId",  func(listingId int) {

		},
		quoteSink,100)


	c.Subscribe(1)
	c.Subscribe(2)

	if s:= <-subscribes; s != "A" {
		t.Errorf("expected subscribe to symbol A")
	}

	if s:= <-subscribes; s != "B" {
		t.Errorf("expected subscribe to symbol B")
	}

}

func Test_connectionDroppedWhenBufferExceeded(t *testing.T) {

}