package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/fix/common"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/fix/marketdata"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/fixsim"
	"google.golang.org/grpc"
)

type listingIdSymbol struct {
	listingId int
	symbol    string
}

type mdServerConnection struct {
	gatewayName       string
	listingIdToSymbol map[int]string
	subscribeChan     chan int
	symbolLookupChan  chan listingIdSymbol
	incRefreshChan    chan *marketdata.MarketDataIncrementalRefresh
	log               *log.Logger
}

type refresh marketdata.MarketDataIncrementalRefresh

func NewMdServerConnection(address string, gatewayName string) (*mdServerConnection, error) {

	m := &mdServerConnection{
		gatewayName,
		make(map[int]string), make(chan int),
		make(chan listingIdSymbol),
		make(chan *marketdata.MarketDataIncrementalRefresh),
		log.New(os.Stdout, gatewayName+":", log.LstdFlags)}

	go m.startSubscriptionHandler(gatewayName, address)

	return m, nil
}

func (m *mdServerConnection) startMarketDataServerConnection(address string) {

	log.Println("Connecting to market data server at ", address)
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		m.log.Println("Failed to dial the market data server:", err)
		return
	}
	defer conn.Close()

	r := &fixsim.ConnectRequest{PartyId: m.gatewayName}
	mdClient := fixsim.NewFixSimMarketDataServiceClient(conn)
	stream, err := mdClient.Connect(context.Background(), r)
	if err != nil {
		m.log.Println("Failed to connect to the market data server:", err)
		return
	}

	for {
		incRefresh, err := stream.Recv()
		if err != nil {
			m.log.Println("market data stream error:", err)
			break
		}

		m.incRefreshChan <- incRefresh
		incRefresh = <-m.incRefreshChan
	}

}

type mdupdate struct {
	listingIdToSymbol *listingIdSymbol
	refresh           *refresh
}






func (m *mdServerConnection) startSubscriptionHandler(address string, connectionId string) {
	m.log.Println("subscription handler started")
	log.Println("Connecting to market data server at ", address)
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		m.log.Println("Failed to dial the market data server:", err)
		return
	}
	defer conn.Close()

	simClient := fixsim.NewFixSimMarketDataServiceClient(conn)

	for {
		select {
		case l := <-m.subscribeChan:
			if _, ok := m.listingIdToSymbol[l]; !ok {
				m.getSymbol(l, m.symbolLookupChan)
			}
		case ls := <-m.symbolLookupChan:

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			m.log.Println("subscribing to ", ls.symbol)
			request := &marketdata.MarketDataRequest{Parties: []*common.Parties{&common.Parties{PartyId: connectionId}},
				InstrmtMdReqGrp: []*common.InstrmtMDReqGrp{&common.InstrmtMDReqGrp{Instrument: &common.Instrument{Symbol: ls.symbol}}}}
			_, err := simClient.Subscribe(ctx, request)
			if err != nil {
				m.log.Println("Failed to subscribe to {}, error: {} ", ls.symbol, err)
				return
			} else {
				m.listingIdToSymbol[ls.listingId] = ls.symbol
				m.log.Println("subscribed to ", ls.symbol)
			}
		}
	}
}

func (m *mdServerConnection) Close() {

}

func (m *mdServerConnection) getSymbol(listingId int, resultChan chan<- listingIdSymbol) {
	// TODO goto database
}

func (m *mdServerConnection) addConnection(c *connection) {

}

func (m *mdServerConnection) subscribe(listingId int) {
	m.subscribeChan <- listingId
}
