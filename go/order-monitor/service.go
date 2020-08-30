package main

import (
	"context"
	"fmt"
	api "github.com/ettec/otp-common/api/executionvenue"
	"github.com/ettec/otp-common/k8s"
	"github.com/ettec/otp-common/model"
	"github.com/ettec/otp-common/orderstore"
	"github.com/ettech/open-trading-platform/go/order-monitor/api/ordermonitor"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"

	"k8s.io/client-go/kubernetes"
	"os"
	"strconv"
	"time"

	"github.com/ettec/otp-common/bootstrap"

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

	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":8080", nil)

	kafkaBrokers := strings.Split(kafkaBrokersString, ",")

	now := time.Now()
	ordersAfter := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)

	store, err := orderstore.NewKafkaStore(kafkaBrokers, "")
	if err != nil {
		panic(fmt.Errorf("failed to create order store: %v", err))
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

	updates := make(chan *model.Order, 1000)

	orders, err := store.SubscribeToAllOrders(updates, ordersAfter)
	if err != nil {
		log.Panic("failed to subscribe to orders: ", err)
	}

	listingsToFetch := map[int32]bool{}
	for _, order := range orders {
		listingsToFetch[order.ListingId] = true

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
				log.Print("cancelling all orders for root originator id:", cr.OriginatorId)
				var cancelParams []*api.CancelOrderParams
				for _, order := range orders {
					if order.GetOriginatorId() == cr.OriginatorId &&
						order.Status == model.OrderStatus_LIVE {

						cancelParams = append(cancelParams, &api.CancelOrderParams{
							OrderId:   order.GetId(),
							ListingId: order.GetListingId(),
							OwnerId:   order.GetOwnerId(),
						})

					}
				}

				numOrdersToCancel := len(cancelParams)
				log.Printf("%v cancellable orders found for originator id:%v", numOrdersToCancel, cr.OriginatorId)

				go func() {
					for idx, params := range cancelParams {
						deadline, cancel := context.WithDeadline(context.Background(), time.Now().Add(cancelTimeoutDuration))
						_, err := orderRouter.CancelOrder(deadline, params)
						cancel()
						if err != nil {
							if err != context.DeadlineExceeded {
								log.Printf("Failed to cancel order %v of %v, id: %v, error:%v", idx+1, numOrdersToCancel, params.GetOrderId(), err)
							} else {
								log.Printf("Deadline exceed, failed to cancel order %v of %v, id: %v", idx+1, numOrdersToCancel, params.GetOrderId())
							}
						} else {
							log.Printf("Cancelled order %v of %v, id: %v", idx+1, numOrdersToCancel, params.GetOrderId())
						}
					}
				}()

			case u := <-updates:

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

			}
		}

	}()

	s := grpc.NewServer()
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
