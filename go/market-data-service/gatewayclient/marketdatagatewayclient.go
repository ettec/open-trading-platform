package gatewayclient

import (
	"context"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/api"
	"github.com/ettec/open-trading-platform/go/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"log"
	"os"
	"sync"
)

type marketGatewayClient struct {
	id            string
	client        api.MarketDataGatewayClient
	conn          *grpc.ClientConn
	out           chan<- *model.ClobQuote
	subscribeMux  sync.Mutex
	subscriptions map[int32]bool
	log           *log.Logger
	errLog        *log.Logger
}

type getMarketDataGatewayClientFn = func(targetAddress string) (api.MarketDataGatewayClient, GrpcConnection, error)

type GrpcConnection interface {
	GetState() connectivity.State
	WaitForStateChange(ctx context.Context, sourceState connectivity.State) bool
}

func NewMarketDataGatewayClient(id string, targetAddress string, out chan<- *model.ClobQuote,
	getConnection getMarketDataGatewayClientFn) (*marketGatewayClient, error) {

	n := &marketGatewayClient{
		id:            id,
		out:           out,
		subscriptions: map[int32]bool{},
		log:           log.New(os.Stdout, "target:"+targetAddress+" ", log.Lshortfile|log.Ltime),
		errLog:        log.New(os.Stderr, "target:"+targetAddress+" ", log.Lshortfile|log.Ltime),
	}

	log.Println("connecting to fix sim market data service at:" + targetAddress)

	client, conn, err := getConnection(targetAddress)
	if err != nil {
		return nil, err
	}

	n.client = client

	go func() {

		for {
			state := conn.GetState()
			for state != connectivity.Ready {
				n.log.Printf("waiting for market gateway connection to be ready....")


				conn.WaitForStateChange(context.Background(), state)
				state = conn.GetState()
				n.log.Println("market gateway connection state is:", state)
			}

			stream, err := connect(n.client, id)
			if err != nil {
				n.errLog.Println("failed to connect to quote stream:", err)
				continue
			}

			n.log.Println("connected to quote stream")

			err = n.resubscribeAllListings()
			if err != nil {
				n.errLog.Println("failed to resubscribe to all listings:", err)
				continue
			}

			if len(n.subscriptions) > 0 {
				n.log.Printf("resubscribed to all %v quotes", len(n.subscriptions))
			}

			for {
				incRefresh, err := stream.Recv()
				if err != nil {
					n.errLog.Println("inbound stream error:", err)
					n.subscribeMux.Lock()
					for listingId, subscribed := range n.subscriptions {
						if subscribed {
							n.out <- &model.ClobQuote{
								ListingId:         listingId,
								Bids:              []*model.ClobLine{},
								Offers:            []*model.ClobLine{},
								StreamInterrupted: true,
							}
						}
					}

					n.subscribeMux.Unlock()

					break
				}
				n.out <- incRefresh
			}
		}
	}()

	return n, nil
}

func (mgc *marketGatewayClient) close() error {
	return mgc.conn.Close()
}

func (mgc *marketGatewayClient) resubscribeAllListings() error {
	mgc.subscribeMux.Lock()
	defer mgc.subscribeMux.Unlock()
	for symbol := range mgc.subscriptions {
		err := subscribe(mgc.client, symbol, mgc.id)
		if err != nil {
			return err
		}
	}
	return nil
}

func connect(client api.MarketDataGatewayClient, id string) (api.MarketDataGateway_ConnectClient, error) {


	r := &api.ConnectRequest{SubscriberId: id


	stream, err := client.Connect(context.Background(), r)
	return stream, err
}

func subscribe(client api.MarketDataGatewayClient, listingId int32, id string) error {
	request := &api.SubscribeRequest{
		SubscriberId: id,
		ListingId:    listingId,
	}
	_, err := client.Subscribe(context.Background(), request)

	return err
}

func (mgc *marketGatewayClient) Subscribe(listingId int32) {

	mgc.subscribeMux.Lock()
	defer mgc.subscribeMux.Unlock()
	if !mgc.subscriptions[listingId] {
		mgc.subscriptions[listingId] = true
		err := subscribe(mgc.client, listingId, mgc.id)
		if err != nil {
			mgc.errLog.Printf("failed to subsribe to listing %v, errorr:%v", listingId, err)
		}
	}

}
