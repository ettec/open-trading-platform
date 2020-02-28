package main

import (
	"context"
	"fmt"
	"github.com/ettec/open-trading-platform/go/execution-venue/api"
	"github.com/ettec/open-trading-platform/go/model"
	"github.com/ettec/open-trading-platform/go/execution-venue/internal/ordercache"
	"github.com/ettec/open-trading-platform/go/execution-venue/internal/ordercache/orderstore"
	"github.com/ettec/open-trading-platform/go/execution-venue/internal/ordergateway/fixgateway"
	"github.com/ettec/open-trading-platform/go/execution-venue/internal/ordermanager"
	"github.com/quickfixgo/quickfix"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	UseLocalStoreKey   = "USE_LOCAL_STORE"
	LocalStorePathKey  = "LOCAL_STORE_PATH"
	KafkaOrderTopicKey = "KAFKA_ORDERS_TOPIC"
	KafkaBrokersKey    = "KAFKA_BROKERS"
)

type service struct {
	orderManager ordermanager.OrderManager
}

func NewService(om ordermanager.OrderManager) *service {
	service := service{orderManager: om}
	return &service
}



func (s *service) CreateAndRouteOrder(context context.Context, params *api.CreateAndRouteOrderParams) (*api.OrderId, error) {

	log.Printf("Received  order parameters-> %v", params)

	if params.GetQuantity() == nil {
		return nil, fmt.Errorf("quantity required on params:%v", params)
	}

	if params.GetPrice() == nil {
		return nil, fmt.Errorf("price required on params:%v", params)
	}

	if params.GetListing() == nil {
		return nil, fmt.Errorf("listing required on params:%v", params)
	}

	result, err := s.orderManager.CreateAndRouteOrder(params)
	if err != nil {
		log.Printf("error when creating and routing order:%v", err)
		return nil, err
	}

	log.Printf("created order id:%v", result.OrderId)

	return &api.OrderId{
		OrderId: result.OrderId,
	}, nil
}

func (s *service) CancelOrder(ctx context.Context, id *api.OrderId) (*model.Empty, error) {
	return &model.Empty{}, s.orderManager.CancelOrder(id)
}

func (s *service) Close() {
	if s.orderManager != nil {
		s.orderManager.Close()
	}
}

func main() {

	port := "50551"
	fmt.Println("Starting Execution Venue Service on port:" + port)
	lis, err := net.Listen("tcp", "0.0.0.0:"+port)

	if err != nil {
		log.Fatalf("Error while listening : %v", err)
	}

	s := grpc.NewServer()

	store, err := createOrderStore()
	if err != nil {
		panic(fmt.Errorf("failed to create order store: %v", err))
	}

	orderCache := ordercache.NewOrderCache(store)

	beginString := "FIXT.1.1"
	targetCompID := "EXEC"
	sendCompID := "BANZAI"
	sessionID := quickfix.SessionID{BeginString: string(beginString), TargetCompID: string(targetCompID), SenderCompID: string(sendCompID)}

	gateway := fixgateway.NewFixOrderGateway(sessionID)

	om := ordermanager.NewOrderManager(orderCache, gateway)

	fixServerCloseChan := make(chan struct{})
	err = createFixGateway(fixServerCloseChan, sessionID, om)
	if err != nil {
		panic(fmt.Errorf("failed to create fix gateway: %v", err))
	}

	defer func() { fixServerCloseChan <- struct{}{} }()

	service := NewService(om)
	defer service.Close()

	api.RegisterExecutionVenueServer(s, service)

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("error   while serving : %v", err)
	}

}

