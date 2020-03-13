package main

import (
	"fmt"
	"github.com/ettec/open-trading-platform/go/common/bootstrap"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/actor"
	mdgapi "github.com/ettec/open-trading-platform/go/market-data-gateway/api"
	"github.com/ettec/open-trading-platform/go/model"
	"github.com/ettech/open-trading-platform/go/market-data-service/api"
	"github.com/ettech/open-trading-platform/go/market-data-service/gatewayclient"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"sync"
	"time"
)

type service struct {
	quoteDistributor    actor.QuoteDistributor
	connMux             sync.Mutex
}

func newService(id string, marketGatewayAddress string, maxReconnectInterval time.Duration) (*service, error) {

	mdcToDistributorChan := make(chan *model.ClobQuote, 1000)

	mdcFn := func(targetAddress string) (mdgapi.MarketDataGatewayClient, gatewayclient.GrpcConnection, error) {
		conn, err := grpc.Dial(targetAddress, grpc.WithInsecure(), grpc.WithBackoffMaxDelay(maxReconnectInterval))
		if err != nil {
			return nil, nil, err
		}

		client := mdgapi.NewMarketDataGatewayClient(conn)
		return client, conn, nil
	}

	mdc, err := gatewayclient.NewMarketDataGatewayClient(id, marketGatewayAddress, mdcToDistributorChan, mdcFn)

	if err != nil {
		return nil, err
	}

	qd := actor.NewQuoteDistributor(mdc.Subscribe, mdcToDistributorChan)
	s := &service{ quoteDistributor: qd}

	return s, nil
}

const SubscriberIdKey = "subscriber_id"

func (s *service) Connect(stream api.MarketDataService_ConnectServer) error {

	ctx, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return fmt.Errorf("failed to retrieve call context")
	}

	values := ctx.Get(SubscriberIdKey)
	if len(values) != 1 {
		return fmt.Errorf("must specify string value for %v", SubscriberIdKey)
	}

	fromClientId := values[0]
	subscriberId := fromClientId + ":" + uuid.New().String()

	log.Printf("connect request received for subscriber %v, unique connection id: %v ", fromClientId, subscriberId, )

	out := make(chan *model.ClobQuote, 100)
	cc := actor.NewClientConnection(subscriberId, out, s.quoteDistributor, maxSubscriptions)
	defer cc.Close()

	go func() {
		for {
			subscription, err := stream.Recv()

			if err != nil {
				log.Printf("subscriber:%v inbound stream error:%v ", subscriberId, err)
				break
			} else {
				log.Printf("subscribe request, subscriber id:%v, listing id:%v", subscriberId, subscription.ListingId)
				cc.Subscribe(subscription.ListingId)
			}
		}
	}()

	for mdUpdate := range out {
		if err := stream.Send(mdUpdate); err != nil {
			log.Printf("error on connection for subscriber %v, closing connection, error:%v ", subscriberId, err)
			break
		}
	}

	return nil
}

const (
	ServiceIdKey   = "SERVICE_ID"
	GatewayAddress = "GATEWAY_ADDRESS"
	ConnectRetrySeconds      = "CONNECT_RETRY_SECONDS"

)

var maxSubscriptions = 10000

func main() {

	port := "50561"
	fmt.Println("Starting Market Data Service on port:" + port)
	lis, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		log.Fatalf("Error while listening : %v", err)
	}

	id := bootstrap.GetEnvVar(ServiceIdKey)

	fixSimAddress:= bootstrap.GetEnvVar(GatewayAddress)

	connectRetrySecs := bootstrap.GetOptionalIntEnvVar(ConnectRetrySeconds, 60 )

	s := grpc.NewServer()
	mdcService, err := newService(id, fixSimAddress, time.Duration(connectRetrySecs)*time.Second)
	if err != nil {
		log.Panicf("failed to create market data service:%v", err)
	}

	api.RegisterMarketDataServiceServer(s, mdcService)

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error while serving : %v", err)
	}

}


