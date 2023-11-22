package main

import (
	"github.com/ettec/otp-common/api/executionvenue"
	"github.com/ettec/otp-common/bootstrap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"k8s.io/apimachinery/pkg/types"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var errLog = log.New(os.Stderr, "", log.Ltime|log.Lshortfile)

type execVenue struct {
	podId  types.UID
	client executionvenue.ExecutionVenueClient
	conn   *grpc.ClientConn
}

func main() {

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})))

	maxConnectRetrySecs := bootstrap.GetOptionalIntEnvVar("MAX_CONNECT_RETRY_SECONDS", 60)

	orderRouter, err := NewOrderRouter(maxConnectRetrySecs)
	if err != nil {
		log.Panicf("failed to create order router: %v", err)
	}

	port := "50581"
	slog.Info("Starting Order Router", "port", port)
	lis, err := net.Listen("tcp", "0.0.0.0:"+port)

	if err != nil {
		log.Panicf("Error while listening : %v", err)
	}

	s := grpc.NewServer()

	executionvenue.RegisterExecutionVenueServer(s, orderRouter)

	reflection.Register(s)

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
		log.Panicf("Error while serving : %v", err)
	}
}