func createOrderStore() (orderstore.OrderStore, error) {
	var store orderstore.OrderStore

	useLocalStore := false
	useLocalStoreValue, exists := os.LookupEnv(UseLocalStoreKey)
	if exists {
		var err error
		useLocalStore, err = strconv.ParseBool(useLocalStoreValue)
		if err != nil {
			return nil, fmt.Errorf("failed to parse value for key %v: %w", UseLocalStoreKey, err)
		}
	}

	if useLocalStore {

		path, exists := os.LookupEnv(LocalStorePathKey)
		if !exists {
			log.Fatalf("must specify %v is using the local store", LocalStorePathKey)
		}

		fileStore, err := orderstore.NewFileStore(path)
		if err != nil {
			log.Fatalf("unable to create store: %v", err)
		}
		store = fileStore

	} else {
		ordersTopic, exists := os.LookupEnv(KafkaOrderTopicKey)
		if !exists {
			log.Fatalf("must specify %v for the kafka store", KafkaOrderTopicKey)
		}

		kafkaBrokers, exists := os.LookupEnv(KafkaBrokersKey)
		if !exists {
			log.Fatalf("must specify %v for the kafka store", KafkaBrokersKey)
		}

		store = orderstore.NewKafkaStore(ordersTopic, strings.Split(kafkaBrokers, ","))

	}

	return store, nil
}

func getFixConfig(sessionId quickfix.SessionID) string {

	allRequiredEnvVars := true
	fileLogPath, ok := os.LookupEnv("FIX_LOG_FILE_PATH")
	allRequiredEnvVars = allRequiredEnvVars && ok
	fileStorePath, ok := os.LookupEnv("FIX_FILE_STORE_PATH")
	allRequiredEnvVars = allRequiredEnvVars && ok
	fixPort, ok := os.LookupEnv("FIX_SOCKET_CONNECT_PORT")
	allRequiredEnvVars = allRequiredEnvVars && ok
	fixHost, ok := os.LookupEnv("FIX_SOCKET_CONNECT_HOST")
	allRequiredEnvVars = allRequiredEnvVars && ok

	template :=
		"[DEFAULT]\n" +
			"ConnectionType=initiator\n" +
			"ReconnectInterval=20\n" +
			"SenderCompID=" + sessionId.SenderCompID + "\n" +
			"FileStorePath=" + fileStorePath + "\n" +
			"FileLogPath=" + fileLogPath + "\n" +
			"\n" +
			"[SESSION]\n" +
			"BeginString=" + sessionId.BeginString + "\n" +
			"DefaultApplVerID=FIX.5.0SP2\n" +
			"TransportDataDictionary=./resources/FIXT11.xml\n" +
			"AppDataDictionary=./resources/FIX50SP2.xml\n" +
			"TargetCompID=" + sessionId.TargetCompID + "\n" +
			"StartTime=00:00:00\n" +
			"EndTime=00:00:00\n" +
			"HeartBtInt=20\n" +
			"SocketConnectPort=" + fixPort + "\n" +
			"SocketConnectHost=" + fixHost + "\n"

	return template
}

func createFixGateway(done chan struct{}, id quickfix.SessionID, handler fixgateway.OrderHandler) error {

	fixConfig := getFixConfig(id)

	app := fixgateway.NewFixHandler(id, handler)

	log.Printf("Creating fix engine with config: %v", fixConfig)

	appSettings, err := quickfix.ParseSettings(strings.NewReader(fixConfig))
	if err != nil {
		return fmt.Errorf("failed parse config: %v", err)
	}
	storeFactory := quickfix.NewFileStoreFactory(appSettings)
	logFactory, err := quickfix.NewFileLogFactory(appSettings)
	if err != nil {
		return fmt.Errorf("failed to create logFactory: %v", err)
	}
	initiator, err := quickfix.NewInitiator(app, storeFactory, appSettings, logFactory)
	if err != nil {
		return fmt.Errorf("failed to create initiator: %v", err)
	}

	go func() {
		err = initiator.Start()
		if err != nil {
			panic(fmt.Errorf("failed to start the fix engine: %v", err))
		}

		<-done

		//for condition == true { do something }
		defer initiator.Stop()
	}()

	return nil
}
