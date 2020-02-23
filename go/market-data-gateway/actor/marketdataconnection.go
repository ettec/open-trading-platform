package actor

import (
	"github.com/ettec/open-trading-platform/go/model"
	"log"
	"os"
	"time"
)


type Connection interface {

	Subscribe(listingId int32) error
	Close() error
}

type NewConnectionFn = func(connectionName string, out chan<- *model.ClobQuote) (Connection, error)



type mdServerConnection struct {
	connectionName         string
	reconnectInterval      time.Duration
	subscriptionChan       chan int32
	log                    *log.Logger
	errLog                 *log.Logger
	connectSignalChan      chan bool
	requestedSubscriptions map[int32]bool
	subscriptions          map[int32]bool
	lastQuote              map[int32]*model.ClobQuote
	connection             Connection
	newConnectionFn        NewConnectionFn
	quotesIn               <-chan *model.ClobQuote
	out                    chan<- *model.ClobQuote
}

func NewMdServerConnection( connectionName string,  out chan<- *model.ClobQuote, newConnection NewConnectionFn, reconnectInterval time.Duration) *mdServerConnection {

	m := &mdServerConnection{
		connectionName:         connectionName,
		out:					out,
		reconnectInterval:      reconnectInterval,
		subscriptionChan:       make(chan int32, 10000),
		lastQuote:				map[int32]*model.ClobQuote{},
		log:                    log.New(os.Stdout, connectionName+":", log.Ltime | log.Lshortfile),
		errLog:                 log.New(os.Stderr, connectionName+":", log.Ltime | log.Lshortfile),
		connectSignalChan:      make(chan bool),
		requestedSubscriptions: map[int32]bool{},
		subscriptions:          map[int32]bool{},
		connection:             nil,
		newConnectionFn:        newConnection,
	}



	go func() {
		for {

			select {
			case quote, ok := <-m.quotesIn:
				if ok {
					m.lastQuote[quote.ListingId] = quote
					m.out <- quote
				} else {
					m.log.Println("inbound quote stream has closed")

					// Send empty quotes upstream
					for k := range m.subscriptions {
						emptyQuote := &model.ClobQuote{
							ListingId:            k,
						}
						m.lastQuote[k] = emptyQuote
						m.out <- emptyQuote
					}

					m.log.Printf("will attempt reconnect inChan %v seconds.", m.reconnectInterval)

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
						err := m.connection.Subscribe(l)
						if err != nil {
							m.errLog.Println("failed to subscribe:", err)
						} else {
							m.subscriptions[l] = true
						}

					}
				}
			case <-m.connectSignalChan:
				m.log.Println("attempting to connect...")
				in := make( chan *model.ClobQuote, 10000)
				m.quotesIn = in
				connection, err := m.newConnectionFn(m.connectionName, in)
				if err == nil {
					m.log.Println("connected, sending subscriptions")
					m.connection = connection
					m.subscriptions = map[int32]bool{}
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


func (m *mdServerConnection) Subscribe(listingId int32) {
	m.subscriptionChan <- listingId
}


