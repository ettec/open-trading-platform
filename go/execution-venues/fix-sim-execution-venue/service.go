package main

import (
	"fmt"
	api "github.com/ettec/otp-common/api/executionvenue"
	"github.com/ettec/otp-common/orderstore"
	"github.com/ettec/otp-evcommon/executionvenue"

	"github.com/ettec/otp-common/bootstrap"

	"github.com/ettec/otp-evcommon/ordercache"
	"github.com/ettec/otp-evcommon/ordermanager"

	"github.com/ettec/open-trading-platform/go/execution-venues/fix-sim-execution-venue/internal/fixgateway"

	"github.com/quickfixgo/quickfix"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"strings"
)

const (
	KafkaBrokersKey = "KAFKA_BROKERS"
	ExecVenueMic    = "MIC"
)

func main() {

	kafkaBrokers := bootstrap.GetEnvVar(KafkaBrokersKey)
	execVenueMic := bootstrap.GetEnvVar(ExecVenueMic)

	s := grpc.NewServer()

	store, err := orderstore.NewKafkaStore(strings.Split(kafkaBrokers, ","), execVenueMic)
	if err != nil {
		panic(fmt.Errorf("failed to create order store: %v", err))
	}

	orderCache, err := ordercache.NewOrderCache(store)
	if err != nil {
		log.Fatalf("failed to create order cache:%v", err)
	}

	beginString := "FIXT.1.1"
	targetCompID := "EXEC"
	sendCompID := "BANZAI"
	sessionID := quickfix.SessionID{BeginString: beginString, TargetCompID: targetCompID, SenderCompID: sendCompID}

	gateway := fixgateway.NewFixOrderGateway(sessionID)

	om := ordermanager.NewOrderManager(orderCache, gateway, execVenueMic)

	fixServerCloseChan := make(chan struct{})
	err = createFixGateway(fixServerCloseChan, sessionID, om)
	if err != nil {
		panic(fmt.Errorf("failed to create fix gateway: %v", err))
	}

	defer func() { fixServerCloseChan <- struct{}{} }()

	service := executionvenue.New(om)
	defer service.Close()

	api.RegisterExecutionVenueServer(s, service)

	reflection.Register(s)

	port := "50551"
	fmt.Println("Starting Execution Venue Service on port:" + port)
	fmt.Println("V3")
	lis, err := net.Listen("tcp", "0.0.0.0:"+port)

	if err != nil {
		log.Fatalf("Error while listening : %v", err)
	}

	if err := s.Serve(lis); err != nil {
		log.Fatalf("error   while serving : %v", err)
	}

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

	if tproot, exists := os.LookupEnv("TELEPRESENCE_ROOT"); exists {
		fileLogPath = tproot + fileLogPath
		fileStorePath = tproot + fileStorePath
	}

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
