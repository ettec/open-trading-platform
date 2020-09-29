package main

import (
	"fmt"
	"github.com/ettec/otp-common/api/executionvenue"
	"github.com/ettec/otp-common/bootstrap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"k8s.io/apimachinery/pkg/types"
	"log"
	"net"
	"os"
)


var errLog = log.New(os.Stderr, "", log.Ltime|log.Lshortfile)

type execVenue struct {
	podId  types.UID
	client executionvenue.ExecutionVenueClient
	conn   *grpc.ClientConn
}


func main() {

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime|log.Lshortfile)

	maxConnectRetrySecs := bootstrap.GetOptionalIntEnvVar("MAX_CONNECT_RETRY_SECONDS", 60)

	orderRouter := New(maxConnectRetrySecs)

	port := "50581"
	fmt.Println("Starting Order Router on port:" + port)
	lis, err := net.Listen("tcp", "0.0.0.0:"+port)

	if err != nil {
		log.Fatalf("Error while listening : %v", err)
	}

	s := grpc.NewServer()

	executionvenue.RegisterExecutionVenueServer(s, orderRouter)

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error while serving : %v", err)

	}

}
