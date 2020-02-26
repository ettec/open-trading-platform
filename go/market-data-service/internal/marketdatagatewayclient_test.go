package internal

import (
	"context"
	"fmt"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/api"
	"github.com/ettec/open-trading-platform/go/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/metadata"
	"testing"
)

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
	subsInChan    chan *api.SubscribeRequest
	streamOutChan chan api.MarketDataGateway_ConnectClient
}

func (t testClient) Subscribe(ctx context.Context, in *api.SubscribeRequest, opts ...grpc.CallOption) (*model.Empty, error) {
	t.subsInChan <- in
	return &model.Empty{}, nil
}

func (t testClient) Connect(ctx context.Context, in *api.ConnectRequest, opts ...grpc.CallOption) (api.MarketDataGateway_ConnectClient, error) {
	return <-t.streamOutChan, nil
}

type testClientStream struct {
	refreshChan    chan *model.ClobQuote
	refreshErrChan chan error
}

func (t testClientStream) Recv() (*model.ClobQuote, error) {
	select {
	case r := <-t.refreshChan:
		return r, nil
	case e := <-t.refreshErrChan:
		return nil, e
	}

	return <-t.refreshChan, nil
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

func Test_marketDataGatewayClient_refreshesForwaredToOut(t *testing.T) {

	client, stream, conn, _, out := setup(t)

	conn.getStateChan <- connectivity.Ready

	client.streamOutChan <- stream

	stream.refreshChan <- &model.ClobQuote{}

	<-out
}

func Test_marketDataGatewayClient_streamErrorSendsNilToOutStream(t *testing.T) {

	client, stream, conn, _, out := setup(t)

	conn.getStateChan <- connectivity.Ready

	client.streamOutChan <- stream

	stream.refreshChan <- &model.ClobQuote{}
	<-out

	stream.refreshErrChan <- fmt.Errorf("testerror")
	r := <-out
	if r != nil {
		t.FailNow()
	}
}


func Test_marketDataGatewayClient_testReconnectAfterError(t *testing.T) {

	client, stream, conn, toTest, out := setup(t)

	conn.getStateChan <- connectivity.Ready

	client.streamOutChan <- stream

	toTest.Subscribe(1)



	stream.refreshChan <- &model.ClobQuote{}
	<-out

	stream.refreshErrChan <- fmt.Errorf("testerror")
	r := <-out
	if r != nil {
		t.FailNow()
	}

	conn.getStateChan <- connectivity.TransientFailure
	conn.getStateChan <- connectivity.Ready
	client.streamOutChan <- stream

	// Original subscribe and resubscribe
	<-client.subsInChan
	<-client.subsInChan



}


func Test_marketDataGatewayClient_resubscribedOnConnect(t *testing.T) {

	client, stream, conn, toTest, _ := setup(t)

	toTest.Subscribe(1)
	<-client.subsInChan

	conn.getStateChan <- connectivity.Ready

	client.streamOutChan <- stream

	s := <-client.subsInChan
	if s.SubscriberId != "testId" {
		t.FailNow()
	}

	if s.ListingId != 1 {
		t.FailNow()
	}

}

func setup(t *testing.T) (testClient, testClientStream, testConnection, *marketGatewayClient, chan *model.ClobQuote) {
	out := make(chan *model.ClobQuote)

	client := testClient{
		subsInChan:    make(chan *api.SubscribeRequest, 10),
		streamOutChan: make(chan api.MarketDataGateway_ConnectClient),
	}

	stream := testClientStream{refreshChan: make(chan *model.ClobQuote),
		refreshErrChan: make(chan error)}

	conn := testConnection{
		getStateChan: make(chan connectivity.State),
	}

	c, err := NewMarketDataGatewayClient("testId", "testTarget", out,
		func(targetAddress string) (api.MarketDataGatewayClient, GrpcConnection, error) {
			return client, conn, nil
		})

	if err != nil {
		t.FailNow()
	}
	return client, stream, conn, c, out
}
