package actor

import (
	"fmt"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/fix/marketdata"
	"testing"
	"time"
)

type testRefreshSink struct {
	refreshes chan *marketdata.MarketDataIncrementalRefresh
}


func (t *testRefreshSink) SendRefresh(refresh *marketdata.MarketDataIncrementalRefresh) {
	t.refreshes <- refresh
}

type testMarketDataClient struct {
	connect func(connectionId string) (IncRefreshSource, error)
	subscribe func(symbol string, subscriberId string) error
	close func() error
}

func (t *testMarketDataClient) Connect(connectionId string) (IncRefreshSource, error){
	return t.connect(connectionId)
}

func (t *testMarketDataClient) Subscribe(symbol string, subscriberId string) error {
		return t.subscribe(symbol,subscriberId)
}

func (t *testMarketDataClient)Close() error {
	return t.close()
}

type recvTuple struct{
	r	*marketdata.MarketDataIncrementalRefresh
	e error
}

type testIncRefreshSource struct {
	send chan recvTuple
}

func (t *testIncRefreshSource) Recv() (*marketdata.MarketDataIncrementalRefresh, error) {
	select {
		case r := <- t.send:
			return r.r, r.e
	}
}



func TestNewMdServerConnection(t *testing.T) {

	connectedCalled := false
	dial := func(target string) (MarketDataClient, error) {
		return &testMarketDataClient{
			connect: func(connectionId string) (source IncRefreshSource, err error) {
				connectedCalled = true
				return &testIncRefreshSource{}, nil
			},
			subscribe: nil,
			close:     nil,
		}, nil
	}

	mdConn := NewMdServerConnection("testaddress", "testconnection", &testRefreshSink{}, dial,0 )

	invoke(mdConn.readInputChannels,2)

	if !connectedCalled {
		t.Errorf("expected connection method to be called")
	}

	if mdConn.connection == nil {
		t.Errorf("expected the md conn to have a connection")
	}

}

func TestSubscribe(t *testing.T) {

	subscribed := make([]string,0)
	dial := func(target string) (MarketDataClient, error) {
		return &testMarketDataClient{
			connect: func(connectionId string) (source IncRefreshSource, err error) {
				return &testIncRefreshSource{}, nil
			},
			subscribe: func(symbol string, subscriberId string) error {
				subscribed = append(subscribed, symbol)
				return nil
			},
			close:     nil,
		}, nil
	}

	mdConn := NewMdServerConnection("testaddress", "testconnection", &testRefreshSink{}, dial, 0 )

	invoke(mdConn.readInputChannels,2)

	mdConn.Subscribe("A")
	mdConn.Subscribe("B")

	invoke(mdConn.readInputChannels,2)


	if subscribed[0] != "A" {
		t.Error("expected to receive A" )
	}

	if subscribed[1] != "B" {
		t.Error("expected to receive B" )
	}

}

func TestRefreshesAreForwardedToSink(t *testing.T) {

	refreshSource := &testIncRefreshSource{send: make( chan recvTuple, 10)}

	dial := func(target string) (MarketDataClient, error) {
		return &testMarketDataClient{
			connect: func(connectionId string) (source IncRefreshSource, err error) {
				return refreshSource, nil
			},
			subscribe: nil,
			close:     nil,
		}, nil
	}

	refreshSink := &testRefreshSink{refreshes:make(chan *marketdata.MarketDataIncrementalRefresh, 10)}

	mdConn := NewMdServerConnection("testaddress", "testconnection", refreshSink, dial, 0 )

	mdConn.readInputChannels()
	mdConn.readInputChannels()

	refreshSource.send <- recvTuple{r:&marketdata.MarketDataIncrementalRefresh{MdReqId:"1"}}
	refreshSource.send <- recvTuple{r:&marketdata.MarketDataIncrementalRefresh{MdReqId:"2"}}


	r1 := <- refreshSink.refreshes
	if r1.MdReqId != "1" {
		t.Errorf("unexpected refresh")
	}

	r2 := <- refreshSink.refreshes
	if r2.MdReqId != "2" {
		t.Errorf("unexpected refresh")
	}

}


