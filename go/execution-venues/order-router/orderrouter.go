package main

import (
	"context"
	api "github.com/ettec/otp-common/api/executionvenue"
	"github.com/ettec/otp-common/loadbalancing"
	"github.com/ettec/otp-common/model"
	"google.golang.org/grpc"
	"time"

	"fmt"
)

type orderRouter struct {
	micToExecVenue     map[string]map[int]*execVenue
	ownerIdToExecVenue map[string]*execVenue
}

func New(connectRetrySecs int) *orderRouter {

	ownerIdToExecVenue := map[string]*execVenue{}
	micToExecVenue := map[string]map[int]*execVenue{}
	micToBalancingPods, err := loadbalancing.GetMicToStatefulPodAddresses("execution-venue")

	if err != nil {
		log.Panicf("failed to get execution venue balancing pods: %v", err)
	}

	for mic, balancingPods := range micToBalancingPods {
		for _, balancingPod := range balancingPods {

			client, err := createExecVenueConnection(time.Duration(connectRetrySecs)*time.Second, balancingPod.TargetAddress)
			if err != nil {
				errLog.Printf("failed to create connection to execution venue service at %v, error: %v", balancingPod.TargetAddress, err)
				continue
			}

			if _, ok := micToExecVenue[mic]; !ok {
				micToExecVenue[mic] = map[int]*execVenue{}
			}

			micToExecVenue[mic][balancingPod.Ordinal] = client

			ownerIdToExecVenue[balancingPod.Name] = client
			log.Printf("added execution venue for mic: %v,  target address: %v, stateful set ordinal %v", mic, balancingPod, balancingPod.Ordinal)
		}
	}

	return &orderRouter{
		micToExecVenue:     micToExecVenue,
		ownerIdToExecVenue: ownerIdToExecVenue,
	}
}

func (o *orderRouter) GetExecutionParametersMetaData(context.Context, *model.Empty) (*api.ExecParamsMetaDataJson, error) {
	return &api.ExecParamsMetaDataJson{}, nil
}

func (o *orderRouter) CreateAndRouteOrder(c context.Context, p *api.CreateAndRouteOrderParams) (*api.OrderId, error) {

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

func (o *orderRouter) getExecutionVenueForListing(listingId int32, destination string) (*execVenue, error) {
	if evs, ok := o.micToExecVenue[destination]; ok {
		numVenues := int32(len(evs))
		ordinal := loadbalancing.GetBalancingOrdinal(listingId, numVenues)
		return evs[ordinal], nil
	} else {
		return nil, fmt.Errorf("no execution venue found for destination %v", destination)
	}
}

func (o *orderRouter) ModifyOrder(c context.Context, p *api.ModifyOrderParams) (*model.Empty, error) {

	if ev, ok := o.ownerIdToExecVenue[p.OwnerId]; ok {
		_, err := ev.client.ModifyOrder(c, p)

		if err != nil {
			return nil, fmt.Errorf("failed to modify order: %v, error: %v", p.OrderId, err)
		}

		return &model.Empty{}, nil

	} else {
		return nil, fmt.Errorf("failed to find execution venue for owner id:%v", p.OwnerId)
	}
}

func (o *orderRouter) CancelOrder(c context.Context, p *api.CancelOrderParams) (*model.Empty, error) {

	if ev, ok := o.ownerIdToExecVenue[p.OwnerId]; ok {
		_, err := ev.client.CancelOrder(c, p)

		if err != nil {
			return nil, fmt.Errorf("failed to cancel order %v: , error: %v", p.OrderId, err)
		}

		return &model.Empty{}, nil

	} else {
		return nil, fmt.Errorf("failed to find execution venue for owner id:%v", p.OwnerId)
	}

}

func createExecVenueConnection(maxReconnectInterval time.Duration, targetAddress string) (cac *execVenue,
	err error) {

	log.Printf("connecting to execution venue service at: %v", targetAddress)

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
