package main

import "github.com/ettec/open-trading-platform/go/model"

type execution struct {
	id    string
	price *model.Decimal64
	qty   *model.Decimal64
}

type parentOrder struct {
	model.Order
	childOrders map[string]*model.Order
}

func newParentOrder(order model.Order) *parentOrder {

	return &parentOrder{
		order,
		map[string]*model.Order{},
	}
}

func (po *parentOrder) onChildOrderUpdate(childOrder *model.Order) {

	var lastExecSeqNo int32

	if previous, ok := po.childOrders[childOrder.Id]; ok {
		lastExecSeqNo = previous.LastExecSeqNo
	}

	if childOrder.LastExecSeqNo > lastExecSeqNo {
		execId := childOrder.Id + ":" + childOrder.LastExecId

		execution := execution{
			id:    execId,
			price: childOrder.LastExecPrice,
			qty:   childOrder.LastExecQuantity,
		}

		po.AddExecution(*execution.price, *execution.qty, execution.id)
	}

	po.childOrders[childOrder.Id] = childOrder

	exposedQnt := model.IasD(0)
	for _, order := range po.childOrders {
		if !order.IsTerminalState() {
			exposedQnt.Add(order.RemainingQuantity)
		}

		if !po.ExposedQuantity.Equal(exposedQnt) {
			po.ExposedQuantity = exposedQnt
		}
	}

}
