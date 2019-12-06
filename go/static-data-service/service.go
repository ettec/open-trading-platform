package main

import (
	"context"
	"fmt"
	"github.com/ettec/open-trading-platform/go/static-data-service/internal/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type service struct {

}


here - impl this

func(s *service)	GetListingsMatching(context.Context, *model.MatchParameters) (*model.Listings, error) {
}

func(s *service)	GetListing(context.Context, *model.ListingId) (*model.Listing, error) {

}

func(s *service)	GetListings(context.Context, *model.ListingIds) (*model.Listings, error) {

}



func main() {

	port := "50551"
	fmt.Println("Starting static data service on port:" + port)
	lis, err := net.Listen("tcp", "0.0.0.0:" + port)

	if err != nil {
		log.Fatalf("Error while listening : %v", err)
	}

	s := grpc.NewServer()
	model.RegisterStaticDataServiceServer(s, service{})

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error while serving : %v", err)
	}

}
