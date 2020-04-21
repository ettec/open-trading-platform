package main

import (
	api "github.com/ettec/open-trading-platform/go/common/api/executionvenue"
	"github.com/ettec/open-trading-platform/go/common/marketdata"
	"github.com/ettec/open-trading-platform/go/model"
)


type execution struct {
	id string
	price *model.Decimal64
	qty   *model.Decimal64
}

type managedOrder struct {
	modified bool
	order *model.Order
	childOrders map[string] *model.Order
	executions map[string] *execution
}

func newManagedOrder( orderUpdates []*model.Order) *managedOrder {

	executions := map[string] *execution{}
	var lastExecSeqNo int32


	for _, order := range orderUpdates {
		if order.LastExecSeqNo > lastExecSeqNo {
			executions[order.LastExecId] = &execution{
				id:    order.LastExecId,
				price: order.LastExecPrice,
				qty:   order.LastExecQuantity,
			}
		}
	}

	return &managedOrder{
		modified: false,
		order:    orderUpdates[len(orderUpdates)-1],
		childOrders: map[string] *model.Order{},
		executions: executions,
	}
}

func(m *managedOrder) IsTerminalState() bool {
	return m.order.IsTerminalState()
}

func(m *managedOrder) UpdateChildOrder( order *model.Order ) {

	var lastExecSeqNo int32

	if previous, ok := m.childOrders[order.Id]; ok {
		lastExecSeqNo = previous.LastExecSeqNo
	}

	if order.LastExecSeqNo > lastExecSeqNo {
		execId := order.Id +":" +order.LastExecId

		if  _, ok := m.executions[execId]; !ok {
			execution := execution{
				id:    execId,
				price: order.LastExecPrice,
				qty:   order.LastExecQuantity,
			}

			m.executions[execId] = &execution

			order.AddExecution(*execution.price, *execution.qty, execution.id)
			m.modified = true
		}
	}

	m.childOrders[order.Id] = order

}



func NewOrderManager(originalOrder []*model.Order, terminalState chan<-model.Order, store func(model.Order) error, orderRouter api.ExecutionVenueClient,
	quoteStream marketdata.MdsQuoteStream, childOrders []*model.Order, childOrderUpdates <-chan *model.Order) {

	order := newManagedOrder(originalOrder)


	go func() {

		for _, childOrder := range childOrders {
			order.UpdateChildOrder(childOrder)
			persist(order, store)
		}


		

		for {
			if order.IsTerminalState() {
				terminalState <-*order.order
				break
			}

			persist(order, store)

		}

	}()

}

func persist(order *managedOrder, store func(model.Order) error) {
	if order.modified {
		order.order.Version = order.order.Version + 1
		store(*order.order)
		order.modified = false
	}
}
