package main

import (
	"context"
	api "github.com/ettec/open-trading-platform/go/common/api/executionvenue"
	"github.com/ettec/open-trading-platform/go/common/k8s"
	"github.com/ettec/open-trading-platform/go/model"
	"google.golang.org/grpc"
	v1 "k8s.io/api/core/v1"
	"strconv"
	"time"

	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type orderRouter struct {
	micToExecVenue map[string]*execVenue
}

func New(external bool, maxConnectRetrySecs int) *orderRouter {

	orderRouter := orderRouter{
		micToExecVenue: map[string]*execVenue{},
	}

	clientSet := k8s.GetK8sClientSet(external)

	namespace := "default"
	list, err := clientSet.CoreV1().Services(namespace).List(metav1.ListOptions{
		LabelSelector: "app=execution-venue",
	})

	if err != nil {
		panic(err)
	}

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

		orderRouter.micToExecVenue[mic] = client
		log.Printf("added execution venue for mic: %v, venue service name: %v, target address: %v", mic, service.Name, targetAddress)
	}

	return &orderRouter
}

func (o *orderRouter) CreateAndRouteOrder(c context.Context, p *api.CreateAndRouteOrderParams) (*api.OrderId, error) {
	mic := p.Listing.Market.Mic
	if ev, ok := o.micToExecVenue[mic]; ok {
		id, err := ev.client.CreateAndRouteOrder(c, &api.CreateAndRouteOrderParams{
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

func (o *orderRouter) CancelOrder(c context.Context, p *api.CancelOrderParams) (*model.Empty, error) {
	mic := p.Listing.Market.Mic
	if ev, ok := o.micToExecVenue[mic]; ok {
		_, err := ev.client.CancelOrder(c, p)

		if err != nil {
			return nil, fmt.Errorf("failed to cancel order on market: %v, error: %v", mic, err)
		}

		return &model.Empty{}, nil

	} else {
		return nil, fmt.Errorf("no execution venue found for mic:%v", mic)
	}

}

func createExecVenueConnection(service *v1.Service, maxReconnectInterval time.Duration, targetAddress string) (cac *execVenue,
	err error) {

	log.Printf("connecting to execution venue service %v at: %v", service.Name, targetAddress)

	conn, err := grpc.Dial(targetAddress, grpc.WithInsecure(), grpc.WithBackoffMaxDelay(maxReconnectInterval))

	if err != nil {
		return nil, err
	}

	client := api.NewExecutionVenueClient(conn)

	return &execVenue{
		client: client,
		conn:   conn,
	}, nil
}
