package actor

import (
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/connections"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/model"
	"log"
	"os"
	"time"
)


type newConnectionFn = func(connectionName string, out chan<- *model.ClobQuote) (connections.Connection, error)

type mdServerConnection struct {
	connectionName         string
	reconnectInterval      time.Duration
	subscriptionChan       chan int
	log                    *log.Logger
	errLog                 *log.Logger
	connectSignalChan      chan bool
	requestedSubscriptions map[int]bool
	subscriptions          map[int]bool
	connection             connections.Connection
	newConnectionFn        newConnectionFn
	quotesIn               <-chan *model.ClobQuote
	out                    chan<- *model.ClobQuote
}

func NewMdServerConnection( connectionName string,  out chan<- *model.ClobQuote, newConnection newConnectionFn, reconnectInterval time.Duration) *mdServerConnection {

	m := &mdServerConnection{
		connectionName:         connectionName,
		out:					out,
		reconnectInterval:      reconnectInterval,
		subscriptionChan:       make(chan int, 10000),
		log:                    log.New(os.Stdout, connectionName+":", log.Ltime | log.Lshortfile),
		errLog:                 log.New(os.Stderr, connectionName+":", log.Ltime | log.Lshortfile),
		connectSignalChan:      make(chan bool),
		requestedSubscriptions: map[int]bool{},
		subscriptions:          map[int]bool{},
		connection:             nil,
		newConnectionFn:        newConnection,
	}



	go func() {
		for {

			select {
			case quote, ok := <-m.quotesIn:
				if ok {
					m.out <- quote
				} else {
					m.log.Printf("inbound quote stream has closed, will attempt reconnect inChan %v seconds.", m.reconnectInterval)
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
				in := make( chan *model.ClobQuote, 10000)
				m.quotesIn = in
				connection, err := m.newConnectionFn(m.connectionName, in)
				if err == nil {
					m.connection = connection
					m.subscriptions = map[int]bool{}
					for s := range m.requestedSubscriptions {
						m.subscriptionChan <- s
					}
				} else {
					m.log.Printf("failed to Connect, will attempt reconnect inChan %v seconds.  Connection error:%v" , m.reconnectInterval.Seconds() ,err)
					go func() {
						time.Sleep(m.reconnectInterval)
						m.connectSignalChan <- true
					}()
				}
			}
		}
	}()

	m.connectSignalChan <- true

	return m
}


func (m *mdServerConnection) Subscribe(listingId int) {
	m.subscriptionChan <- listingId
}


