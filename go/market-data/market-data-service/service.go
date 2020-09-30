package main

import (
	"context"
	"fmt"
	api "github.com/ettec/otp-common/api/marketdataservice"

	"github.com/ettec/otp-common/bootstrap"
	"github.com/ettec/otp-common/k8s"
	"github.com/ettec/otp-common/loadbalancing"
	"github.com/ettec/otp-common/model"
	"github.com/ettec/otp-common/staticdata"
	"github.com/ettech/open-trading-platform/go/market-data/market-data-service/marketdatasource"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	v12 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

var connections = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "mds_active_connections",
	Help: "The number of active subscribers to the mds",
})

var quotesSent = promauto.NewCounter(prometheus.CounterOpts{
	Name: "mds_quotes_sent",
	Help: "The number of quotes sent across all clients",
})


type service struct {
	micToSources map[string]map[int]*marketdatasource.MdsConnection
	getListingFn func(listingId int32, result chan<- *model.Listing)
	sourceMutex  sync.Mutex
	subscribers  map[string] chan *model.ClobQuote
	toClientBufferSize int
}

func (s *service) Subscribe(_ context.Context, r *api.MdsSubscribeRequest) (*model.Empty, error) {

	log.Printf("Subscribe request received for subscriber id: %v, listing id:%v, retrieving listing....", r.SubscriberId, r.ListingId)
	listingChan := make(chan *model.Listing)
	s.getListingFn(r.ListingId, listingChan)
	listing := <-listingChan
	log.Printf("listing %v received", listing.Id)

	mic := listing.Market.Mic
	sources := s.getSourcesForMic(mic)
	if len(sources) > 0 {
		numGateways := int32(len(sources))
		ordinal := loadbalancing.GetBalancingOrdinal(r.ListingId, numGateways)

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

func( s *service) getSourcesForMic(mic string) map[int]*marketdatasource.MdsConnection {
	s.sourceMutex.Lock()
	defer s.sourceMutex.Unlock()

	result := map[int]*marketdatasource.MdsConnection{}
	for k,v := range s.micToSources[mic] {
		result[k] = v
	}

	return result
}

func (s *service) addSubscriber( subscriberId string, out chan *model.ClobQuote ) {
	s.sourceMutex.Lock()
	s.sourceMutex.Unlock()

	s.subscribers[subscriberId] = out
	for mic, gateways := range s.micToSources {

		for _, gateway := range gateways {
			gateway.AddConnection(subscriberId, out)
		}
		log.Printf("connected subscriber %v to %v market data sources for mic %v", subscriberId, len(gateways), mic)
	}

}


func (s *service) removeSubscriber( subscriberId string ) {
	s.sourceMutex.Lock()
	s.sourceMutex.Unlock()

	delete(s.subscribers, subscriberId)
}


func (s *service) Connect(request *api.MdsConnectRequest, stream api.MarketDataService_ConnectServer) error {

	subscriberId := request.GetSubscriberId()

	log.Println("connect request received for subscriber: ", subscriberId)

	out := make(chan *model.ClobQuote, s.toClientBufferSize )

	s.addSubscriber(subscriberId, out)


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


func (s *service) addMarketDataConnection(bsp *loadbalancing.BalancingStatefulPod, connection *marketdatasource.MdsConnection) {
	s.sourceMutex.Lock()
	defer s.sourceMutex.Unlock()

	if _, ok := s.micToSources[bsp.Mic]; !ok {
		s.micToSources[bsp.Mic] = map[int]*marketdatasource.MdsConnection{}
	}

	if _, ok := s.micToSources[bsp.Mic][bsp.Ordinal]; !ok {

		s.micToSources[bsp.Mic][bsp.Ordinal] = connection
		log.Printf("added market data source for mic: %v,  target address: %v, stateful set ordinal %v", bsp.Mic, bsp, bsp.Ordinal)

		for subscriberId, conn := range s.subscribers {
			connection.AddConnection(subscriberId, conn)
			log.Printf("add subscriber %v to market data source", subscriberId)
		}
	}

}


var errLog = log.New(os.Stderr, "", log.Flags())

func main() {

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime|log.Lshortfile)

	id := bootstrap.GetEnvVar("MDS_ID")
	connectRetrySecs := bootstrap.GetOptionalIntEnvVar("CONNECT_RETRY_SECONDS", 60)
	maxSubscriptions := bootstrap.GetOptionalIntEnvVar("MAX_SUBSCRIPTIONS", 10000)
	toClientBufferSize := bootstrap.GetOptionalIntEnvVar("TO_CLIENT_BUFFER_SIZE", 1000)

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

	mdService := service{map[string]map[int]*marketdatasource.MdsConnection{}, sds.GetListing,
		 sync.Mutex{},  map[string]chan *model.ClobQuote{},
		toClientBufferSize,
	}

	go func() {

		namespace := "default"
		clientSet := k8s.GetK8sClientSet(false)
		serviceType := "market-data-gateway"
		pods, err := clientSet.CoreV1().Pods(namespace).Watch(v1.ListOptions{
			LabelSelector: "servicetype=" + serviceType,
		})

		if err != nil {
			panic(err)
		}

		for e := range pods.ResultChan() {
			pod := e.Object.(*v12.Pod)
			bsp, err := loadbalancing.GetBalancingStatefulPod(*pod)
			if err != nil {
				panic(err)
			}

			if e.Type == watch.Added {

				client, err := marketdatasource.NewMdsConnection(id, bsp.TargetAddress, time.Duration(connectRetrySecs)*time.Second,
					maxSubscriptions)
				if err != nil {
					errLog.Printf("failed to create connection to market data source at %v, error: %v", bsp, err)
					continue
				}

				mdService.addMarketDataConnection(bsp, client)
			}
		}
	}()

	port := "50551"
	fmt.Println("starting market data service on port:" + port)
	lis, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		log.Panicf("Error while listening : %v", err)
	}

	s := grpc.NewServer()

	api.RegisterMarketDataServiceServer(s, &mdService)

	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Panicf("Error while serving : %v", err)
	}
}
