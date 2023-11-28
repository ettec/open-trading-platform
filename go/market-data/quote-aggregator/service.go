package main

import (
	"context"
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
	"log/slog"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})))

	id := bootstrap.GetEnvVar("GATEWAY_ID")
	toClientBufferSize := bootstrap.GetOptionalIntEnvVar("TO_CLIENT_BUFFER_SIZE", 1000)
	inboundQuoteBufferSize := bootstrap.GetOptionalIntEnvVar("INBOUND_QUOTE_BUFFER_SIZE", 1000)
	maxConnectRetry := time.Duration(bootstrap.GetOptionalIntEnvVar("MAX_CONNECT_RETRY_SECONDS", 60)) * time.Second
	inboundListingsBufferSize := bootstrap.GetOptionalIntEnvVar("INBOUND_LISTINGS_BUFFER_SIZE", 1000)

	http.Handle("/metrics", promhttp.Handler())
	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Panicf("failed to start metrics server:%v", err)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sds, err := staticdata.NewStaticDataSource(ctx)
	if err != nil {
		log.Panicf("failed to get static data source:%v", err)
	}

	mdsAddress, err := k8s.GetServiceAddress("market-data-service")
	if err != nil {
		log.Panicf("failed to get market data service address:%v", err)
	}

	mdsQuoteStream, err := marketdata.NewQuoteStreamFromMarketDataService(ctx, id, mdsAddress, maxConnectRetry, inboundQuoteBufferSize)
	if err != nil {
		log.Panicf("failed to get quote stream from market data service:%v", err)
	}

	quoteAggregator := quoteaggregator.New(ctx, sds.GetListingsWithSameInstrument, mdsQuoteStream, inboundListingsBufferSize)

	mdSource := marketdata.NewMarketDataSource(marketdata.NewQuoteDistributor(ctx, quoteAggregator, toClientBufferSize))

	port := "50551"
	slog.Info("Starting Quote Aggregator", "port", port)
	lis, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		log.Panicf("Error while listening : %v", err)
	}

	s := grpc.NewServer()

	marketdatasource.RegisterMarketDataSourceServer(s, mdSource)

	reflection.Register(s)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh,
		syscall.SIGKILL,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		<-sigCh
		s.GracefulStop()
	}()

	if err := s.Serve(lis); err != nil {
		log.Panicf("Error while serving : %v", err)
	}

}
