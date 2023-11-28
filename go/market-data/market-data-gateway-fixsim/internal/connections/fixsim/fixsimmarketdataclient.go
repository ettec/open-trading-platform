package fixsim

import (
	"context"
	"github.com/ettec/open-trading-platform/go/market-data/market-data-gateway-fixsim/internal/fix/common"
	"github.com/ettec/open-trading-platform/go/market-data/market-data-gateway-fixsim/internal/fix/marketdata"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/metadata"
	"log/slog"
)

type fixSimMarketDataClient struct {
	subscriptionsChan chan string
	out               chan *marketdata.MarketDataIncrementalRefresh
}

func (fsc *fixSimMarketDataClient) Subscribe(symbol string) error {
	fsc.subscriptionsChan <- symbol
	return nil
}

func (fsc *fixSimMarketDataClient) Chan() <-chan *marketdata.MarketDataIncrementalRefresh {
	return fsc.out
}

type GrpcConnection interface {
	GetState() connectivity.State
	WaitForStateChange(ctx context.Context, sourceState connectivity.State) bool
}

func NewFixSimMarketDataClient(ctx context.Context, id string, client FixSimMarketDataServiceClient, conn GrpcConnection,
	outBufferSize int) (*fixSimMarketDataClient, error) {

	mdClient := &fixSimMarketDataClient{
		subscriptionsChan: make(chan string, 100),
		out:               make(chan *marketdata.MarketDataIncrementalRefresh, outBufferSize),
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
					slog.Info("new stream connected, resubscribing to all listings")
					for symbol := range subscriptions {
						err := stream.Send(&marketdata.MarketDataRequest{Parties: []*common.Parties{{PartyId: id}},
							InstrmtMdReqGrp: []*common.InstrmtMDReqGrp{{Instrument: &common.Instrument{Symbol: symbol}}}})

						if err != nil {
							slog.Error("failed to resubscribe to quote", "symbol", symbol, "error", err)
							break
						}
					}

					slog.Info("resubscribed to all quotes", "numSubscriptions", len(subscriptions))
				}
			case symbol := <-mdClient.subscriptionsChan:
				if !subscriptions[symbol] {
					subscriptions[symbol] = true
					if stream != nil {
						err := stream.Send(&marketdata.MarketDataRequest{Parties: []*common.Parties{{PartyId: id}},
							InstrmtMdReqGrp: []*common.InstrmtMDReqGrp{{Instrument: &common.Instrument{Symbol: symbol}}}})

						if err != nil {
							slog.Error("failed so subscribe to quote", "symbol", symbol, "error", err)
						}
					}
				}

			}
		}

	}()

	go func() {
		defer close(mdClient.out)
		for {
			state := conn.GetState()
			for state != connectivity.Ready {
				slog.Info("waiting for fix sim market data connection to be ready....")

				conn.WaitForStateChange(ctx, state)
				state = conn.GetState()
				slog.Info("market gateway connection state updated", "newState", state)
			}

			stream, err := client.Connect(metadata.AppendToOutgoingContext(ctx, "subscriber_id", id))
			if err != nil {
				slog.Error("failed to connect to fix market simulator", "error", err)
				continue
			}

			slog.Info("connected to fix market simulator")

			streamChan <- stream

			for {
				incRefresh, err := stream.Recv()
				if err != nil {
					slog.Error("error receiving from inbound stream", "error", err)
					mdClient.out <- nil
					break
				} else {
					mdClient.out <- incRefresh
				}

			}

			streamChan <- nil
		}
	}()

	return mdClient, nil
}
