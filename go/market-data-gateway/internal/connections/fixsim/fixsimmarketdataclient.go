package fixsim

import (
	"context"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/fix/common"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/fix/marketdata"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"log"
	"os"
	"sync"
)

type fixSimMarketDataClient struct {
	id            string
	client        FixSimMarketDataServiceClient
	conn          *grpc.ClientConn
	out           chan<- *marketdata.MarketDataIncrementalRefresh
	subscribeMux  sync.Mutex
	subscriptions map[string]bool
	log           *log.Logger
	errLog        *log.Logger
}

type getMarketSimConnectionFn = func (targetAddress string) (FixSimMarketDataServiceClient, GrpcConnection, error)

type GrpcConnection interface {
	GetState() connectivity.State
	WaitForStateChange(ctx context.Context, sourceState connectivity.State) bool
}

func NewFixSimMarketDataClient(id string, targetAddress string, out chan<- *marketdata.MarketDataIncrementalRefresh,
	getConnection getMarketSimConnectionFn) (*fixSimMarketDataClient, error) {

	n := &fixSimMarketDataClient{
		id:     id,
		out:    out,
		subscriptions: map[string]bool{},
		log:    log.New(os.Stdout, "target:" + targetAddress + " ", log.Lshortfile|log.Ltime),
		errLog: log.New(os.Stderr, "target:" + targetAddress + " ", log.Lshortfile|log.Ltime),
	}

	log.Println("connecting to fix sim market data service at:" + targetAddress)

	client, conn,  err := getConnection(targetAddress)
	if err != nil {
		return nil, err
	}

	n.client = client

	go func() {

		for {
			state := conn.GetState()
			for state != connectivity.Ready {
				n.log.Printf("waiting for fix market sim connection to be ready....")
				conn.WaitForStateChange(context.Background(), state)
				state = conn.GetState()
				n.log.Println("market sim connection state is:", state)
			}

			stream, err := connect(n.client, id)
			if err != nil {
				n.errLog.Println("failed to connect to quote stream:", err)
				continue
			}

			n.log.Println("connected to quote stream")


			err = n.resubscribeAllSymbols()
			if err != nil {
				n.errLog.Println("failed to resubscribe to all symbols:", err)
				continue
			}


			if len(n.subscriptions) >0 {
				n.log.Printf("resubscribed to all %v quotes", len(n.subscriptions))
			}

			for {
				incRefresh, err := stream.Recv()
				if err != nil {
					n.errLog.Println("inbound stream error:", err)
					n.out <- nil
					break
				}
				n.out <- incRefresh
			}
		}
	}()

	return n, nil
}

func (fsc *fixSimMarketDataClient) close() error {
	return fsc.conn.Close()
}

func (fsc *fixSimMarketDataClient) resubscribeAllSymbols() error {
	fsc.subscribeMux.Lock()
	defer fsc.subscribeMux.Unlock()
	for symbol := range fsc.subscriptions {
		err := subscribe(fsc.client, symbol, fsc.id)
		if err != nil {
			return err
		}
	}
	return nil
}

func connect(client FixSimMarketDataServiceClient, id string) (FixSimMarketDataService_ConnectClient, error) {
	r := &Party{PartyId: id}
	stream, err := client.Connect(context.Background(), r)
	return stream, err
}


func subscribe(client FixSimMarketDataServiceClient, symbol string, id string ) error {
	request := &marketdata.MarketDataRequest{Parties: []*common.Parties{{PartyId: id}},
		InstrmtMdReqGrp: []*common.InstrmtMDReqGrp{{Instrument: &common.Instrument{Symbol: symbol}}}}
	_, err := client.Subscribe(context.Background(), request)

	return err
}

func (fsc *fixSimMarketDataClient) subscribe(symbol string) error {

	fsc.subscribeMux.Lock()
	defer fsc.subscribeMux.Unlock()
	if !fsc.subscriptions[symbol] {
		fsc.subscriptions[symbol] = true
		return subscribe(fsc.client, symbol, fsc.id)
	}

	return nil
}

