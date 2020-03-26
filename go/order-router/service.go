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
	watch2 "k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/rest"
	"log"
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

type clientAndConn struct {
	client api2.ExecutionVenueClient
	conn   *grpc.ClientConn
}

type orderRouter struct {
	micToExecVenue    map[string]*clientAndConn
	micToExecVenueMux sync.RWMutex
	errLog            *log.Logger
}

func (o orderRouter) getExecVenue(market string) (*clientAndConn, bool) {
	o.micToExecVenueMux.RLock()
	defer o.micToExecVenueMux.RUnlock()
	c, ok := o.micToExecVenue[market]
	return c, ok
}

func (o orderRouter) putExecVenue(market string, client *clientAndConn) error {
	o.micToExecVenueMux.Lock()
	defer o.micToExecVenueMux.Unlock()
	if _, ok := o.micToExecVenue[market]; ok {
		return fmt.Errorf("market to execution venue map already contains an entry for market:%v", market)
	}

	o.micToExecVenue[market] = client

	return nil
}

func (o orderRouter) deleteExecVenue(market string) (*clientAndConn, error) {
	o.micToExecVenueMux.Lock()
	defer o.micToExecVenueMux.Unlock()
	if result, ok := o.micToExecVenue[market]; ok {
		delete(o.micToExecVenue, market)
		return result, nil
	} else {
		return nil, fmt.Errorf("no execution venue entry for market:%v", market)
	}
}

func (o orderRouter) CreateAndRouteOrder(c context.Context, p *api.CreateAndRouteOrderParams) (*api.OrderId, error) {
	mic := p.Listing.Market.Mic
	if ev, ok := o.getExecVenue(mic); ok {
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
	if ev, ok := o.getExecVenue(mic); ok {
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
		micToExecVenue:    map[string]*clientAndConn{},
		micToExecVenueMux: sync.RWMutex{},
		errLog:            log.New(os.Stderr, "", log.Lshortfile|log.Ltime),
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

	pods, err := clientSet.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

	namespace := "default"
	watch, err := clientSet.CoreV1().Pods(namespace).Watch(metav1.ListOptions{
		LabelSelector: "app=execution-venue",
		Watch:         true,
	})

	go func() {
		for event := range watch.ResultChan() {

			pod, ok := event.Object.(*v1.Pod)
			if !ok {
				log.Panic("unexpected type")
			}

			if _, ok := pod.Labels["market"]; !ok {
				orderRouter.errLog.Printf("ignoring execution venue pod as it does not have a market label, pod: %v", pod)
				continue
			}

			market := pod.Labels["market"]

			switch event.Type {
			case watch2.Added:
				client, targetAddress, err := createExecVenueConnection(pod, time.Duration(maxConnectRetrySecs)*time.Second)
				if err != nil {
					orderRouter.errLog.Printf("failed to create connection to execution venue at %v, error: %v", targetAddress, err)
					continue
				}

				orderRouter.putExecVenue(market, client)
				log.Printf("added execution venue for market: %v, venue target address: %v", market, targetAddress)

			case watch2.Deleted:
				client, err := orderRouter.deleteExecVenue(market)
				if err != nil {
					orderRouter.errLog.Printf("failed to delete connection as no execution venue for market %v error: %v", market, err)
					continue
				}

				client.conn.Close()
				log.Printf("removed execution venue for market: %v, venue: %v", market, client)
			}

		}
	}()

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

func createExecVenueConnection(pod *v1.Pod, maxReconnectInterval time.Duration) (cac *clientAndConn,  targetAddress string,
	err error) {

	podIp := pod.Status.PodIP
	var podPort int32
	for _, port := range pod.Spec.Containers[0].Ports {
		if port.Name == "api" {
			podPort = port.ContainerPort
		}
	}

	if podPort == 0 {
		return nil, "", fmt.Errorf("execution venue pod does not have a port named api, pod: %v", pod)
	}

	targetAddress = podIp + ":" + strconv.Itoa(int(podPort))
	log.Printf("connecting to execution venue at:%v", targetAddress)

	conn, err := grpc.Dial(targetAddress, grpc.WithInsecure(), grpc.WithBackoffMaxDelay(maxReconnectInterval))
	client := api2.NewExecutionVenueClient(conn)

	return &clientAndConn{
		client: client,
		conn:   conn,
	}, targetAddress, nil
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
