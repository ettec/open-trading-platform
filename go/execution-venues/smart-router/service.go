package main

import (
	"context"
	"fmt"
	"github.com/ettec/open-trading-platform/go/common"
	api "github.com/ettec/open-trading-platform/go/common/api/executionvenue"
	"github.com/ettec/open-trading-platform/go/common/k8s"
	"github.com/ettec/open-trading-platform/go/common/marketdata"
	"github.com/ettec/open-trading-platform/go/model"
	"github.com/ettech/open-trading-platform/go/smart-router/internal"
	"github.com/google/uuid"
	"k8s.io/client-go/kubernetes"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/ettec/open-trading-platform/go/common/bootstrap"

	"github.com/ettec/open-trading-platform/go/execution-venues/common/ordercache"
	"github.com/ettec/open-trading-platform/go/execution-venues/common/orderstore"

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
	ExecVenueMic           = "MIC"
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
	orderCache                   *ordercache.OrderCache
	orderRouter                  api.ExecutionVenueClient
	quoteDistributor             marketdata.QuoteDistributor
	getListingsFn                internal.GetListingsWithSameInstrument
	doneChan                     chan string
	orders                       sync.Map
	childOrderUpdatesDistributor *childOrderUpdatesDistributor
}

func (s *smartRouter) CreateAndRouteOrder(ctx context.Context, params *api.CreateAndRouteOrderParams) (*api.OrderId, error) {

	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	om, err := internal.NewOrderManager(id.String(), params, s.id, s.getListingsFn, s.doneChan, func(order model.Order) error {
		return s.orderCache.Store(&order)
	}, s.orderRouter,
		s.quoteDistributor.GetNewQuoteStream(), s.childOrderUpdatesDistributor.NewOrderStream(id.String(), 1000))

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

func main() {

	id := bootstrap.GetOptionalEnvVar(Id, "smart-router")
	maxConnectRetry := time.Duration(bootstrap.GetOptionalIntEnvVar(MaxConnectRetrySeconds, 60)) * time.Second
	external := bootstrap.GetOptionalBoolEnvVar(External, false)
	kafkaBrokersString := bootstrap.GetEnvVar(KafkaBrokersKey)
	execVenueMic := bootstrap.GetEnvVar(ExecVenueMic)

	s := grpc.NewServer()

	kafkaBrokers := strings.Split(kafkaBrokersString, ",")

	store, err := orderstore.NewKafkaStore(kafkaBrokers, execVenueMic)
	if err != nil {
		panic(fmt.Errorf("failed to create order store: %v", err))
	}

	orderCache, err := ordercache.NewOrderCache(store)
	if err != nil {
		log.Fatalf("failed to create order cache:%v", err)
	}

	sds, err := common.NewStaticDataSource(common.STATIC_DATA_SERVICE_ADDRESS)
	if err != nil {
		panic(err)
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

	mdsQuoteStream, err := marketdata.NewMdsQuoteStream(id, targetAddress, maxConnectRetry, 1000)
	qd := marketdata.NewQuoteDistributor(mdsQuoteStream, 100)

	if err != nil {
		panic(err)
	}

	orderRouter, err := getOrderRouter(clientSet, maxConnectRetry)
	if err != nil {
		panic(err)
	}

	_, childOrderUpdates, err := orderstore.GetChildOrders(id, kafkaBrokers)
	childUpdatesDistributor := newChildOrderUpdatesDistributor(childOrderUpdates)

	sr := &smartRouter{
		id:                           id,
		orderCache:                   orderCache,
		orderRouter:                  orderRouter,
		quoteDistributor:             qd,
		getListingsFn:                sds.GetListingsWithSameInstrument,
		doneChan:                     make(chan string, 100),
		orders:                       sync.Map{},
		childOrderUpdatesDistributor: childUpdatesDistributor,
	}

	// recovery behaviour

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

type childOrderStream struct {
	parentOrderId string
	orderChan     chan *model.Order
	distributor   *childOrderUpdatesDistributor
}

func newChildOrderStream(parentOrderId string, bufferSize int, d *childOrderUpdatesDistributor) *childOrderStream {
	stream := &childOrderStream{parentOrderId: parentOrderId, orderChan: make(chan *model.Order, bufferSize), distributor: d}
	d.openOrderChan <- parentIdAndChan{
		parentId:  parentOrderId,
		orderChan: stream.orderChan,
	}
	return stream
}

func (c *childOrderStream) GetStream() <-chan *model.Order {
	return c.orderChan
}

func (c *childOrderStream) Close() {
	c.distributor.closeOrderChan <- c.parentOrderId
}

type parentIdAndChan struct {
	parentId  string
	orderChan chan *model.Order
}

type childOrderUpdatesDistributor struct {
	openOrderChan  chan parentIdAndChan
	closeOrderChan chan string
}

func (d *childOrderUpdatesDistributor) NewOrderStream(parentOrderId string, bufferSize int) *childOrderStream {
	return newChildOrderStream(parentOrderId, bufferSize, d)
}

func newChildOrderUpdatesDistributor(updates <-chan orderstore.ChildOrder) *childOrderUpdatesDistributor {

	idToChan := map[string]chan *model.Order{}

	d := &childOrderUpdatesDistributor{
		openOrderChan:  make(chan parentIdAndChan),
		closeOrderChan: make(chan string),
	}

	go func() {

		for {
			select {
			case u := <-updates:
				if orderChan, ok := idToChan[u.ParentOrderId]; ok {
					select {
					case orderChan <- u.Child:
					default:
						log.Printf("slow consumer, closing child order update channel, parent order id %v", u.ParentOrderId)
						close(orderChan)
						delete(idToChan, u.ParentOrderId)
					}

				}
			case o := <-d.openOrderChan:
				idToChan[o.parentId] = o.orderChan
			case c := <-d.closeOrderChan:
				delete(idToChan, c)

			}
		}

	}()

	return d
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
