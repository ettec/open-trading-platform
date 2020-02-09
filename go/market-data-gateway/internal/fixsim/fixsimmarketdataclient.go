package fixsim

import (
	"context"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/fix/common"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/fix/marketdata"
	"google.golang.org/grpc"
)

type IncRefreshSource interface {
	Recv() (*marketdata.MarketDataIncrementalRefresh, error)
}

type fixSimMarketDataClientImpl struct {
	client FixSimMarketDataServiceClient
	conn   *grpc.ClientConn
}

func (fsc *fixSimMarketDataClientImpl) Close() error {
	return fsc.conn.Close()
}

func (fsc *fixSimMarketDataClientImpl) Subscribe(symbol string, subscriberId string) error {
	request := &marketdata.MarketDataRequest{Parties: []*common.Parties{{PartyId: subscriberId}},
		InstrmtMdReqGrp: []*common.InstrmtMDReqGrp{{Instrument: &common.Instrument{Symbol: symbol}}}}
	_, err := fsc.client.Subscribe(context.Background(), request)
	return err
}

func (fsc *fixSimMarketDataClientImpl) Connect(connectionId string) (IncRefreshSource, error) {
	r := &ConnectRequest{PartyId: connectionId}
	stream, err := fsc.client.Connect(context.Background(), r)
	return stream, err
}

