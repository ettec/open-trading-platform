package internal

import "github.com/ettec/open-trading-platform/go/model"

type parentOrder struct {
	model.Order
	childOrders          map[string]*model.Order
	executions           map[string]*model.Execution
	childOrderRefs       map[string]model.Ref
	childOrdersRecovered bool
}

func newParentOrder(order model.Order) *parentOrder {

	childOrderRefs := map[string]model.Ref{}
	for _, ref := range order.ChildOrdersRefs {
		childOrderRefs[ref.Id] = *ref
	}

	return &parentOrder{
		order,
		map[string]*model.Order{},
		map[string]*model.Execution{},
		childOrderRefs,
		false,
	}
}

func (po *parentOrder) onChildOrderUpdate(childOrder *model.Order) bool {

	po.childOrders[childOrder.Id] = childOrder

	var newExecution *model.Execution

	if childOrder.LastExecId != "" {
		if _, exists := po.executions[childOrder.LastExecId]; !exists {
			newExecution = &model.Execution{
				Id:    childOrder.LastExecId,
				Price: *childOrder.LastExecPrice,
				Qty:   *childOrder.LastExecQuantity,
			}

			po.executions[childOrder.LastExecId] = newExecution
		}
	}

	if !po.childOrdersRecovered {
		po.childOrdersRecovered = true
		for _, persistedRef := range po.ChildOrdersRefs {
			if order, exists := po.childOrders[persistedRef.Id]; !exists || persistedRef.Version > order.Version {
				po.childOrdersRecovered = false
				break
			}
		}
	}

	if ref, exists := po.childOrderRefs[childOrder.Id]; exists {
		if childOrder.Version <= ref.Version {
			return po.childOrdersRecovered
		} else {
			newRef := model.Ref{Id: childOrder.Id, Version: childOrder.Version}
			po.childOrderRefs[childOrder.Id] = newRef
			foundIdx := -1
			for idx, existingRef := range po.ChildOrdersRefs {
				if existingRef.Id == newRef.Id {
					foundIdx = idx
					break
				}
			}

			po.ChildOrdersRefs[foundIdx] = &newRef
		}
	} else {
		newRef := model.Ref{Id: childOrder.Id, Version: childOrder.Version}
		po.childOrderRefs[childOrder.Id] = newRef
		po.ChildOrdersRefs = append(po.ChildOrdersRefs, &newRef)
	}

	if newExecution != nil {
		po.AddExecution(*newExecution)
	}

	exposedQnt := model.IasD(0)
	for _, order := range po.childOrders {
		if !order.IsTerminalState() {
			exposedQnt.Add(order.RemainingQuantity)
		}
	}

	if !po.ExposedQuantity.Equal(exposedQnt) {
		po.ExposedQuantity = exposedQnt
	}

	if po.GetTargetStatus() == model.OrderStatus_CANCELLED {
		if po.GetExposedQuantity().Equal(zero) {
			po.SetTargetStatus(model.OrderStatus_NONE)
			po.SetStatus(model.OrderStatus_CANCELLED)
		}
	}

	return po.childOrdersRecovered

}
