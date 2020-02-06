package actor

import (
	"context"
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
	actorImpl
	gatewayName            string
	address                string
	reconnectInterval      time.Duration
	subscriptionChan       chan string
	sink                   RefreshSink
	log                    *log.Logger
	connectionChan         chan fixsim.FixSimMarketDataServiceClient
	connectSignalChan      chan bool
	requestedSubscriptions map[string]bool
	subscriptions          map[string]bool
	connection             fixsim.FixSimMarketDataServiceClient
}

func NewMdServerConnection(address string, gatewayName string, sink RefreshSink) *mdServerConnection {

	m := &mdServerConnection{
		gatewayName:            gatewayName,
		address:                address,
		reconnectInterval:      20 * time.Second,
		subscriptionChan:       make(chan string, 10000),
		sink:                   sink,
		log:                    log.New(os.Stdout, gatewayName+":", log.LstdFlags),
		connectionChan:         make(chan fixsim.FixSimMarketDataServiceClient),
		connectSignalChan:      make(chan bool, 1),
		requestedSubscriptions: map[string]bool{},
		subscriptions:          map[string]bool{},
		connection:             nil,
	}

	m.connectSignalChan <- true

	m.actorImpl = newActorImpl("mdServerConnection for  "+address, m.readInputChannels)

	return m
}

func (m *mdServerConnection) Subscribe(symbol string) {
	m.subscriptionChan <- symbol
}

func (m *mdServerConnection) readInputChannels() (chan<- bool, error) {

	select {
	case m.connection = <-m.connectionChan:
		if m.connection == nil {
			m.subscriptions = map[string]bool{}
			go func() {
				time.Sleep(m.reconnectInterval)
				m.connectSignalChan <- true
			}()
		} else {
			for s, _ := range m.requestedSubscriptions {
				m.subscriptionChan <- s
			}
		}
	case s := <-m.subscriptionChan:
		if !m.requestedSubscriptions[s] {
			m.requestedSubscriptions[s] = true
			if m.connection != nil {
				if err := subscribe(s, m.gatewayName, m.connection); err == nil {
					m.subscriptions[s] = true
				} else {
					m.log.Printf("failed to subscribe to symbol %v, error:%v", s, err)
				}
			}
		}
	case <-m.connectSignalChan:
		go m.connect(m.connectionChan)
	case d := <-m.actorImpl.closeChan:
		return d, nil
	}

	return nil, nil
}

func subscribe(symbol string,  subscriberId string, connection fixsim.FixSimMarketDataServiceClient) error {
	request := &marketdata.MarketDataRequest{Parties: []*common.Parties{{PartyId: subscriberId}},
		InstrmtMdReqGrp: []*common.InstrmtMDReqGrp{{Instrument: &common.Instrument{Symbol: symbol}}}}
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
	defer func() {
		conn.Close()
		connectionChan <- nil
	}()

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

}
