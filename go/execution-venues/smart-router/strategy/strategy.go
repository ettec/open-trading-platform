package strategy

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

type Strategy struct {
	// These two channels must be checked and handled as part of the event loop
	CancelChan           chan string
	ChildOrderUpdateChan <-chan *model.Order

	ExecVenueId string
	ParentOrder *executionvenue.ParentOrder
	Log         *logger.Logger
	ErrLog      *logger.Logger

	lastStoredOrder  []byte
	store            func(*model.Order) error
	orderRouter      api.ExecutionVenueClient
	childOrderStream executionvenue.ChildOrderStream

	doneChan chan<- string
}

func NewStrategyFromCreateParams(parentOrderId string, params *api.CreateAndRouteOrderParams, execVenueId string,
	store func(*model.Order) error, orderRouter api.ExecutionVenueClient,
	childOrderStream executionvenue.ChildOrderStream, doneChan chan<- string) (*Strategy, error) {

	initialState := model.NewOrder(parentOrderId, params.OrderSide, params.Quantity, params.Price, params.Listing.Id,
		params.OriginatorId, params.OriginatorRef, params.RootOriginatorId, params.RootOriginatorRef)

	err := initialState.SetTargetStatus(model.OrderStatus_LIVE)

	if err != nil {
		return nil, err
	}

	om := NewStrategyFromParentOrder(initialState, store, execVenueId, orderRouter, childOrderStream, doneChan)
	return om, nil
}

func NewStrategyFromParentOrder(initialState *model.Order, store func(*model.Order) error, execVenueId string, orderRouter api.ExecutionVenueClient, childOrderStream executionvenue.ChildOrderStream, doneChan chan<- string) *Strategy {
	po := executionvenue.NewParentOrder(*initialState)
	return &Strategy{
		lastStoredOrder:      nil,
		CancelChan:           make(chan string, 1),
		store:                store,
		ParentOrder:          po,
		ExecVenueId:          execVenueId,
		orderRouter:          orderRouter,
		childOrderStream:     childOrderStream,
		ChildOrderUpdateChan: childOrderStream.GetStream(),
		doneChan:             doneChan,
		Log:                  logger.New(os.Stdout, "order:"+po.Id+" ", logger.Lshortfile|logger.Ltime),
		ErrLog:               logger.New(os.Stderr, "order:"+po.Id+" ", logger.Lshortfile|logger.Ltime),
	}
}

func (om *Strategy) Cancel() {
	om.CancelChan <- ""
}

func (om *Strategy) GetParentOrderId() string {
	return om.ParentOrder.GetId()
}

func (om *Strategy) CancelParentOrder(listingSource func(int32) *model.Listing) error {
	if !om.ParentOrder.IsTerminalState() {
		om.Log.Print("cancelling order")
		err := om.ParentOrder.SetTargetStatus(model.OrderStatus_CANCELLED)

		if err != nil {
			return fmt.Errorf("failed to cancel order:%w", err)
		}

		pendingChildOrderCancels := false
		for _, co := range om.ParentOrder.ChildOrders {
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
			err := om.ParentOrder.SetStatus(model.OrderStatus_CANCELLED)
			if err != nil {
				return fmt.Errorf("failed to set status of parent order: %w", err)
			}

		}

	}

	return nil
}

func (om *Strategy) SendChildOrder(side model.Side, quantity *model.Decimal64, price *model.Decimal64, listing *model.Listing) error {

	if quantity.GreaterThan(om.ParentOrder.GetAvailableQty()) {
		return fmt.Errorf("cannot send child order for %v as it exceeds the available quantity on the parent order: %v", quantity,
			om.ParentOrder.GetAvailableQty())
	}

	params := &api.CreateAndRouteOrderParams{
		OrderSide:         side,
		Quantity:          quantity,
		Price:             price,
		Listing:           listing,
		OriginatorId:      om.ExecVenueId,
		OriginatorRef:     om.GetParentOrderId(),
		RootOriginatorId:  om.ParentOrder.RootOriginatorId,
		RootOriginatorRef: om.ParentOrder.RootOriginatorRef,
	}

	id, err := om.orderRouter.CreateAndRouteOrder(context.Background(), params)

	if err != nil {
		return fmt.Errorf("failed to submit child order:%w", err)
	}

	pendingOrder := model.NewOrder(id.OrderId, params.OrderSide, params.Quantity, params.Price, params.Listing.GetId(),
		om.ExecVenueId, om.GetParentOrderId(), om.ParentOrder.RootOriginatorId, om.ParentOrder.RootOriginatorRef)

	// First persisted orders start at version 0, this is a placeholder until the first child order update is received
	pendingOrder.Version = -1

	om.ParentOrder.OnChildOrderUpdate(pendingOrder)

	return nil
}

func (om *Strategy) CheckIfDone() (done bool, err error) {
	done = false
	err = om.persistParentOrderChanges()
	if err != nil {
		return false, fmt.Errorf("failed to persist parent order changes:%w", err)
	}

	if om.ParentOrder.IsTerminalState() {
		done = true
		om.childOrderStream.Close()
		om.doneChan <- om.ParentOrder.GetId()
	}
	return done, nil
}

func (om *Strategy) OnChildOrderUpdate(ok bool, co *model.Order) {
	if ok {
		om.ParentOrder.OnChildOrderUpdate(co)
	} else {
		om.ErrLog.Printf("child order update chan unexpectedly closed, cancelling order")
		om.Cancel()
	}
}

func (om *Strategy) persistParentOrderChanges() error {

	orderAsBytes, err := proto.Marshal(&om.ParentOrder.Order)

	if bytes.Compare(om.lastStoredOrder, orderAsBytes) != 0 {

		if om.lastStoredOrder != nil {
			om.ParentOrder.Version = om.ParentOrder.Version + 1
		}

		toStore, err := proto.Marshal(&om.ParentOrder.Order)
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

func (om *Strategy) CancelOrderWithErrorMsg(msg string) {
	om.CancelChan <- msg
}
