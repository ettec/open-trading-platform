package main

import (
	"fmt"
	common "github.com/ettec/otp-common"
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
	log.SetFlags(log.Ltime|log.Lshortfile)

	id := bootstrap.GetEnvVar("ID")

	maxConnectRetry := time.Duration(bootstrap.GetOptionalIntEnvVar("MAX_CONNECT_RETRY_SECONDS", 60)) * time.Second
	kafkaBrokersString := bootstrap.GetEnvVar("KAFKA_BROKERS")

	s := grpc.NewServer()

	kafkaBrokers := strings.Split(kafkaBrokersString, ",")

	sds, err := staticdata.NewStaticDataSource(false)
	if err != nil {
		log.Fatalf("failed to create static data source:%v", err)
	}

	mdsAddress, err := k8s.GetServiceAddress("market-data-service")
	if err != nil {
		panic(err)
	}

	mdsQuoteStream, err := marketdata.NewQuoteStreamFromMdService(id, mdsAddress, maxConnectRetry, 1000)
	if err != nil {
		panic(err)
	}

	qd := marketdata.NewQuoteDistributor(mdsQuoteStream, 100)

	orderRouter, err := strategy.GetOrderRouter(k8s.GetK8sClientSet(false), maxConnectRetry)
	if err != nil {
		panic(err)
	}

	executeFn := func(om *strategy.Strategy) {
		ExecuteAsSmartRouterStrategy(om, sds.GetListingsWithSameInstrument, qd.GetNewQuoteStream())
	}


	store, err := orderstore.NewKafkaStore(orderstore.DefaultReaderConfig(common.ORDERS_TOPIC, kafkaBrokers),
		orderstore.DefaultWriterConfig(common.ORDERS_TOPIC, kafkaBrokers), id)

	if err != nil {
		panic(fmt.Errorf("failed to create order store: %v", err))
	}

	childOrderUpdates, err := ordermanagement.GetChildOrders(id, orderstore.DefaultReaderConfig(common.ORDERS_TOPIC, kafkaBrokers),
	 bootstrap.GetOptionalIntEnvVar("SMARTROUTER_CHILD_ORDER_UPDATES_BUFFER_SIZE", 1000))
	if err != nil {
		panic(err)
	}

	distributor := ordermanagement.NewChildOrderUpdatesDistributor(childOrderUpdates)

	sm := strategy.NewStrategyManager(id, store, distributor, orderRouter, executeFn)

	executionvenue.RegisterExecutionVenueServer(s, sm)

	reflection.Register(s)

	port := "50551"
	fmt.Println("Starting Smart Router on port:" + port)
	lis, err := net.Listen("tcp", "0.0.0.0:"+port)

	if err != nil {
		log.Panicf("Error while listening : %v", err)
	}

	if err := s.Serve(lis); err != nil {
		log.Panicf("error   while serving : %v", err)
	}

}
