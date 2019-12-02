package main

import (
	"fmt"
	"github.com/ettec/open-trading-platform/view-service/internal/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"context"
)

type server struct{}

func (*server) AddSubscription(context context.Context, subscription *cmds.Subscription) (*cmds.AddSubscriptionResponse, error) {

	//return nil, status.Errorf(codes.NotFound, "No subscriber found for id %v", subscription.SubscriberId)


	here implement this along with subscriber to kafka etc

	return &cmds.AddSubscriptionResponse{
		Message: "Subscription for " + subscription.ListingId + " added for subscriber " + subscription.SubscriberId,
	}, nil
}


func main() {

	port := "50551"
	fmt.Println("Starting Client Market Data Server on the port:" + port)
	lis, err := net.Listen("tcp", "0.0.0.0:" + port)

	if err != nil {
		log.Fatalf("Error while listening : %v", err)
	}

	s := grpc.NewServer()
	model.RegisterViewServiceServer(s, &server{})

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error while serving : %v", err)
	}

}