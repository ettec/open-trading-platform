package actor

import (
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/fix/marketdata"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/fixsim"
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


here - where we need to change this so it take a quote clob source, therefore normalisation done in fix sim package


type MarketDataClient interface {
	Connect(connectionId string) (IncRefreshSource, error)
	Subscribe(symbol string, subscriberId string) error
	Close() error
}



type mdServerConnection struct {
	actorImpl
	connectionName         string
	address                string
	reconnectInterval      time.Duration
	subscriptionChan       chan string
	sink                   RefreshSink
	log                    *log.Logger
	connectionChan         chan MarketDataClient
	connectSignalChan      chan bool
	requestedSubscriptions map[string]bool
	subscriptions          map[string]bool
	connection             MarketDataClient
	dial                   fixsim.dial
}

func NewMdServerConnection(address string, connectionName string, sink RefreshSink, connectionDial fixsim.dial, reconnectInterval time.Duration) *mdServerConnection {

	m := &mdServerConnection{
		connectionName:         connectionName,
		address:                address,
		reconnectInterval:      reconnectInterval,
		subscriptionChan:       make(chan string, 10000),
		sink:                   sink,
		log:                    log.New(os.Stdout, connectionName+":", log.LstdFlags),
		connectionChan:         make(chan MarketDataClient),
		connectSignalChan:      make(chan bool, 1),
		requestedSubscriptions: map[string]bool{},
		subscriptions:          map[string]bool{},
		connection:             nil,
		dial:                   connectionDial,
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
			for s := range m.requestedSubscriptions {
				m.subscriptionChan <- s
			}
		}
	case s := <-m.subscriptionChan:
		m.requestedSubscriptions[s] = true
		if !m.subscriptions[s]  {

			if m.connection != nil {
				if err := m.connection.Subscribe(s, m.connectionName); err == nil {
					m.subscriptions[s] = true
				} else {
					m.log.Printf("failed to subscribe to symbol %v, error:%v", s, err)
				}
			}
		}
	case <-m.connectSignalChan:
		go m.connect(m.connectionChan, m.dial)
	case d := <-m.actorImpl.closeChan:
		if m.connection != nil {
			m.connection.Close()
		}
		return d, nil
	}

	return nil, nil
}

func (m *mdServerConnection) connect(connectionChan chan MarketDataClient, dial fixsim.dial) {

	log.Println("Connecting to market data server at ", m.address)
	mdClient, err := dial(m.address)
	if err != nil {
		m.log.Println("Failed to dial the market data server:", err)
		return
	}
	defer func() {
		if err := mdClient.Close(); err != nil {
			m.log.Println("error whilst closing:", err)
		}

		connectionChan <- nil
	}()

	stream, err := mdClient.Connect(m.connectionName)
	if err != nil {
		m.log.Println("Failed to connect to the market data server:", err)
		return
	}

	connectionChan <- mdClient

	for {
		incRefresh, err := stream.Recv()
		if err != nil {
			m.log.Println("market data sendQuoteFn error:", err)
			break
		}

		m.sink.SendRefresh(incRefresh)
	}

}
