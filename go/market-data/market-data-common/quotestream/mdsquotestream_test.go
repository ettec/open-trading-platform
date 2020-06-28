package quotestream

import (
	"context"
	"fmt"
	"github.com/ettec/otp-common/api/marketdatasource"
	"github.com/ettec/otp-model"

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
	streamOutChan chan marketdatasource.MarketDataSource_ConnectClient
}

func (t testClient) Connect(ctx context.Context, opts ...grpc.CallOption) (marketdatasource.MarketDataSource_ConnectClient, error) {
	return <-t.streamOutChan, nil
}

type testClientStream struct {
	refreshChan    chan *model.ClobQuote
	refreshErrChan chan error
	subscribeChan  chan *marketdatasource.SubscribeRequest
}

func (t testClientStream) Send(request *marketdatasource.SubscribeRequest) error {
	t.subscribeChan <- request
	return nil
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

func Test_marketDataGatewayClient_sendsEmptyQuoteForAllListingsOnConnectionError(t *testing.T) {

	client, stream, conn, mdc, out := setup(t)

	conn.getStateChan <- connectivity.Ready

	mdc.Subscribe(1)
	mdc.Subscribe(2)

	client.streamOutChan <- stream

	stream.refreshChan <- &model.ClobQuote{}
	<-out

	stream.refreshErrChan <- fmt.Errorf("testerror")
	r := <-out

	if !r.StreamInterrupted {
		t.FailNow()
	}

	if r.ListingId != 1 && r.ListingId != 2 {
		t.FailNow()
	}

	if len(r.Bids) != 0 || len(r.Offers) != 0 {
		t.FailNow()
	}

	r = <-out
	if r.ListingId != 1 && r.ListingId != 2 {
		t.FailNow()
	}

	if len(r.Bids) != 0 || len(r.Offers) != 0 {
		t.FailNow()
	}

	if !r.StreamInterrupted {
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
	if r.StreamInterrupted != true {
		t.FailNow()
	}

	conn.getStateChan <- connectivity.TransientFailure
	conn.getStateChan <- connectivity.Ready
	client.streamOutChan <- stream

	<-stream.subscribeChan

}

func Test_marketDataGatewayClient_resubscribedOnConnect(t *testing.T) {

	client, stream, conn, toTest, _ := setup(t)

	toTest.Subscribe(1)

	conn.getStateChan <- connectivity.Ready

	client.streamOutChan <- stream

	s := <-stream.subscribeChan

	if s.ListingId != 1 {
		t.FailNow()
	}

}

func setup(t *testing.T) (testClient, testClientStream, testConnection, *mdsQuoteStream, chan *model.ClobQuote) {
	out := make(chan *model.ClobQuote)

	client := testClient{

		streamOutChan: make(chan marketdatasource.MarketDataSource_ConnectClient),
	}

	stream := testClientStream{refreshChan: make(chan *model.ClobQuote),
		refreshErrChan: make(chan error),
		subscribeChan:  make(chan *marketdatasource.SubscribeRequest, 10)}

	conn := testConnection{
		getStateChan: make(chan connectivity.State),
	}

	c, err := NewMdsQuoteStreamFromFn("testId", "testTarget", out,
		func(targetAddress string) (marketdatasource.MarketDataSourceClient, GrpcConnection, error) {
			return client, conn, nil
		})

	if err != nil {
		t.FailNow()
	}
	return client, stream, conn, c, out
}
