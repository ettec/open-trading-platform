package main

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
	"time"
)

var zero *model.Decimal64

func init() {
	zero = &model.Decimal64{}
}

type orderManager struct {
	Id                 string
	lastStoredOrder    []byte
	cancelChan         chan bool
	store              func(*model.Order) error
	managedOrder       *executionvenue.ParentOrder
	orderRouter        api.ExecutionVenueClient
	log                *logger.Logger
	errLog             *logger.Logger
}

func (om *orderManager) GetManagedOrderId() string {
	return om.managedOrder.GetId()
}

func (om *orderManager) Cancel() {
	om.cancelChan <- true
}

func (om *orderManager) persistManagedOrderChanges() error {

	orderAsBytes, err := proto.Marshal(&om.managedOrder.Order)

	if bytes.Compare(om.lastStoredOrder, orderAsBytes) != 0 {

		if om.lastStoredOrder != nil {
			om.managedOrder.Version = om.managedOrder.Version + 1
		}

		toStore, err := proto.Marshal(&om.managedOrder.Order)
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

func NewOrderManagerFromParams(id string, params *api.CreateAndRouteOrderParams,
	orderManagerId string, buckets []bucket, doneChan chan<- string, store func(*model.Order) error, orderRouter api.ExecutionVenueClient,
	childOrderStream executionvenue.ChildOrderStream) (*orderManager, error) {

	initialState := model.NewOrder(id, params.OrderSide, params.Quantity, params.Price, params.Listing.Id,
		params.OriginatorId, params.OriginatorRef, params.RootOriginatorId, params.RootOriginatorRef)
	err := initialState.SetTargetStatus(model.OrderStatus_LIVE)
	initialState.ExecParametersJson = params.ExecParametersJson
	if err != nil {
		return nil, err
	}

	om := NewOrderManager(initialState, store, orderManagerId, orderRouter,  buckets, params.Listing, childOrderStream, doneChan)

	return om, nil
}

func NewOrderManager(initialState *model.Order, store func(*model.Order) error, orderManagerId string, orderRouter api.ExecutionVenueClient,
	buckets []bucket, listing *model.Listing, childOrderStream executionvenue.ChildOrderStream,
	doneChan chan<- string) *orderManager {
	po := executionvenue.NewParentOrder(*initialState)

	om := newOrderManager(store, po, orderManagerId, orderRouter)

	go func() {

		om.log.Println("order initialised")

		ticker := time.NewTicker(1 * time.Second)

		for {
			err := om.persistManagedOrderChanges()
			if err != nil {
				errLog.Printf("failed to persist order changes: %v", err )
				om.cancelOrderWithErrorMsg("failed to persist order changes")
			}

			if po.IsTerminalState() {
				childOrderStream.Close()
				doneChan <- po.GetId()
				break
			}

			select {
			case <-ticker.C:
				nowUtc := time.Now().Unix()
				shouldHaveSentQty := &model.Decimal64{}
				for  i :=0; i< len(buckets); i++ {
					if buckets[i].utcStartTimeSecs <= nowUtc {
						shouldHaveSentQty.Add(&buckets[i].quantity)
					}
				}

				sentQty  := &model.Decimal64{}
				sentQty.Add(om.managedOrder.GetTradedQuantity())
				sentQty.Add(om.managedOrder.GetExposedQuantity())

				if sentQty.LessThan(shouldHaveSentQty) {
					shouldHaveSentQty.Sub(sentQty)
					om.sendChildOrder(om.managedOrder.Side,  shouldHaveSentQty, om.managedOrder.Price, listing)
				}

			case <-om.cancelChan:

				if !po.IsTerminalState() {
					om.log.Print("cancelling order")
					err := po.SetTargetStatus(model.OrderStatus_CANCELLED)
					if err != nil {
						errLog.Printf("failed to set target status:%v", err)
					}

					pendingChildOrderCancels := false
					for _, co := range po.ChildOrders {
						if !co.IsTerminalState() {
							pendingChildOrderCancels = true
							_, err := orderRouter.CancelOrder(context.Background(), &api.CancelOrderParams{
								OrderId: co.Id,
								Listing: listing,
							})
							if err != nil {
								errLog.Printf("failed to cancel order:%v", err)
							}

						}

					}

					if !pendingChildOrderCancels {
						err := po.SetStatus(model.OrderStatus_CANCELLED)
						if err != nil {
							errLog.Printf("failed to set status: %v", err)
						}

					}

				}
			case co, ok := <-childOrderStream.GetStream():
				if ok {
					po.OnChildOrderUpdate(co)
				} else {
					om.errLog.Printf("child order update chan unexpectedly closed, cancelling order")
					om.Cancel()
				}

			}

		}

		ticker.Stop()
	}()


	return om
}

func newOrderManager(store func(*model.Order) error, po *executionvenue.ParentOrder, execVenueId string, orderRouter api.ExecutionVenueClient) *orderManager {
	om := &orderManager{
		lastStoredOrder: nil,
		cancelChan:      make(chan bool, 1),
		store:           store,
		managedOrder:    po,
		Id:              execVenueId,
		orderRouter:     orderRouter,
		log:             logger.New(os.Stdout, "order:"+po.Id, logger.Lshortfile|logger.Ltime),
		errLog:          logger.New(os.Stderr, "order:"+po.Id, logger.Lshortfile|logger.Ltime),
	}
	return om
}


func (om *orderManager) sendChildOrder(side model.Side, quantity *model.Decimal64, price *model.Decimal64, listing *model.Listing) {
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

func (om *orderManager) cancelOrderWithErrorMsg(msg string) {
	om.managedOrder.ErrorMessage = msg
	om.cancelChan <- true
}
