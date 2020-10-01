package main

import (
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
	"net"
	"os"
	"strings"
	"time"
)

func main() {

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.Lshortfile)

	id := bootstrap.GetEnvVar("ID")
	maxConnectRetry := time.Duration(bootstrap.GetOptionalIntEnvVar("MAX_CONNECT_RETRY_SECONDS", 60)) * time.Second
	kafkaBrokersString := bootstrap.GetEnvVar("KAFKA_BROKERS")

	log.Print("Starting vwap strategy")

	s := grpc.NewServer()

	kafkaBrokers := strings.Split(kafkaBrokersString, ",")

	sds, err := staticdata.NewStaticDataSource(false)
	if err != nil {
		log.Fatalf("failed to create static data source:%v", err)
	}

	clientSet := k8s.GetK8sClientSet(false)

	orderRouter, err := api.GetOrderRouter(clientSet, maxConnectRetry)
	if err != nil {
		panic(err)
	}

	executeFn := func(om *strategy.Strategy) {

		om.Log.Printf("execute strategy for params: %v", om.ParentOrder.GetExecParametersJson())

		vwapParamsJson := om.ParentOrder.GetExecParametersJson()

		listingIn := make(chan *model.Listing)
		om.Log.Printf("fetching listing %v.....", om.ParentOrder.ListingId)

		sds.GetListing(om.ParentOrder.ListingId, listingIn)
		listing := <-listingIn
		om.Log.Printf("got listing %v", listing)

		quantity := om.ParentOrder.Quantity

		vwapParams := &vwapParameters{}
		err := json.Unmarshal([]byte(vwapParamsJson), vwapParams)
		if err != nil {
			om.CancelChan <- fmt.Sprintf("failed to parse parameters:%v", err)
		}

		if vwapParams.UtcStartTimeSecs <= 0 || vwapParams.UtcEndTimeSecs <= 0 {
			om.CancelChan <-"invalid start or end time specified"
		}

		if vwapParams.UtcStartTimeSecs >= vwapParams.UtcEndTimeSecs {
			om.CancelChan <- "start time must be before end time"
		}

		if vwapParams.UtcEndTimeSecs < time.Now().Unix() {
			om.CancelChan <-"end time has already passed"
		}

		if model.IasD(vwapParams.Buckets).GreaterThan(quantity) {
			om.CancelChan <-"num Buckets must be less than or equal to the quantity"
		}

		buckets, err := getBucketsFromParamsString(vwapParamsJson, *quantity, listing)
		if err != nil {
			om.CancelChan <-fmt.Sprintf("failed to get Buckets from params:%v", err)
		}

		executeAsVwapStrategy(om, buckets, listing)
	}

	store, err := orderstore.NewKafkaStore(orderstore.DefaultReaderConfig(common.ORDERS_TOPIC, kafkaBrokers),
		orderstore.DefaultWriterConfig(common.ORDERS_TOPIC, kafkaBrokers), id)

	if err != nil {
		panic(fmt.Errorf("failed to create order store: %v", err))
	}

	childOrderUpdates, err := ordermanagement.GetChildOrders(id, orderstore.DefaultReaderConfig(common.ORDERS_TOPIC, kafkaBrokers),
		bootstrap.GetOptionalIntEnvVar("VWAPSTRATEGY_CHILD_ORDER_UPDATES_BUFFER_SIZE", 1000))

	if err != nil {
		panic(err)
	}

	distributor := ordermanagement.NewChildOrderUpdatesDistributor(childOrderUpdates)

	sm := strategy.NewStrategyManager(id, store, distributor, orderRouter, executeFn)

	executionvenue.RegisterExecutionVenueServer(s, sm)

	reflection.Register(s)

	port := "50551"
	fmt.Println("Starting Execution Venue Service on port:" + port)
	lis, err := net.Listen("tcp", "0.0.0.0:"+port)

	if err != nil {
		log.Fatalf("Error while listening : %v", err)
	}

	if err := s.Serve(lis); err != nil {
		log.Fatalf("error   while serving : %v", err)
	}

}
