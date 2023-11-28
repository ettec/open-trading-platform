package main

import (
	"context"
	"encoding/json"
	"fmt"
	common "github.com/ettec/otp-common"
	"github.com/ettec/otp-common/api"
	"github.com/ettec/otp-common/api/executionvenue"
	"github.com/ettec/otp-common/bootstrap"
	"github.com/ettec/otp-common/k8s"
	"github.com/ettec/otp-common/model"
	"github.com/ettec/otp-common/ordermanagement"
	"github.com/ettec/otp-common/orderstore"
	"github.com/ettec/otp-common/staticdata"
	"github.com/ettec/otp-common/strategy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})))

	id := bootstrap.GetEnvVar("ID")
	maxConnectRetry := time.Duration(bootstrap.GetOptionalIntEnvVar("MAX_CONNECT_RETRY_SECONDS", 60)) * time.Second
	kafkaBrokersString := bootstrap.GetEnvVar("KAFKA_BROKERS")

	slog.Info("Starting vwap strategy")

	kafkaBrokers := strings.Split(kafkaBrokersString, ",")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sds, err := staticdata.NewStaticDataSource(ctx)
	if err != nil {
		log.Panicf("failed to create static data source:%v", err)
	}

	clientSet := k8s.GetK8sClientSet(false)

	orderRouter, err := api.GetOrderRouter(clientSet, maxConnectRetry)
	if err != nil {
		panic(err)
	}

	executeFn := func(om *strategy.Strategy) {

		om.Log.Info("executing strategy", "params", om.ParentOrder.GetExecParametersJson())

		vwapParamsJson := om.ParentOrder.GetExecParametersJson()

		listingIn := make(chan staticdata.ListingResult)
		om.Log.Info("fetching listing", "listingId", om.ParentOrder.ListingId)

		sds.GetListing(ctx, om.ParentOrder.ListingId, listingIn)
		listingResult := <-listingIn
		if listingResult.Err != nil {
			om.CancelChan <- fmt.Sprintf("failed to get listing %v", listingResult.Err)
		}

		om.Log.Info("got listing", "listingId", listingResult.Listing)
		close(listingIn)

		quantity := om.ParentOrder.Quantity

		vwapParams := &vwapParameters{}
		err := json.Unmarshal([]byte(vwapParamsJson), vwapParams)
		if err != nil {
			om.CancelChan <- fmt.Sprintf("failed to parse parameters:%v", err)
		}

		if vwapParams.UtcStartTimeSecs <= 0 || vwapParams.UtcEndTimeSecs <= 0 {
			om.CancelChan <- "invalid start or end time specified"
		}

		if vwapParams.UtcStartTimeSecs >= vwapParams.UtcEndTimeSecs {
			om.CancelChan <- "start time must be before end time"
		}

		if vwapParams.UtcEndTimeSecs < time.Now().Unix() {
			om.CancelChan <- "end time has already passed"
		}

		if model.IasD(vwapParams.Buckets).GreaterThan(quantity) {
			om.CancelChan <- "num Buckets must be less than or equal to the quantity"
		}

		buckets, err := getBucketsFromParamsString(vwapParamsJson, *quantity, listingResult.Listing)
		if err != nil {
			om.CancelChan <- fmt.Sprintf("failed to get Buckets from params:%v", err)
		}

		executeAsVwapStrategy(ctx, om, buckets, listingResult.Listing)
	}

	store, err := orderstore.NewKafkaStore(orderstore.DefaultReaderConfig(common.ORDERS_TOPIC, kafkaBrokers),
		orderstore.DefaultWriterConfig(common.ORDERS_TOPIC, kafkaBrokers), id)

	if err != nil {
		log.Panicf("failed to create order store: %v", err)
	}

	childOrderUpdates, err := ordermanagement.GetChildOrders(ctx, id, orderstore.DefaultReaderConfig(common.ORDERS_TOPIC, kafkaBrokers),
		bootstrap.GetOptionalIntEnvVar("VWAPSTRATEGY_CHILD_ORDER_UPDATES_BUFFER_SIZE", 1000))

	if err != nil {
		log.Panicf("failed to create child order updates channel:%v", err)
	}

	distributor := ordermanagement.NewChildOrderUpdatesDistributor(childOrderUpdates, 10000)

	sm, err := strategy.NewStrategyManager(ctx, id, store, distributor, orderRouter, executeFn)
	if err != nil {
		log.Panicf("failed to create strategy manager:%v", err)
	}

	s := grpc.NewServer()
	executionvenue.RegisterExecutionVenueServer(s, sm)

	reflection.Register(s)

	port := "50551"
	slog.Info("Starting Execution Venue Service", "port", port)
	lis, err := net.Listen("tcp", "0.0.0.0:"+port)

	if err != nil {
		log.Panicf("Error while listening : %v", err)
	}

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
		log.Panicf("error   while serving : %v", err)
	}

}
