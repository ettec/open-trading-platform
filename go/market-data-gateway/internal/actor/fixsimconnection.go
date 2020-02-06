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
	"time"
)

type RefreshSink interface {
	SendRefresh(refresh *marketdata.MarketDataIncrementalRefresh)
}

type MdServerConnection interface {
	Actor
	Subscribe(symbol string)
}

type mdServerConnection struct {
	gatewayName       string
	address           string
	subscriptionChan  chan string
	sink              RefreshSink
	log               *log.Logger
}

func NewMdServerConnection(address string, gatewayName string, sink RefreshSink) *mdServerConnection {

	m := &mdServerConnection{
		gatewayName:       gatewayName,
		address:           address,
		subscriptionChan:  make(chan string, 10000),
		sink:              sink,
		log:               log.New(os.Stdout, gatewayName+":", log.LstdFlags)}

	return m
}

func (m *mdServerConnection) Start() {

	m.connect()

}

func (m *mdServerConnection) Close(chan<- bool) {
	if m.clientConn != nil {
		m.clientConn.Close()
	}
}

func (m *mdServerConnection) Subscribe(symbol string)  {
	m.subscriptionChan<-symbol
}

func (m *mdServerConnection) run() {

	connectionChan := make(chan fixsim.FixSimMarketDataServiceClient)

	go m.connect(connectionChan)

	subscriptions := map[string]bool{}
	subscribed := map[string]bool{}

	connection := <- connectionChan

	reconnectTimer := time.NewTicker(10 * time.Second)

	for {
		select {
		case connection = <-connectionChan:
			if connection == nil {
				subscribed = map[string]bool{}
			} else {
				for s, _:= range subscriptions {
					m.subscriptionChan<-s
				}
			}
		case s := <-m.subscriptionChan:
			if !subscriptions[s] {
				subscriptions[s] = true
				if connection != nil {
					if err := subscribe( s, connection); err == nil {
						subscribed[s]=true
					} else {
						m.log.Printf("failed to subscribe to %v, error:%v", s, err)
					}
				}
			}
		case <-reconnectTimer.C:
			if connection == nil {
				here
			}


		}
	}

}

func subscribe(s string, connection fixsim.FixSimMarketDataServiceClient) error {
	request := &marketdata.MarketDataRequest{Parties: []*common.Parties{{PartyId: m.gatewayName}},
		InstrmtMdReqGrp: []*common.InstrmtMDReqGrp{{Instrument: &common.Instrument{Symbol: s}}}}
	_, err := connection.Subscribe(context.Background(), request)
	return err

}

func (m *mdServerConnection) connect(connectionChan chan fixsim.FixSimMarketDataServiceClient) {

	log.Println("Connecting to market data server at ", m.address)
	conn, err := grpc.Dial(m.address, grpc.WithInsecure(), grpc.WithBlock())
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

	connectionChan <- mdClient

	for {
		incRefresh, err := stream.Recv()
		if err != nil {
			m.log.Println("market data stream error:", err)
			break
		}

		m.sink.SendRefresh(incRefresh)
	}

	connectionChan <- nil

}
