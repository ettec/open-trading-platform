package main

import (
	"fmt"
	api "github.com/ettec/open-trading-platform/go/common/api/executionvenue"
	"github.com/ettec/open-trading-platform/go/common/k8s"
	"github.com/ettec/open-trading-platform/go/common/marketdata"
	"github.com/ettec/open-trading-platform/go/execution-venues/common/ordergateway"
	"github.com/ettec/open-trading-platform/go/model"
	"k8s.io/client-go/kubernetes"
	"os"
	"strconv"
	"time"

	"github.com/ettec/open-trading-platform/go/execution-venues/common/executionvenue"

	"github.com/ettec/open-trading-platform/go/common/bootstrap"

	"github.com/ettec/open-trading-platform/go/execution-venues/common/ordercache"
	"github.com/ettec/open-trading-platform/go/execution-venues/common/orderstore"

	"github.com/ettec/open-trading-platform/go/execution-venues/common/ordermanager"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	logger "log"
	"net"
	"strings"
)

const (
	Id                     = "ID"
	KafkaBrokersKey        = "KAFKA_BROKERS"
	ExecVenueMic           = "MIC"
	External               = "EXTERNAL"
	MaxConnectRetrySeconds = "MAX_CONNECT_RETRY_SECONDS"
)

var log = logger.New(os.Stdout, "", logger.Ltime|logger.Lshortfile)
var errLog = logger.New(os.Stderr, "", logger.Ltime|logger.Lshortfile)

func main() {

	id := bootstrap.GetOptionalEnvVar(Id, "smart-router")
	maxConnectRetry := time.Duration(bootstrap.GetOptionalIntEnvVar(MaxConnectRetrySeconds, 60)) * time.Second
	external := bootstrap.GetOptionalBoolEnvVar(External, false)
	kafkaBrokers := bootstrap.GetEnvVar(KafkaBrokersKey)
	execVenueMic := bootstrap.GetEnvVar(ExecVenueMic)

	s := grpc.NewServer()

	store, err := orderstore.NewKafkaStore(strings.Split(kafkaBrokers, ","), execVenueMic)
	if err != nil {
		panic(fmt.Errorf("failed to create order store: %v", err))
	}

	orderCache, err := ordercache.NewOrderCache(store)
	if err != nil {
		log.Fatalf("failed to create order cache:%v", err)
	}

	clientSet := k8s.GetK8sClientSet(external)

	namespace := "default"
	xosrServiceLabelSelector := "app=market-data-source,mic=XOSR"
	list, err := clientSet.CoreV1().Services(namespace).List(metav1.ListOptions{
		LabelSelector: xosrServiceLabelSelector,
	})

	if err != nil {
		panic(err)
	}

	if len(list.Items) != 1 {
		log.Panicf("no service found for selector: %v", xosrServiceLabelSelector)
	}

	service := list.Items[0]

	var podPort int32
	for _, port := range service.Spec.Ports {
		if port.Name == "api" {
			podPort = port.Port
		}
	}

	if podPort == 0 {
		log.Panic("aggregate quote source does not have an 'api' port")
	}

	targetAddress := service.Name + ":" + strconv.Itoa(int(podPort))

	mdsQuoteStream, err := marketdata.NewMdsQuoteStream(id, targetAddress, maxConnectRetry, 1000)

	if err != nil {
		panic(err)
	}

	client, err := getOrderRouter(clientSet, maxConnectRetry)
	if err != nil {
		panic(err)
	}

	om := ordermanager.NewOrderManager(orderCache, NewSmartRouterOrderGateway(client, mdsQuoteStream,

		func(listingId int32, listingGroupsIn chan<- []*model.Listing) {}), execVenueMic)

	sr := executionvenue.New(om)
	defer sr.Close()

	api.RegisterExecutionVenueServer(s, sr)

	reflection.Register(s)

	port := "50551"
	fmt.Println("Starting Execution Venue Service on port:" + port)
	lis, err := net.Listen("tcp", "0.0.0.0:"+port)

	if err != nil {
		log.Fatalf("Error while listening : %v", err)
	}

	if err := s.Serve(lis); err != nil {
		log.Fatalf("error   while serving : %v", err)
	}

}

