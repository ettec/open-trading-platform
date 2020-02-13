package actor

import (
	"context"
	"fmt"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/actor"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/connections"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/connections/fixsim"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/fix/marketdata"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/model"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"sync"
	"time"
)

type service struct {
	partyIdToConnection map[string]clientConnection
	quoteDistributor    actor.QuoteDistributor
	connMux             sync.Mutex
}

type clientConnection struct {
	id string
	connection actor.ClientConnection
	distributorChan chan *model.ClobQuote
}

const maxSubscriptions = 10000

func newService(address, name string) (*service, error) {

	listingIdToSymbol := map[int]string{1:"A", 2:"B", 3:"C", 4:"D"}
	newConnection := func(connectionName string, out chan<- *model.ClobQuote) (connections.Connection, error) {


		newMarketDataClient := func(id string, out chan<- *marketdata.MarketDataIncrementalRefresh) (fixsim.MarketDataClient,error) {
			return fixsim.NewFixSimMarketDataClient(id, address, out)
		}


		return fixsim.NewFixSimConnection(newMarketDataClient, connectionName, func(listingId int) (s string, err error) {
			if sym, ok := listingIdToSymbol[listingId]; ok {
				return sym, nil
			} else {
				return "", fmt.Errorf("symbol not found for listing id %v", listingId)
			}
		}, out)

	}

	serverToDistributorChan := make( chan *model.ClobQuote, 1000 )

	mdConnection := actor.NewMdServerConnection(name, serverToDistributorChan, newConnection,  20 * time.Second )
	qd := actor.NewQuoteDistributor(mdConnection.Subscribe, serverToDistributorChan)

	s := &service{partyIdToConnection: make(map[string]clientConnection), quoteDistributor: qd}

	return s, nil
}



func (s *service) getConnection(partyId string) (actor.ClientConnection, bool) {
	s.connMux.Lock()
	defer s.connMux.Unlock()

	con, ok := s.partyIdToConnection[partyId]
	return con, ok
}

func (s *service) addConnection(subscriberId string, stream model.MarketDataGateway_ConnectServer) (actor.ClientConnection, error) {
	s.connMux.Lock()
	defer s.connMux.Unlock()

	if _, ok := s.partyIdToConnection[subscriberId];ok {
		return nil, fmt.Errorf("connection already exists for subscriber id " + subscriberId)
	}

	distributorToConnectionChan := make(chan *model.ClobQuote, 1000)
	connection := actor.NewClientConnection(subscriberId, stream.Send, s.quoteDistributor.Subscribe,  distributorToConnectionChan, maxSubscriptions)
	cc := clientConnection{
		id:              subscriberId,
		connection:      connection,
		distributorChan: distributorToConnectionChan,
	}

	s.partyIdToConnection[subscriberId] = cc
	s.quoteDistributor.AddOutQuoteChan(distributorToConnectionChan)
}

func (s *service) removeConnection(subscriberId string) error {
	s.connMux.Lock()
	defer s.connMux.Unlock()


	if cc, ok := s.partyIdToConnection[subscriberId];!ok {
		return  fmt.Errorf("no connection exists for subscriber id %v" , subscriberId)
	} else {
		s.quoteDistributor.RemoveOutQuoteChan(cc.distributorChan)
		cc.connection.Close()
		delete(s.partyIdToConnection, subscriberId)
	}

	return nil
}


func (s *service) Subscribe(c context.Context, r *model.SubscribeRequest) (*empty.Empty, error) {

	if conn, ok := s.getConnection(r.SubscriberId); ok {
		conn.Subscribe(int(r.ListingId))
	} else {
		return nil, fmt.Errorf("failed to subscribe, no connection exists for subscriber " + r.SubscriberId)
	}

}

func (s *service) Connect(request *model.ConnectRequest, stream model.MarketDataGateway_ConnectServer) error {


	here - wire this up and then test


	stream.

	subscriberId := request.GetSubscriberId()

	if conn, ok := s.getConnection(subscriberId); ok {
		log.Printf("connection for client %v already exists, closing existing connection.", subscriberId)
		s.removeConnection(subscriberId)
		done := make( chan bool, 1)
		conn.Close(done)
		<-done
		log.Print("connection closed:", subscriberId)
	}

	stream.Context()

	log.Println("creating client connection for ", subscriberId)
	cc := actor.NewClientConnection(subscriberId,  stream, s.quoteDistributor, 1000)
	cc.start()
	s.quoteDistributor.AddConnection(cc)
	s.addConnection(subscriberId, cc)

	return nil
}

func main() {

	port := "50551"
	fmt.Println("Starting Client Market Data Gateway on port:" + port)
	lis, err := net.Listen("tcp", "0.0.0.0:"+port)

	if err != nil {
		log.Fatalf("Error while listening : %v", err)
	}

	s := grpc.NewServer()
	model.RegisterMarketDataGatewayServer(s, newService())

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error while serving : %v", err)
	}

}
