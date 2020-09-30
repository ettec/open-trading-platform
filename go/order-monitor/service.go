package main

import (
	"context"
	"fmt"
	common "github.com/ettec/otp-common"
	"github.com/ettec/otp-common/api/executionvenue"
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
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"net"
	"strings"
)


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

type orderMonitor struct {
	cancelOrdersForOriginatorChan chan *ordermonitor.CancelAllOrdersForOriginatorIdParams
}

func (s *orderMonitor) CancelAllOrdersForOriginatorId(_ context.Context, params *ordermonitor.CancelAllOrdersForOriginatorIdParams) (*model.Empty, error) {
	s.cancelOrdersForOriginatorChan <- params
	return &model.Empty{}, nil
}

func main() {

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime|log.Lshortfile)

	logAllOrderUpdates := bootstrap.GetOptionalBoolEnvVar("LOG_ORDER_UPDATES", false)
	maxConnectRetry := time.Duration(bootstrap.GetOptionalIntEnvVar("MAX_CONNECT_RETRY_SECONDS", 60)) * time.Second
	external := bootstrap.GetOptionalBoolEnvVar("EXTERNAL", false)
	kafkaBrokersString := bootstrap.GetEnvVar("KAFKA_BROKERS")
	cancelTimeoutDuration := time.Duration(bootstrap.GetOptionalIntEnvVar("CANCEL_TIMEOUT_SECS", 5)) * time.Second
	orderUpdatesBufSize := bootstrap.GetOptionalIntEnvVar("INBOUND_ORDER_UPDATES_BUFFER_SIZE", 1000)


	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":8080", nil)

	kafkaBrokers := strings.Split(kafkaBrokersString, ",")

	now := time.Now()
	ordersAfter := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)

	store, err := orderstore.NewKafkaStore(orderstore.DefaultReaderConfig(common.ORDERS_TOPIC, kafkaBrokers),
		orderstore.DefaultWriterConfig(common.ORDERS_TOPIC, kafkaBrokers), "")

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

	updates := make(chan *model.Order, orderUpdatesBufSize)

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
			log.Print("failed to update status gauge:", err)
		}
		gauge.Inc()

		if logAllOrderUpdates {
			log.Printf("Order Update:%v\n\n", order)
		}
	}

	go func() {

		for {
			select {

			case cr := <-cancelChan:
				log.Print("cancelling all orders for root originator id:", cr.OriginatorId)
				var cancelParams []*executionvenue.CancelOrderParams
				for _, order := range orders {
					if order.GetOriginatorId() == cr.OriginatorId &&
						order.Status == model.OrderStatus_LIVE {

						cancelParams = append(cancelParams, &executionvenue.CancelOrderParams{
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

				if logAllOrderUpdates {
					log.Printf("Order Update:%v\n\n", u)
				}

				if order, exists := orders[u.Id]; exists {
					orders[u.Id] = u
					g, err := getStatusGauge(order)
					if err != nil {
						log.Print("failed to update status gauge:", err)
						continue
					}
					g.Dec()

					g, err = getStatusGauge(u)
					if err != nil {
						log.Print("failed to update status gauge:", err)
						continue
					}
					g.Inc()

				} else {
					totalOrders.Inc()
					orders[u.Id] = u
					g, err := getStatusGauge(u)
					if err != nil {
						log.Print("failed to update status gauge:", err)
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

func getOrderRouter(clientSet *kubernetes.Clientset, maxConnectRetrySecs time.Duration) (executionvenue.ExecutionVenueClient, error) {
	namespace := "default"
	list, err := clientSet.CoreV1().Services(namespace).List(v1.ListOptions{
		LabelSelector: "app=order-router",
	})

	if err != nil {
		panic(err)
	}

	var client executionvenue.ExecutionVenueClient

	for _, service := range list.Items {

		var podPort int32
		for _, port := range service.Spec.Ports {
			if port.Name == "executionvenue" {
				podPort = port.Port
			}
		}

		if podPort == 0 {
			log.Printf("ignoring order router service as it does not have a port named executionvenue, service: %v", service)
			continue
		}

		targetAddress := service.Name + ":" + strconv.Itoa(int(podPort))

		log.Printf("connecting to order router service %v at: %v", service.Name, targetAddress)

		conn, err := grpc.Dial(targetAddress, grpc.WithInsecure(), grpc.WithBackoffMaxDelay(maxConnectRetrySecs))

		if err != nil {
			panic(err)
		}

		client = executionvenue.NewExecutionVenueClient(conn)
		break
	}

	if client == nil {
		return nil, fmt.Errorf("failed to find order router")
	}

	return client, nil
}
