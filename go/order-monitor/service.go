package main

import (
	"context"
	"fmt"
	"github.com/ettec/open-trading-platform/go/common"
	api "github.com/ettec/open-trading-platform/go/common/api/executionvenue"
	"github.com/ettec/open-trading-platform/go/common/k8s"
	"github.com/ettec/open-trading-platform/go/common/orderstore"
	"github.com/ettec/open-trading-platform/go/common/staticdata"
	"github.com/ettec/open-trading-platform/go/model"
	"github.com/ettech/open-trading-platform/go/order-monitor/api/ordermonitor"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"k8s.io/client-go/kubernetes"
	"os"
	"strconv"
	"time"

	"github.com/ettec/open-trading-platform/go/common/bootstrap"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	logger "log"
	"net"
	"strings"
)

var log = logger.New(os.Stdout, "", logger.Ltime|logger.Lshortfile)
var errLog = logger.New(os.Stderr, "", logger.Ltime|logger.Lshortfile)

var totalOrders = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "total_orders",
	Help: "The total number of orders",
})

var liveOrders = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "live_orders",
	Help: "The number of live orders",
})

var cancelledOrders = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "cancelled_orders",
	Help: "The number of cancelled orders",
})

var filledOrders = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "filled_orders",
	Help: "The number of filled orders",
})

var noneStatusOrders = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "none_status_orders",
	Help: "The number of orders with none status",
})

var pendingLiveOrders = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "pending_live_orders",
	Help: "The number of pending live orders",
})

var pendingCancelOrders = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "pending_cancel_orders",
	Help: "The number of pending cancel orders",
})

const (
	KafkaBrokersKey        = "KAFKA_BROKERS"
	External               = "EXTERNAL"
	MaxConnectRetrySeconds = "MAX_CONNECT_RETRY_SECONDS"
	CancelTimeoutSecs      = "CANCEL_TIMEOUT_SECS"
)

type orderMonitor struct {
	cancelOrdersForOriginatorChan chan *ordermonitor.CancelAllOrdersForOriginatorIdParams
}

func (s *orderMonitor) CancelAllOrdersForOriginatorId(ctx context.Context, params *ordermonitor.CancelAllOrdersForOriginatorIdParams) (*model.Empty, error) {
	s.cancelOrdersForOriginatorChan <- params
	return &model.Empty{}, nil
}

