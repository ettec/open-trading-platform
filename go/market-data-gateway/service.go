package main

import (
	"context"
	"fmt"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/actor"
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
	fixsimConn          actor.MdServerConnection
	quoteDistributor    actor.QuoteDistributor
	actors              []actor.Actor
}



type connection struct {
	QuoteChan     chan *quote
	stream        api.MarketDataGateway_ConnectServer
	subscriptions map[int]bool
	closeChan     chan bool
}

func newService(address , name string) *service {

	var actors []actor.Actor
	qd := actor.NewQuoteDistributor()
	actors = append(actors, qd)
	quoteNormaliser := actor.NewClobQuoteNormaliser(qd)
	actors = append(actors, qd)
	fixsimConn := actor.NewMdServerConnection(address, name, quoteNormaliser)
	actors = append(actors, fixsimConn)


	for _, actor := range actors {
		actor.Start()
	}

	return &service{partyIdToConnection: make(map[string]connection), fixsimConn: fixsimConn, quoteDistributor: qd,
		actors: actors}
}

func (*service) Subscribe(c context.Context, r *marketdata.MarketDataRequest) (*empty.Empty, error) {

	return nil, nil
}

type clientConnection struct {
	subscriptions actor.SubscriptionHandler

}


func (s *service) Connect(request *api.ConnectRequest, stream api.MarketDataGateway_ConnectServer) error {




	//here is where we chain it all together when a new connection is received

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
	api.RegisterMarketDataGatewayServer(s, newService())

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error while serving : %v", err)
	}

}
