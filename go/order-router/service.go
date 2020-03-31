package main

import (
	"context"
	"github.com/ettec/open-trading-platform/go/common/bootstrap"
	api2 "github.com/ettec/open-trading-platform/go/execution-venue/api"
	"github.com/ettec/open-trading-platform/go/model"
	"github.com/ettec/open-trading-platform/go/order-router/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	logger "log"
	"net"
	"strconv"
	"sync"
	"time"

	"flag"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
)

var log = logger.New(os.Stdout, "", logger.Ltime|logger.Lshortfile)
var errLog = logger.New(os.Stderr, "", logger.Ltime|logger.Lshortfile)

type execVenue struct {
	podId  types.UID
	client api2.ExecutionVenueClient
	conn   *grpc.ClientConn
}

type orderRouter struct {
	micToExecVenue    map[string]*execVenue
	micToExecVenueMux sync.RWMutex
}

func (o orderRouter) getConnectedExecVenue(market string) (venue *execVenue, ok bool) {
	o.micToExecVenueMux.RLock()
	defer o.micToExecVenueMux.RUnlock()
	venue, ok = o.micToExecVenue[market]
	return venue, ok
}

func (o orderRouter) putExecVenue(market string, client *execVenue)  {
	o.micToExecVenueMux.Lock()
	defer o.micToExecVenueMux.Unlock()
	o.micToExecVenue[market] =  client
}



func (o orderRouter) CreateAndRouteOrder(c context.Context, p *api.CreateAndRouteOrderParams) (*api.OrderId, error) {
	mic := p.Listing.Market.Mic
	if ev, ok := o.getConnectedExecVenue(mic); ok {
		id, err := ev.client.CreateAndRouteOrder(c, &api2.CreateAndRouteOrderParams{
			OrderSide: p.OrderSide,
			Quantity:  p.Quantity,
			Price:     p.Price,
			Listing:   p.Listing,
		})

		if err != nil {
			return nil, fmt.Errorf("failed to route order:%v", err)
		}

		return &api.OrderId{
			OrderId: id.OrderId,
		}, nil

	} else {
		return nil, fmt.Errorf("no execution venue found for mic:%v", mic)
	}
}

func (o orderRouter) CancelOrder(c context.Context, p *api.CancelOrderParams) (*model.Empty, error) {
	mic := p.Listing.Market.Mic
	if ev, ok := o.getConnectedExecVenue(mic); ok {
		_, err := ev.client.CancelOrder(c, &api2.OrderId{
			OrderId: p.OrderId,
		})

		if err != nil {
			return nil, fmt.Errorf("failed to cancel order on market: %v, error: %v", mic, err)
		}

		return &model.Empty{}, nil

	} else {
		return nil, fmt.Errorf("no execution venue found for mic:%v", mic)
	}

}

const (
	MaxConnectRetrySeconds = "MAX_CONNECT_RETRY_SECONDS"
	External               = "EXTERNAL"
)

func main() {

	maxConnectRetrySecs := bootstrap.GetOptionalIntEnvVar(MaxConnectRetrySeconds, 60)
	external := bootstrap.GetOptionalBoolEnvVar(External, false)

	orderRouter := orderRouter{
		micToExecVenue:    map[string]*execVenue{},
		micToExecVenueMux: sync.RWMutex{},
	}

	var clientSet *kubernetes.Clientset
	if external {
		var kubeconfig *string
		if home := homeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		flag.Parse()

		// use the current context in kubeconfig
		config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			panic(err.Error())
		}

		// create the clientSet
		clientSet, err = kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}

	} else {
		config, err := rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
		// creates the clientSet
		clientSet, err = kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}
	}


	namespace := "default"
	list, err := clientSet.CoreV1().Services(namespace).List(metav1.ListOptions{
		LabelSelector: "app=execution-venue",
	})

	for _, service := range list.Items {
		const micLabel = "mic"
		if _, ok := service.Labels[micLabel]; !ok {
			errLog.Printf("ignoring execution venue service as it does not have a mic label, service: %v", service)
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
			log.Printf("ignoring execution venue service as it does not have a port named api, service: %v", service)
			continue
		}

		targetAddress := service.Name + ":" + strconv.Itoa(int(podPort))

		client, err := createExecVenueConnection(&service, time.Duration(maxConnectRetrySecs)*time.Second, targetAddress)
		if err != nil {
			errLog.Printf("failed to create connection to execution venue service at %v, error: %v", targetAddress, err)
			continue
		}

		orderRouter.putExecVenue(mic, client)
		log.Printf("added execution venue for mic: %v, venue service name: %v, target address: %v", mic, service.Name, targetAddress)
	}

	port := "50581"
	fmt.Println("Starting Order Router on port:" + port)
	lis, err := net.Listen("tcp", "0.0.0.0:"+port)

	if err != nil {
		log.Fatalf("Error while listening : %v", err)
	}

	s := grpc.NewServer()

	api.RegisterOrderRouterServer(s, &orderRouter)

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error while serving : %v", err)

	}

}

func createExecVenueConnection(service *v1.Service, maxReconnectInterval time.Duration, targetAddress string) (cac *execVenue,
	err error) {

	log.Printf("connecting to execution venue service %v at: %v", service.Name, targetAddress)

	conn, err := grpc.Dial(targetAddress, grpc.WithInsecure(), grpc.WithBackoffMaxDelay(maxReconnectInterval))
	client := api2.NewExecutionVenueClient(conn)

	conn.GetState()

	return &execVenue{
		client: client,
		conn:   conn,
	}, nil
}



func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
