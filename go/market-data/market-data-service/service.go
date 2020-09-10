package main

import (
	"context"
	"fmt"
	api "github.com/ettec/otp-common/api/marketdataservice"
	"github.com/ettec/otp-common/bootstrap"
	"github.com/ettec/otp-common/loadbalancing"
	"github.com/ettec/otp-common/model"
	"github.com/ettec/otp-common/staticdata"
	"github.com/ettech/open-trading-platform/go/market-data/market-data-service/marketdatasource"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	logger "log"
	"net"
	"net/http"
	"os"
	"time"
)

var connections = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "mds_active_connections",
	Help: "The number of active connections to the mds",
})

var quotesSent = promauto.NewCounter(prometheus.CounterOpts{
	Name: "mds_quotes_sent",
	Help: "The number of quotes sent across all clients",
})

type service struct {
	micToSources map[string]map[int]*marketdatasource.MdsConnection
	getListingFn func(listingId int32, result chan<- *model.Listing)
}

func (s *service) Subscribe(_ context.Context, r *api.MdsSubscribeRequest) (*model.Empty, error) {

	log.Printf("Subscribe request received for subscriber id: %v, listing id:%v, retrieving listing....", r.SubscriberId, r.ListingId)
	listingChan := make(chan *model.Listing)
	s.getListingFn(r.ListingId, listingChan)
	listing := <-listingChan
	log.Printf("listing %v received", listing.Id)

	mic := listing.Market.Mic
	if sources, ok := s.micToSources[mic]; ok {
		numGateways := int32(len(sources))
		ordinal := loadbalancing.GetBalancingOrdinal(listing, numGateways)

		if source, ok := sources[ordinal]; ok {
			if conn, ok := source.GetConnection(r.SubscriberId); ok {

				if err := conn.Subscribe(listing.Id); err != nil {
					return nil, err
				}

				return &model.Empty{}, nil
			} else {
				return nil, fmt.Errorf("failed  to subscribe, no connection exists for subscriber " + r.SubscriberId)
			}
		} else {
			return nil, fmt.Errorf("no market source exists for stateful set ordinal %v and mic %v", ordinal, mic)
		}

	} else {
		return nil, fmt.Errorf("no market data source exists for mic %v", mic)
	}

}

func (s *service) Connect(request *api.MdsConnectRequest, stream api.MarketDataService_ConnectServer) error {

	subscriberId := request.GetSubscriberId()

	log.Println("connect request received for subscriber: ", subscriberId)

	out := make(chan *model.ClobQuote, 100)

	for mic, gateways := range s.micToSources {

		for _, gateway := range gateways {
			gateway.AddConnection(subscriberId, out)
		}
		log.Printf("connected subscriber %v to %v market data sources for mic %v", subscriberId, len(gateways), mic)
	}

	connections.Inc()

	for mdUpdate := range out {

		if err := stream.Send(mdUpdate); err != nil {
			log.Printf("error on connection for subscriber %v, closing connection, error:%v", subscriberId, err)
			break
		}

		quotesSent.Inc()
	}

	connections.Dec()

	return nil
}

var maxSubscriptions = 10000
var log = logger.New(os.Stdout, "", logger.Ltime|logger.Lshortfile)
var errLog = logger.New(os.Stderr, "", logger.Ltime|logger.Lshortfile)

func main() {

	id := bootstrap.GetEnvVar("MDS_ID")
	connectRetrySecs := bootstrap.GetOptionalIntEnvVar("CONNECT_RETRY_SECONDS", 60)

	http.Handle("/metrics", promhttp.Handler())
	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			errLog.Printf("failed to listen on metrics server port:%v", err)
		}
	}()

	sds, err := staticdata.NewStaticDataSource(false)
	if err != nil {
		log.Fatalf("failed to get static data source:%v", err)
	}

	mdService := service{micToSources: map[string]map[int]*marketdatasource.MdsConnection{}, getListingFn: sds.GetListing}

	micToBalancingPods, err := loadbalancing.GetMicToStatefulPodAddresses("market-data-gateway")

	if err != nil {
		log.Panicf("failed to get market data gateway pod balancingPods: %v", err)
	}

	for mic, balancingPods := range micToBalancingPods {
		for _, balancingPod := range balancingPods {

			client, err := marketdatasource.NewMdsConnection(id, balancingPod.TargetAddress, time.Duration(connectRetrySecs)*time.Second,
				maxSubscriptions)
			if err != nil {
				errLog.Printf("failed to create connection to market data source at %v, error: %v", balancingPod, err)
				continue
			}

			if _, ok := mdService.micToSources[mic]; !ok {
				mdService.micToSources[mic] = map[int]*marketdatasource.MdsConnection{}
			}

			mdService.micToSources[mic][balancingPod.Ordinal] = client
			log.Printf("added market data source for mic: %v,  target address: %v, stateful set ordinal %v", mic, balancingPod, balancingPod.Ordinal)
		}
	}

	port := "50551"
	fmt.Println("starting market data service on port:" + port)
	lis, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		log.Fatalf("Error while listening : %v", err)
	}

	s := grpc.NewServer()

	api.RegisterMarketDataServiceServer(s, &mdService)

	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error while serving : %v", err)
	}
}
