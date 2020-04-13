package main

import (
	"fmt"
	api "github.com/ettec/open-trading-platform/go/common/api/executionvenue"
	"github.com/ettec/open-trading-platform/go/common/bootstrap"
	"github.com/ettec/open-trading-platform/go/common/executionvenue"
	"github.com/ettec/open-trading-platform/go/common/topics"
	"github.com/ettec/open-trading-platform/go/execution-venue/internal/ordercache"
	"github.com/ettec/open-trading-platform/go/execution-venue/internal/ordercache/orderstore"
	"github.com/ettec/open-trading-platform/go/execution-venue/internal/ordergateway/fixgateway"
	"github.com/ettec/open-trading-platform/go/execution-venue/internal/ordermanager"
	"github.com/quickfixgo/quickfix"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
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

	store, err := orderstore.NewKafkaStore(topics.GetOrdersTopic(execVenueMic), strings.Split(kafkaBrokers, ","), execVenueMic)
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
	lis, err := net.Listen("tcp", "0.0.0.0:"+port)

	if err != nil {
		log.Fatalf("Error while listening : %v", err)
	}

	if err := s.Serve(lis); err != nil {
		log.Fatalf("error   while serving : %v", err)
	}

}
