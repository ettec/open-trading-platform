package main

import (
	"context"
	"fmt"
	common "github.com/ettec/otp-common"
	"github.com/ettec/otp-common/api"
	"github.com/ettec/otp-common/api/executionvenue"
	"github.com/ettec/otp-common/k8s"
	"github.com/ettec/otp-common/marketdata"
	"github.com/ettec/otp-common/ordermanagement"
	"github.com/ettec/otp-common/staticdata"
	"github.com/ettec/otp-common/strategy"
	"os"
	"time"

	"github.com/ettec/otp-common/bootstrap"

	"github.com/ettec/otp-common/orderstore"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"strings"
)

func main() {

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.Lshortfile)

	id := bootstrap.GetEnvVar("ID")

	maxConnectRetry := time.Duration(bootstrap.GetOptionalIntEnvVar("MAX_CONNECT_RETRY_SECONDS", 60)) * time.Second
	kafkaBrokersString := bootstrap.GetEnvVar("KAFKA_BROKERS")

	s := grpc.NewServer()

	kafkaBrokers := strings.Split(kafkaBrokersString, ",")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sds, err := staticdata.NewStaticDataSource(ctx)
	if err != nil {
		log.Fatalf("failed to create static data source:%v", err)
	}

	mdsAddress, err := k8s.GetServiceAddress("market-data-service")
	if err != nil {
		log.Panicf("failed to market data service address: %v", err)
	}

	mdsQuoteStream, err := marketdata.NewQuoteStreamFromMarketDataService(ctx, id, mdsAddress, maxConnectRetry,
		bootstrap.GetOptionalIntEnvVar("SMARTROUTER_INBOUND_QUOTE_BUFFER_SIZE", 1000))
	if err != nil {
		log.Panicf("failed to create quote stream from market data service: %v", err)
	}

	qd := marketdata.NewQuoteDistributor(ctx, mdsQuoteStream, bootstrap.GetOptionalIntEnvVar("SMARTROUTER_QUOTE_DISTRIBUTOR_BUFFER_SIZE", 1000))

	orderRouter, err := api.GetOrderRouter(k8s.GetK8sClientSet(false), maxConnectRetry)
	if err != nil {
		log.Panicf("failed to get order router: %v", err)
	}

	executeFn := func(om *strategy.Strategy) {
		ExecuteAsSmartRouterStrategy(ctx, om, sds.GetListingsWithSameInstrument, qd.NewQuoteStream())
	}

	store, err := orderstore.NewKafkaStore(orderstore.DefaultReaderConfig(common.ORDERS_TOPIC, kafkaBrokers),
		orderstore.DefaultWriterConfig(common.ORDERS_TOPIC, kafkaBrokers), id)

	if err != nil {
		log.Panicf("failed to create order store: %v", err)
	}

	childOrderUpdates, err := ordermanagement.GetChildOrders(ctx, id, orderstore.DefaultReaderConfig(common.ORDERS_TOPIC, kafkaBrokers),
		bootstrap.GetOptionalIntEnvVar("SMARTROUTER_CHILD_ORDER_UPDATES_BUFFER_SIZE", 1000))
	if err != nil {
		log.Panicf("failed to get child order updates channel: %v", err)
	}

	distributor := ordermanagement.NewChildOrderUpdatesDistributor(childOrderUpdates, 10000)

	sm, err := strategy.NewStrategyManager(ctx, id, store, distributor, orderRouter, executeFn)
	if err != nil {
		log.Panicf("failed to create strategy manager: %v", err)
	}

	executionvenue.RegisterExecutionVenueServer(s, sm)

	reflection.Register(s)

	port := "50551"
	fmt.Println("Starting Smart Router on port:" + port)
	lis, err := net.Listen("tcp", "0.0.0.0:"+port)

	if err != nil {
		log.Panicf("failed to listen on port %s : %v", port, err)
	}

	if err := s.Serve(lis); err != nil {
		log.Panicf("error while serving: %v", err)
	}

}
