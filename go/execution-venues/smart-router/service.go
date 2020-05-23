package main

import (
	"context"
	"fmt"
	"github.com/ettec/open-trading-platform/go/common"
	api "github.com/ettec/open-trading-platform/go/common/api/executionvenue"
	"github.com/ettec/open-trading-platform/go/common/k8s"
	"github.com/ettec/open-trading-platform/go/common/marketdata"
	"github.com/ettec/open-trading-platform/go/common/marketdata/quotestream"
	"github.com/ettec/open-trading-platform/go/common/staticdata"
	"github.com/ettec/open-trading-platform/go/model"
	"github.com/ettech/open-trading-platform/go/smart-router/internal"
	"github.com/google/uuid"
	"k8s.io/client-go/kubernetes"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/ettec/open-trading-platform/go/common/bootstrap"

	"github.com/ettec/open-trading-platform/go/common/orderstore"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	logger "log"
	"net"
	"strings"
)

const (
	Id                     = "ID"
	KafkaBrokersKey        = "KAFKA_BROKERS"
	External               = "EXTERNAL"
	MaxConnectRetrySeconds = "MAX_CONNECT_RETRY_SECONDS"
)

var log = logger.New(os.Stdout, "", logger.Ltime|logger.Lshortfile)
var errLog = logger.New(os.Stderr, "", logger.Ltime|logger.Lshortfile)

type OrderManager interface {
	GetManagedOrderId() string
	Cancel()
}

type smartRouter struct {
	id                           string
	store                        orderstore.OrderStore
	orderRouter                  api.ExecutionVenueClient
	quoteDistributor             marketdata.QuoteDistributor
	getListingsFn                internal.GetListingsWithSameInstrument
	doneChan                     chan string
	orders                       sync.Map
	childOrderUpdatesDistributor ChildOrderUpdatesDistributor
}

type ChildOrderUpdatesDistributor interface {
	NewOrderStream(parentOrderId string, bufferSize int) internal.ChildOrderStream
	Start()
}

func (s *smartRouter) CreateAndRouteOrder(ctx context.Context, params *api.CreateAndRouteOrderParams) (*api.OrderId, error) {

	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	om, err := internal.NewOrderManagerFromParams(id.String(), params, s.id, s.getListingsFn, s.doneChan, s.store.Write, s.orderRouter,
		s.quoteDistributor.GetNewQuoteStream(), s.childOrderUpdatesDistributor.NewOrderStream(id.String(), ChildUpdatesBufferSize))

	if err != nil {
		return nil, err
	}

	s.orders.Store(om.GetManagedOrderId(), om)

	return &api.OrderId{
		OrderId: om.Id,
	}, nil
}

func (s *smartRouter) CancelOrder(ctx context.Context, params *api.CancelOrderParams) (*model.Empty, error) {

	if val, exists := s.orders.Load(params.OrderId); exists {
		om := val.(OrderManager)
		om.Cancel()
		return &model.Empty{}, nil
	} else {
		return nil, fmt.Errorf("no order found for id:%v", params.OrderId)
	}
}

const ChildUpdatesBufferSize = 1000

func main() {

	id := bootstrap.GetOptionalEnvVar(Id, "smart-router")
	maxConnectRetry := time.Duration(bootstrap.GetOptionalIntEnvVar(MaxConnectRetrySeconds, 60)) * time.Second
	external := bootstrap.GetOptionalBoolEnvVar(External, false)
	kafkaBrokersString := bootstrap.GetEnvVar(KafkaBrokersKey)

	s := grpc.NewServer()

	kafkaBrokers := strings.Split(kafkaBrokersString, ",")

	store, err := orderstore.NewKafkaStore(kafkaBrokers, common.SR_MIC)
	if err != nil {
		panic(fmt.Errorf("failed to create order store: %v", err))
	}

	sds, err := staticdata.NewStaticDataSource(common.STATIC_DATA_SERVICE_ADDRESS)
	if err != nil {
		log.Fatalf("failed to create static data source:%v", err)
	}

	clientSet := k8s.GetK8sClientSet(external)

	namespace := "default"
	xosrServiceLabelSelector := "app=market-data-source,mic=" + common.SR_MIC
	list, err := clientSet.CoreV1().Services(namespace).List(metav1.ListOptions{
		LabelSelector: xosrServiceLabelSelector,
	})

	if err != nil {
		panic(err)
	}

	if len(list.Items) != 1 {
		log.Panicf("no service found for selector: %v", xosrServiceLabelSelector)
	}

	service := list.Items[0]

	var podPort int32
	for _, port := range service.Spec.Ports {
		if port.Name == "api" {
			podPort = port.Port
		}
	}

	if podPort == 0 {
		log.Panic("aggregate quote source does not have an 'api' port")
	}

	targetAddress := service.Name + ":" + strconv.Itoa(int(podPort))

	mdsQuoteStream, err := quotestream.NewMdsQuoteStream(id, targetAddress, maxConnectRetry, 1000)
	if err != nil {
		panic(err)
	}

	qd := marketdata.NewQuoteDistributor(mdsQuoteStream, 100)

	orderRouter, err := getOrderRouter(clientSet, maxConnectRetry)
	if err != nil {
		panic(err)
	}

	childOrderUpdates, err := internal.GetChildOrders(id, kafkaBrokers, ChildUpdatesBufferSize)
	if err != nil {
		panic(err)
	}

	childUpdatesDistributor := internal.NewChildOrderUpdatesDistributor(childOrderUpdates)

	sr := &smartRouter{
		id:                           id,
		store:                        store,
		orderRouter:                  orderRouter,
		quoteDistributor:             qd,
		getListingsFn:                sds.GetListingsWithSameInstrument,
		doneChan:                     make(chan string, 100),
		orders:                       sync.Map{},
		childOrderUpdatesDistributor: childUpdatesDistributor,
	}

	go func() {
		id := <-sr.doneChan
		sr.orders.Delete(id)
		log.Printf("order %v is done, remove from orders cache", id)
	}()

	parentOrders, err := store.RecoverInitialCache()
	if err != nil {
		panic(err)
	}

	for _, order := range parentOrders {
		if !order.IsTerminalState() {
			om := internal.NewOrderManager(order, sr.store.Write, sr.id, sr.orderRouter, sr.getListingsFn,
				sr.quoteDistributor.GetNewQuoteStream(), sr.childOrderUpdatesDistributor.NewOrderStream(order.Id, 1000),
				sr.doneChan)
			sr.orders.Store(om.GetManagedOrderId(), om)
		}
	}

	childUpdatesDistributor.Start()

	api.RegisterExecutionVenueServer(s, sr)

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

func getOrderRouter(clientSet *kubernetes.Clientset, maxConnectRetrySecs time.Duration) (api.ExecutionVenueClient, error) {
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

		conn, err := grpc.Dial(targetAddress, grpc.WithInsecure(), grpc.WithBackoffMaxDelay(maxConnectRetrySecs))

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
