package main

import (
	"fmt"
	"github.com/ettec/otp-common/api/marketdatasource"
	"github.com/ettec/otp-common/bootstrap"

	"github.com/ettec/open-trading-platform/go/market-data/market-data-gateway-fixsim/internal/connections/fixsim"
	"github.com/ettec/open-trading-platform/go/market-data/market-data-gateway-fixsim/internal/fix/marketdata"
	"github.com/ettec/otp-common/staticdata"
	md "github.com/ettec/otp-mdcommon"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	"time"
)

func newService(id string, fixSimAddress string, maxReconnectInterval time.Duration) (marketdatasource.MarketDataSourceServer, error) {

	listingSrc, err := staticdata.NewStaticDataSource(false)
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

	fixSimConn, err := fixsim.NewFixSimAdapter(newMarketDataClientFn, id, listingSrc.GetListing, 1000)
	if err != nil {
		return nil, err
	}

	qd := md.NewQuoteDistributor(fixSimConn, 100)

	s := md.NewMarketDataSource(qd)

	return s, nil
}

const (
	GatewayIdKey        = "GATEWAY_ID"
	FixSimAddress       = "FIX_SIM_ADDRESS"
	ConnectRetrySeconds = "CONNECT_RETRY_SECONDS"
)

var maxSubscriptions = 10000

func main() {

	port := "50551"
	fmt.Println("Starting Market Data Gateway on port:" + port)
	lis, err := net.Listen("tcp", "0.0.0.0:"+port)

	id := bootstrap.GetEnvVar(GatewayIdKey)
	fixSimAddress := bootstrap.GetEnvVar(FixSimAddress)
	connectRetrySecs := bootstrap.GetOptionalIntEnvVar(ConnectRetrySeconds, 60)

	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatalf("Error while listening : %v", err)
	}

	s := grpc.NewServer()

	service, err := newService(id, fixSimAddress, time.Duration(connectRetrySecs)*time.Second)
	if err != nil {
		log.Fatalf("error creating service: %v", err)
	}

	marketdatasource.RegisterMarketDataSourceServer(s, service)

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error while serving : %v", err)
	}

}