func TestSubscribesSentWhilstNotConnectedAreResentOnConnect(t *testing.T) {
	subscribed := make([]string,0)
	dial := func(target string) (MarketDataClient, error) {
		return &testMarketDataClient{
			connect: func(connectionId string) (source IncRefreshSource, err error) {
				return &testIncRefreshSource{}, nil
			},
			subscribe: func(symbol string, subscriberId string) error {
				subscribed = append(subscribed, symbol)
				return nil
			},
			close:     nil,
		}, nil
	}

	mdConn := NewMdServerConnection("testaddress", "testconnection", &testRefreshSink{}, dial, 0 )


	mdConn.Subscribe("A")
	mdConn.Subscribe("B")

	invoke(mdConn.readInputChannels,6)


	if subscribed[0] != "A" {
		t.Error("expected to receive A" )
	}

	if subscribed[1] != "B" {
		t.Error("expected to receive B" )
	}
}


func TestReconnectOccursAfterConnectionFailure(t *testing.T) {

	refreshSource := &testIncRefreshSource{send: make( chan recvTuple, 10)}

	closeChan := make(chan bool)

	subscriptions := make( chan string,10)

	dial := func(target string) (MarketDataClient, error) {
		return &testMarketDataClient{
			connect: func(connectionId string) (source IncRefreshSource, err error) {
				return refreshSource, nil
			},
			subscribe: func(symbol string, subscriberId string) error {
				subscriptions <-symbol
				return nil
			},
			close: func() error {
				closeChan<-true
				return nil
			},
		}, nil
	}

	refreshSink := &testRefreshSink{refreshes:make(chan *marketdata.MarketDataIncrementalRefresh, 10)}

	reconnectInterval := 1 * time.Second

	mdConn := NewMdServerConnection("testaddress", "testconnection", refreshSink, dial, reconnectInterval )

	mdConn.readInputChannels()
	mdConn.readInputChannels()

	mdConn.Subscribe("A")
	mdConn.Subscribe("B")

	mdConn.readInputChannels()
	mdConn.readInputChannels()

	<-subscriptions
	<-subscriptions


	refreshSource.send <- recvTuple{r: &marketdata.MarketDataIncrementalRefresh{MdReqId:"1"}}

	r1 := <- refreshSink.refreshes
	if r1.MdReqId != "1" {
		t.Errorf("unexpected refresh")
	}

	refreshSource.send <- recvTuple{e: fmt.Errorf("test error") }
	<-closeChan

	mdConn.readInputChannels()
	if mdConn.connection != nil {
		t.Errorf("expected connection to be nil after error")
	}

	time.Sleep(reconnectInterval/2)

	if mdConn.connection != nil {
		t.Errorf("expected connection to be still be nil")
	}

	time.Sleep(reconnectInterval)
	mdConn.readInputChannels()
	mdConn.readInputChannels()
	mdConn.readInputChannels()
	mdConn.readInputChannels()


	if mdConn.connection == nil {
		t.Errorf("expected connection to be live")
	}

	if s := <-subscriptions; s != "A" {
		t.Errorf("expected subscription to be resent for A")
	}

	if s := <-subscriptions; s != "B" {
		t.Errorf("expected subscription to be resent for B")
	}

}

func TestConnectionIsClosedWhenMarketDataConnectionActorIsClosed(t *testing.T) {

	connCloseChan := make(chan bool, 10)

	dial := func(target string) (MarketDataClient, error) {
		return &testMarketDataClient{
			connect: func(connectionId string) (source IncRefreshSource, err error) {
				return &testIncRefreshSource{}, nil
			},
			subscribe: nil,
			close: func() error {
				connCloseChan <-true
				return nil
			},
		}, nil
	}

	refreshSink := &testRefreshSink{refreshes:make(chan *marketdata.MarketDataIncrementalRefresh, 10)}


	mdConn := NewMdServerConnection("testaddress", "testconnection", refreshSink, dial, 0 )
	mdConn.closeChan = make(chan chan<-bool,1)
	mdConn.readInputChannels()
	mdConn.readInputChannels()

	done := make(chan bool)
	mdConn.Close(done)
	mdConn.readInputChannels()
	<-connCloseChan

}
