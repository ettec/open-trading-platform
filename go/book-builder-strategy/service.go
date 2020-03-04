package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ettec/open-trading-platform/go/book-builder-strategy/api"
	"github.com/ettec/open-trading-platform/go/book-builder-strategy/depth"
	"github.com/ettec/open-trading-platform/go/book-builder-strategy/orderentryapi"
	"github.com/ettec/open-trading-platform/go/common"
	"github.com/ettec/open-trading-platform/go/common/bootstrap"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/actor"
	mdgapi "github.com/ettec/open-trading-platform/go/market-data-gateway/api"
	"github.com/ettec/open-trading-platform/go/market-data-service/gatewayclient"
	"github.com/ettec/open-trading-platform/go/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"io/ioutil"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

const (
	ServiceId                = "SERVICE_ID"
	MarketDataGatewayAddress = "MARKET_DATA_GATEWAY_ADDRESS"
	OrderEntryServiceAddress = "ORDER_ENTRY_SERVICE_ADDRESS"
	StaticDataServiceAddress = "STATIC_DATA_SERVICE_ADDRESS"
	ConnectRetrySeconds      = "CONNECT_RETRY_SECONDS"
)

func main() {

	port := "50551"
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

	s := grpc.NewServer()
	bbs, err := newService(serviceId, mdGatewayAddr, orderEntryAddr, staticDataServiceAddr,
		time.Duration(connectRetrySecs)*time.Second)

	if err != nil {
		log.Panicf("failed to create book builder strategy service:%v", err)
	}

	api.RegisterBookBuilderStrategyServer(s, bbs)

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error while serving : %v", err)
	}

	body,err := ioutil.ReadFile("./resources/depth.json")
	if err != nil {
		log.Fatalf("failed to read depths data:%v", err)
	}

	depths := depth.Depths{}
	err = json.Unmarshal(body, &depths)
	if err != nil {
		log.Fatalf("failed to unmarshall depths data:%v", err)
	}

	here - setup up books from data

}

type service struct {
	id                string
	quoteDistributor  actor.QuoteDistributor
	orderEntryService orderentryapi.OrderEntryServiceClient
	listingSource     common.ListingSource
	books             map[int32]*book
	booksMx           sync.Mutex
}

func newService(id string, mdGatewayAddr string, orderEntryAddr string, staticDataServiceAddr string,
	maxReconnectInterval time.Duration) (*service, error) {

	ls, err :=common.NewListingSource(staticDataServiceAddr)
	if err != nil {
		return nil, err
	}

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

	return &service{id: id, quoteDistributor: qd, orderEntryService: oec, listingSource:ls}, nil
}

func (s *service) BuildBookForListing(c context.Context, p *api.BuildBookForListingParams) (*model.Empty, error) {

	s.booksMx.Lock()

	if s.books[p.ListingId] == nil {
		s.books[p.ListingId] = newBook(p.ListingId, s.listingSource, s.quoteDistributor)
	} else {
		return nil, fmt.Errorf("book already exists for listing id:%v", p.ListingId)
	}

	return &model.Empty{}, nil
}

type book struct {
	log *log.Logger
}

func newBook(listingId int32, ls common.ListingSource, distributor actor.QuoteDistributor) *book {

	b := &book{
		log : log.New(os.Stdout, fmt.Sprintf(" book: %v ", listingId), log.Ltime),
	}

	b.log.Println("fetching listing")

	lc := make( chan *model.Listing, 1 )
	ls.GetListing(listingId, lc )

	listing := <-lc

	quotesIn := make( chan *model.ClobQuote, 1000 )

	distributor.AddOutQuoteChan(quotesIn)
	distributor.Subscribe(listingId, quotesIn)




	go func() {

		select {
			case q := <- quotesIn:



		}


	}()




	return b
}
