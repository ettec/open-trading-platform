package internal

import (
	"bytes"
	"context"
	"fmt"
	api "github.com/ettec/otp-common/api/executionvenue"
	"github.com/ettec/otp-common/executionvenue"
	"github.com/ettec/otp-common/model"
	"github.com/golang/protobuf/proto"
	logger "log"
	"os"
)

type commonOrderManager struct {
	lastStoredOrder  []byte
	cancelChan       chan bool
	store            func(*model.Order) error
	managedOrder     *executionvenue.ParentOrder
	Id               string
	orderRouter      api.ExecutionVenueClient
	childOrderStream executionvenue.ChildOrderStream
	doneChan         chan<- string
	log              *logger.Logger
	errLog           *logger.Logger
}

func newCommonOrderManager(store func(*model.Order) error, po *executionvenue.ParentOrder, execVenueId string, orderRouter api.ExecutionVenueClient,
	childOrderStream executionvenue.ChildOrderStream, doneChan chan<- string) *commonOrderManager {
	return &commonOrderManager{
		lastStoredOrder:  nil,
		cancelChan:       make(chan bool, 1),
		store:            store,
		managedOrder:     po,
		Id:               execVenueId,
		orderRouter:      orderRouter,
		childOrderStream: childOrderStream,
		doneChan:         doneChan,
		log:              logger.New(os.Stdout, "order:"+po.Id, logger.Lshortfile|logger.Ltime),
		errLog:           logger.New(os.Stderr, "order:"+po.Id, logger.Lshortfile|logger.Ltime),
	}
}


func (om *commonOrderManager) Cancel() {
	om.cancelChan <- true
}

func (om *commonOrderManager) GetManagedOrderId() string {
	return om.managedOrder.GetId()
}


func (om *commonOrderManager) cancelOrderWithErrorMsg(msg string) {
	om.managedOrder.ErrorMessage = msg
	om.cancelChan <- true
}

func (om *commonOrderManager) cancelOrder( orderRouter api.ExecutionVenueClient, listingSource func( int32) *model.Listing) {
	if !om.managedOrder.IsTerminalState() {
		om.log.Print("cancelling order")
		om.managedOrder.SetTargetStatus(model.OrderStatus_CANCELLED)

		pendingChildOrderCancels := false
		for _, co := range om.managedOrder.ChildOrders {
			if !co.IsTerminalState() {
				pendingChildOrderCancels = true
				orderRouter.CancelOrder(context.Background(), &api.CancelOrderParams{
					OrderId: co.Id,
					Listing: listingSource(co.ListingId),
				})
			}

		}

		if !pendingChildOrderCancels {
			om.managedOrder.SetStatus(model.OrderStatus_CANCELLED)
		}

	}
}

func (om *commonOrderManager) sendChildOrder(side model.Side, quantity *model.Decimal64, price *model.Decimal64, listing *model.Listing) {
	params := &api.CreateAndRouteOrderParams{
		OrderSide:         side,
		Quantity:          quantity,
		Price:             price,
		Listing:           listing,
		OriginatorId:      om.Id,
		OriginatorRef:     om.GetManagedOrderId(),
		RootOriginatorId:  om.managedOrder.RootOriginatorId,
		RootOriginatorRef: om.managedOrder.RootOriginatorRef,
	}

	id, err := om.orderRouter.CreateAndRouteOrder(context.Background(), params)

	if err != nil {
		msg := fmt.Sprintf("failed to submit child order:%v", err)
		om.cancelOrderWithErrorMsg(msg)
		return
	}

	pendingOrder := model.NewOrder(id.OrderId, params.OrderSide, params.Quantity, params.Price, params.Listing.GetId(),
		om.Id, om.GetManagedOrderId(), om.managedOrder.RootOriginatorId, om.managedOrder.RootOriginatorRef)

	// First persisted orders start at version 0, this is a placeholder until the first child order update is received
	pendingOrder.Version = -1

	om.managedOrder.OnChildOrderUpdate(pendingOrder)

}


func (om *commonOrderManager) persistChanges() bool {
	close := false
	om.persistManagedOrderChanges()
	if om.managedOrder.IsTerminalState() {
		close = true
		om.childOrderStream.Close()
		om.doneChan <- om.managedOrder.GetId()
	}
	return close
}

func (om *commonOrderManager) persistManagedOrderChanges() error {

	orderAsBytes, err := proto.Marshal(&om.managedOrder.Order)

	if bytes.Compare(om.lastStoredOrder, orderAsBytes) != 0 {

		if om.lastStoredOrder != nil {
			om.managedOrder.Version = om.managedOrder.Version + 1
		}

		toStore, err := proto.Marshal(&om.managedOrder.Order)

		om.lastStoredOrder = toStore

		orderCopy := &model.Order{}
		proto.Unmarshal(toStore, orderCopy)
		om.store(orderCopy)

		if err != nil {
			return err
		}

	}

	return err
}

func (om *commonOrderManager) onChildOrderUpdate(ok bool,  co *model.Order) {
	if ok {
		om.managedOrder.OnChildOrderUpdate(co)
	} else {
		om.errLog.Printf("child order update chan unexpectedly closed, cancelling order")
		om.Cancel()
	}
}