package ordermanager

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

type OrderManager struct {
	lastStoredOrder      []byte
	CancelChan           chan bool
	store                func(*model.Order) error
	ManagedOrder         *executionvenue.ParentOrder
	Id                   string
	orderRouter          api.ExecutionVenueClient
	childOrderStream     executionvenue.ChildOrderStream
	ChildOrderUpdateChan <-chan *model.Order
	doneChan             chan<- string
	Log                  *logger.Logger
	ErrLog               *logger.Logger
}

func NewOrderManagerFromCreateParams(id string, params *api.CreateAndRouteOrderParams, orderManagerId string,
	store func(*model.Order) error, orderRouter api.ExecutionVenueClient,
	childOrderStream executionvenue.ChildOrderStream, doneChan chan<- string) (*OrderManager, error) {
	initialState := model.NewOrder(id, params.OrderSide, params.Quantity, params.Price, params.Listing.Id,
		params.OriginatorId, params.OriginatorRef, params.RootOriginatorId, params.RootOriginatorRef)
	err := initialState.SetTargetStatus(model.OrderStatus_LIVE)
	if err != nil {
		return nil, err
	}

	om := NewCommonOrderManagerFromState(initialState, store, orderManagerId, orderRouter, childOrderStream, doneChan)
	return om, nil
}

func NewCommonOrderManagerFromState(initialState *model.Order, store func(*model.Order) error, orderManagerId string, orderRouter api.ExecutionVenueClient, childOrderStream executionvenue.ChildOrderStream, doneChan chan<- string) *OrderManager {
	po := executionvenue.NewParentOrder(*initialState)
	return &OrderManager{
		lastStoredOrder:      nil,
		CancelChan:           make(chan bool, 1),
		store:                store,
		ManagedOrder:         po,
		Id:                   orderManagerId,
		orderRouter:          orderRouter,
		childOrderStream:     childOrderStream,
		ChildOrderUpdateChan: childOrderStream.GetStream(),
		doneChan:             doneChan,
		Log:                  logger.New(os.Stdout, "order:"+po.Id, logger.Lshortfile|logger.Ltime),
		ErrLog:               logger.New(os.Stderr, "order:"+po.Id, logger.Lshortfile|logger.Ltime),
	}
}

func (om *OrderManager) Cancel() {
	om.CancelChan <- true
}

func (om *OrderManager) GetManagedOrderId() string {
	return om.ManagedOrder.GetId()
}

func (om *OrderManager) cancelOrderWithErrorMsg(msg string) {
	om.ManagedOrder.ErrorMessage = msg
	om.CancelChan <- true
}

func (om *OrderManager) CancelOrder(listingSource func(int32) *model.Listing) error {
	if !om.ManagedOrder.IsTerminalState() {
		om.Log.Print("cancelling order")
		err := om.ManagedOrder.SetTargetStatus(model.OrderStatus_CANCELLED)

		if err != nil {
			return fmt.Errorf("failed to cancel order:%w", err)
		}

		pendingChildOrderCancels := false
		for _, co := range om.ManagedOrder.ChildOrders {
			if !co.IsTerminalState() {
				pendingChildOrderCancels = true
				_, err := om.orderRouter.CancelOrder(context.Background(), &api.CancelOrderParams{
					OrderId: co.Id,
					Listing: listingSource(co.ListingId),
				})

				if err != nil {
					return fmt.Errorf("failed to cancel child order:%w", err)
				}

			}

		}

		if !pendingChildOrderCancels {
			err := om.ManagedOrder.SetStatus(model.OrderStatus_CANCELLED)
			if err != nil {
				return fmt.Errorf("failed to set status of managed order: %w", err)
			}

		}

	}

	return nil
}

func (om *OrderManager) SendChildOrder(side model.Side, quantity *model.Decimal64, price *model.Decimal64, listing *model.Listing) {
	params := &api.CreateAndRouteOrderParams{
		OrderSide:         side,
		Quantity:          quantity,
		Price:             price,
		Listing:           listing,
		OriginatorId:      om.Id,
		OriginatorRef:     om.GetManagedOrderId(),
		RootOriginatorId:  om.ManagedOrder.RootOriginatorId,
		RootOriginatorRef: om.ManagedOrder.RootOriginatorRef,
	}

	id, err := om.orderRouter.CreateAndRouteOrder(context.Background(), params)

	if err != nil {
		msg := fmt.Sprintf("failed to submit child order:%v", err)
		om.cancelOrderWithErrorMsg(msg)
		return
	}

	pendingOrder := model.NewOrder(id.OrderId, params.OrderSide, params.Quantity, params.Price, params.Listing.GetId(),
		om.Id, om.GetManagedOrderId(), om.ManagedOrder.RootOriginatorId, om.ManagedOrder.RootOriginatorRef)

	// First persisted orders start at version 0, this is a placeholder until the first child order update is received
	pendingOrder.Version = -1

	om.ManagedOrder.OnChildOrderUpdate(pendingOrder)

}

func (om *OrderManager) PersistChanges() (finished bool, err error) {
	finished = false
	err = om.persistManagedOrderChanges()
	if err != nil {
		return false, fmt.Errorf("failed to persist managed order changes:%w",err)
	}

	if om.ManagedOrder.IsTerminalState() {
		finished = true
		om.childOrderStream.Close()
		om.doneChan <- om.ManagedOrder.GetId()
	}
	return finished, nil
}

func (om *OrderManager) persistManagedOrderChanges() error {

	orderAsBytes, err := proto.Marshal(&om.ManagedOrder.Order)

	if bytes.Compare(om.lastStoredOrder, orderAsBytes) != 0 {

		if om.lastStoredOrder != nil {
			om.ManagedOrder.Version = om.ManagedOrder.Version + 1
		}

		toStore, err := proto.Marshal(&om.ManagedOrder.Order)
		if err != nil {
			return err
		}

		om.lastStoredOrder = toStore

		orderCopy := &model.Order{}
		err = proto.Unmarshal(toStore, orderCopy)
		if err != nil {
			return err
		}

		err = om.store(orderCopy)
		if err != nil {
			return err
		}


	}

	return err
}

func (om *OrderManager) OnChildOrderUpdate(ok bool, co *model.Order) {
	if ok {
		om.ManagedOrder.OnChildOrderUpdate(co)
	} else {
		om.ErrLog.Printf("child order update chan unexpectedly closed, cancelling order")
		om.Cancel()
	}
}
