package model

import (
	"fmt"
)

var noneStateValidTargetStates = map[OrderStatus]bool{OrderStatus_LIVE: true, OrderStatus_CANCELLED: true}
var liveStateValidTargetStates = map[OrderStatus]bool{OrderStatus_LIVE: true, OrderStatus_CANCELLED: true, OrderStatus_NONE: true}
var cancelledStateValidTargetStates = map[OrderStatus]bool{OrderStatus_NONE: true}
var filledStateValidTargetStates = map[OrderStatus]bool{OrderStatus_NONE: true}

var noneTargetStateValidStates = map[OrderStatus]bool{OrderStatus_NONE: true, OrderStatus_FILLED: true}
var liveTargetStateValidStates = map[OrderStatus]bool{OrderStatus_LIVE: true, OrderStatus_FILLED: true}
var cancelledTargetStateValidStates = map[OrderStatus]bool{OrderStatus_CANCELLED: true, OrderStatus_NONE: true, OrderStatus_LIVE: true, OrderStatus_FILLED: true}

func (ord *Order) SetStatus(status OrderStatus) error {

	if ord.Status == OrderStatus_FILLED {
		return nil
	}

	switch ord.TargetStatus {
	case OrderStatus_NONE:
		if !noneTargetStateValidStates[status] {
			return ord.createStatusTransitionError(status)
		}
	case OrderStatus_LIVE:
		if !liveTargetStateValidStates[status] {
			return ord.createStatusTransitionError(status)
		}
	case OrderStatus_CANCELLED:
		if !cancelledTargetStateValidStates[status] {
			return ord.createStatusTransitionError(status)
		}
	}

	ord.Status = status

	if ord.TargetStatus == status {
		ord.TargetStatus = OrderStatus_NONE
	}

	if status == OrderStatus_FILLED {
		ord.TargetStatus = OrderStatus_NONE
	}

	return nil

}

func (ord *Order) SetTargetStatus(targetStatus OrderStatus) error {

	if ord.TargetStatus != OrderStatus_NONE {
		return ord.createTargetStatusTransitionError(targetStatus)
	}

	if targetStatus != OrderStatus_NONE {
		switch ord.Status {
		case OrderStatus_NONE:
			if !noneStateValidTargetStates[targetStatus] {
				return ord.createTargetStatusTransitionError(targetStatus)
			}
		case OrderStatus_LIVE:
			if !liveStateValidTargetStates[targetStatus] {
				return ord.createTargetStatusTransitionError(targetStatus)
			}
		case OrderStatus_CANCELLED:
			if !cancelledStateValidTargetStates[targetStatus] {
				return ord.createTargetStatusTransitionError(targetStatus)
			}
		case OrderStatus_FILLED:
			if !filledStateValidTargetStates[targetStatus] {
				return ord.createTargetStatusTransitionError(targetStatus)
			}
		}
	}

	ord.TargetStatus = targetStatus
	return nil
}

func (ord *Order) createTargetStatusTransitionError(targetStatus OrderStatus) error {
	return fmt.Errorf("requested transition to target status %v is invalid for an order with status %v and target status %v",
		targetStatus.String(), ord.Status.String(), ord.TargetStatus.String())
}

func (ord *Order) createStatusTransitionError(status OrderStatus) error {
	return fmt.Errorf("requested transition to status %v is invalid for an order with status %v and target status %v",
		status.String(), ord.Status.String(), ord.TargetStatus.String())
}
