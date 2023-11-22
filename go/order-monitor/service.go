package main

import (
	"context"
	"errors"
	"fmt"
	common "github.com/ettec/otp-common"
	"github.com/ettec/otp-common/api"
	"github.com/ettec/otp-common/api/executionvenue"
	"github.com/ettec/otp-common/k8s"
	"github.com/ettec/otp-common/model"
	"github.com/ettec/otp-common/orderstore"
	"github.com/ettech/open-trading-platform/go/order-monitor/api/ordermonitor"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os/signal"
	"syscall"

	"os"
	"time"

	"github.com/ettec/otp-common/bootstrap"

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
	kafkaBrokers                  []string
	orderUpdatesBufSize           int
	ordersAfter                   time.Time
}

func (m *orderMonitor) CancelAllOrdersForOriginatorId(_ context.Context, params *ordermonitor.CancelAllOrdersForOriginatorIdParams) (*model.Empty, error) {
	m.cancelOrdersForOriginatorChan <- params
	return &model.Empty{}, nil
}

func main() {

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})))

	logAllOrderUpdates := bootstrap.GetOptionalBoolEnvVar("LOG_ORDER_UPDATES", true)
	maxConnectRetry := time.Duration(bootstrap.GetOptionalIntEnvVar("MAX_CONNECT_RETRY_SECONDS", 60)) * time.Second
	kafkaBrokersString := bootstrap.GetEnvVar("KAFKA_BROKERS")
	cancelTimeoutDuration := time.Duration(bootstrap.GetOptionalIntEnvVar("CANCEL_TIMEOUT_SECS", 5)) * time.Second
	orderUpdatesBufSize := bootstrap.GetOptionalIntEnvVar("INBOUND_ORDER_UPDATES_BUFFER_SIZE", 1000)

	cancelAllCmdChan := make(chan *ordermonitor.CancelAllOrdersForOriginatorIdParams)
	now := time.Now()
	ordersAfter := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)

	om := &orderMonitor{
		cancelOrdersForOriginatorChan: cancelAllCmdChan,
		kafkaBrokers:                  strings.Split(kafkaBrokersString, ","),
		orderUpdatesBufSize:           orderUpdatesBufSize,
		ordersAfter:                   ordersAfter,
	}

	http.Handle("/metrics", promhttp.Handler())
	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			slog.Error("failed to start metrics server", "error", err)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := om.startOrderStatusMonitoring(ctx, logAllOrderUpdates); err != nil {
		log.Panicf("failed to start order status monitoring: %v", err)
	}

	if err := om.startCancelAllHandler(ctx, maxConnectRetry, cancelTimeoutDuration); err != nil {
		log.Panicf("failed to start cancel all handler: %v", err)
	}

	s := grpc.NewServer()
	ordermonitor.RegisterOrderMonitorServer(s, om)
	reflection.Register(s)

	port := "50551"
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

	slog.Info("Starting Order Monitor Service", "port", port)

	if err := s.Serve(lis); err != nil {
		log.Panicf("error while serving : %v", err)
	}

}

