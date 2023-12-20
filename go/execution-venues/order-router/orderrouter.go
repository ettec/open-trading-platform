package main

import (
	"context"
	"fmt"
	"github.com/ettec/otp-common/api/executionvenue"
	"github.com/ettec/otp-common/k8s"
	"github.com/ettec/otp-common/loadbalancing"
	"github.com/ettec/otp-common/model"
	"google.golang.org/grpc"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"log"
	"log/slog"
	"sync"
	"time"
)

type orderRouter struct {
	micToExecVenue     map[string]map[int]*execVenue
	ownerIdToExecVenue map[string]*execVenue
	mux                sync.Mutex
}

func NewOrderRouter(_ context.Context, connectRetrySecs int) (*orderRouter, error) {

	router := &orderRouter{
		micToExecVenue:     map[string]map[int]*execVenue{},
		ownerIdToExecVenue: map[string]*execVenue{},
		mux:                sync.Mutex{},
	}

	namespace := "default"
	clientSet := k8s.GetK8sClientSet(false)
	labelSelector := "servicetype in (execution-venue, execution-venue-and-market-data-gateway)"
	pods, err := clientSet.CoreV1().Pods(namespace).Watch(v1.ListOptions{
		LabelSelector: labelSelector,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to watch pods with label selector \"%s\", error: %v", labelSelector, err)
	}

	go func() {

		for e := range pods.ResultChan() {
			pod := e.Object.(*corev1.Pod)
			bsp, err := loadbalancing.GetBalancingStatefulPod(*pod)
			if err != nil {
				log.Panicf("failed to get balancing stateful pod: %v", err)
			}

			if e.Type == watch.Added {
				client, err := createExecVenueConnection(time.Duration(connectRetrySecs)*time.Second, bsp.TargetAddress)
				if err != nil {
					slog.Error("failed to create connection to execution venue service", "targetAddress", bsp.TargetAddress, "error", err)
					continue
				}

				router.addExecVenue(bsp, client)
			} else if e.Type == watch.Deleted {
				router.removeExecVenue(bsp)
			}
		}
	}()

	return router, nil
}

func (o *orderRouter) removeExecVenue(bsp *loadbalancing.BalancingStatefulPod) {
	o.mux.Lock()
	defer o.mux.Unlock()

	if ordinalToExecutionVenue, ok := o.micToExecVenue[bsp.Mic]; ok {
		delete(ordinalToExecutionVenue, bsp.Ordinal)
	}

	delete(o.ownerIdToExecVenue, bsp.Name)

	slog.Info("removed execution venue", "mic", bsp.Mic,
		"targetAddress", bsp.TargetAddress, "ordinal", bsp.Ordinal)
}

func (o *orderRouter) addExecVenue(bsp *loadbalancing.BalancingStatefulPod, ev *execVenue) {
	o.mux.Lock()
	defer o.mux.Unlock()

	if _, ok := o.micToExecVenue[bsp.Mic]; !ok {
		o.micToExecVenue[bsp.Mic] = map[int]*execVenue{}
	}

	o.micToExecVenue[bsp.Mic][bsp.Ordinal] = ev

	o.ownerIdToExecVenue[bsp.Name] = ev

	slog.Info("added execution venue", "mic", bsp.Mic,
		"targetAddress", bsp.TargetAddress, "ordinal", bsp.Ordinal)
}

func (o *orderRouter) GetExecutionParametersMetaData(context.Context, *model.Empty) (*executionvenue.ExecParamsMetaDataJson, error) {
	return &executionvenue.ExecParamsMetaDataJson{}, nil
}

func (o *orderRouter) CreateAndRouteOrder(c context.Context, p *executionvenue.CreateAndRouteOrderParams) (*executionvenue.OrderId, error) {

	ev, err := o.getExecutionVenueForListing(p.ListingId, p.Destination)

	if err != nil {
		return nil, fmt.Errorf("failed to get execution venue for listing ID %d and destination %s, error: %w",
			p.ListingId, p.Destination, err)
	}

	id, err := ev.client.CreateAndRouteOrder(c, p)
	if err != nil {
		slog.Error("failed to route create order request", "request ", p, "error", err)
		return nil, fmt.Errorf("failed to route order:%w", err)
	}

	slog.Info("routed create order request", "request", p, "executionVenue", ev.podId, "orderId", id)

	return id, nil

}

func (o *orderRouter) getExecutionVenueForOwnerId(ownerId string) (*execVenue, error) {
	o.mux.Lock()
	defer o.mux.Unlock()
	if ev, ok := o.ownerIdToExecVenue[ownerId]; ok {
		return ev, nil

	} else {
		return nil, fmt.Errorf("failed to find execution venue for owner id:%v", ownerId)
	}

}

func (o *orderRouter) getExecutionVenueForListing(listingId int32, destination string) (*execVenue, error) {
	o.mux.Lock()
	defer o.mux.Unlock()
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

	slog.Info("connecting to execution venue service", "targetAddress", targetAddress)

	conn, err := grpc.Dial(targetAddress, grpc.WithInsecure(), grpc.WithBackoffMaxDelay(maxReconnectInterval))

	if err != nil {
		return nil, fmt.Errorf("failed to dial execution venue service: %w", err)
	}

	client := executionvenue.NewExecutionVenueClient(conn)

	return &execVenue{
		client: client,
		conn:   conn,
	}, nil
}
