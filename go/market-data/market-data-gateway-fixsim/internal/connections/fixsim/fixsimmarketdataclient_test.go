package fixsim

import (
	"context"
	"fmt"
	"github.com/ettec/open-trading-platform/go/market-data/market-data-gateway-fixsim/internal/fix/marketdata"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/metadata"
	"testing"
)

func TestMdIncRefreshesAreForwarded(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, stream, conn, toTest := setup(t, ctx)

	conn.getStateChan <- connectivity.Ready

	client.streamOutChan <- stream

	stream.refreshChan <- &marketdata.MarketDataIncrementalRefresh{}

	<-toTest.Chan()
}

func TestNilSentToOutChanWhenErrorOccurs(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, stream, conn, toTest := setup(t, ctx)

	conn.getStateChan <- connectivity.Ready

	client.streamOutChan <- stream

	stream.refreshChan <- &marketdata.MarketDataIncrementalRefresh{}
	<-toTest.Chan()

	stream.refreshErrChan <- fmt.Errorf("testerror")
	r := <-toTest.Chan()
	if r != nil {
		t.FailNow()
	}
}

func TestReconnectsAfterError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, stream, conn, toTest := setup(t, ctx)

	conn.getStateChan <- connectivity.Ready

	client.streamOutChan <- stream

	err := toTest.Subscribe("A")
	assert.NoError(t, err)

	<-stream.subsInChan

	stream.refreshChan <- &marketdata.MarketDataIncrementalRefresh{}
	<-toTest.Chan()

	stream.refreshErrChan <- fmt.Errorf("testerror")
	r := <-toTest.Chan()
	if r != nil {
		t.FailNow()
	}

	conn.getStateChan <- connectivity.TransientFailure
	conn.getStateChan <- connectivity.Ready
	client.streamOutChan <- stream

	// resubscribe
	<-stream.subsInChan

}

func TestResubscribesToListingsOnConnect(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, stream, conn, toTest := setup(t, ctx)

	err := toTest.Subscribe("A")
	assert.NoError(t, err)

	conn.getStateChan <- connectivity.Ready

	client.streamOutChan <- stream

	s := <-stream.subsInChan
	if s.Parties[0].PartyId != "testId" {
		t.FailNow()
	}

	if s.InstrmtMdReqGrp[0].Instrument.Symbol != "A" {
		t.FailNow()
	}

}

func setup(t *testing.T, ctx context.Context) (testClient, testClientStream, testConnection, *fixSimMarketDataClient) {

	client := testClient{
		streamOutChan: make(chan FixSimMarketDataService_ConnectClient),
	}

	stream := testClientStream{refreshChan: make(chan *marketdata.MarketDataIncrementalRefresh),
		refreshErrChan: make(chan error),
		subsInChan:     make(chan *marketdata.MarketDataRequest),
	}

	conn := testConnection{
		getStateChan: make(chan connectivity.State),
	}

	c, err := NewFixSimMarketDataClient(ctx, "testId", client, conn, 100)

	assert.NoError(t, err)

	return client, stream, conn, c
}

type testConnection struct {
	state        connectivity.State
	getStateChan chan connectivity.State
}

func (t testConnection) GetState() connectivity.State {
	t.state = <-t.getStateChan
	return t.state
}

func (t testConnection) WaitForStateChange(ctx context.Context, sourceState connectivity.State) bool {

	for {
		if t.state != sourceState {
			return true
		}
		t.state = <-t.getStateChan
	}

}

type testClient struct {
	streamOutChan chan FixSimMarketDataService_ConnectClient
}

func (t testClient) Connect(ctx context.Context, opts ...grpc.CallOption) (FixSimMarketDataService_ConnectClient, error) {
	return <-t.streamOutChan, nil
}

type testClientStream struct {
	refreshChan    chan *marketdata.MarketDataIncrementalRefresh
	refreshErrChan chan error
	subsInChan     chan *marketdata.MarketDataRequest
}

func (t testClientStream) Recv() (*marketdata.MarketDataIncrementalRefresh, error) {
	select {
	case r := <-t.refreshChan:
		return r, nil
	case e := <-t.refreshErrChan:
		return nil, e
	}
}

func (t testClientStream) Send(m *marketdata.MarketDataRequest) error {
	t.subsInChan <- m
	return nil
}

func (t testClientStream) Header() (metadata.MD, error) {
	panic("implement me")
}

func (t testClientStream) Trailer() metadata.MD {
	panic("implement me")
}

func (t testClientStream) CloseSend() error {
	panic("implement me")
}

func (t testClientStream) Context() context.Context {
	panic("implement me")
}

func (t testClientStream) SendMsg(m interface{}) error {
	panic("implement me")
}

func (t testClientStream) RecvMsg(m interface{}) error {
	panic("implement me")
}
