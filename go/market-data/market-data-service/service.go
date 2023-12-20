package main

import (
	"context"
	"fmt"
	api "github.com/ettec/otp-common/api/marketdataservice"
	"github.com/ettec/otp-common/bootstrap"
	"github.com/ettec/otp-common/k8s"
	"github.com/ettec/otp-common/loadbalancing"
	"github.com/ettec/otp-common/marketdata"
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
	"log/slog"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"syscall"
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

type NewQuoteStreamFromMdSourceFunc func(ctx context.Context, id string, targetAddress string, maxReconnectInterval time.Duration,
	quoteBufferSize int) (marketdata.QuoteStream, error)

func (f NewQuoteStreamFromMdSourceFunc) NewQuoteStreamFromMdSource(ctx context.Context, id string, targetAddress string, maxReconnectInterval time.Duration,
	quoteBufferSize int) (marketdata.QuoteStream, error) {
	return f(ctx, id, targetAddress, maxReconnectInterval, quoteBufferSize)
}

type connectionFactory interface {
	Connect(ctx context.Context, subscriberId string) marketdata.QuoteStream
	AddMarketDataGateway(gateway marketdatasource.MarketDataGateway) error
}

type service struct {
	connectionFactory        connectionFactory
	subscriberIdToConnection map[string]marketdata.QuoteStream
	mutex                    sync.Mutex
}

func (s *service) Subscribe(_ context.Context, r *api.MdsSubscribeRequest) (*model.Empty, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if quoteStream, ok := s.subscriberIdToConnection[r.SubscriberId]; ok {
		if err := quoteStream.Subscribe(r.ListingId); err != nil {
			return nil, fmt.Errorf("failed to subscribe, subscriber %v, listing %v, error: %w", r.SubscriberId, r.ListingId, err)
		}
	} else {
		return nil, fmt.Errorf("failed to subscribe, no connection exists for subscriber " + r.SubscriberId)
	}

	return nil, nil
}

func (s *service) Connect(request *api.MdsConnectRequest, stream api.MarketDataService_ConnectServer) error {
	subscriberId := request.GetSubscriberId()
	slog.Info("connect request received", "subscriberId", subscriberId)

	s.mutex.Lock()
	defer s.mutex.Unlock()
	if connection, ok := s.subscriberIdToConnection[subscriberId]; ok {
		slog.Info("connection already exists, closing existing connection", "subscriberId", subscriberId)
		connection.Close()
		connections.Dec()
	}

	connection := s.connectionFactory.Connect(stream.Context(), subscriberId)
	connections.Inc()
	defer func() {
		connection.Close()
		connections.Dec()
	}()

	s.subscriberIdToConnection[subscriberId] = connection
	s.mutex.Unlock()

	for {
		select {
		case <-stream.Context().Done():
			slog.Info("connection closed", "subscriberId", subscriberId)
			return nil
		case quote := <-connection.Chan():
			if err := stream.Send(quote); err != nil {
				slog.Error("failed to send quote, closing connection", "subscriberId", subscriberId, "error", err)
				return fmt.Errorf("failed to send quote: %w", err)
			}
			quotesSent.Inc()
		}
	}
}

func main() {

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	id := bootstrap.GetEnvVar("MDS_ID")
	connectRetrySecs := bootstrap.GetOptionalIntEnvVar("CONNECT_RETRY_SECONDS", 60)
	maxSubscriptions := bootstrap.GetOptionalIntEnvVar("MAX_SUBSCRIPTIONS", 10000)
	toClientBufferSize := bootstrap.GetOptionalIntEnvVar("TO_CLIENT_BUFFER_SIZE", 1000)

	http.Handle("/metrics", promhttp.Handler())
	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			slog.Error("failed to listen on metrics server port", "error", err)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sds, err := staticdata.NewStaticDataSource(ctx)
	if err != nil {
		log.Panicf("failed to get static data source:%v", err)
	}

	cf := marketdatasource.NewMarketDataService(ctx, id,
		NewQuoteStreamFromMdSourceFunc(marketdata.NewQuoteStreamFromMdSource),
		sds.GetListing, toClientBufferSize, connectRetrySecs, maxSubscriptions)

	namespace := "default"
	clientSet := k8s.GetK8sClientSet(false)

	labelSelector := "servicetype in (market-data-gateway, execution-venue-and-market-data-gateway)"
	pods, err := clientSet.CoreV1().Pods(namespace).Watch(v1.ListOptions{
		LabelSelector: labelSelector,
	})

	if err != nil {
		log.Panicf("Error watching pods: %v", err)
	}

	go func() {
		for e := range pods.ResultChan() {
			pod := e.Object.(*v12.Pod)
			bsp, err := loadbalancing.GetBalancingStatefulPod(*pod)
			if err != nil {
				slog.Error("failed to get balancing stateful pod", "pod", pod, "error", err)
				continue
			}

			if e.Type == watch.Added {
				if err = cf.AddMarketDataGateway(marketDataService{bsp: bsp}); err != nil {
					slog.Error("failed to add new gateway", "balancingStatefulPod", bsp, "error", err)
				}
			}
		}
	}()

	port := "50551"
	slog.Info("starting market data service", "port", port)
	lis, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		log.Panicf("Error while listening : %v", err)
	}

	service := &service{connectionFactory: cf, subscriberIdToConnection: map[string]marketdata.QuoteStream{}}

	s := grpc.NewServer()

	api.RegisterMarketDataServiceServer(s, service)

	reflection.Register(s)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh,
		syscall.SIGKILL,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		<-sigCh
		s.GracefulStop()
	}()

	if err := s.Serve(lis); err != nil {
		log.Panicf("Error while serving : %v", err)
	}
}

type marketDataService struct {
	bsp *loadbalancing.BalancingStatefulPod
}

func (m marketDataService) GetAddress() string {
	return m.bsp.TargetAddress
}
func (m marketDataService) GetOrdinal() int {
	return m.bsp.Ordinal
}
func (m marketDataService) GetMarketMic() string {
	return m.bsp.Mic
}
