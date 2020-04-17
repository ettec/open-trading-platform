package main

import (
	"fmt"
	api "github.com/ettec/open-trading-platform/go/common/api/executionvenue"
	"github.com/ettec/open-trading-platform/go/common/bootstrap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"k8s.io/apimachinery/pkg/types"
	logger "log"
	"net"
	"os"
)

var log = logger.New(os.Stdout, "", logger.Ltime|logger.Lshortfile)
var errLog = logger.New(os.Stderr, "", logger.Ltime|logger.Lshortfile)

type execVenue struct {
	podId  types.UID
	client api.ExecutionVenueClient
	conn   *grpc.ClientConn
}

const (
	MaxConnectRetrySeconds = "MAX_CONNECT_RETRY_SECONDS"
	External               = "EXTERNAL"
)

func main() {

	maxConnectRetrySecs := bootstrap.GetOptionalIntEnvVar(MaxConnectRetrySeconds, 60)
	external := bootstrap.GetOptionalBoolEnvVar(External, false)

	orderRouter := New(external, maxConnectRetrySecs)

	port := "50581"
	fmt.Println("Starting Order Router on port:" + port)
	lis, err := net.Listen("tcp", "0.0.0.0:"+port)

	if err != nil {
		log.Fatalf("Error while listening : %v", err)
	}

	s := grpc.NewServer()

	api.RegisterExecutionVenueServer(s, orderRouter)

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error while serving : %v", err)

	}

}
