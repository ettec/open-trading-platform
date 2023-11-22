package main

import (
	"context"
	"errors"
	"fmt"
	api "github.com/ettec/open-trading-platform/go/order-data-service/api/orderdataservice"
	"github.com/ettec/otp-common"
	"github.com/ettec/otp-common/bootstrap"
	"github.com/ettec/otp-common/model"
	"github.com/ettec/otp-common/orderstore"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

type orderAndWriteTime struct {
	order     *model.Order
	writeTime time.Time
}

const maxInitialOrderConflationInterval = 500 * time.Millisecond

type service struct {
	orderSubscriptions sync.Map
	kafkaBrokers       []string
	ordersAfter        time.Time
	toClientBufferSize int
}

func (s *service) SubscribeToOrdersWithRootOriginatorId(request *api.SubscribeToOrdersWithRootOriginatorIdArgs, stream api.OrderDataService_SubscribeToOrdersWithRootOriginatorIdServer) error {
	username, appInstanceId, err := getMetaData(stream.Context())
	if err != nil {
		return fmt.Errorf("failed to get metadata, error:%w", err)
	}
	streamLog := slog.With("appInstanceId", appInstanceId, "topic", common.ORDERS_TOPIC)

	ordersAfter := &model.Timestamp{Seconds: s.ordersAfter.Unix()}

	_, exists := s.orderSubscriptions.LoadOrStore(appInstanceId, appInstanceId)
	if !exists {
		defer func() {
			s.orderSubscriptions.Delete(appInstanceId)
			slog.Info("unsubscribed from order updates")
		}()

		slog.Info("subscribing to order updates", "username", username, "request", request)
		orderUpdatesChan := s.getOrderUpdatesChan(stream.Context(), streamLog, request.RootOriginatorId, ordersAfter)
		if err = sendOrderUpdates(stream.Context(), orderUpdatesChan, stream.Send); err != nil {
			return fmt.Errorf("failed to send order updates:%w", err)
		}
	} else {
		slog.Warn("subscription already exists, ignoring subscription request")
		return fmt.Errorf("subscription to orders already exists for application instance id %s", appInstanceId)
	}

	return nil
}

func (s *service) GetOrderHistory(ctx context.Context, args *api.GetOrderHistoryArgs) (*api.OrderHistory, error) {
	_, _, err := getMetaData(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get metadata, error:%w", err)
	}

	reader := NewKafkaMessageSource(orderstore.DefaultReaderConfig(common.ORDERS_TOPIC, s.kafkaBrokers))
	defer func() {
		if err := reader.Close(); err != nil {
			slog.Error("error when closing kafka message source", "error", err)
		}
	}()

	var updates []*api.OrderUpdate
	for {
		key, value, writeTime, err := reader.ReadMessage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to read message:%w", err)
		}

		orderId := string(key)
		if orderId == args.OrderId {
			order := &model.Order{}
			err = proto.Unmarshal(value, order)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal order:%w", err)
			}
			updates = append(updates, &api.OrderUpdate{
				Order: order,
				Time:  model.NewTimeStamp(writeTime),
			})
			if order.Version >= args.ToVersion {
				return &api.OrderHistory{Updates: updates}, nil
			}
		}
	}
}

func sendOrderUpdates(ctx context.Context, in <-chan orderAndWriteTime, send func(*model.Order) error) error {

	firstOrder := true
	startTime := time.Time{}
	conflatedOrders := map[string]*model.Order{}
	conflatingInitialOrderState := true

	ticker := time.NewTicker(maxInitialOrderConflationInterval)
	defer ticker.Stop()
	tickerChan := ticker.C
	var lastReceivedTime time.Time
	var lastReceivedOrder *orderAndWriteTime

	for {

		select {
		case <-ctx.Done():
			return nil
		case oat, ok := <-in:
			if !ok {
				return errors.New("inbound channel closed")
			}

			if conflatingInitialOrderState {
				if firstOrder {
					firstOrder = false
					startTime = time.Now()
				}

				lastReceivedOrder = &oat
				lastReceivedTime = time.Now()

				conflatedOrders[oat.order.Id] = oat.order

				conflatingInitialOrderState = conflatingOrders(startTime, lastReceivedOrder, lastReceivedTime)

				if !conflatingInitialOrderState {
					err := sendConflatedOrders(send, conflatedOrders)
					if err != nil {
						return err
					}
				}

				continue
			}

			if err := send(oat.order); err != nil {
				return fmt.Errorf("failed to send order:%w", err)
			}
		case <-tickerChan:
			conflatingInitialOrderState = conflatingOrders(startTime, lastReceivedOrder, lastReceivedTime)
			if !conflatingInitialOrderState {
				ticker.Stop()
				if err := sendConflatedOrders(send, conflatedOrders); err != nil {
					return fmt.Errorf("failed to send conflated orders: %w", err)
				}
			}
		}
	}

}

