package main

import (
	"fmt"
	"github.com/ettec/open-trading-platform/go/common"
	"github.com/ettec/open-trading-platform/go/common/bootstrap"
	md "github.com/ettec/open-trading-platform/go/common/marketdata"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/api/marketdatasource"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/connections/fixsim"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/fix/marketdata"
	"github.com/ettec/open-trading-platform/go/model"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

type service struct {
	quoteDistributor md.QuoteDistributor
	connMux          sync.Mutex
}

func newService(id string, fixSimAddress string, staticDataServiceAddress string, maxReconnectInterval time.Duration) (*service, error) {

	listingSrc, err := common.NewStaticDataSource(staticDataServiceAddress)
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

	qd := md.NewQuoteDistributor(fixSimConn.Subscribe, serverToDistributorChan)

	s := &service{quoteDistributor: qd}

	return s, nil
}

const SubscriberIdKey = "subscriber_id"

func (s *service) Connect(stream marketdatasource.MarketDataSource_ConnectServer) error {

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

	log.Printf("connect request received for subscriber %v, unique connection id: %v ", fromClientId, subscriberId)

	out := make(chan *model.ClobQuote, 100)
	cc := md.NewConflatedQuoteConnection(subscriberId, out, s.quoteDistributor, maxSubscriptions)
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
	fmt.Println("Starting Market Data Gateway on port:" + port)
	lis, err := net.Listen("tcp", "0.0.0.0:"+port)

	id := bootstrap.GetEnvVar(GatewayIdKey)
	fixSimAddress := bootstrap.GetEnvVar(FixSimAddress)
	staticDataServiceAddress := bootstrap.GetEnvVar(StaticDataServiceAddress)
	connectRetrySecs := bootstrap.GetOptionalIntEnvVar(ConnectRetrySeconds, 60)

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

	marketdatasource.RegisterMarketDataSourceServer(s, service)

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error while serving : %v", err)
	}

}
