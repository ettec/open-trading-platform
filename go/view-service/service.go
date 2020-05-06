package main

import (
	"context"
	"fmt"
	"github.com/ettec/open-trading-platform/go/common"
	"github.com/ettec/open-trading-platform/go/model"
	api "github.com/ettec/open-trading-platform/go/view-service/api/viewservice"

	"github.com/ettec/open-trading-platform/go/view-service/internal/messagesource"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	logger "log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	KafkaBrokersKey = "KAFKA_BROKERS"
	External        = "EXTERNAL"
)

var log = logger.New(os.Stdout, "", logger.Ltime|logger.Lshortfile)
var errLog = logger.New(os.Stderr, "", logger.Ltime|logger.Lshortfile)

type service struct {
	orderSubscriptions sync.Map
	kafkaBrokers       []string
}

func (s *service) SubscribeToOrdersWithRootOriginatorId(request *api.SubscribeToOrdersWithRootOriginatorIdArgs, stream api.ViewService_SubscribeToOrdersWithRootOriginatorIdServer) error {
	username, appInstanceId, err := getMetaData(stream.Context())
	if err != nil {
		return err
	}

	after := request.After
	if after == nil {
		after = &model.Timestamp{Seconds: time.Now().Unix()}
	}

	log.Printf("received order subscription request from app instance id:%v, user:%v", appInstanceId, username)

	_, exists := s.orderSubscriptions.LoadOrStore(appInstanceId, appInstanceId)
	if !exists {

		defer s.orderSubscriptions.Delete(appInstanceId)

		out := make(chan *model.Order, 1000)

		source := messagesource.NewKafkaMessageSource(common.ORDERS_TOPIC, s.kafkaBrokers)
		go streamOrderTopic(common.ORDERS_TOPIC, source, appInstanceId, out, after, request.RootOriginatorId)

		for order := range out {
			err = stream.Send(order)
			if err != nil {
				errLog.Printf("closing connection as failed to send order to %v, error:%v", appInstanceId, err)
				close(out)
				return nil
			}
		}

	} else {
		return fmt.Errorf("subscription to orders already exists for app instance id %v", appInstanceId)
	}

	return nil

}

func (s *service) GetOrderHistory(ctx context.Context, args *api.GetOrderHistoryArgs) (*api.Orders, error) {
	_, _, err := getMetaData(ctx)
	if err != nil {
		return nil, err
	}

	reader := messagesource.NewKafkaMessageSource(common.ORDERS_TOPIC, s.kafkaBrokers)
	defer reader.Close()

	var orders []*model.Order
	for {
		key, value, err := reader.ReadMessage(ctx)
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
			orders = append(orders, order)
			if order.Version == args.ToVersion {
				return &api.Orders{Orders: orders}, nil
			}
		}
	}
}

func newService() *service {
	s := service{}

	kafkaBrokers, exists := os.LookupEnv(KafkaBrokersKey)
	if !exists {
		log.Fatalf("must specify %v for the kafka store", KafkaBrokersKey)
	}
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

func streamOrderTopic(topic string, reader messagesource.Source, appInstanceId string,
	out chan<- *model.Order, after *model.Timestamp, rootOriginatorId string) {

	defer reader.Close()

	for {
		_, value, err := reader.ReadMessage(context.Background())

		if err != nil {
			logTopicReadError(appInstanceId, topic, err)
			return
		}

		order := model.Order{}
		err = proto.Unmarshal(value, &order)
		if err != nil {
			logTopicReadError(appInstanceId, topic, err)
			return
		}

		if order.Created != nil && order.Created.After(after) && order.RootOriginatorId == rootOriginatorId {
			out <- &order
		}
	}
}

func logTopicReadError(appInstanceId string, topic string, err error) {
	log.Printf("AppInstanceId: %v, Topic: %v, error occurred whilt attempting to stream message: %v", appInstanceId,
		topic, err)
}

func main() {

	port := "50551"
	fmt.Println("Starting view service on port:" + port)
	lis, err := net.Listen("tcp", "0.0.0.0:"+port)

	if err != nil {
		log.Fatalf("Error while listening : %v", err)
	}

	s := grpc.NewServer()
	api.RegisterViewServiceServer(s, newService())

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error while serving : %v", err)
	}

}
