package main

import (
	"context"
	"fmt"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/actor"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/connections"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/connections/fixsim"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/fix/marketdata"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/model"
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
	partyIdToConnection map[string]*clientConnection
	quoteDistributor    actor.QuoteDistributor
	connMux             sync.Mutex
}

type clientConnection struct {
	id              string
	connection      actor.ClientConnection
	distributorChan chan *model.ClobQuote
	subscriptionCnt int
	out             chan<- *model.ClobQuote
}

func newService(id string, fixSimAddress string) *service {

	listingIdToSymbol := map[int]string{1: "A", 2: "B", 3: "C", 4: "D"}

	newConnection := func(connectionName string, out chan<- *model.ClobQuote) (connections.Connection, error) {

		newMarketDataClient := func(id string, out chan<- *marketdata.MarketDataIncrementalRefresh) (fixsim.MarketDataClient, error) {
			return fixsim.NewFixSimMarketDataClient(id, fixSimAddress, out)
		}

		return fixsim.NewFixSimConnection(newMarketDataClient, connectionName, func(listingId int) (s string, err error) {
			if sym, ok := listingIdToSymbol[listingId]; ok {
				return sym, nil
			} else {
				return "", fmt.Errorf("symbol not found for listing id %v", listingId)
			}
		}, out)

	}

	serverToDistributorChan := make(chan *model.ClobQuote, 1000)

	mdConnection := actor.NewMdServerConnection(id, serverToDistributorChan, newConnection, 20*time.Second)
	qd := actor.NewQuoteDistributor(mdConnection.Subscribe, serverToDistributorChan)

	s := &service{partyIdToConnection: make(map[string]*clientConnection), quoteDistributor: qd}

	return s
}

func (s *service) getConnection(partyId string) (*clientConnection, bool) {
	s.connMux.Lock()
	defer s.connMux.Unlock()

	con, ok := s.partyIdToConnection[partyId]
	return con, ok
}

func (s *service) addConnection(subscriberId string, out chan<- *model.ClobQuote) (*clientConnection, error) {
	s.connMux.Lock()
	defer s.connMux.Unlock()

	if conn, ok := s.partyIdToConnection[subscriberId]; ok {
		log.Printf("connection for client %v already exists, closing existing connection.", subscriberId)
		conn.connection.Close()
		s.quoteDistributor.RemoveOutQuoteChan(conn.distributorChan)
		close(conn.out)
		log.Print("connection closed:", subscriberId)
	}

	distributorToConnectionChan := make(chan *model.ClobQuote, 1000)
	connection := actor.NewClientConnection(subscriberId, out, s.quoteDistributor.Subscribe, distributorToConnectionChan, maxSubscriptions)
	cc := &clientConnection{
		id:              subscriberId,
		connection:      connection,
		distributorChan: distributorToConnectionChan,
		out:			 out,
	}

	s.partyIdToConnection[subscriberId] = cc
	s.quoteDistributor.AddOutQuoteChan(distributorToConnectionChan)

	return cc, nil
}

func (s *service) Subscribe(c context.Context, r *model.SubscribeRequest) (*model.Empty, error) {

	if conn, ok := s.getConnection(r.SubscriberId); ok {

		if conn.subscriptionCnt >= maxSubscriptions {
			return nil, fmt.Errorf("maximum subscription count %v exceeded", maxSubscriptions)
		}

		conn.connection.Subscribe(int(r.ListingId))
		conn.subscriptionCnt++
		return &model.Empty{}, nil
	} else {
		return nil, fmt.Errorf("failed to subscribe, no connection exists for subscriber " + r.SubscriberId)
	}

}

func (s *service) Connect(request *model.ConnectRequest, stream model.MarketDataGateway_ConnectServer) error {

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
	GatewayIdKey  = "GATEWAY_ID"
	FixSimAddress = "FIX_SIM_ADDRESS"

	// The maximum number of listing subscriptions per connection
	MaxSubscriptionsKey = "MAX_SUBSCRIPTIONS"
)

var maxSubscriptions = 10000

func main() {

	port := "50551"
	fmt.Println("Starting Client Market Data Gateway on port:" + port)
	lis, err := net.Listen("tcp", "0.0.0.0:"+port)

	id, ok := os.LookupEnv(GatewayIdKey)
	if !ok {
		log.Panicf("expected %v env var to be set", GatewayIdKey)
	}

	fixSimAddress, ok := os.LookupEnv(FixSimAddress)
	if !ok {
		log.Panicf("expected %v env var to be set", FixSimAddress)
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
	model.RegisterMarketDataGatewayServer(s, newService(id, fixSimAddress))

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error while serving : %v", err)
	}

}
