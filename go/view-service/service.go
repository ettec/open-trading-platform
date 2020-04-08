package main

import (
	"context"
	"fmt"
	bootstrap2 "github.com/ettec/open-trading-platform/go/common/bootstrap"
	"github.com/ettec/open-trading-platform/go/common/k8s"
	"github.com/ettec/open-trading-platform/go/common/topics"
	"github.com/ettec/open-trading-platform/go/model"
	api "github.com/ettec/open-trading-platform/go/view-service/api"
	"github.com/ettec/open-trading-platform/go/view-service/internal/messagesource"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	orderTopics        []string
}

func newService() *service {
	s := service{}

	kafkaBrokers, exists := os.LookupEnv(KafkaBrokersKey)
	if !exists {
		log.Fatalf("must specify %v for the kafka store", KafkaBrokersKey)
	}
	s.kafkaBrokers = strings.Split(kafkaBrokers, ",")

	clientSet := k8s.GetK8sClientSet(bootstrap2.GetOptionalBoolEnvVar(External, false))

	namespace := "default"
	list, err := clientSet.CoreV1().Services(namespace).List(v1.ListOptions{
		LabelSelector: "app=execution-venue",
	})

	if err != nil {
		panic(err)
	}

	log.Printf("found %v execution venues", len(list.Items))

	for _, service := range list.Items {
		const micLabel = "mic"
		if _, ok := service.Labels[micLabel]; !ok {
			errLog.Printf("ignoring execution venue as it does not have a mic label, service: %v", service)
			continue
		}

		mic := service.Labels[micLabel]

		topic := topics.GetOrdersTopic(mic)
		s.orderTopics = append(s.orderTopics, topic)
		log.Printf("added order topic:, %v", topic)
	}

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
		after = &model.Timestamp{Seconds: time.Now().Unix()}
	}

	if err != nil {
		return err
	}

	log.Printf("received order subscription request from app instance id:%v, user:%v", appInstanceId, username)

	_, exists := s.orderSubscriptions.LoadOrStore(appInstanceId, appInstanceId)
	if !exists {

		defer s.orderSubscriptions.Delete(appInstanceId)

		out := make(chan *model.Order, 1000)
		for _, topic := range s.orderTopics {
			source := messagesource.NewKafkaMessageSource(topic, s.kafkaBrokers)
			go streamOrderTopic(topic, source, appInstanceId, out, after)
		}

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

func streamOrderTopic(topic string, reader messagesource.Source, appInstanceId string,
	out chan<- *model.Order, after *model.Timestamp) {

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
