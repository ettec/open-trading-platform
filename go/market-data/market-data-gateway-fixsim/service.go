package main

import (
	"fmt"
	"github.com/ettec/otp-common/api/marketdatasource"
	"github.com/ettec/otp-common/bootstrap"
	"os"

	"github.com/ettec/open-trading-platform/go/market-data/market-data-gateway-fixsim/internal/connections/fixsim"
	"github.com/ettec/open-trading-platform/go/market-data/market-data-gateway-fixsim/internal/fix/marketdata"
	md "github.com/ettec/otp-common/marketdata"
	"github.com/ettec/otp-common/staticdata"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	"time"
)

func newService(id string, fixSimAddress string, maxReconnectInterval time.Duration,
	inboundQuoteBufferSize int, clientQuoteBufferSize int) (marketdatasource.MarketDataSourceServer, error) {

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

	fixSimConn, err := fixsim.NewFixSimAdapter(newMarketDataClientFn, id, listingSrc.GetListing,
		inboundQuoteBufferSize)
	if err != nil {
		return nil, err
	}

	qd := md.NewQuoteDistributor(fixSimConn, clientQuoteBufferSize)

	s := md.NewMarketDataSource(qd)

	return s, nil
}

func main() {

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.Lshortfile)

	id := bootstrap.GetEnvVar("GATEWAY_ID")
	fixSimAddress := bootstrap.GetEnvVar("FIX_SIM_ADDRESS")
	maxConnectRetrySecs := bootstrap.GetOptionalIntEnvVar("MAX_CONNECT_RETRY_SECONDS", 60)
	inboundQuoteBufferSize := bootstrap.GetOptionalIntEnvVar("INBOUND_QUOTE_BUFFER_SIZE", 1000)
	clientQuoteBufferSize := bootstrap.GetOptionalIntEnvVar("CLIENT_QUOTE_BUFFER_SIZE", 1000)

	port := "50551"
	fmt.Println("Starting Market Data Gateway on port:" + port)
	lis, err := net.Listen("tcp", "0.0.0.0:"+port)

	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Panicf("Error while listening : %v", err)
	}

	s := grpc.NewServer()

	service, err := newService(id, fixSimAddress, time.Duration(maxConnectRetrySecs)*time.Second, inboundQuoteBufferSize,
		clientQuoteBufferSize)
	if err != nil {
		log.Panicf("error creating service: %v", err)
	}

	marketdatasource.RegisterMarketDataSourceServer(s, service)

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Panicf("Error while serving : %v", err)
	}

}