func (m *orderMonitor) startCancelAllHandler(ctx context.Context, maxConnectRetry time.Duration, cancelTimeoutDuration time.Duration) error {
	clientSet := k8s.GetK8sClientSet(false)

	orderRouter, err := api.GetOrderRouter(clientSet, maxConnectRetry)
	if err != nil {
		return fmt.Errorf("failed to get order router: %w", err)
	}
	go func() {

		store, err := orderstore.NewKafkaStore(orderstore.DefaultReaderConfig(common.ORDERS_TOPIC, m.kafkaBrokers),
			orderstore.DefaultWriterConfig(common.ORDERS_TOPIC, m.kafkaBrokers), "")
		if err != nil {
			slog.Error("failed to create order store", "error", err)
			return
		}

		orders, updates, err := store.SubscribeToAllOrders(ctx, m.ordersAfter, m.orderUpdatesBufSize)
		if err != nil {
			slog.Error("failed to create order store", "error", err)
			return
		}

		nonTerminatedOrders := make(map[string]*model.Order)
		for _, order := range orders {
			if !order.IsTerminalState() {
				nonTerminatedOrders[order.GetId()] = order
			}
		}

		for {
			select {
			case <-ctx.Done():
				return
			case update := <-updates:
				slog.Info("got order update", "order", update)
				if update.IsTerminalState() {
					delete(nonTerminatedOrders, update.GetId())
				} else {
					nonTerminatedOrders[update.GetId()] = update
				}
			case cr := <-m.cancelOrdersForOriginatorChan:
				slog.Info("cancelling all orders for root originator", "originatorId", cr.OriginatorId)
				var cancelParams []*executionvenue.CancelOrderParams
				for _, order := range nonTerminatedOrders {
					if order.GetOriginatorId() == cr.OriginatorId {
						cancelParams = append(cancelParams, &executionvenue.CancelOrderParams{
							OrderId:   order.GetId(),
							ListingId: order.GetListingId(),
							OwnerId:   order.GetOwnerId(),
						})
					}
				}

				numOrdersToCancel := len(cancelParams)
				slog.Info("cancellable orders found for originator id", "numCancellableOrders", numOrdersToCancel,
					"originatorId", cr.OriginatorId)

				go func() {
					for _, params := range cancelParams {
						deadline, cancel := context.WithDeadline(ctx, time.Now().Add(cancelTimeoutDuration))
						_, err = orderRouter.CancelOrder(deadline, params)
						cancel()
						if err != nil {
							if !errors.Is(err, context.DeadlineExceeded) {
								slog.Error("Failed to cancel order", "orderId", params.GetOrderId(), "error", err)
							} else {
								slog.Error("Deadline exceed, failed to cancel order", "orderId", params.GetOrderId())
							}
						} else {
							slog.Info("Cancelled order", "orderId", params.GetOrderId())
						}
					}
				}()
			}
		}
	}()
	return err
}

func (m *orderMonitor) startOrderStatusMonitoring(ctx context.Context, logAllOrderUpdates bool) error {

	store, err := orderstore.NewKafkaStore(orderstore.DefaultReaderConfig(common.ORDERS_TOPIC, m.kafkaBrokers),
		orderstore.DefaultWriterConfig(common.ORDERS_TOPIC, m.kafkaBrokers), "")
	if err != nil {
		log.Panicf("failed to create order store: %v", err)
	}

	orders, updates, err := store.SubscribeToAllOrders(ctx, m.ordersAfter, m.orderUpdatesBufSize)
	if err != nil {
		return fmt.Errorf("failed to subscribe to all orders:%w", err)
	}

	// Setup initial order stats
	for _, order := range orders {
		totalOrders.Inc()
		gauge, err := getOrderStatusGauge(order)
		if err != nil {
			slog.Error("failed to get status gauge for order", "order", order, "error", err)
		}
		gauge.Inc()

		if logAllOrderUpdates {
			slog.Info("Initial state", "order", order)
		}
	}

	// Listening for order status counts
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case update := <-updates:
				if logAllOrderUpdates {
					slog.Info("Updated state", "order", update)
				}

				if order, exists := orders[update.Id]; exists {
					orders[update.Id] = update
					g, err := getOrderStatusGauge(order)
					if err != nil {
						slog.Error("failed to update status gauge", "error", err)
						continue
					}
					g.Dec()
				} else {
					totalOrders.Inc()
					orders[update.Id] = update
				}

				g, err := getOrderStatusGauge(update)
				if err != nil {
					slog.Error("failed to update status gauge", "error", err)
					continue
				}
				g.Inc()
			}
		}

	}()

	return nil
}

func getOrderStatusGauge(order *model.Order) (prometheus.Gauge, error) {

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
