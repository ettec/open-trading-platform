package internal

import (
	"context"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/api"
	"github.com/ettec/open-trading-platform/go/model"
	"google.golang.org/grpc"
	"log"
	"os"
)

type  receiveQuote = func() (*model.ClobQuote, error)

type marketDataClient struct {
	id string
	client api.MarketDataGatewayClient
	conn   *grpc.ClientConn
	out            chan<- *model.ClobQuote
	errLog         *log.Logger
}

func NewMarketDataClient(id string, targetAddress string, out chan<- *model.ClobQuote) (*marketDataClient, error) {

	n := &marketDataClient{
		id:				id,
		out:            out,
		errLog:         log.New(os.Stderr, targetAddress, log.Lshortfile | log.Ltime),
	}

	conn, err := grpc.Dial(targetAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	n.conn = conn
	n.client = api.NewMarketDataGatewayClient(conn)

	stream, err := n.connect()

	go func() {
		defer func() {
			close(n.out)
			if err := n.Close(); err != nil {
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



func (mdc *marketDataClient) Close() error {
	return mdc.conn.Close()
}

func (mdc *marketDataClient) Subscribe(listingId int32) error  {

	request := &api.SubscribeRequest{
		SubscriberId:         mdc.id,
		ListingId:            int32(listingId),
	}

	_, err := mdc.client.Subscribe(context.Background(), request)

	return err
}

func (mdc *marketDataClient) connect() (receiveQuote, error) {
	r := &api.ConnectRequest{
		SubscriberId:         mdc.id,
	}
	stream, err := mdc.client.Connect(context.Background(), r)
	return stream.Recv, err
}
