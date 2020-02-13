package fixsim

import (
	"context"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/fix/common"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/fix/marketdata"
	"google.golang.org/grpc"
	"log"
	"os"
)




type  receiveIncRefreshFn = func() (*marketdata.MarketDataIncrementalRefresh, error)

type fixSimMarketDataClient struct {
	id string
	client FixSimMarketDataServiceClient
	conn   *grpc.ClientConn
	out            chan<- *marketdata.MarketDataIncrementalRefresh
	errLog         *log.Logger
}

func NewFixSimMarketDataClient(id string, targetAddress string, out chan<- *marketdata.MarketDataIncrementalRefresh) (*fixSimMarketDataClient, error) {

	n := &fixSimMarketDataClient{
		id:				id,
		out:            out,
		errLog:         log.New(os.Stderr, targetAddress, log.Lshortfile | log.Ltime),
	}

	conn, err := grpc.Dial(targetAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	n.conn = conn
	n.client = NewFixSimMarketDataServiceClient(conn)

	stream, err := n.connect()

	go func() {
		defer func() {
			close(n.out)
			if err := n.close(); err != nil {
				n.errLog.Println("error whilst closing:", err)
			}
		}()

		for {
			incRefresh, err := stream()
			if err != nil {
				n.errLog.Println("inbound stream error:", err)
				return
			}

			n.out <- incRefresh
		}
	}()

	return n, nil
}



func (fsc *fixSimMarketDataClient) close() error {
	return fsc.conn.Close()
}

func (fsc *fixSimMarketDataClient) subscribe(symbol string) error {
	request := &marketdata.MarketDataRequest{Parties: []*common.Parties{{PartyId: fsc.id}},
		InstrmtMdReqGrp: []*common.InstrmtMDReqGrp{{Instrument: &common.Instrument{Symbol: symbol}}}}
	_, err := fsc.client.Subscribe(context.Background(), request)
	return err
}

func (fsc *fixSimMarketDataClient) connect() (receiveIncRefreshFn, error) {
	r := &ConnectRequest{PartyId: fsc.id}
	stream, err := fsc.client.Connect(context.Background(), r)
	return stream.Recv, err
}

