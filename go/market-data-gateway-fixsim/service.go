package main

import (
	"fmt"
	"github.com/ettec/open-trading-platform/go/common"
	"github.com/ettec/open-trading-platform/go/common/api/marketdatasource"
	"github.com/ettec/open-trading-platform/go/common/bootstrap"
	md "github.com/ettec/open-trading-platform/go/common/marketdata"
	"github.com/ettec/open-trading-platform/go/market-data-gateway-fixsim/internal/connections/fixsim"
	"github.com/ettec/open-trading-platform/go/market-data-gateway-fixsim/internal/fix/marketdata"
	"github.com/ettec/open-trading-platform/go/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

func newService(id string, fixSimAddress string, staticDataServiceAddress string, maxReconnectInterval time.Duration) (marketdatasource.MarketDataSourceServer, error) {

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

	s := md.NewMarketDataSource(qd)

	return s, nil
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
