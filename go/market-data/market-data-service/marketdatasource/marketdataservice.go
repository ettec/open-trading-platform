package marketdatasource

import (
	"context"
	"fmt"
	"github.com/ettec/otp-common/loadbalancing"
	"github.com/ettec/otp-common/marketdata"
	"github.com/ettec/otp-common/model"
	"github.com/ettec/otp-common/staticdata"
	"log/slog"
	"slices"
	"sync"
	"time"
)

type getListingFn func(ctx context.Context, listingId int32, result chan<- staticdata.ListingResult)

type MarketDataGateway interface {
	GetAddress() string
	GetOrdinal() int
	GetMarketMic() string
}

//go:generate go run github.com/golang/mock/mockgen -destination mocks/gatewaystreamsource.go -package mocks github.com/ettech/open-trading-platform/go/market-data/market-data-service/marketdatasource GatewayStreamSource
type GatewayStreamSource interface {
	NewQuoteStreamFromMdSource(ctx context.Context, id string, targetAddress string, maxReconnectInterval time.Duration,
		quoteBufferSize int) (marketdata.QuoteStream, error)
}

type MarketDataService struct {
	ctx                       context.Context
	id                        string
	gatewayStreamSource       GatewayStreamSource
	micToSources              map[string]map[int]*MdsConnection
	getListing                getListingFn
	subscribers               map[string]chan *model.ClobQuote
	bufferSize                int
	retryConnectSeconds       int
	maxSubscriptionsPerClient int

	sourceMutex sync.Mutex

	subscriberIdToConn        map[string]*connection
	gatewayToQuoteDistributor map[MarketDataGateway]*marketdata.QuoteDistributor
}

func NewMarketDataService(ctx context.Context, id string,
	gatewayStreamSource GatewayStreamSource,
	getListing getListingFn,
	toClientBufferSize int, retryConnectSeconds int, maxSubscriptionsPerClient int) *MarketDataService {
	return &MarketDataService{
		ctx:                       ctx,
		id:                        id,
		gatewayStreamSource:       gatewayStreamSource,
		micToSources:              map[string]map[int]*MdsConnection{},
		getListing:                getListing,
		subscribers:               map[string]chan *model.ClobQuote{},
		bufferSize:                toClientBufferSize,
		retryConnectSeconds:       retryConnectSeconds,
		maxSubscriptionsPerClient: maxSubscriptionsPerClient,

		subscriberIdToConn:        map[string]*connection{},
		gatewayToQuoteDistributor: map[MarketDataGateway]*marketdata.QuoteDistributor{},
	}

}

func (f *MarketDataService) AddMarketDataGateway(gateway MarketDataGateway) error {
	f.sourceMutex.Lock()
	defer f.sourceMutex.Unlock()

	mdgQuoteStream, err := f.gatewayStreamSource.NewQuoteStreamFromMdSource(f.ctx, f.id, gateway.GetAddress(), time.Duration(f.retryConnectSeconds)*time.Second,
		f.maxSubscriptionsPerClient)
	if err != nil {
		return fmt.Errorf("failed to create connection to market data source at %v, error: %w", gateway.GetAddress(), err)
	}

	qd := marketdata.NewQuoteDistributor(f.ctx, mdgQuoteStream, f.bufferSize)
	f.gatewayToQuoteDistributor[gateway] = qd

	for _, conn := range f.subscriberIdToConn {
		conn.addGateway(gateway, qd)
	}

	return nil
}

func (f *MarketDataService) Connect(ctx context.Context, subscriberId string) marketdata.QuoteStream {
	f.sourceMutex.Lock()
	defer f.sourceMutex.Unlock()

	conn := newConnection(ctx, subscriberId, f.getListing, f.bufferSize)
	f.subscriberIdToConn[subscriberId] = conn

	for gateway, quoteDistributor := range f.gatewayToQuoteDistributor {
		conn.addGateway(gateway, quoteDistributor)
	}

	return conn
}

type connection struct {
	ctx                  context.Context
	log                  *slog.Logger
	cancel               context.CancelFunc
	subscriberId         string
	getListingFn         getListingFn
	gatewayToQuoteStream map[MarketDataGateway]marketdata.QuoteStream
	out                  chan *model.ClobQuote

	mutex sync.Mutex
}

func newConnection(parentCtx context.Context, subscriberId string, getListingFn getListingFn,
	bufferSize int) *connection {

	ctx, cancel := context.WithCancel(parentCtx)

	conn := &connection{ctx: ctx, cancel: cancel, subscriberId: subscriberId,
		getListingFn:         getListingFn,
		gatewayToQuoteStream: map[MarketDataGateway]marketdata.QuoteStream{},
		out:                  make(chan *model.ClobQuote, bufferSize),
		log:                  slog.With("subsriberId", subscriberId),
	}

	return conn
}

func (c *connection) addGateway(gateway MarketDataGateway, quoteDistributor *marketdata.QuoteDistributor) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	stream := quoteDistributor.NewQuoteStream()
	c.gatewayToQuoteStream[gateway] = stream

	go func() {
		defer stream.Close()
		for {
			select {
			case <-c.ctx.Done():
				return
			case quote, ok := <-stream.Chan():
				if !ok {
					return
				}
				c.out <- quote
			}
		}
	}()
}

func (c *connection) Subscribe(listingId int32) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.log.Info("subscription request", "listingId", listingId)

	listingChan := make(chan staticdata.ListingResult, 1)
	c.getListingFn(c.ctx, listingId, listingChan)
	listingResult := <-listingChan
	if listingResult.Err != nil {
		return fmt.Errorf("failed to get listing %v, error: %w", listingId, listingResult.Err)
	}

	mic := listingResult.Listing.Market.Mic
	var gatewaysForMic []MarketDataGateway
	for gateway, _ := range c.gatewayToQuoteStream {
		if gateway.GetMarketMic() == mic {
			gatewaysForMic = append(gatewaysForMic, gateway)
		}
	}

	slices.SortFunc(gatewaysForMic, func(i, j MarketDataGateway) int {
		return i.GetOrdinal() - j.GetOrdinal()
	})

	if len(gatewaysForMic) > 0 {
		numGateways := int32(len(gatewaysForMic))
		ordinal := loadbalancing.GetBalancingOrdinal(listingId, numGateways)
		gateway := gatewaysForMic[ordinal]
		stream := c.gatewayToQuoteStream[gateway]
		if err := stream.Subscribe(listingResult.Listing.Id); err != nil {
			return fmt.Errorf("failed to subscribe to market quote for subscriber %v, listing %v, error: %w", c.subscriberId, listingResult.Listing.Id, err)
		}
	} else {
		return fmt.Errorf("no market data gateway found for mic %v", mic)
	}

	return nil
}

func (c *connection) Chan() <-chan *model.ClobQuote {
	return c.out
}

func (c *connection) Close() {
	c.cancel()
}
