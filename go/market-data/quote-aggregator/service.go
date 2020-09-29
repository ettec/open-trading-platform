package main

import (
	"github.com/ettec/open-trading-platform/go/market-data/quote-aggregator/quoteaggregator"
	"github.com/ettec/otp-common/api/marketdatasource"
	"github.com/ettec/otp-common/bootstrap"
	"github.com/ettec/otp-common/k8s"
	"github.com/ettec/otp-common/marketdata"
	"github.com/ettec/otp-common/staticdata"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)



func main() {

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime|log.Lshortfile)

	id := bootstrap.GetEnvVar("GATEWAY_ID")
	toClientBufferSize := bootstrap.GetOptionalIntEnvVar("TO_CLIENT_BUFFER_SIZE", 1000)
	inboundQuoteBufferSize := bootstrap.GetOptionalIntEnvVar("INBOUND_QUOTE_BUFFER_SIZE", 1000)
	maxConnectRetry := time.Duration(bootstrap.GetOptionalIntEnvVar("MAX_CONNECT_RETRY_SECONDS", 60)) * time.Second

	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":8080", nil)

	sds, err := staticdata.NewStaticDataSource(false)
	if err != nil {
		panic(err)
	}

	mdsAddress, err := k8s.GetServiceAddress("market-data-service")
	if err != nil {
		panic(err)
	}

	mdsQuoteStream, err := marketdata.NewQuoteStreamFromMdService(id, mdsAddress, maxConnectRetry, inboundQuoteBufferSize)
	if err != nil {
		panic(err)
	}

	quoteAggregator := quoteaggregator.New(sds.GetListingsWithSameInstrument, mdsQuoteStream)

	mdSource := marketdata.NewMarketDataSource(marketdata.NewQuoteDistributor(quoteAggregator, toClientBufferSize))

	port := "50551"
	log.Println("Starting Quote Aggregator on port:" + port)
	lis, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		log.Panicf("Error while listening : %v", err)
	}

	s := grpc.NewServer()

	marketdatasource.RegisterMarketDataSourceServer(s, mdSource)

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Panicf("Error while serving : %v", err)
	}

}
