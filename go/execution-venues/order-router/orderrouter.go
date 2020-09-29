package main

import (
	"context"
	"fmt"
	"github.com/ettec/otp-common/api/executionvenue"
	"github.com/ettec/otp-common/k8s"
	"github.com/ettec/otp-common/loadbalancing"
	"github.com/ettec/otp-common/model"
	"google.golang.org/grpc"
	v12 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"log"
	"sync"
	"time"
)

type orderRouter struct {
	micToExecVenue     map[string]map[int]*execVenue
	ownerIdToExecVenue map[string]*execVenue
	venueMux           sync.Mutex
}

func New(connectRetrySecs int) *orderRouter {

	router := &orderRouter{
		micToExecVenue:     map[string]map[int]*execVenue{},
		ownerIdToExecVenue: map[string]*execVenue{},
		venueMux:           sync.Mutex{},
	}

	go func() {

		namespace := "default"
		clientSet := k8s.GetK8sClientSet(false)
		serviceType := "execution-venue"
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

				client, err := createExecVenueConnection(time.Duration(connectRetrySecs)*time.Second, bsp.TargetAddress)
				if err != nil {
					errLog.Printf("failed to create connection to execution venue service at %v, error: %v", bsp.TargetAddress, err)
					continue
				}

				router.addExecVenue(bsp, client)
			} else if e.Type == watch.Deleted {
				router.removeExecVenue(bsp)
			}
		}
	}()

	return router
}

func (o *orderRouter) removeExecVenue(bsp *loadbalancing.BalancingStatefulPod) {
	o.venueMux.Lock()
	defer o.venueMux.Unlock()

	if ordToEv, ok := o.micToExecVenue[bsp.Mic]; ok {
		delete(ordToEv, bsp.Ordinal)
	}

	delete(o.ownerIdToExecVenue, bsp.Name)

	log.Printf("removed execution venue for mic: %v,  target address: %v, stateful set ordinal %v", bsp.Mic, bsp.TargetAddress, bsp.Ordinal)
}

func (o *orderRouter) addExecVenue(bsp *loadbalancing.BalancingStatefulPod, ev *execVenue) {
	o.venueMux.Lock()
	defer o.venueMux.Unlock()

	if _, ok := o.micToExecVenue[bsp.Mic]; !ok {
		o.micToExecVenue[bsp.Mic] = map[int]*execVenue{}
	}

	o.micToExecVenue[bsp.Mic][bsp.Ordinal] = ev

	o.ownerIdToExecVenue[bsp.Name] = ev

	log.Printf("added execution venue for mic: %v,  target address: %v, stateful set ordinal %v", bsp.Mic, bsp.TargetAddress, bsp.Ordinal)

}

func (o *orderRouter) GetExecutionParametersMetaData(context.Context, *model.Empty) (*executionvenue.ExecParamsMetaDataJson, error) {
	return &executionvenue.ExecParamsMetaDataJson{}, nil
}

func (o *orderRouter) CreateAndRouteOrder(c context.Context, p *executionvenue.CreateAndRouteOrderParams) (*executionvenue.OrderId, error) {

	ev, err := o.getExecutionVenueForListing(p.ListingId, p.Destination)

	if err != nil {
		return nil, err
	}

	id, err := ev.client.CreateAndRouteOrder(c, p)

	log.Printf("routed create order request %v to execution venue %v, returned order id %v", p, ev, id)

	if err != nil {
		return nil, fmt.Errorf("failed to route order:%v", err)
	}

	return id, nil

}

func (o *orderRouter) getExecutionVenueForOwnerId(ownerId string) (*execVenue, error) {
	o.venueMux.Lock()
	defer o.venueMux.Unlock()
	if ev, ok := o.ownerIdToExecVenue[ownerId]; ok {
		return ev, nil

	} else {
		return nil, fmt.Errorf("failed to find execution venue for owner id:%v", ownerId)
	}

}

func (o *orderRouter) getExecutionVenueForListing(listingId int32, destination string) (*execVenue, error) {
	o.venueMux.Lock()
	defer o.venueMux.Unlock()
	if evs, ok := o.micToExecVenue[destination]; ok {
		numVenues := int32(len(evs))
		ordinal := loadbalancing.GetBalancingOrdinal(listingId, numVenues)
		return evs[ordinal], nil
	} else {
		return nil, fmt.Errorf("no execution venue found for destination %v", destination)
	}
}

func (o *orderRouter) ModifyOrder(c context.Context, p *executionvenue.ModifyOrderParams) (*model.Empty, error) {

	ev, err := o.getExecutionVenueForOwnerId(p.OwnerId)

	if err != nil {
		return nil, fmt.Errorf("failed to modify order %v, error: %v", p.OrderId, err)
	}

	_, err = ev.client.ModifyOrder(c, p)

	if err != nil {
		return nil, fmt.Errorf("failed to modify order: %v, error: %v", p.OrderId, err)
	}

	return &model.Empty{}, nil
}

func (o *orderRouter) CancelOrder(c context.Context, p *executionvenue.CancelOrderParams) (*model.Empty, error) {

	ev, err := o.getExecutionVenueForOwnerId(p.OwnerId)

	if err != nil {
		return nil, fmt.Errorf("failed to cancel order %v, error: %v", p.OrderId, err)
	}

	_, err = ev.client.CancelOrder(c, p)

	if err != nil {
		return nil, fmt.Errorf("failed to cancel order %v: , error: %v", p.OrderId, err)
	}

	return &model.Empty{}, nil
}

func createExecVenueConnection(maxReconnectInterval time.Duration, targetAddress string) (cac *execVenue,
	err error) {

	log.Printf("connecting to execution venue service at: %v", targetAddress)

	conn, err := grpc.Dial(targetAddress, grpc.WithInsecure(), grpc.WithBackoffMaxDelay(maxReconnectInterval))

	if err != nil {
		return nil, err
	}

	client := executionvenue.NewExecutionVenueClient(conn)

	return &execVenue{
		client: client,
		conn:   conn,
	}, nil
}