func sendConflatedOrders(send func(*model.Order) error, conflatedOrders map[string]*model.Order) error {
	for id, order := range conflatedOrders {
		err := send(order)
		if err != nil {
			return fmt.Errorf("failed to send order, closing connection, error:%v", err)
		}
		delete(conflatedOrders, id)
	}

	return nil
}

func conflatingOrders(startTime time.Time, lastReceivedOrder *orderAndWriteTime, lastReceivedTime time.Time) bool {
	if lastReceivedOrder == nil {
		return true
	}
	now := time.Now()

	return lastReceivedOrder.writeTime.Before(startTime) && now.Sub(lastReceivedTime) < maxInitialOrderConflationInterval
}

func newService(toClientBufferSize int) *service {

	now := time.Now()
	ordersAfter := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	s := service{
		ordersAfter:        ordersAfter,
		toClientBufferSize: toClientBufferSize,
	}

	kafkaBrokers := bootstrap.GetEnvVar("KAFKA_BROKERS")
	s.kafkaBrokers = strings.Split(kafkaBrokers, ",")

	return &s
}

func getMetaData(ctx context.Context) (username string, appInstanceId string, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", "", fmt.Errorf("failed to read metadata from the context")
	}

	appInstanceIds := md.Get("app-instance-id")
	if len(appInstanceIds) != 1 {
		return "", "", fmt.Errorf("unable to retrieve app-instance-id from metadata")
	}
	appInstanceId = appInstanceIds[0]

	usernames := md.Get("user-name")
	if len(usernames) != 1 {
		return "", "", fmt.Errorf("unable to retrieve user-name from metadata")
	}
	username = usernames[0]

	return username, appInstanceId, nil
}

func (s *service) getOrderUpdatesChan(ctx context.Context, streamLog *slog.Logger, forRootOriginatorId string, after *model.Timestamp) <-chan orderAndWriteTime {

	out := make(chan orderAndWriteTime, s.toClientBufferSize)

	go func() {
		defer close(out)
		reader := NewKafkaMessageSource(orderstore.DefaultReaderConfig(common.ORDERS_TOPIC, s.kafkaBrokers))
		defer func() {
			if err := reader.Close(); err != nil {
				slog.Error("error when closing kafka message source", "error", err)
			}
		}()

		for {
			_, value, writeTime, err := reader.ReadMessage(ctx)
			if err != nil {
				streamLog.Error("failed to read message", "error", err)
				return
			}

			order := &model.Order{}
			err = proto.Unmarshal(value, order)
			if err != nil {
				streamLog.Error("failed to unmarshal order", "error", err)
				return
			}

			if order.Created != nil && order.Created.After(after) && order.RootOriginatorId == forRootOriginatorId {
				out <- orderAndWriteTime{
					order:     order,
					writeTime: writeTime,
				}
			}
		}
	}()

	return out
}

func main() {

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})))

	toClientBufferSize := bootstrap.GetOptionalIntEnvVar("TO_CLIENT_BUFFER_SIZE", 1000)

	port := "50551"
	slog.Info("Starting order data service", "port", port)
	listener, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		log.Panicf("Error while listening : %v", err)
	}

	defer func() {
		if err := listener.Close(); err != nil {
			slog.Error("error when closing listener", "error", err)
		}
	}()

	s := grpc.NewServer()
	api.RegisterOrderDataServiceServer(s, newService(toClientBufferSize))
	reflection.Register(s)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh,
		syscall.SIGKILL,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		<-sigCh
		s.GracefulStop()
	}()

	if err := s.Serve(listener); err != nil {
		log.Panicf("Error while serving : %v", err)
	}

}
