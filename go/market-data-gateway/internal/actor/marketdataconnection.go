package actor

import (
	"fmt"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/connections"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/fix/marketdata"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/model"
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

type newConnection = func(connectionName string) connections.Connection

type mdServerConnection struct {
	actorImpl
	connectionName         string
	reconnectInterval      time.Duration
	subscriptionChan       chan int
	log                    *log.Logger
	connectSignalChan      chan bool
	requestedSubscriptions map[int]bool
	subscriptions          map[int]bool
	connection             connections.Connection
	newConnection          newConnection
	quotesIn               <-chan *model.ClobQuote
	out                    chan *model.ClobQuote
}

func NewMdServerConnection(address string, connectionName string, newConnection newConnection, reconnectInterval time.Duration) *mdServerConnection {

	m := &mdServerConnection{
		connectionName:         connectionName,
		reconnectInterval:      reconnectInterval,
		subscriptionChan:       make(chan int, 10000),
		log:                    log.New(os.Stdout, connectionName+":", log.LstdFlags),
		connectSignalChan:      make(chan bool, 1),
		requestedSubscriptions: map[int]bool{},
		subscriptions:          map[int]bool{},
		connection:             nil,
		newConnection:          newConnection,
	}

	m.actorImpl = newActorImpl("mdServerConnection for  "+address, m.readInputChannels)

	return m
}

func (m *mdServerConnection) Connect() (<-chan *model.ClobQuote, error) {

	if m.out != nil {
		return nil, fmt.Errorf("connect has already been called")
	}

	m.connectSignalChan <- true
	m.out = make(chan *model.ClobQuote, 1000)
	return m.out, nil
}


func (m *mdServerConnection) Subscribe(listingId int) {
	m.subscriptionChan <- listingId
}

func (m *mdServerConnection) readInputChannels() (chan<- bool, error) {

	select {
	case quote, ok := <-m.quotesIn:
		if ok {
			m.out <- quote
		} else {
			log.Printf("inbound quote stream has closed, will attempt reconnect in %v seconds.", m.reconnectInterval)
			m.quotesIn = nil
			go func() {
				time.Sleep(m.reconnectInterval)
				m.connectSignalChan <- true
			}()
		}

	case l := <-m.subscriptionChan:
		m.requestedSubscriptions[l] = true
		if !m.subscriptions[l] {
			if m.connection != nil {
				m.connection.Subscribe(l)
				m.subscriptions[l] = true
			}
		}
	case <-m.connectSignalChan:
		m.connection = m.newConnection(m.connectionName)
		in, err := m.connection.Connect()
		if err == nil {
			m.subscriptions = map[int]bool{}
			m.quotesIn = in
			for s := range m.requestedSubscriptions {
				m.subscriptionChan <- s
			}
		} else {
			m.log.Printf("failed to connect, will attempt reconnect in %v seconds.  Connection error:%v" , m.reconnectInterval.Seconds() ,err)
			go func() {
				time.Sleep(m.reconnectInterval)
				m.connectSignalChan <- true
			}()
		}

	case d := <-m.actorImpl.closeChan:
		if m.connection != nil {
			err := m.connection.Close()
			if err != nil {
				log.Printf("error whilst closing connection:%v", err)
			}
		}
		if m.out != nil {
			close(m.out)
		}

		return d, nil
	}

	return nil, nil
}

