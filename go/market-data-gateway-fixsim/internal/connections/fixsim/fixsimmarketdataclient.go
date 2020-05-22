package fixsim

import (
	"context"
	"github.com/ettec/open-trading-platform/go/market-data-gateway-fixsim/internal/fix/common"
	"github.com/ettec/open-trading-platform/go/market-data-gateway-fixsim/internal/fix/marketdata"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/metadata"
	"log"
	"os"
)

var outboundSubscriptions = promauto.NewCounter(prometheus.CounterOpts{
	Name: "outbound_subscriptions",
	Help: "The number of outbound subscriptions",
})

var quotesReceived = promauto.NewCounter(prometheus.CounterOpts{
	Name: "quotes_received",
	Help: "The number of quotes received from all streams",
})

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

	n := &fixSimMarketDataClient{
		subscriptionsChan: make(chan string, 100),
		log:               log.New(os.Stdout, "target:"+targetAddress+" ", log.Lshortfile|log.Ltime),
		errLog:            log.New(os.Stderr, "target:"+targetAddress+" ", log.Lshortfile|log.Ltime),
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
					log.Printf("new stream connected, resubscribing to all listings")
					for symbol := range subscriptions {
						err := stream.Send(&marketdata.MarketDataRequest{Parties: []*common.Parties{{PartyId: id}},
							InstrmtMdReqGrp: []*common.InstrmtMDReqGrp{{Instrument: &common.Instrument{Symbol: symbol}}}})

						if err != nil {
							n.errLog.Printf("failed to resubscribe to all quotes using new stream: %v", err)
							break
						}
					}

					n.log.Printf("resubscribed to all %v quotes", len(subscriptions))

				}
			case symbol := <-n.subscriptionsChan:
				if !subscriptions[symbol] {
					subscriptions[symbol] = true
					outboundSubscriptions.Inc()
					if stream != nil {
						err := stream.Send(&marketdata.MarketDataRequest{Parties: []*common.Parties{{PartyId: id}},
							InstrmtMdReqGrp: []*common.InstrmtMDReqGrp{{Instrument: &common.Instrument{Symbol: symbol}}}})

						if err != nil {
							n.errLog.Printf("failed so subscribe to symbol %v, error:%v", symbol, err)
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
				n.log.Printf("waiting for fix sim market data connection to be ready....")

				conn.WaitForStateChange(context.Background(), state)
				state = conn.GetState()
				n.log.Println("market gateway connection state is:", state)
			}

			stream, err := client.Connect(metadata.AppendToOutgoingContext(context.Background(), "subscriber_id", id))
			if err != nil {
				n.errLog.Println("failed to connect to quote stream:", err)
				continue
			}

			n.log.Println("connected to quote stream")

			streamChan <- stream

			for {
				incRefresh, err := stream.Recv()
				if err != nil {
					n.errLog.Println("inbound stream error:", err)
					out <- nil
					break
				}
				out <- incRefresh
				quotesReceived.Inc()
			}

			streamChan <- nil
		}
	}()

	return n, nil
}

func (fsc *fixSimMarketDataClient) close() error {
	return fsc.conn.Close()
}

func (fsc *fixSimMarketDataClient) subscribe(symbol string) error {
	fsc.subscriptionsChan <- symbol
	return nil
}
