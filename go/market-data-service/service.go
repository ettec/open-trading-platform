package main

import (
	"context"
	"fmt"
	"github.com/ettec/open-trading-platform/go/common/bootstrap"
	"github.com/ettec/open-trading-platform/go/common/k8s"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/actor"
	mdgapi "github.com/ettec/open-trading-platform/go/market-data-gateway/api"
	"github.com/ettec/open-trading-platform/go/model"
	"github.com/ettech/open-trading-platform/go/market-data-service/api"
	"github.com/ettech/open-trading-platform/go/market-data-service/gatewayclient"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	logger "log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

type service struct {
	micToGateway map[string]*gatewayConnection
}

type gatewayConnection struct {
	partyIdToConnection map[string]actor.ClientConnection
	quoteDistributor    actor.QuoteDistributor
	connMux             sync.Mutex
}

func newGateway(id string, marketGatewayAddress string, maxReconnectInterval time.Duration) (*gatewayConnection, error) {

	mdcToDistributorChan := make(chan *model.ClobQuote, 1000)

	mdcFn := func(targetAddress string) (mdgapi.MarketDataGatewayClient, gatewayclient.GrpcConnection, error) {
		conn, err := grpc.Dial(targetAddress, grpc.WithInsecure(), grpc.WithBackoffMaxDelay(maxReconnectInterval))
		if err != nil {
			return nil, nil, err
		}

		client := mdgapi.NewMarketDataGatewayClient(conn)
		return client, conn, nil
	}

	mdc, err := gatewayclient.NewMarketDataGatewayClient(id, marketGatewayAddress, mdcToDistributorChan, mdcFn)

	if err != nil {
		return nil, err
	}

	qd := actor.NewQuoteDistributor(mdc.Subscribe, mdcToDistributorChan)
	gateway := &gatewayConnection{partyIdToConnection: make(map[string]actor.ClientConnection), quoteDistributor: qd}

	return gateway, nil
}

func (s *gatewayConnection) getConnection(partyId string) (actor.ClientConnection, bool) {
	s.connMux.Lock()
	defer s.connMux.Unlock()

	con, ok := s.partyIdToConnection[partyId]
	return con, ok
}

func (s *gatewayConnection) addConnection(subscriberId string, out chan<- *model.ClobQuote) actor.ClientConnection {
	s.connMux.Lock()
	defer s.connMux.Unlock()

	if conn, ok := s.partyIdToConnection[subscriberId]; ok {
		log.Printf("connection for client %v already exists, closing existing connection.", subscriberId)
		conn.Close()
		log.Print("connection closed:", subscriberId)
	}

	cc := actor.NewClientConnection(subscriberId, out, s.quoteDistributor, maxSubscriptions)

	s.partyIdToConnection[subscriberId] = cc

	return cc
}

func (s *service) Subscribe(_ context.Context, r *api.MdsSubscribeRequest) (*model.Empty, error) {

	mic := r.Listing.Market.Mic
	if gateway, ok := s.micToGateway[mic]; ok {
		if conn, ok := gateway.getConnection(r.SubscriberId); ok {

			if err := conn.Subscribe(r.Listing.Id); err != nil {
				return nil, err
			}

			return &model.Empty{}, nil
		} else {
			return nil, fmt.Errorf("failed to subscribe, no connection exists for subscriber " + r.SubscriberId)
		}

	} else {
		return nil, fmt.Errorf("no gateway exists for mic %v", mic)
	}

}

func (s *service) Connect(request *api.MdsConnectRequest, stream api.MarketDataService_ConnectServer) error {

	subscriberId := request.GetSubscriberId()

	log.Println("connect request received for subscriber ", subscriberId)

	out := make(chan *model.ClobQuote, 100)

	for mic, gateway := range s.micToGateway {
		gateway.addConnection(subscriberId, out)
		log.Printf("connected subscriber %v to gatway for mic %ve", subscriberId, mic)
	}

	for mdUpdate := range out {
		if err := stream.Send(mdUpdate); err != nil {
			log.Printf("error on connection for subscriber %v, closing connection, error:%v", subscriberId, err)
			break
		}
	}

	return nil
}

const (
	ServiceIdKey        = "SERVICE_ID"
	ConnectRetrySeconds = "CONNECT_RETRY_SECONDS"
	External            = "EXTERNAL"
)

var maxSubscriptions = 10000

var log = logger.New(os.Stdout, "", logger.Ltime|logger.Lshortfile)
var errLog = logger.New(os.Stderr, "", logger.Ltime|logger.Lshortfile)

func main() {
	
	id := bootstrap.GetEnvVar(ServiceIdKey)

	connectRetrySecs := bootstrap.GetOptionalIntEnvVar(ConnectRetrySeconds, 60)

	external := bootstrap.GetOptionalBoolEnvVar(External, false)

	mdService := service{micToGateway: map[string]*gatewayConnection{}}

	clientSet := k8s.GetK8sClientSet(external)

	namespace := "default"
	list, err := clientSet.CoreV1().Services(namespace).List(v1.ListOptions{
		LabelSelector: "app=market-data-service",
	})

	for _, service := range list.Items {
		const micLabel = "mic"
		if _, ok := service.Labels[micLabel]; !ok {
			errLog.Printf("ignoring market data service as it does not have a mic label, service: %v", service)
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
			log.Printf("ignoring market data service as it does not have a port named api, service: %v", service)
			continue
		}

		targetAddress := service.Name + ":" + strconv.Itoa(int(podPort))

		client, err := newGateway(id, targetAddress, time.Duration(connectRetrySecs)*time.Second)
		if err != nil {
			errLog.Printf("failed to create connection to execution venue service at %v, error: %v", targetAddress, err)
			continue
		}

		mdService.micToGateway[mic] = client

		log.Printf("added market data service for mic: %v, service name: %v, target address: %v", mic, service.Name, targetAddress)
	}

	port := "50551"
	fmt.Println("Starting Market Data Service on port:" + port)
	lis, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		log.Fatalf("Error while listening : %v", err)
	}

	s := grpc.NewServer()
	if err != nil {
		log.Panicf("failed to create market data service:%v", err)
	}

	api.RegisterMarketDataServiceServer(s, &mdService)

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error while serving : %v", err)
	}

}
