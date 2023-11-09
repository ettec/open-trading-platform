package marketdatasource

import (
	"context"
	"fmt"
	"github.com/ettec/otp-common/marketdata"
	"log"
	"sync"
	"time"
)

type MdsConnection struct {
	partyIdToConnection map[string]marketdata.QuoteStream
	quoteDistributor    *marketdata.QuoteDistributor

	connMux          sync.Mutex
	maxSubscriptions int
}

func NewMdsConnection(ctx context.Context, id string, marketGatewayAddress string, maxReconnectInterval time.Duration,
	maxSubscriptions int) (*MdsConnection, error) {

	stream, err := marketdata.NewQuoteStreamFromMdSource(ctx, id, marketGatewayAddress, maxReconnectInterval, 1000)
	if err != nil {
		return nil, fmt.Errorf("failed to create quote stream: %w", err)
	}

	qd := marketdata.NewQuoteDistributor(ctx, stream, 1000)
	gateway := &MdsConnection{partyIdToConnection: make(map[string]marketdata.QuoteStream), quoteDistributor: qd,
		maxSubscriptions: maxSubscriptions}

	return gateway, nil
}

func (s *MdsConnection) GetConnection(partyId string) (marketdata.QuoteStream, bool) {
	s.connMux.Lock()
	defer s.connMux.Unlock()

	con, ok := s.partyIdToConnection[partyId]
	return con, ok
}

func (s *MdsConnection) AddConnection(subscriberId string) marketdata.QuoteStream {
	s.connMux.Lock()
	defer s.connMux.Unlock()

	if conn, ok := s.partyIdToConnection[subscriberId]; ok {
		log.Printf("connection for client %v already exists, closing existing connection.", subscriberId)
		conn.Close()
		log.Print("connection closed:", subscriberId)
	}

	quoteStream := marketdata.NewConflatedQuoteStream(subscriberId, s.quoteDistributor.NewQuoteStream(),
		s.maxSubscriptions)

	s.partyIdToConnection[subscriberId] = quoteStream

	return quoteStream
}
