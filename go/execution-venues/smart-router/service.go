package main

import (
	"fmt"
	api "github.com/ettec/open-trading-platform/go/common/api/executionvenue"
	"github.com/ettec/open-trading-platform/go/common/k8s"
	"github.com/ettec/open-trading-platform/go/execution-venues/common/ordergateway"
	"k8s.io/client-go/kubernetes"
	"os"
	"strconv"
	"time"

	"github.com/ettec/open-trading-platform/go/execution-venues/common/executionvenue"

	"github.com/ettec/open-trading-platform/go/common/bootstrap"

	"github.com/ettec/open-trading-platform/go/common/topics"
	"github.com/ettec/open-trading-platform/go/execution-venues/common/ordercache"
	"github.com/ettec/open-trading-platform/go/execution-venues/common/orderstore"

	"github.com/ettec/open-trading-platform/go/execution-venues/common/ordermanager"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	logger "log"
	"net"
	"strings"
)

const (
	KafkaBrokersKey = "KAFKA_BROKERS"
	ExecVenueMic    = "MIC"
	External        = "EXTERNAL"
	MaxConnectRetrySeconds = "MAX_CONNECT_RETRY_SECONDS"
)

var log = logger.New(os.Stdout, "", logger.Ltime|logger.Lshortfile)
var errLog = logger.New(os.Stderr, "", logger.Ltime|logger.Lshortfile)

func main() {

	maxConnectRetrySecs := bootstrap.GetOptionalIntEnvVar(MaxConnectRetrySeconds, 60)
	external := bootstrap.GetOptionalBoolEnvVar(External, false)
	kafkaBrokers := bootstrap.GetEnvVar(KafkaBrokersKey)
	execVenueMic := bootstrap.GetEnvVar(ExecVenueMic)

	s := grpc.NewServer()

	store, err := orderstore.NewKafkaStore(topics.GetOrdersTopic(execVenueMic), strings.Split(kafkaBrokers, ","), execVenueMic)
	if err != nil {
		panic(fmt.Errorf("failed to create order store: %v", err))
	}

	orderCache, err := ordercache.NewOrderCache(store)
	if err != nil {
		log.Fatalf("failed to create order cache:%v", err)
	}

	clientSet := k8s.GetK8sClientSet(external)

	namespace := "default"
	list, err := clientSet.CoreV1().Services(namespace).List(v1.ListOptions{
		LabelSelector: "app=market-data-source",
		
	})

	if err != nil {
		panic(err)
	}

	log.Printf("found %v market data sources", len(list.Items))

	for _, service := range list.Items {
		const micLabel = "mic"
		if _, ok := service.Labels[micLabel]; !ok {
			errLog.Printf("ignoring market data source as it does not have a mic label, marketDataSource: %v", service)
			continue
		}

		mic := service.Labels[micLabel]

		var podPort int32
		for _, port := range service.Spec.Ports {
			if port.Name == "api" {
				podPort = port.Port
			}
		}

		if podPort == 0 {
			log.Printf("ignoring market data marketDataSource as it does not have a port named api, marketDataSource: %v", service)
			continue
		}

		targetAddress := service.Name + ":" + strconv.Itoa(int(podPort))

		micToMdsAddress[mic] = targetAddress

		log.Printf("found market data source for mic: %v, marketDataSource name: %v, target address: %v", mic, service.Name, targetAddress)
	}






	client, err := getOrderRouter(clientSet, maxConnectRetrySecs)
	if err != nil {
		panic(err)
	}

	om := ordermanager.NewOrderManager(orderCache, NewSmartRouterOrderGateway(client), execVenueMic)

	service := executionvenue.New(om)
	defer service.Close()

	api.RegisterExecutionVenueServer(s, service)

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

func getOrderRouter(clientSet *kubernetes.Clientset, maxConnectRetrySecs int) (api.ExecutionVenueClient, error) {
	namespace := "default"
	list, err := clientSet.CoreV1().Services(namespace).List(metav1.ListOptions{
		LabelSelector: "app=order-router",
	})

	if err != nil {
		panic(err)
	}

	var client api.ExecutionVenueClient

	for _, service := range list.Items {

		var podPort int32
		for _, port := range service.Spec.Ports {
			if port.Name == "api" {
				podPort = port.Port
			}
		}

		if podPort == 0 {
			log.Printf("ignoring order router service as it does not have a port named api, service: %v", service)
			continue
		}

		targetAddress := service.Name + ":" + strconv.Itoa(int(podPort))

		log.Printf("connecting to order router service %v at: %v", service.Name, targetAddress)

		conn, err := grpc.Dial(targetAddress, grpc.WithInsecure(), grpc.WithBackoffMaxDelay(time.Duration(maxConnectRetrySecs)*time.Second))

		if err != nil {
			panic(err)
		}

		client = api.NewExecutionVenueClient(conn)
		break
	}

	if client == nil {
		return nil, fmt.Errorf("failed to find order router")
	}

	return client, nil
}

func NewSmartRouterOrderGateway(orderRouter api.ExecutionVenueClient) ordergateway.OrderGateway{
	return nil
}
