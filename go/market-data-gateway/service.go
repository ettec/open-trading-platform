package actor

import (
	"context"
	"fmt"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/actor"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/connections"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/connections/fixsim"
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
	partyIdToConnection map[string]actor.ClientConnection
	quoteDistributor    actor.QuoteDistributor
	connMux             sync.Mutex
	actors []actor.Actor
}

func newService(address, name string) (*service, error) {

	listingIdToSymbol := map[int]string{1:"A", 2:"B", 3:"C", 4:"D"}
	newConnection := func(connectionName string) connections.Connection {

		return fixsim.NewFixSimConnection(address, connectionName, func(listingId int) (s string, err error) {
			if sym, ok := listingIdToSymbol[listingId]; ok {
				return sym, nil
			} else {
				return "", fmt.Errorf("symbol not found for listing id %v", listingId)
			}
		})

	}

	var actors []actor.Actor
	mdConnection := actor.NewMdServerConnection(name, newConnection, 20 * time.Second )
	actors = append(actors, mdConnection)

	qd := actor.NewQuoteDistributor(mdConnection)
	actors = append(actors, qd)

	s := &service{partyIdToConnection: make(map[string]actor.ClientConnection), quoteDistributor: qd, actors: actors}


	for _, actor := range actors {
		actor.Start()
	}

	return s, nil
}



func (s *service) getConnection(partyId string) (actor.ClientConnection, bool) {
	s.connMux.Lock()
	defer s.connMux.Unlock()

	con, ok := s.partyIdToConnection[partyId]
	return con, ok
}

func (s *service) addConnection(partyId string, connection actor.ClientConnection) {
	s.connMux.Lock()
	defer s.connMux.Unlock()

	s.partyIdToConnection[partyId] = connection
}

func (s *service) removeConnection(subscriberId string) {
	s.connMux.Lock()
	defer s.connMux.Unlock()

	delete(s.partyIdToConnection, subscriberId)
}


func (s *service) Subscribe(c context.Context, r *model.SubscribeRequest) (*empty.Empty, error) {

	if conn, ok := s.getConnection(r.SubscriberId); ok {
		conn.Subscribe(int(r.ListingId))
	} else {
		return nil, fmt.Errorf("failed to subscribe, no connection exists for subscriber " + r.SubscriberId)
	}

}

func (s *service) Connect(request *model.ConnectRequest, stream model.MarketDataGateway_ConnectServer) error {

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
