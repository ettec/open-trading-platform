package marketdata

import (
	"context"
	"github.com/ettec/open-trading-platform/go/common/api/marketdatasource"
	"github.com/ettec/open-trading-platform/go/model"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/metadata"
	"log"
	"os"
)

type MdsQuoteStream interface {
	Subscribe(listingId int32)
	GetStream() <-chan *model.ClobQuote
}

type mdsQuoteStream struct {
	conn              *grpc.ClientConn
	out               chan *model.ClobQuote
	subscriptionsChan chan int32
	log               *log.Logger
	errLog            *log.Logger
}

type GetMdsClientFn = func(targetAddress string) (marketdatasource.MarketDataSourceClient, GrpcConnection, error)

type GrpcConnection interface {
	GetState() connectivity.State
	WaitForStateChange(ctx context.Context, sourceState connectivity.State) bool
}

func NewMdsQuoteStream(id string, targetAddress string, maxReconnectInterval time.Duration,
	quoteBufferSize int) (MdsQuoteStream, error) {

	out := make(chan *model.ClobQuote, quoteBufferSize)

	mdcFn := func(targetAddress string) (marketdatasource.MarketDataSourceClient, GrpcConnection, error) {
		conn, err := grpc.Dial(targetAddress, grpc.WithInsecure(), grpc.WithBackoffMaxDelay(maxReconnectInterval))
		if err != nil {
			return nil, nil, err
		}

		client := marketdatasource.NewMarketDataSourceClient(conn)
		return client, conn, nil
	}

	return NewMdsQuoteStreamFromFn(id, targetAddress, out, mdcFn)
}

func NewMdsQuoteStreamFromFn(id string, targetAddress string, out chan *model.ClobQuote,
	getConnection GetMdsClientFn) (MdsQuoteStream, error) {

	n := &mdsQuoteStream{
		subscriptionsChan: make(chan int32, 100),
		out: out,
		log:               log.New(os.Stdout, "target:"+targetAddress+" ", log.Lshortfile|log.Ltime),
		errLog:            log.New(os.Stderr, "target:"+targetAddress+" ", log.Lshortfile|log.Ltime),
	}

	log.Println("connecting to market data source at " + targetAddress)

	client, conn, err := getConnection(targetAddress)
	if err != nil {
		return nil, err
	}

	streamChan := make(chan marketdatasource.MarketDataSource_ConnectClient, 1)

	go func() {
		subscriptions := map[int32]bool{}

		var stream marketdatasource.MarketDataSource_ConnectClient
		for {

			select {
			case newStream := <-streamChan:
				stream = newStream
				if stream != nil {
					log.Printf("new stream connected, resubscribing to all listings")
					for listingId := range subscriptions {
						err := stream.Send(&marketdatasource.SubscribeRequest{
							ListingId: listingId,
						})

						if err != nil {
							n.errLog.Printf("failed to resubscribe to all quotes using new stream: %v", err)
							break
						}
					}

					n.log.Printf("resubscribed to all %v quotes", len(subscriptions))

				} else {
					log.Printf("stream connection lost, sending empty quotes to all subscriptions")
					for listingId, subscribed := range subscriptions {
						if subscribed {
							out <- &model.ClobQuote{
								ListingId:         listingId,
								Bids:              []*model.ClobLine{},
								Offers:            []*model.ClobLine{},
								StreamInterrupted: true,
								StreamStatusMsg:   "market data source client stream interrupted",
							}
						}
					}

				}
			case listingId := <-n.subscriptionsChan:
				if !subscriptions[listingId] {
					subscriptions[listingId] = true
					if stream != nil {
						err := stream.Send(&marketdatasource.SubscribeRequest{
							ListingId: listingId,
						})

						if err != nil {
							n.errLog.Printf("failed so subscribe to listing %v, error:%v", listingId, err)
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
				n.log.Printf("waiting for market data source connection to be ready....")

				conn.WaitForStateChange(context.Background(), state)
				state = conn.GetState()
				n.log.Println("market data source connection state is:", state)
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
					break
				}
				out <- incRefresh
			}

			streamChan <- nil
		}
	}()

	return n, nil
}

func (mgc *mdsQuoteStream) GetStream() <-chan *model.ClobQuote {
	return mgc.out
}

func (mgc *mdsQuoteStream) Subscribe(listingId int32) {
	mgc.subscriptionsChan <- listingId
}
