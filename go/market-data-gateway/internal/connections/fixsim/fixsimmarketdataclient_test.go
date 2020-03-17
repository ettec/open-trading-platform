package fixsim

import (
	"context"
	"fmt"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/fix/marketdata"
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

	streamOutChan chan FixSimMarketDataService_ConnectClient
}



func (t testClient) Connect(ctx context.Context, opts ...grpc.CallOption) (FixSimMarketDataService_ConnectClient, error) {
	return <-t.streamOutChan, nil
}

type testClientStream struct {
	refreshChan    chan *marketdata.MarketDataIncrementalRefresh
	refreshErrChan chan error
	subsInChan    chan *marketdata.MarketDataRequest
}

func (t testClientStream) Recv() (*marketdata.MarketDataIncrementalRefresh, error) {
	select {
	case r := <-t.refreshChan:
		return r, nil
	case e := <-t.refreshErrChan:
		return nil, e
	}

	return <-t.refreshChan, nil
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

func Test_fixSimMarketDataClient_refreshesForwaredToOut(t *testing.T) {

	client, stream, conn, _, out := setup(t)

	conn.getStateChan <- connectivity.Ready

	client.streamOutChan <- stream

	stream.refreshChan <- &marketdata.MarketDataIncrementalRefresh{}

	<-out
}

func Test_fixSimMarketDataClient_streamErrorSendsNilToOutStream(t *testing.T) {

	client, stream, conn, _, out := setup(t)

	conn.getStateChan <- connectivity.Ready

	client.streamOutChan <- stream

	stream.refreshChan <- &marketdata.MarketDataIncrementalRefresh{}
	<-out

	stream.refreshErrChan <- fmt.Errorf("testerror")
	r := <-out
	if r != nil {
		t.FailNow()
	}
}


func Test_fixSimMarketDataClient_testReconnectAfterError(t *testing.T) {

	client, stream, conn, toTest, out := setup(t)

	conn.getStateChan <- connectivity.Ready

	client.streamOutChan <- stream

	toTest.subscribe("A")



	stream.refreshChan <- &marketdata.MarketDataIncrementalRefresh{}
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
	<-stream.subsInChan
	<-stream.subsInChan



}


func Test_fixSimMarketDataClient_resubscribedOnConnect(t *testing.T) {

	client, stream, conn, toTest, _ := setup(t)

	toTest.subscribe("A")

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

func setup(t *testing.T) (testClient, testClientStream, testConnection, *fixSimMarketDataClient, chan *marketdata.MarketDataIncrementalRefresh) {
	out := make(chan *marketdata.MarketDataIncrementalRefresh)

	client := testClient{

		streamOutChan: make(chan FixSimMarketDataService_ConnectClient),

	}

	stream := testClientStream{refreshChan: make(chan *marketdata.MarketDataIncrementalRefresh),
		refreshErrChan: make(chan error),
		subsInChan: make(chan *marketdata.MarketDataRequest),
	}

	conn := testConnection{
		getStateChan: make(chan connectivity.State),
	}

	c, err := NewFixSimMarketDataClient("testId", "testTarget", out,
		func(targetAddress string) (FixSimMarketDataServiceClient, GrpcConnection, error) {
			return client, conn, nil
		})

	if err != nil {
		t.FailNow()
	}
	return client, stream, conn, c, out
}
