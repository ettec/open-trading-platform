package main

import (
	"context"
	"fmt"
	"github.com/ettec/otp-common/api/marketdatasource"
	"github.com/ettec/otp-common/bootstrap"
	"os"

	"github.com/ettec/open-trading-platform/go/market-data/market-data-gateway-fixsim/internal/connections/fixsim"
	md "github.com/ettec/otp-common/marketdata"
	"github.com/ettec/otp-common/staticdata"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func newService(ctx context.Context, id string, fixSimAddress string, maxReconnectInterval time.Duration,
	inboundQuoteBufferSize int, clientQuoteBufferSize int) (marketdatasource.MarketDataSourceServer, error) {

	listingSrc, err := staticdata.NewStaticDataSource(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create static data source: %w", err)
	}

	conn, err := grpc.Dial(fixSimAddress, grpc.WithInsecure(), grpc.WithBackoffMaxDelay(maxReconnectInterval))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to fix sim at %s: %w", fixSimAddress, err)
	}

	grpcClient := fixsim.NewFixSimMarketDataServiceClient(conn)

	fixSimClient, err := fixsim.NewFixSimMarketDataClient(ctx, id, grpcClient, conn, clientQuoteBufferSize)
	if err != nil {
		return nil, fmt.Errorf("failed to create fix simulator market data client: %w", err)
	}

	fixSimQuoteStream, err := fixsim.NewQuoteStreamFromFixClient(ctx, fixSimClient, id, listingSrc.GetListing,
		inboundQuoteBufferSize)
	if err != nil {
		return nil, fmt.Errorf("failed to create fix quote stream: %w", err)
	}

	qd := md.NewQuoteDistributor(ctx, fixSimQuoteStream, clientQuoteBufferSize)

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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service, err := newService(ctx, id, fixSimAddress, time.Duration(maxConnectRetrySecs)*time.Second, inboundQuoteBufferSize,
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
