package main

import (
	"context"
	"github.com/ettec/open-trading-platform/go/common/bootstrap"
	api2 "github.com/ettec/open-trading-platform/go/execution-venue/api"
	"github.com/ettec/open-trading-platform/go/model"
	"github.com/ettec/open-trading-platform/go/order-router/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/reflection"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	watch2 "k8s.io/apimachinery/pkg/watch"
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
	micToExecVenue    map[string][]*execVenue
	micToExecVenueMux sync.RWMutex
}

func (o orderRouter) getConnectedExecVenue(market string) (*execVenue, bool) {
	o.micToExecVenueMux.RLock()
	defer o.micToExecVenueMux.RUnlock()
	venues, ok := o.micToExecVenue[market]
	for _, venue := range venues {

		if venue.conn.GetState() == connectivity.Ready {
			return venue, ok
		}

	}

	return nil, false
}

func (o orderRouter) putExecVenue(market string, client *execVenue) error {
	o.micToExecVenueMux.Lock()
	defer o.micToExecVenueMux.Unlock()
	o.micToExecVenue[market] = append(o.micToExecVenue[market], client)

	return nil
}

func (o orderRouter) hasExecVenue(podId types.UID) bool {

	for _, venues := range o.micToExecVenue {
		for _, venue := range venues {
			if venue.podId == podId {
				return true
			}
		}
	}

	return false
}

func (o orderRouter) deleteExecVenue(podId types.UID) (*execVenue, error) {
	o.micToExecVenueMux.Lock()
	defer o.micToExecVenueMux.Unlock()

	var removed *execVenue

	micToExecVenues := map[string][]*execVenue{}
	for mic, venues := range o.micToExecVenue {
		var newVenues []*execVenue
		for _, venue := range venues {
			if venue.podId != podId {
				newVenues = append(newVenues, venue)
			} else {
				removed = venue
			}
		}
		micToExecVenues[mic] = newVenues
	}

	if removed != nil {
		return removed, nil
	} else {
		return nil, fmt.Errorf("no exec venue found for podId:%v", podId)
	}

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
		micToExecVenue:    map[string][]*execVenue{},
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
				errLog.Printf("ignoring execution venue pod as it does not have a market label, pod: %v", pod)
				continue
			}

			market := pod.Labels["market"]

			switch event.Type {
			case watch2.Added:
				fallthrough
			case watch2.Modified:

				if !orderRouter.hasExecVenue(pod.UID) {
					if targetAddress, ok := getTargetAddress(pod); ok {
						client, err := createExecVenueConnection(pod, time.Duration(maxConnectRetrySecs)*time.Second, targetAddress)
						if err != nil {
							errLog.Printf("failed to create connection to execution venue at %v, error: %v", targetAddress, err)
							continue
						}

						orderRouter.putExecVenue(market, client)
						log.Printf("added execution venue for market: %v, venue pod name: %v, target address: %v", market, pod.Name, targetAddress)
					}
				}

			case watch2.Deleted:
				client, err := orderRouter.deleteExecVenue(pod.UID)
				if err != nil {
					errLog.Printf("failed to delete connection as no execution venue for market %v error: %v", market, err)
					continue
				}

				client.conn.Close()
				log.Printf("removed execution venue for market: %v, venue pod name: %v", market, pod.Name)
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

func createExecVenueConnection(pod *v1.Pod, maxReconnectInterval time.Duration, targetAddress string) (cac *execVenue,
	err error) {

	log.Printf("connecting to execution venue pod %v at: %v", pod.Name, targetAddress)

	conn, err := grpc.Dial(targetAddress, grpc.WithInsecure(), grpc.WithBackoffMaxDelay(maxReconnectInterval))
	client := api2.NewExecutionVenueClient(conn)

	conn.GetState()

	return &execVenue{
		podId:  pod.UID,
		client: client,
		conn:   conn,
	}, nil
}

func getTargetAddress(pod *v1.Pod) (targetAddress string, ok bool) {
	podIp := pod.Status.PodIP

	if podIp == "" {
		return "", false
	}

	var podPort int32
	for _, port := range pod.Spec.Containers[0].Ports {
		if port.Name == "api" {
			podPort = port.ContainerPort
		}
	}

	if podPort == 0 {
		log.Printf("execution venue pod does not have a port named api, pod: %v", pod)
		return "", false

	}

	targetAddress = podIp + ":" + strconv.Itoa(int(podPort))
	return targetAddress, true
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
