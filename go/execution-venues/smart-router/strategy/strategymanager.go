package strategy

import (
	"context"
	"fmt"
	api "github.com/ettec/otp-common/api/executionvenue"
	"github.com/ettec/otp-common/executionvenue"
	"github.com/ettec/otp-common/model"
	"github.com/ettec/otp-common/orderstore"
	"github.com/google/uuid"
	logger "log"
	"os"
	"sync"
)

const ChildUpdatesBufferSize = 1000

var log = logger.New(os.Stdout, "", logger.Ltime|logger.Lshortfile)

type ChildOrderUpdates interface {
	Start()
	NewOrderStream(parentOrderId string, bufferSize int) executionvenue.ChildOrderStream
}



type strategyManager struct {
	id                string
	store             orderstore.OrderStore
	orderRouter       api.ExecutionVenueClient
	doneChan          chan string
	orders            sync.Map
	childOrderUpdates ChildOrderUpdates
	executeFn         func(om *Strategy)
}

func NewStrategyManager(id string, parentOrderStore orderstore.OrderStore, childOrderUpdates ChildOrderUpdates,
	orderRouter api.ExecutionVenueClient, executeFn func(om *Strategy)) *strategyManager {

	sm := &strategyManager{
		id:                id,
		store:             parentOrderStore,
		orderRouter:       orderRouter,
		doneChan:          make(chan string, 100),
		orders:            sync.Map{},
		childOrderUpdates: childOrderUpdates,
	}

	sm.executeFn = executeFn

	go func() {
		id := <-sm.doneChan
		sm.orders.Delete(id)
		log.Printf("order %v done", id)
	}()

	parentOrders, err := sm.store.RecoverInitialCache()
	if err != nil {
		panic(err)
	}

	for _, order := range parentOrders {
		if !order.IsTerminalState() {

			om := NewStrategyFromParentOrder(order, sm.store.Write, sm.id, sm.orderRouter,
				sm.childOrderUpdates.NewOrderStream(order.Id, 1000),
				sm.doneChan)
			sm.orders.Store(om.GetParentOrderId(), om)

			sm.executeFn(om)
		}
	}

	sm.childOrderUpdates.Start()
	return sm
}

func (s *strategyManager) GetExecutionParametersMetaData(ctx context.Context, empty *model.Empty) (*api.ExecParamsMetaDataJson, error) {
	panic("implement me")
}

func (s *strategyManager) CreateAndRouteOrder(ctx context.Context, params *api.CreateAndRouteOrderParams) (*api.OrderId, error) {

	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	om, err := NewStrategyFromCreateParams(id.String(), params, s.id, s.store.Write, s.orderRouter,
		s.childOrderUpdates.NewOrderStream(id.String(), ChildUpdatesBufferSize), s.doneChan)

	if err != nil {
		return nil, err
	}

	s.orders.Store(om.GetParentOrderId(), om)

	s.executeFn(om)

	return &api.OrderId{
		OrderId: om.ExecVenueId,
	}, nil
}

func (s *strategyManager) ModifyOrder(ctx context.Context, params *api.ModifyOrderParams) (*model.Empty, error) {
	return nil, fmt.Errorf("order modification not supported")
}

func (s *strategyManager) CancelOrder(ctx context.Context, params *api.CancelOrderParams) (*model.Empty, error) {

	if val, exists := s.orders.Load(params.OrderId); exists {
		om := val.(Strategy)
		om.Cancel()
		return &model.Empty{}, nil
	} else {
		return nil, fmt.Errorf("no order found for id:%v", params.OrderId)
	}
}
