package main

import (
	"encoding/json"
	"fmt"
	common "github.com/ettec/otp-common"
	api "github.com/ettec/otp-common/api/executionvenue"
	"github.com/ettec/otp-common/bootstrap"
	"github.com/ettec/otp-common/executionvenue"
	"github.com/ettec/otp-common/k8s"
	"github.com/ettec/otp-common/model"
	"github.com/ettec/otp-common/orderstore"
	"github.com/ettec/otp-common/staticdata"
	"github.com/ettec/otp-common/strategy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	logger "log"
	"net"
	"os"
	"strings"
	"time"
)

const (
	KafkaBrokersKey        = "KAFKA_BROKERS"
	MaxConnectRetrySeconds = "MAX_CONNECT_RETRY_SECONDS"
)

var log = logger.New(os.Stdout, "", logger.Ltime|logger.Lshortfile)

func main() {

	id := common.SR_MIC
	maxConnectRetry := time.Duration(bootstrap.GetOptionalIntEnvVar(MaxConnectRetrySeconds, 60)) * time.Second
	kafkaBrokersString := bootstrap.GetEnvVar(KafkaBrokersKey)

	s := grpc.NewServer()

	kafkaBrokers := strings.Split(kafkaBrokersString, ",")

	sds, err := staticdata.NewStaticDataSource(false)
	if err != nil {
		log.Fatalf("failed to create static data source:%v", err)
	}

	clientSet := k8s.GetK8sClientSet(false)

	orderRouter, err := strategy.GetOrderRouter(clientSet, maxConnectRetry)
	if err != nil {
		panic(err)
	}

	executeFn := func(om *strategy.Strategy) {

		vwapParamsJson := om.ParentOrder.GetExecParametersJson()

		listingIn := make(chan *model.Listing)
		listing := <-listingIn
		sds.GetListing(om.ParentOrder.ListingId, listingIn )

		quantity := om.ParentOrder.Quantity

		vwapParams := &vwapParameters{}
		err := json.Unmarshal([]byte(vwapParamsJson), vwapParams)
		if err != nil {
			om.CancelOrderWithErrorMsg(fmt.Sprintf("failed to parse parameters:%v", err))
		}

		if vwapParams.utcStartTimeSecs <= 0 || vwapParams.utcEndTimeSecs <= 0 {
			om.CancelOrderWithErrorMsg("invalid start or end time specified")
		}

		if vwapParams.utcStartTimeSecs >= vwapParams.utcEndTimeSecs {
			om.CancelOrderWithErrorMsg("start time must be before end time")
		}

		if model.IasD(vwapParams.buckets).GreaterThan(quantity) {
			om.CancelOrderWithErrorMsg("num Buckets must be less than or equal to the quantity")
		}

		buckets, err := getBucketsFromParamsString(vwapParamsJson, quantity, listing)
		if err != nil {
			om.CancelOrderWithErrorMsg(fmt.Sprintf("failed to get buckets from params:%v", err))
		}


		executeAsVwapStrategy(om, buckets, listing)
	}

	store, err := orderstore.NewKafkaStore(kafkaBrokers, id)
	if err != nil {
		panic(fmt.Errorf("failed to create order store: %v", err))
	}

	childOrderUpdates, err := executionvenue.GetChildOrders(id, kafkaBrokers, strategy.ChildUpdatesBufferSize)
	if err != nil {
		panic(err)
	}

	distributor := executionvenue.NewChildOrderUpdatesDistributor(childOrderUpdates)

	sm := strategy.NewStrategyManager(id, store, distributor, orderRouter, executeFn)

	api.RegisterExecutionVenueServer(s, sm)

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
