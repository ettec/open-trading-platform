package main

import (
	"context"
	"fmt"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/actor"
	"github.com/ettech/open-trading-platform/go/market-data-service/api"

	"github.com/ettec/open-trading-platform/go/model"
	"github.com/ettech/open-trading-platform/go/market-data-service/internal"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

type service struct {
	partyIdToConnection map[string]actor.ClientConnection
	quoteDistributor    actor.QuoteDistributor
	connMux             sync.Mutex
}

func newService(id string, marketGatewayAddress string) *service {

	newConnection := func(connectionName string, out chan<- *model.ClobQuote) (actor.Connection, error) {
		return internal.NewMarketDataClient(connectionName, marketGatewayAddress, out )
	}

	serverToDistributorChan := make(chan *model.ClobQuote, 1000)

	mdConnection := actor.NewMdServerConnection(id, serverToDistributorChan, newConnection, 20*time.Second)
	qd := actor.NewQuoteDistributor(mdConnection.Subscribe, serverToDistributorChan)

	s := &service{partyIdToConnection: make(map[string]actor.ClientConnection), quoteDistributor: qd}

	return s
}

func (s *service) getConnection(partyId string) (actor.ClientConnection, bool) {
	s.connMux.Lock()
	defer s.connMux.Unlock()

	con, ok := s.partyIdToConnection[partyId]
	return con, ok
}

func (s *service) addConnection(subscriberId string, out chan<- *model.ClobQuote) (actor.ClientConnection, error) {
	s.connMux.Lock()
	defer s.connMux.Unlock()

	if conn, ok := s.partyIdToConnection[subscriberId]; ok {
		log.Printf("connection for client %v already exists, closing existing connection.", subscriberId)
		conn.Close()
		log.Print("connection closed:", subscriberId)
	}

	cc := actor.NewClientConnection(subscriberId, out, s.quoteDistributor, maxSubscriptions)

	s.partyIdToConnection[subscriberId] = cc

	return cc, nil
}

func (s *service) Subscribe(c context.Context, r *api.MdsSubscribeRequest) (*model.Empty, error) {

	if conn, ok := s.getConnection(r.SubscriberId); ok {

		if err := conn.Subscribe(r.ListingId); err == nil {
			return &model.Empty{}, nil
		} else {
			return nil, err
		}

	} else {
		return nil, fmt.Errorf("failed to subscribe, no connection exists for subscriber " + r.SubscriberId)
	}

}

func (s *service) Connect(request *api.MdsConnectRequest, stream api.MarketDataService_ConnectServer) error {

	subscriberId := request.GetSubscriberId()

	log.Println("connect request received for subscriber ", subscriberId)

	out := make(chan *model.ClobQuote, 100)

	s.addConnection(subscriberId, out)

	for mdUpdate := range out {
		if err := stream.Send(mdUpdate); err != nil {
			log.Printf("error on connection for subscriber %v, closing connection, error:%v", subscriberId, err)
			break
		}
	}

	return nil
}

const (
	ServiceIdKey   = "SERVICE_ID"
	GatewayAddress = "GATEWAY_ADDRESS"

	// The maximum number of listing subscriptions per connection
	MaxSubscriptionsKey = "MAX_SUBSCRIPTIONS"
)

var maxSubscriptions = 10000

func main() {

	port := "50651"
	fmt.Println("Starting Market Data Service on port:" + port)
	lis, err := net.Listen("tcp", "0.0.0.0:"+port)

	id, ok := os.LookupEnv(ServiceIdKey)
	if !ok {
		log.Panicf("expected %v env var to be set", ServiceIdKey)
	}

	fixSimAddress, ok := os.LookupEnv(GatewayAddress)
	if !ok {
		log.Panicf("expected %v env var to be set", GatewayAddress)
	}

	maxSubsEnv, ok := os.LookupEnv(MaxSubscriptionsKey)
	if ok {
		maxSubscriptions, err = strconv.Atoi(maxSubsEnv)
		if err != nil {
			log.Panicf("cannot parse %v, error: %v", MaxSubscriptionsKey, err)
		}
	}

	if err != nil {
		log.Fatalf("Error while listening : %v", err)
	}

	s := grpc.NewServer()
	api.RegisterMarketDataServiceServer(s, newService(id, fixSimAddress))

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error while serving : %v", err)
	}

}