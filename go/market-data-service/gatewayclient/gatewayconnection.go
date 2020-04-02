package gatewayclient

import (
	"github.com/ettec/open-trading-platform/go/common/marketdata"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/actor"
	"github.com/ettech/open-trading-platform/go/market-data-service/api/marketdatasource"

	"github.com/ettec/open-trading-platform/go/model"
	"google.golang.org/grpc"
	"log"
	"sync"
	"time"
)

type GatewayConnection struct {
	partyIdToConnection map[string]actor.ClientConnection
	quoteDistributor    marketdata.QuoteDistributor
	connMux             sync.Mutex
	maxSubscriptions    int
}

func NewGatewayConnection(id string, marketGatewayAddress string, maxReconnectInterval time.Duration,
	maxSubscriptions int) (*GatewayConnection, error) {

	mdcToDistributorChan := make(chan *model.ClobQuote, 1000)

	mdcFn := func(targetAddress string) (marketdatasource.MarketDataSourceClient, GrpcConnection, error) {
		conn, err := grpc.Dial(targetAddress, grpc.WithInsecure(), grpc.WithBackoffMaxDelay(maxReconnectInterval))
		if err != nil {
			return nil, nil, err
		}

		client := marketdatasource.NewMarketDataSourceClient(conn)
		return client, conn, nil
	}

	mdc, err := NewMarketDataGatewayClient(id, marketGatewayAddress, mdcToDistributorChan, mdcFn)

	if err != nil {
		return nil, err
	}

	qd := marketdata.NewQuoteDistributor(mdc.Subscribe, mdcToDistributorChan)
	gateway := &GatewayConnection{partyIdToConnection: make(map[string]actor.ClientConnection), quoteDistributor: qd,
		maxSubscriptions: maxSubscriptions}

	return gateway, nil
}

func (s *GatewayConnection) GetConnection(partyId string) (actor.ClientConnection, bool) {
	s.connMux.Lock()
	defer s.connMux.Unlock()

	con, ok := s.partyIdToConnection[partyId]
	return con, ok
}

func (s *GatewayConnection) AddConnection(subscriberId string, out chan<- *model.ClobQuote) actor.ClientConnection {
	s.connMux.Lock()
	defer s.connMux.Unlock()

	if conn, ok := s.partyIdToConnection[subscriberId]; ok {
		log.Printf("connection for client %v already exists, closing existing connection.", subscriberId)
		conn.Close()
		log.Print("connection closed:", subscriberId)
	}

	cc := actor.NewClientConnection(subscriberId, out, s.quoteDistributor, s.maxSubscriptions)

	s.partyIdToConnection[subscriberId] = cc

	return cc
}
