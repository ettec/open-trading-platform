package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ettec/open-trading-platform/go/book-builder-strategy/api"
	"github.com/ettec/open-trading-platform/go/book-builder-strategy/depth"
	"github.com/ettec/open-trading-platform/go/book-builder-strategy/orderentryapi"
	"github.com/ettec/open-trading-platform/go/book-builder-strategy/strategy"
	"github.com/ettec/open-trading-platform/go/common"
	"github.com/ettec/open-trading-platform/go/common/bootstrap"
	services "github.com/ettec/open-trading-platform/go/common/services"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/actor"
	mdgapi "github.com/ettec/open-trading-platform/go/market-data-gateway/api"
	"github.com/ettec/open-trading-platform/go/market-data-service/gatewayclient"
	"github.com/ettec/open-trading-platform/go/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"io/ioutil"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

const (
	ServiceId                = "SERVICE_ID"
	MarketDataGatewayAddress = "MARKET_DATA_GATEWAY_ADDRESS"
	OrderEntryServiceAddress = "ORDER_ENTRY_SERVICE_ADDRESS"
	StaticDataServiceAddress = "STATIC_DATA_SERVICE_ADDRESS"
	ConnectRetrySeconds      = "CONNECT_RETRY_SECONDS"
	BookScanIntervalMillis   = "BOOK_SCAN_INTERVAL_MILLIS"
	TradeProbability         = "TRADE_PROBABILITY"
	Variation                = "VARIATION"
	MinQtyPercent            = "MIN_QTY_PERCENT"
	SymbolsToRun             = "SYMBOLS_TO_RUN"
)

func main() {

	port := "50571"
	fmt.Println("Starting Market Data Service on port:" + port)
	lis, err := net.Listen("tcp", "0.0.0.0:"+port)

	if err != nil {
		log.Fatalf("Error while listening : %v", err)
	}

	serviceId := bootstrap.GetEnvVar(ServiceId)
	mdGatewayAddr := bootstrap.GetEnvVar(MarketDataGatewayAddress)
	orderEntryAddr := bootstrap.GetEnvVar(OrderEntryServiceAddress)
	staticDataServiceAddr := bootstrap.GetEnvVar(StaticDataServiceAddress)
	connectRetrySecs := bootstrap.GetOptionalIntEnvVar(ConnectRetrySeconds, 60)
	bookScanInterval := time.Duration(bootstrap.GetOptionalIntEnvVar(BookScanIntervalMillis, 1000)) * time.Millisecond
	tradeProbability := bootstrap.GetOptionalFloatEnvVar(TradeProbability, 0.1)
	variation := bootstrap.GetOptionalFloatEnvVar(Variation, 0.005)
	minQty := bootstrap.GetOptionalFloatEnvVar(MinQtyPercent, 0.9)
	symbolsToRunArg := bootstrap.GetOptionalEnvVar(SymbolsToRun, "*")



	ls, err := common.NewListingSource(staticDataServiceAddr)
	if err != nil {
		log.Panicf("failed to create listing source service:%v", err)
	}

	s := grpc.NewServer()
	bbs, err := newService(serviceId, mdGatewayAddr, orderEntryAddr, ls,
		time.Duration(connectRetrySecs)*time.Second)

	if err != nil {
		log.Panicf("failed to create book builder strategy service:%v", err)
	}



	body, err := ioutil.ReadFile("./resources/depth.json")
	if err != nil {
		log.Fatalf("failed to read depths data:%v", err)
	}

	depths := depth.Depths{}
	err = json.Unmarshal(body, &depths)
	if err != nil {
		log.Fatalf("failed to unmarshall depths data:%v", err)
	}

	symToDepths := map[string]depth.Depth{}

	for _, depth := range depths {
		symToDepths[depth.Symbol] = depth
	}

	var symbolsToRun []string
	if symbolsToRunArg == "*" {
		for k,_ := range symToDepths {
			symbolsToRun = append(symbolsToRun, k)
		}
	} else {
		symbolsToRun = strings.Split(symbolsToRunArg, ",")

	}

	listingChan := make(chan *model.Listing, 1)
	for _, sym := range symbolsToRun {

		ls.GetListingMatching(&services.MatchParameters{SymbolMatch: sym}, listingChan)
		listing := <-listingChan
		if listing != nil {
			book, err := strategy.NewBookBuilder(listing, bbs.quoteDistributor, symToDepths[sym], bbs.orderEntryService,
				bookScanInterval,  tradeProbability, variation, minQty)
			if err != nil {
				log.Printf("failed to start strategy for listing:%v, error:%v", listing, err)
			} else {
				book.Start()
			}


		}

	}

	api.RegisterBookBuilderStrategyServer(s, bbs)

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error while serving : %v", err)

	}

}

type service struct {
	id                string
	quoteDistributor  actor.QuoteDistributor
	orderEntryService orderentryapi.OrderEntryServiceClient
	listingSource     common.ListingSource
	books             map[int32]*strategy.BookBuilder
	booksMx           sync.Mutex
}

func newService(id string, mdGatewayAddr string, orderEntryAddr string, ls common.ListingSource,
	maxReconnectInterval time.Duration) (*service, error) {

	mdcFn := func(targetAddress string) (mdgapi.MarketDataGatewayClient, gatewayclient.GrpcConnection, error) {
		conn, err := grpc.Dial(targetAddress, grpc.WithInsecure(), grpc.WithBackoffMaxDelay(maxReconnectInterval))
		if err != nil {
			return nil, nil, err
		}

		client := mdgapi.NewMarketDataGatewayClient(conn)
		return client, conn, nil
	}

	mdcToDistributorChan := make(chan *model.ClobQuote, 1000)

	mdc, err := gatewayclient.NewMarketDataGatewayClient(id, mdGatewayAddr, mdcToDistributorChan, mdcFn)
	if err != nil {
		return nil, err
	}

	qd := actor.NewQuoteDistributor(mdc.Subscribe, mdcToDistributorChan)

	conn, err := grpc.Dial(orderEntryAddr, grpc.WithInsecure(), grpc.WithBackoffMaxDelay(maxReconnectInterval))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to order entry end point:%w", err)
	}

	oec := orderentryapi.NewOrderEntryServiceClient(conn)

	return &service{id: id, quoteDistributor: qd, orderEntryService: oec, listingSource: ls}, nil
}

func (s *service) Start(c context.Context, p *api.ListingId) (*model.Empty, error) {
	s.booksMx.Lock()
	defer s.booksMx.Unlock()

	if book, exists := s.books[p.ListingId]; exists {
		err :=  book.Start()
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("no book is setup for for listing id:%v", p.ListingId)
	}

	return &model.Empty{}, nil

}

func (s *service) Stop(c context.Context, p *api.ListingId) (*model.Empty, error) {

	s.booksMx.Lock()
	defer s.booksMx.Unlock()

	if book, exists := s.books[p.ListingId]; exists {
		err :=  book.Stop()
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("no book is setup for for listing id:%v", p.ListingId)
	}

	return &model.Empty{}, nil
}
