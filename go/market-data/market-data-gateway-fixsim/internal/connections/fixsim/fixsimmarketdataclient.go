package fixsim

import (
	"context"
	"github.com/ettec/open-trading-platform/go/market-data/market-data-gateway-fixsim/internal/fix/common"
	"github.com/ettec/open-trading-platform/go/market-data/market-data-gateway-fixsim/internal/fix/marketdata"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/metadata"
	"log"
	"os"
)

type fixSimMarketDataClient struct {
	conn              *grpc.ClientConn
	subscriptionsChan chan string
	log               *log.Logger
	errLog            *log.Logger
}

type getMarketSimConnectionFn = func(targetAddress string) (FixSimMarketDataServiceClient, GrpcConnection, error)

type GrpcConnection interface {
	GetState() connectivity.State
	WaitForStateChange(ctx context.Context, sourceState connectivity.State) bool
}

func NewFixSimMarketDataClient(id string, targetAddress string, out chan<- *marketdata.MarketDataIncrementalRefresh,
	getConnection getMarketSimConnectionFn) (*fixSimMarketDataClient, error) {

	fixClient := &fixSimMarketDataClient{
		subscriptionsChan: make(chan string, 100),
		log:               log.New(log.Writer(), "target:"+targetAddress+" ", log.Flags()),
		errLog:            log.New(os.Stderr, "target:"+targetAddress+" ", log.Flags()),
	}

	log.Println("connecting to fix sim market data service at:" + targetAddress)

	client, conn, err := getConnection(targetAddress)
	if err != nil {
		return nil, err
	}

	streamChan := make(chan FixSimMarketDataService_ConnectClient, 1)

	go func() {
		subscriptions := map[string]bool{}

		var stream FixSimMarketDataService_ConnectClient
		for {

			select {
			case newStream := <-streamChan:
				stream = newStream
				if stream != nil {
					fixClient.log.Printf("new stream connected, resubscribing to all listings")
					for symbol := range subscriptions {
						err := stream.Send(&marketdata.MarketDataRequest{Parties: []*common.Parties{{PartyId: id}},
							InstrmtMdReqGrp: []*common.InstrmtMDReqGrp{{Instrument: &common.Instrument{Symbol: symbol}}}})

						if err != nil {
							fixClient.errLog.Printf("failed to resubscribe to all quotes using new stream: %v", err)
							break
						}
					}

					fixClient.log.Printf("resubscribed to all %v quotes", len(subscriptions))

				}
			case symbol := <-fixClient.subscriptionsChan:
				if !subscriptions[symbol] {
					subscriptions[symbol] = true
					if stream != nil {
						err := stream.Send(&marketdata.MarketDataRequest{Parties: []*common.Parties{{PartyId: id}},
							InstrmtMdReqGrp: []*common.InstrmtMDReqGrp{{Instrument: &common.Instrument{Symbol: symbol}}}})

						if err != nil {
							fixClient.errLog.Printf("failed so subscribe to symbol %v, error:%v", symbol, err)
						}
					}
				}

			}
		}

	}()

	go func() {

		for {
			state := conn.GetState()
			for state != connectivity.Ready {
				fixClient.log.Printf("waiting for fix sim market data connection to be ready....")

				conn.WaitForStateChange(context.Background(), state)
				state = conn.GetState()
				fixClient.log.Println("market gateway connection state is:", state)
			}

			stream, err := client.Connect(metadata.AppendToOutgoingContext(context.Background(), "subscriber_id", id))
			if err != nil {
				fixClient.errLog.Println("failed to connect to quote stream:", err)
				continue
			}

			fixClient.log.Println("connected to quote stream")

			streamChan <- stream

			for {
				incRefresh, err := stream.Recv()
				if err != nil {
					fixClient.errLog.Println("inbound stream error:", err)
					out <- nil
					break
				}
				out <- incRefresh
			}

			streamChan <- nil
		}
	}()

	return fixClient, nil
}

func (fsc *fixSimMarketDataClient) close() error {
	return fsc.conn.Close()
}

func (fsc *fixSimMarketDataClient) subscribe(symbol string) error {
	fsc.subscriptionsChan <- symbol
	return nil
}