func getOrderRouter(clientSet *kubernetes.Clientset, maxConnectRetrySecs time.Duration) (api.ExecutionVenueClient, error) {
	namespace := "default"
	list, err := clientSet.CoreV1().Services(namespace).List(metav1.ListOptions{
		LabelSelector: "app=order-router",
	})

	if err != nil {
		panic(err)
	}

	var client api.ExecutionVenueClient

	for _, service := range list.Items {

		var podPort int32
		for _, port := range service.Spec.Ports {
			if port.Name == "api" {
				podPort = port.Port
			}
		}

		if podPort == 0 {
			log.Printf("ignoring order router service as it does not have a port named api, service: %v", service)
			continue
		}

		targetAddress := service.Name + ":" + strconv.Itoa(int(podPort))

		log.Printf("connecting to order router service %v at: %v", service.Name, targetAddress)

		conn, err := grpc.Dial(targetAddress, grpc.WithInsecure(), grpc.WithBackoffMaxDelay(maxConnectRetrySecs))

		if err != nil {
			panic(err)
		}

		client = api.NewExecutionVenueClient(conn)
		break
	}

	if client == nil {
		return nil, fmt.Errorf("failed to find order router")
	}

	return client, nil
}

const QuoteAggregatorMic = "XOSR"

type gateway struct {
	sendChan   chan sendArgs
	cancelChan chan cancelArgs
	orderState chan childOrderUpdates
}

type getListingsWithSameInstrument = func(listingId int32, listingGroupsIn chan<- []*model.Listing)

type childOrderUpdates struct {
	orderId      string
	updatedChild *model.Order
}

func NewSmartRouterOrderGateway(orderRouter api.ExecutionVenueClient, quoteStream marketdata.MdsQuoteStream,
	getListings getListingsWithSameInstrument) ordergateway.OrderGateway {
	gateway := gateway{
		sendChan:   make(chan sendArgs, 100),
		cancelChan: make(chan cancelArgs, 100),
	}

	srListingToListings := map[int32][]*model.Listing{}

	listingToLastQuote := map[int32][]*model.ClobQuote{}

	listingGroupsIn := make(chan []*model.Listing, 1000)

	sendsPendingQuotes := map[int32][]sendArgs{}

	go func() {

		for {
			select {
			case s := <-gateway.sendChan:
				if listings, ok := srListingToListings[s.listing.Id]; ok {

					quotesPending := false
					for _, listing := range listings {
						if _, ok := listingToLastQuote[listing.Id]; !ok {
							quotesPending = true
							sendsPendingQuotes[listing.Id] = append(sendsPendingQuotes[listing.Id], s)
							break
						}
					}

					if !quotesPending {

						// here - pick the best listing, split the order?

					}

				} else {
					getListings(s.listing.Id, listingGroupsIn)
				}
			case listings := <-listingGroupsIn:
				for _, listing := range listings {
					var toListings []*model.Listing
					var listingId int32
					if listing.Market.Mic == QuoteAggregatorMic {
						listingId = listing.Id
					} else {
						toListings = append(toListings, listing)
					}

					srListingToListings[listingId] = toListings
				}

			}

		}

	}()

	return &gateway
}

type sendArgs struct {
	order   *model.Order
	listing *model.Listing
}

func (s *gateway) Send(order *model.Order, listing *model.Listing) error {
	s.sendChan <- sendArgs{
		order:   order,
		listing: listing,
	}

	return nil
}

type cancelArgs struct {
	order *model.Order
}

func (s *gateway) Cancel(order *model.Order) error {

	s.cancelChan <- cancelArgs{
		order: order,
	}

	return nil

}
