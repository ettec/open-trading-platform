package main

import (
	"context"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/fix/marketdata"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/fixsim"
	"google.golang.org/grpc"
	"log"
	"os"
)

type listingIdSymbol struct {
	listingId int
	symbol    string
}

type mdServerConnection struct {
	gatewayName       string
	listingIdToSymbol map[int]string

	incRefreshChan    chan *marketdata.MarketDataIncrementalRefresh
	log               *log.Logger
}

type refresh marketdata.MarketDataIncrementalRefresh

func NewMdServerConnection(address string, gatewayName string) (*mdServerConnection, error) {

	m := &mdServerConnection{
		gatewayName,
		make(map[int]string),
		make(chan *marketdata.MarketDataIncrementalRefresh),
		log.New(os.Stdout, gatewayName+":", log.LstdFlags)}

	//go m.startReadLoop(gatewayName, address)

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
	}

}

type mdupdate struct {
	listingIdToSymbol *listingIdSymbol
	refresh           *refresh
}





func (m *mdServerConnection) Close() {

}

func (m *mdServerConnection) fetchSymbol(listingId int, resultChan chan<- listingIdSymbol) {
	// TODO goto database
}

func (m *mdServerConnection) addConnection(c *connection) {

}

func (m *mdServerConnection) subscribe(listingId int) {
	//m.subscribeChan <- listingId
}
