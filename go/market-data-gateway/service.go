package main

import (
	"context"
	"fmt"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/api"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/fix/marketdata"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type service struct {
	partyIdToConnection map[string]connection
}

func newService() *service {
	return &service{partyIdToConnection: make(map[string]connection)}
}

func (*service) Subscribe(c context.Context, r *marketdata.MarketDataRequest) (*empty.Empty, error) {

	return nil, nil
}

func (s *service) Connect(request *api.ConnectRequest, stream api.MarketDataGateway_ConnectServer) error {

	/*
		partyId := request.GetPartyId()

		log
			con, ok = s.partyIdToConnection[partyId]
		if ok {
			return fmt.Errorf("Connection for part")
		}*/

	return nil
}

func main() {

	port := "50551"
	fmt.Println("Starting Client Market Data Gateway on port:" + port)
	lis, err := net.Listen("tcp", "0.0.0.0:"+port)

	if err != nil {
		log.Fatalf("Error while listening : %v", err)
	}

	s := grpc.NewServer()
	api.RegisterMarketDataGatewayServer(s, &service{})

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error while serving : %v", err)
	}

}
