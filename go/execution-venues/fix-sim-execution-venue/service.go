package main

import (
	"context"
	"fmt"
	"github.com/ettec/open-trading-platform/go/execution-venues/fix-sim-execution-venue/internal/executionvenue"
	common "github.com/ettec/otp-common"
	api "github.com/ettec/otp-common/api/executionvenue"
	"github.com/ettec/otp-common/ordermanagement"
	"github.com/ettec/otp-common/orderstore"
	"github.com/ettec/otp-common/staticdata"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/ettec/otp-common/bootstrap"

	"github.com/ettec/open-trading-platform/go/execution-venues/fix-sim-execution-venue/internal/fixgateway"

	"github.com/quickfixgo/quickfix"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"strings"
)

func main() {

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})))

	kafkaBrokers := bootstrap.GetEnvVar("KAFKA_BROKERS")
	id := bootstrap.GetEnvVar("ID")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sds, err := staticdata.NewStaticDataSource(ctx)
	if err != nil {
		log.Panicf("failed to create static data source:%v", err)
	}

	brokers := strings.Split(kafkaBrokers, ",")
	store, err := orderstore.NewKafkaStore(orderstore.DefaultReaderConfig(common.ORDERS_TOPIC, brokers),
		orderstore.DefaultWriterConfig(common.ORDERS_TOPIC, brokers), id)

	if err != nil {
		log.Panicf("failed to create order store: %v", err)
	}

	orderCache, err := ordermanagement.NewOwnerOrderCache(ctx, id, store)

	if err != nil {
		log.Panicf("failed to create order cache:%v", err)
	}

	beginString := "FIXT.1.1"
	targetCompID := "EXEC"
	sendCompID := id
	sessionID := quickfix.SessionID{BeginString: beginString, TargetCompID: targetCompID, SenderCompID: sendCompID}

	gateway := fixgateway.NewFixOrderGateway(sessionID)

	om := executionvenue.NewOrderManager(ctx, orderCache, gateway, sds.GetListing,
		bootstrap.GetOptionalIntEnvVar("ORDER_MANAGER_CMD_BUFFER_SIZE", 100))

	closeFixGatewayFn, err := createFixGateway(sessionID, om)
	if err != nil {
		log.Panicf("failed to create fix gateway: %v", err)
	}
	defer closeFixGatewayFn()

	service := executionvenue.New(om)

	s := grpc.NewServer()
	api.RegisterExecutionVenueServer(s, service)

	reflection.Register(s)

	port := "50551"
	slog.Info("Starting Execution Venue Service", "port", port)

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

	if err := s.Serve(lis); err != nil {
		log.Panicf("error   while serving : %v", err)
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

func createFixGateway(id quickfix.SessionID, handler fixgateway.OrderHandler) (close func(), err error) {

	fixConfig := getFixConfig(id)

	app := fixgateway.NewFixHandler(id, handler)

	slog.Info("Creating fix engine", "config", fixConfig)

	appSettings, err := quickfix.ParseSettings(strings.NewReader(fixConfig))
	if err != nil {
		return nil, fmt.Errorf("failed parse config: %w", err)
	}
	storeFactory := quickfix.NewFileStoreFactory(appSettings)
	logFactory, err := quickfix.NewFileLogFactory(appSettings)
	if err != nil {
		return nil, fmt.Errorf("failed to create logFactory: %w", err)
	}
	initiator, err := quickfix.NewInitiator(app, storeFactory, appSettings, logFactory)
	if err != nil {
		return nil, fmt.Errorf("failed to create initiator: %w", err)
	}

	if err = initiator.Start(); err != nil {
		return nil, fmt.Errorf("failed to start the fix engine: %w", err)
	}

	return initiator.Stop, nil
}
