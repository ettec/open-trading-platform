package main

import (
	"context"
	"fmt"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/actor"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/api"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/connections/fixsim"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/fix/marketdata"
	"github.com/ettec/open-trading-platform/go/model"
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

func newService(id string, fixSimAddress string, staticDataServiceAddress string, maxReconnectInterval time.Duration) (*service, error) {

	listingSrc, err := actor.NewListingSource(staticDataServiceAddress)
	if err != nil {
		return nil, err
	}

	newMarketDataClientFn := func(id string, out chan<- *marketdata.MarketDataIncrementalRefresh) (fixsim.MarketDataClient, error) {
		return fixsim.NewFixSimMarketDataClient(id, fixSimAddress, out, func(targetAddress string) (fixsim.FixSimMarketDataServiceClient, fixsim.GrpcConnection, error) {
			conn, err := grpc.Dial(targetAddress, grpc.WithInsecure(), grpc.WithBackoffMaxDelay(maxReconnectInterval))
			if err != nil {
				return nil, nil, err
			}

			client := fixsim.NewFixSimMarketDataServiceClient(conn)
			return client, conn, nil
		})
	}

	serverToDistributorChan := make(chan *model.ClobQuote, 1000)
	fixSimConn, err := fixsim.NewFixSimAdapter(newMarketDataClientFn, id, listingSrc.GetListing, serverToDistributorChan)
	if err != nil {
		return nil, err
	}

	qd := actor.NewQuoteDistributor(fixSimConn.Subscribe, serverToDistributorChan)

	s := &service{partyIdToConnection: make(map[string]actor.ClientConnection), quoteDistributor: qd}

	return s, nil
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

func (s *service) Subscribe(c context.Context, r *api.SubscribeRequest) (*model.Empty, error) {

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

func (s *service) Connect(request *api.ConnectRequest, stream api.MarketDataGateway_ConnectServer) error {

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
	GatewayIdKey             = "GATEWAY_ID"
	FixSimAddress            = "FIX_SIM_ADDRESS"
	StaticDataServiceAddress = "STATIC_DATA_SERVICE_ADDRESS"
	ConnectRetrySeconds      = "CONNECT_RETRY_SECONDS"

	// The maximum number of listing subscriptions per connection
	MaxSubscriptionsKey = "MAX_SUBSCRIPTIONS"
)

var maxSubscriptions = 10000

func main() {

	port := "50551"
	fmt.Println("Starting Client Market Data Gateway on port:" + port)
	lis, err := net.Listen("tcp", "0.0.0.0:"+port)

	id := getBootstrapEnvVar(GatewayIdKey)
	fixSimAddress := getBootstrapEnvVar(FixSimAddress)
	staticDataServiceAddress := getBootstrapEnvVar(StaticDataServiceAddress)
	connectRetrySecs := getOptionalBootstrapIntEnvVar(ConnectRetrySeconds, 60 )

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

	service, err := newService(id, fixSimAddress, staticDataServiceAddress, time.Duration(connectRetrySecs)*time.Second)
	if err != nil {
		log.Fatalf("error creating service: %v", err)
	}

	api.RegisterMarketDataGatewayServer(s, service)

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error while serving : %v", err)
	}

}

func getBootstrapEnvVar(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("missing required env var %v", key)
	}

	log.Printf("%v set to %v", key, value)

	return value
}

func getOptionalBootstrapIntEnvVar(key string, def int) int {
	strValue, exists := os.LookupEnv(key)
	result := def
	if exists {
		var err error
		result, err = strconv.Atoi(strValue)
		if err != nil {
			log.Panicf("cannot parse %v, error: %v", key, err)
		}
	}

	log.Printf("%v set to %v", key, result)

	return result
}
