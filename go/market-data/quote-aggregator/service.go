package main

import (
	"github.com/ettec/otp-common/api/marketdatasource"
	"github.com/ettec/otp-common/bootstrap"
	"github.com/ettec/otp-common/k8s"
	marketdata "github.com/ettec/otp-mdcommon"
	"github.com/ettec/otp-mdcommon/quotestream"
	"github.com/ettec/otp-mdcommon/source"

	"github.com/ettec/open-trading-platform/go/market-data/quote-aggregator/quoteaggregator"
	"github.com/ettec/otp-common/staticdata"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	logger "log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	ServiceIdKey             = "SERVICE_ID"
	ConnectRetrySeconds      = "CONNECT_RETRY_SECONDS"
	StaticDataServiceAddress = "STATIC_DATA_SERVICE_ADDRESS"
	External                 = "EXTERNAL"
)

var log = logger.New(os.Stdout, "", logger.Ltime|logger.Lshortfile)
var errLog = logger.New(os.Stderr, "", logger.Ltime|logger.Lshortfile)

func main() {

	id := bootstrap.GetOptionalEnvVar(ServiceIdKey, "quoteaggregator")

	connectRetrySecs := bootstrap.GetOptionalIntEnvVar(ConnectRetrySeconds, 60)

	external := bootstrap.GetOptionalBoolEnvVar(External, false)

	staticDataServiceAddress := bootstrap.GetEnvVar(StaticDataServiceAddress)

	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":8080", nil)

	micToMdsAddress := map[string]string{}

	clientSet := k8s.GetK8sClientSet(external)

	namespace := "default"
	list, err := clientSet.CoreV1().Services(namespace).List(v1.ListOptions{
		LabelSelector: "app=market-data-source",
	})

	if err != nil {
		panic(err)
	}

	log.Printf("found %v market data sources", len(list.Items))

	for _, service := range list.Items {
		const micLabel = "mic"
		if _, ok := service.Labels[micLabel]; !ok {
			errLog.Printf("ignoring market data source as it does not have a mic label, marketDataSource: %v", service)
			continue
		}

		mic := service.Labels[micLabel]

		var podPort int32
		for _, port := range service.Spec.Ports {
			if port.Name == "api" {
				podPort = port.Port
			}
		}

		if podPort == 0 {
			log.Printf("ignoring market data marketDataSource as it does not have a port named api, marketDataSource: %v", service)
			continue
		}

		targetAddress := service.Name + ":" + strconv.Itoa(int(podPort))

		micToMdsAddress[mic] = targetAddress

		log.Printf("found market data source for mic: %v, marketDataSource name: %v, target address: %v", mic, service.Name, targetAddress)
	}

	mdcFn := func(targetAddress string) (marketdatasource.MarketDataSourceClient, quotestream.GrpcConnection, error) {
		conn, err := grpc.Dial(targetAddress, grpc.WithInsecure(), grpc.WithBackoffMaxDelay(time.Duration(connectRetrySecs)*time.Second))
		if err != nil {
			return nil, nil, err
		}

		client := marketdatasource.NewMarketDataSourceClient(conn)
		return client, conn, nil
	}

	sds, err := staticdata.NewStaticDataSource(staticDataServiceAddress)
	if err != nil {
		panic(err)
	}

	quoteAggregator := quoteaggregator.New(id, sds.GetListingsWithSameInstrument,
		micToMdsAddress, 1000, mdcFn)

	mdSource := source.NewMarketDataSource(marketdata.NewQuoteDistributor(quoteAggregator, 1000))

	port := "50551"
	log.Println("Starting Market Data Service on port:" + port)
	lis, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		log.Fatalf("Error while listening : %v", err)
	}

	s := grpc.NewServer()
	if err != nil {
		log.Panicf("failed to create market data marketDataSource:%v", err)
	}

	marketdatasource.RegisterMarketDataSourceServer(s, mdSource)

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error while serving : %v", err)
	}

}
