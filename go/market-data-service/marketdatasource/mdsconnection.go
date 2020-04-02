package marketdatasource

import (
	"github.com/ettec/open-trading-platform/go/common/marketdata"
	"github.com/ettec/open-trading-platform/go/model"
	"github.com/ettech/open-trading-platform/go/market-data-service/api/marketdatasource"
	"google.golang.org/grpc"
	"log"
	"sync"
	"time"
)

type MdsConnection struct {
	partyIdToConnection map[string]marketdata.ConflatedQuoteConnection
	quoteDistributor    marketdata.QuoteDistributor






	connMux          sync.Mutex
	maxSubscriptions int
}

func NewGatewayConnection(id string, marketGatewayAddress string, maxReconnectInterval time.Duration,
	maxSubscriptions int) (*MdsConnection, error) {

	mdcToDistributorChan := make(chan *model.ClobQuote, 1000)

	mdcFn := func(targetAddress string) (marketdatasource.MarketDataSourceClient, GrpcConnection, error) {
		conn, err := grpc.Dial(targetAddress, grpc.WithInsecure(), grpc.WithBackoffMaxDelay(maxReconnectInterval))
		if err != nil {
			return nil, nil, err
		}

		client := marketdatasource.NewMarketDataSourceClient(conn)
		return client, conn, nil
	}

	mdc, err := NewMarketDataSourceClient(id, marketGatewayAddress, mdcToDistributorChan, mdcFn)

	if err != nil {
		return nil, err
	}

	qd := marketdata.NewQuoteDistributor(mdc.Subscribe, mdcToDistributorChan)
	gateway := &MdsConnection{partyIdToConnection: make(map[string]marketdata.ConflatedQuoteConnection), quoteDistributor: qd,
		maxSubscriptions: maxSubscriptions}

	return gateway, nil
}

func (s *MdsConnection) GetConnection(partyId string) (marketdata.ConflatedQuoteConnection, bool) {
	s.connMux.Lock()
	defer s.connMux.Unlock()

	con, ok := s.partyIdToConnection[partyId]
	return con, ok
}

func (s *MdsConnection) AddConnection(subscriberId string, out chan<- *model.ClobQuote) marketdata.ConflatedQuoteConnection {
	s.connMux.Lock()
	defer s.connMux.Unlock()

	if conn, ok := s.partyIdToConnection[subscriberId]; ok {
		log.Printf("connection for client %v already exists, closing existing connection.", subscriberId)
		conn.Close()
		log.Print("connection closed:", subscriberId)
	}

	cc := marketdata.NewConflatedQuoteConnection(subscriberId, out, s.quoteDistributor, s.maxSubscriptions)

	s.partyIdToConnection[subscriberId] = cc

	return cc
}
