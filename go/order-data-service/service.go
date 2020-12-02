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
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

var logFlags = log.Ltime|log.Lshortfile
var errLog = log.New(os.Stderr,"", logFlags)

type service struct {
	orderSubscriptions sync.Map
	kafkaBrokers       []string
	ordersAfter        time.Time
	toClientBufferSize int
}

type orderAndWriteTime struct {
	order     *model.Order
	writeTime time.Time
}


type Source interface {
	ReadMessage(ctx context.Context) (key []byte, value []byte, writeTime time.Time, err error)
	Close() error
}

const maxInitialOrderConflationInterval = 500 * time.Millisecond

func (s *service) SubscribeToOrdersWithRootOriginatorId(request *api.SubscribeToOrdersWithRootOriginatorIdArgs, stream api.OrderDataService_SubscribeToOrdersWithRootOriginatorIdServer) error {
	username, appInstanceId, err := getMetaData(stream.Context())
	if err != nil {
		return err
	}
	after := &model.Timestamp{Seconds: s.ordersAfter.Unix()}

	log.Printf("subscribe to orders from app instance id:%v, user:%v, args:%v", appInstanceId, username, request)

	_, exists := s.orderSubscriptions.LoadOrStore(appInstanceId, appInstanceId)
	if !exists {

		defer s.orderSubscriptions.Delete(appInstanceId)

		out := make(chan orderAndWriteTime, 1000)

		
		go func() {
			source := NewKafkaMessageSource(orderstore.DefaultReaderConfig(common.ORDERS_TOPIC, s.kafkaBrokers))
			defer func() {
				if err := source.Close(); err != nil {
					errLog.Printf("error when closing kafka message source:%v", err)
				}
			}()
			streamOrderTopic(common.ORDERS_TOPIC, source, appInstanceId, out, after, request.RootOriginatorId)
			
		}()

		err := sendUpdates(out, stream.Send)
		if err != nil {
			return err
		}

	} else {
		return fmt.Errorf("subscription to orders already exists for application instance id %v", appInstanceId)
	}

	return nil

}

func sendUpdates(out chan orderAndWriteTime, send func(*model.Order) error) error {

	firstOrder := true
	startTime := time.Time{}
	conflatedOrders := map[string]*model.Order{}
	conflatingInitialOrderState := true

	ticker := time.NewTicker(maxInitialOrderConflationInterval)
	tickerChan := ticker.C
	var lastReceivedTime time.Time
	var lastReceivedOrder *orderAndWriteTime

	for {

		select {
		case oat, ok := <-out:
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

			err := send(oat.order)
			if err != nil {
				return fmt.Errorf("failed to send order, closing connection, error:%v", err)
			}
		case <-tickerChan:
			conflatingInitialOrderState = conflatingOrders(startTime, lastReceivedOrder, lastReceivedTime)
			if !conflatingInitialOrderState {
				ticker.Stop()
				err := sendConflatedOrders(send, conflatedOrders)
				if err != nil {
					return err
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

func (s *service) GetOrderHistory(ctx context.Context, args *api.GetOrderHistoryArgs) (*api.OrderHistory, error) {
	_, _, err := getMetaData(ctx)
	if err != nil {
		return nil, err
	}

	reader := NewKafkaMessageSource(orderstore.DefaultReaderConfig(common.ORDERS_TOPIC, s.kafkaBrokers))
	defer func() {
		if err := reader.Close(); err != nil {
			errLog.Printf("error when closing kafka message source:%v", err)
		}
	}()

	var updates []*api.OrderUpdate
	for {
		key, value, writeTime, err := reader.ReadMessage(ctx)
		if err != nil {
			return nil, err
		}

		orderId := string(key)
		if orderId == args.OrderId {
			order := &model.Order{}
			err = proto.Unmarshal(value, order)
			if err != nil {
				return nil, err
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

func streamOrderTopic(topic string, reader Source, appInstanceId string,
	out chan<- orderAndWriteTime, after *model.Timestamp, rootOriginatorId string) {

	for {
		_, value, writeTime, err := reader.ReadMessage(context.Background())

		if err != nil {
			logTopicReadError(appInstanceId, topic, err)
			return
		}

		order := &model.Order{}
		err = proto.Unmarshal(value, order)
		if err != nil {
			logTopicReadError(appInstanceId, topic, err)
			close(out)
			return
		}

		if order.Created != nil && order.Created.After(after) && order.RootOriginatorId == rootOriginatorId {
			out <- orderAndWriteTime{
				order:     order,
				writeTime: writeTime,
			}
		}
	}
}

func logTopicReadError(appInstanceId string, topic string, err error) {
	errLog.Printf("AppInstanceId: %v, Topic: %v, error occurred whilt attempting to stream message: %v", appInstanceId,
		topic, err)
}

func main() {

	log.SetOutput(os.Stdout)
	log.SetFlags(logFlags)

	toClientBufferSize := bootstrap.GetOptionalIntEnvVar("TO_CLIENT_BUFFER_SIZE", 1000)

	port := "50551"
	fmt.Println("Starting view service on port:" + port)
	listener, err := net.Listen("tcp", "0.0.0.0:"+port)

	defer func() {
		if err := listener.Close(); err != nil {
			errLog.Printf("error when closing listener:%v", err)
		}
	}()

	if err != nil {
		log.Panicf("Error while listening : %v", err)
	}

	s := grpc.NewServer()
	defer s.Stop()

	api.RegisterOrderDataServiceServer(s, newService(toClientBufferSize))

	reflection.Register(s)
	if err := s.Serve(listener); err != nil {
		log.Panicf("Error while serving : %v", err)
	}

}
