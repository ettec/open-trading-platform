package actor

import (
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/connections"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/fix/marketdata"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/model"
	"log"
	"testing"
	"time"
)

type testMarketDataClient struct {
	connect   func() (<-chan *model.ClobQuote, error)
	subscribe func(listingId int)
	close     func() error
}

func (t *testMarketDataClient) Connect() (<-chan *model.ClobQuote, error) {
	return t.connect()
}

func (t *testMarketDataClient) Subscribe(listingId int) {
	t.subscribe(listingId)
}

func (t *testMarketDataClient) Close() error {
	return t.close()
}

type recvTuple struct {
	r *marketdata.MarketDataIncrementalRefresh
	e error
}

type testIncRefreshSource struct {
	send chan recvTuple
}

func (t *testIncRefreshSource) Recv() (*marketdata.MarketDataIncrementalRefresh, error) {
	select {
	case r := <-t.send:
		return r.r, r.e
	}
}

func TestNewMdServerConnection(t *testing.T) {

	connectedCalled := false

	dial := func(target string) connections.Connection {
		return &testMarketDataClient{
			connect: func() (source <-chan *model.ClobQuote, err error) {
				connectedCalled = true
				return make(chan *model.ClobQuote), nil
			},
			subscribe: nil,
			close:     nil,
		}
	}

	mdConn := NewMdServerConnection("testaddress", "testconnection", dial, 0)

	mdConn.Connect()

	invoke(mdConn.readInputChannels, 1)

	if !connectedCalled {
		t.Errorf("expected connection method to be called")
	}

	if mdConn.connection == nil {
		t.Errorf("expected the md conn to have a connection")
	}

}

func invoke(f func() (chan<- bool, error), times int) {
	for i := 0; i < times; i++ {
		f()
	}
}

func TestSubscribe(t *testing.T) {

	subscribed := make([]int, 0)

	dial := func(target string) connections.Connection {
		return &testMarketDataClient{
			connect: func() (source <-chan *model.ClobQuote, err error) {

				return make(chan *model.ClobQuote), nil
			},
			subscribe: func(listingId int) {
				subscribed = append(subscribed, listingId)
			},
			close: nil,
		}
	}

	mdConn := NewMdServerConnection("testaddress", "testconnection", dial, 0)

	mdConn.Connect()

	invoke(mdConn.readInputChannels, 1)

	mdConn.Subscribe(1)
	mdConn.Subscribe(2)

	invoke(mdConn.readInputChannels, 2)

	if subscribed[0] != 1 {
		t.Error("expected to receive 1")
	}

	if subscribed[1] != 2 {
		t.Error("expected to receive 2")
	}

}

func TestRefreshesAreForwardedToSink(t *testing.T) {

	clobSource := make(chan *model.ClobQuote, 10)

	dial := func(target string) connections.Connection {
		return &testMarketDataClient{
			connect: func() (source <-chan *model.ClobQuote, err error) {
				return clobSource, nil
			},
			subscribe: nil,
			close:     nil,
		}
	}

	mdConn := NewMdServerConnection("testaddress", "testconnection", dial, 0)

	out, _ := mdConn.Connect()

	mdConn.readInputChannels()

	clobSource <- &model.ClobQuote{ListingId: 1}
	clobSource <- &model.ClobQuote{ListingId: 2}

	mdConn.readInputChannels()
	mdConn.readInputChannels()

	r1 := <-out
	if r1.ListingId != 1 {
		t.Errorf("unexpected refresh")
	}

	r2 := <-out
	if r2.ListingId != 2 {
		t.Errorf("unexpected refresh")
	}

}

func TestSubscribesSentWhilstNotConnectedAreResentOnConnect(t *testing.T) {

	subscribed := make([]int, 0)
	clobSource := make(chan *model.ClobQuote, 10)
	dial := func(target string) connections.Connection {
		return &testMarketDataClient{
			connect: func() (source <-chan *model.ClobQuote, err error) {
				return clobSource, nil
			},
			subscribe: func(listingId int) {
				subscribed = append(subscribed, listingId)
			},
			close: nil,
		}
	}

	mdConn := NewMdServerConnection("testaddress", "testconnection", dial, 0)

	mdConn.Subscribe(1)
	mdConn.Subscribe(2)

	invoke(mdConn.readInputChannels, 2)

	mdConn.Connect()

	mdConn.readInputChannels()
	mdConn.readInputChannels()
	mdConn.readInputChannels()

	if subscribed[0] != 1 {
		t.Error("expected to receive 1")
	}

	if subscribed[1] != 2 {
		t.Error("expected to receive 2")
	}
}

func TestReconnectOccursAfterConnectionFailure(t *testing.T) {

	clobSource := make(chan *model.ClobQuote, 10)

	subscriptions := make(chan int, 10)

	connectCalled := false

	dial := func(target string) connections.Connection {
		return &testMarketDataClient{
			connect: func() (source <-chan *model.ClobQuote, err error) {
				log.Println("connect called")
				connectCalled = true
				return clobSource, nil
			},
			subscribe: func(listingId int) {
				subscriptions <- listingId
			},
			close: func() error {
				return nil
			},
		}
	}

	reconnectInterval := 1 * time.Second

	mdConn := NewMdServerConnection("testaddress", "testconnection", dial, reconnectInterval)

	out, _ := mdConn.Connect()

	mdConn.readInputChannels()

	if !connectCalled {
		t.Errorf("expected connect to be called")
	}
	connectCalled = false

	mdConn.Subscribe(1)
	mdConn.Subscribe(2)

	mdConn.readInputChannels()
	mdConn.readInputChannels()

	<-subscriptions
	<-subscriptions

	clobSource <- &model.ClobQuote{ListingId: 1}

	mdConn.readInputChannels()

	r1 := <-out
	if r1.ListingId != 1 {
		t.Errorf("unexpected refresh")
	}

	close(clobSource)

	mdConn.readInputChannels()

	clobSource = make(chan *model.ClobQuote, 10)

	mdConn.readInputChannels()

	time.Sleep(reconnectInterval / 2)


	time.Sleep(reconnectInterval)
	mdConn.readInputChannels()
	mdConn.readInputChannels()


	if !connectCalled {
		t.Errorf("expected connect to be called")
	}
	connectCalled = false

	if s := <-subscriptions; s != 1 {
		t.Errorf("expected subscription to be resent for 1")
	}

	if s := <-subscriptions; s != 2 {
		t.Errorf("expected subscription to be resent for 2")
	}

}

func TestConnectionIsClosedWhenMarketDataConnectionActorIsClosed(t *testing.T) {

	clobSource := make(chan *model.ClobQuote, 10)
	closeChan := make(chan bool, 10)

	dial := func(target string) connections.Connection {
		return &testMarketDataClient{
			connect: func() (source <-chan *model.ClobQuote, err error) {

				return clobSource, nil
			},
			subscribe: func(listingId int) {

			},
			close: func() error {
				closeChan <- true
				return nil
			},
		}
	}

	mdConn := NewMdServerConnection("testaddress", "testconnection", dial, 0)
	mdConn.Start()
	mdConn.Connect()

	d:= make(chan bool, 10)

	go mdConn.Close(d)
	<-d

	<-closeChan

}
