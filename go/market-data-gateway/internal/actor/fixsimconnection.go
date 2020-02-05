package actor

import (
	"context"
	"fmt"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/fix/common"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/fix/marketdata"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/fixsim"
	"google.golang.org/grpc"
	"log"
	"os"
)



type RefreshSink interface {
	SendRefresh(refresh *marketdata.MarketDataIncrementalRefresh)
}

type MdServerConnection interface {
	Start() Actor
	Close(chan<- bool)
	Subscribe(symbol string) error
}

type mdServerConnection struct {
	gatewayName       string
	address           string
	listingIdToSymbol map[int]string
	sink              RefreshSink
	log               *log.Logger
	client            fixsim.FixSimMarketDataServiceClient
	clientConn		  *grpc.ClientConn
}

func NewMdServerConnection(address string, gatewayName string, sink RefreshSink) *mdServerConnection {

	m := &mdServerConnection{
		gatewayName:       gatewayName,
		address:           address,
		listingIdToSymbol: make(map[int]string),
		sink:              sink,
		log:               log.New(os.Stdout, gatewayName+":", log.LstdFlags)}

	return m
}

func (m *mdServerConnection) Start() Actor {
	m.connect()
	return m
}

func (m *mdServerConnection) Close(chan<-bool)  {
	if m.clientConn != nil {
		m.clientConn.Close()
	}
}


func (m *mdServerConnection) Subscribe(symbol string) error {

	request := &marketdata.MarketDataRequest{Parties: []*common.Parties{{PartyId: m.gatewayName}},
		InstrmtMdReqGrp: []*common.InstrmtMDReqGrp{{Instrument: &common.Instrument{Symbol: symbol}}}}
	_, err := m.client.Subscribe(context.Background(), request)
	if err != nil {
		return fmt.Errorf("Failed to Subscribe to %v, error: %w ", symbol, err)
	}

	return nil
}

func (m *mdServerConnection) connect() {

	log.Println("Connecting to market data server at ", m.address)
	conn, err := grpc.Dial(m.address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		m.log.Println("Failed to dial the market data server:", err)
		return
	}

	m.clientConn = conn

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

		m.sink.SendRefresh(incRefresh)
	}

}
