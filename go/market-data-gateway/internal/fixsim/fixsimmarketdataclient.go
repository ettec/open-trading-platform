package fixsim

import (
	"context"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/fix/common"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/fix/marketdata"
	"google.golang.org/grpc"
)


type  receiveIncRefreshFn = func() (*marketdata.MarketDataIncrementalRefresh, error)

type fixSimMarketDataClientImpl struct {
	client FixSimMarketDataServiceClient
	conn   *grpc.ClientConn
}

func (fsc *fixSimMarketDataClientImpl) close() error {
	return fsc.conn.Close()
}

func (fsc *fixSimMarketDataClientImpl) subscribe(symbol string, subscriberId string) error {
	request := &marketdata.MarketDataRequest{Parties: []*common.Parties{{PartyId: subscriberId}},
		InstrmtMdReqGrp: []*common.InstrmtMDReqGrp{{Instrument: &common.Instrument{Symbol: symbol}}}}
	_, err := fsc.client.Subscribe(context.Background(), request)
	return err
}

func (fsc *fixSimMarketDataClientImpl) connect(connectionId string) (receiveIncRefreshFn, error) {
	r := &ConnectRequest{PartyId: connectionId}
	stream, err := fsc.client.Connect(context.Background(), r)
	return stream.Recv, err
}