func main() {

	maxConnectRetry := time.Duration(bootstrap.GetOptionalIntEnvVar(MaxConnectRetrySeconds, 60)) * time.Second
	external := bootstrap.GetOptionalBoolEnvVar(External, false)
	kafkaBrokersString := bootstrap.GetEnvVar(KafkaBrokersKey)
	cancelTimeoutDuration := time.Duration(bootstrap.GetOptionalIntEnvVar(CancelTimeoutSecs, 5)) * time.Second

	s := grpc.NewServer()

	kafkaBrokers := strings.Split(kafkaBrokersString, ",")

	store, err := orderstore.NewKafkaStore(kafkaBrokers, "")
	if err != nil {
		panic(fmt.Errorf("failed to create order store: %v", err))
	}

	sds, err := staticdata.NewStaticDataSource(common.STATIC_DATA_SERVICE_ADDRESS)
	if err != nil {
		log.Fatalf("failed to create static data source:%v", err)
	}

	clientSet := k8s.GetK8sClientSet(external)

	orderRouter, err := getOrderRouter(clientSet, maxConnectRetry)
	if err != nil {
		panic(err)
	}

	cancelChan := make(chan *ordermonitor.CancelAllOrdersForOriginatorIdParams)
	om := &orderMonitor{
		cancelOrdersForOriginatorChan: cancelChan,
	}

	listings := map[int32]*model.Listing{}
	updates := make(chan *model.Order, 1000)
	listingsChan := make(chan *model.Listing, 1000)
	orders, err := store.SubscribeToAllOrders(updates)

	for _, order := range orders {
		totalOrders.Inc()
		gauge, err := getStatusGauge(order)
		if err != nil {
			errLog.Print("failed to update status gauge:", err)
		}
		gauge.Inc()
	}

	go func() {
		for {
			select {

			case cr := <-cancelChan:
				log.Print("cancelling all orders for root originator id:", cr.RootOriginatorId)
				var cancelParams []*api.CancelOrderParams
				for _, order := range orders {
					if order.GetRootOriginatorId() == cr.RootOriginatorId &&
						order.Status == model.OrderStatus_LIVE {
						if listing, exists := listings[order.GetListingId()]; exists {
							cancelParams = append(cancelParams, &api.CancelOrderParams{
								OrderId: order.GetId(),
								Listing: listing,
							})
						}
					}
				}

				numOrdersToCancel := len(cancelParams)
				log.Printf("%v cancellable orders found for root originator id:%v", numOrdersToCancel, cr.RootOriginatorId)

				go func() {
					for idx, params := range cancelParams {
						deadline, cancel := context.WithDeadline(context.Background(), time.Now().Add(cancelTimeoutDuration))
						_, err := orderRouter.CancelOrder(deadline, params)
						cancel()
						if err != nil {
							if err != context.DeadlineExceeded {
								log.Printf("Failed to cancel order %v of %v, id: %v", idx+1, numOrdersToCancel, params.GetOrderId())
							} else {
								log.Printf("Deadline exceed, failed to cancel order %v of %v, id: %v", idx+1, numOrdersToCancel, params.GetOrderId())
							}
						} else {
							log.Printf("Cancelled order %v of %v, id: %v", idx+1, numOrdersToCancel, params.GetOrderId())
						}
					}
				}()

			case u := <-updates:
				if _, exists := listings[u.ListingId]; !exists {
					sds.GetListing(u.ListingId, listingsChan)
				}

				if order, exists := orders[u.Id]; exists {
					orders[u.Id] = u
					g, err := getStatusGauge(order)
					if err != nil {
						errLog.Print("failed to update status gauge:", err)
						continue
					}
					g.Dec()

					g, err = getStatusGauge(u)
					if err != nil {
						errLog.Print("failed to update status gauge:", err)
						continue
					}
					g.Inc()

				} else {
					totalOrders.Inc()
					orders[u.Id] = u
					g, err := getStatusGauge(u)
					if err != nil {
						errLog.Print("failed to update status gauge:", err)
						continue
					}
					g.Inc()
				}

			case l := <-listingsChan:
				listings[l.Id] = l

			}
		}

	}()

	ordermonitor.RegisterOrderMonitorServer(s, om)
	reflection.Register(s)

	port := "50551"
	lis, err := net.Listen("tcp", "0.0.0.0:"+port)

	if err != nil {
		log.Fatalf("Error while listening : %v", err)
	}

	if err := s.Serve(lis); err != nil {
		log.Fatalf("error   while serving : %v", err)
	}

	fmt.Println("Started Order Monitor Service on port:" + port)
}

func getStatusGauge(order *model.Order) (prometheus.Gauge, error) {

	if order.TargetStatus == model.OrderStatus_CANCELLED {
		return pendingCancelOrders, nil
	}

	if order.TargetStatus == model.OrderStatus_LIVE {
		return pendingLiveOrders, nil
	}

	if order.Status == model.OrderStatus_LIVE {
		return liveOrders, nil
	}

	if order.Status == model.OrderStatus_CANCELLED {
		return cancelledOrders, nil
	}

	if order.Status == model.OrderStatus_FILLED {
		return filledOrders, nil
	}

	if order.Status == model.OrderStatus_NONE {
		return noneStatusOrders, nil
	}

	return nil, fmt.Errorf("no status gauge for order with status: %v, tartget status: %v, order id:%v",
		order.GetStatus(), order.GetTargetStatus(), order.GetId())
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
