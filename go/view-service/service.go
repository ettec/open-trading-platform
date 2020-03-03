package main

import (
	"context"
	"fmt"
	"github.com/ettec/open-trading-platform/go/model"
	api "github.com/ettec/open-trading-platform/go/view-service/api"
	"github.com/ettec/open-trading-platform/go/view-service/internal/messagesource"
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

const (
	KafkaOrderTopicKey = "KAFKA_ORDERS_TOPIC"
	KafkaBrokersKey    = "KAFKA_BROKERS"
)

type service struct {
	orderSubscriptions sync.Map
	orderTopic         string
	kafkaBrokers       []string
}

func newService() *service {
	s := service{}

	orderTopic, exists := os.LookupEnv(KafkaOrderTopicKey)
	if !exists {
		log.Fatalf("must specify %v for the kafka store", KafkaOrderTopicKey)
	}

	s.orderTopic = orderTopic

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

func (s *service) Subscribe(request *api.SubscribeToOrders, stream api.ViewService_SubscribeServer) error {

	username, appInstanceId, err := getMetaData(stream.Context())


	after := request.After
	if after == nil {
		after = &model.Timestamp{Seconds:time.Now().Unix()}
	}

	if err != nil {
		return err
	}

	log.Printf("received order subscription request from app instance id:%v, user:%v", appInstanceId, username)

	_, exists := s.orderSubscriptions.LoadOrStore(appInstanceId, appInstanceId)
	if !exists {
		source := messagesource.NewKafkaMessageSource(s.orderTopic, s.kafkaBrokers)
		streamTopic(s.orderTopic, source, appInstanceId, &s.orderSubscriptions, stream, after)
	} else {
		return fmt.Errorf("subscription to orders already exists for app instance id %v", appInstanceId)
	}

	return nil
}

func streamTopic(topic string, reader messagesource.Source, appInstanceId string, subscriptionsMap *sync.Map,
	stream api.ViewService_SubscribeServer, after *model.Timestamp) {

	defer subscriptionsMap.Delete(appInstanceId)

	defer reader.Close()

	for {
		msg, err := reader.ReadMessage(context.Background())

		if err != nil {
			logTopicReadError(appInstanceId, topic, err)
			return
		}

		order := model.Order{}
		err = proto.Unmarshal(msg, &order)
		if err != nil {
			logTopicReadError(appInstanceId, topic, err)
			return
		}

		if order.Created != nil && order.Created.After(after) {
			err = stream.Send(&order)
			if err != nil {
				logTopicReadError(appInstanceId, topic, err)
				return
			}
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
