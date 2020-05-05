package model

import (
	"fmt"
	"github.com/shopspring/decimal"
	"time"
)

var zero decimal.Decimal

func init() {
	zero = decimal.New(0, 0)
}

var noneStateValidTargetStates = map[OrderStatus]bool{OrderStatus_LIVE: true, OrderStatus_CANCELLED: true}
var liveStateValidTargetStates = map[OrderStatus]bool{OrderStatus_LIVE: true, OrderStatus_CANCELLED: true, OrderStatus_NONE: true}
var cancelledStateValidTargetStates = map[OrderStatus]bool{OrderStatus_NONE: true}
var filledStateValidTargetStates = map[OrderStatus]bool{OrderStatus_NONE: true}

var noneTargetStateValidStates = map[OrderStatus]bool{OrderStatus_NONE: true, OrderStatus_FILLED: true}
var liveTargetStateValidStates = map[OrderStatus]bool{OrderStatus_LIVE: true, OrderStatus_FILLED: true}
var cancelledTargetStateValidStates = map[OrderStatus]bool{OrderStatus_CANCELLED: true, OrderStatus_NONE: true, OrderStatus_LIVE: true, OrderStatus_FILLED: true}

func (o *Order) SetStatus(status OrderStatus) error {

	if o.Status == OrderStatus_FILLED {
		return nil
	}

	switch o.TargetStatus {
	case OrderStatus_NONE:
		if !noneTargetStateValidStates[status] {
			return o.createStatusTransitionError(status)
		}
	case OrderStatus_LIVE:
		if !liveTargetStateValidStates[status] {
			return o.createStatusTransitionError(status)
		}
	case OrderStatus_CANCELLED:
		if !cancelledTargetStateValidStates[status] {
			return o.createStatusTransitionError(status)
		}
	}

	o.Status = status

	if o.TargetStatus == status {
		o.TargetStatus = OrderStatus_NONE
	}

	if status == OrderStatus_FILLED {
		o.TargetStatus = OrderStatus_NONE
	}

	return nil

}

func (o *Order) SetTargetStatus(targetStatus OrderStatus) error {

	if o.TargetStatus != OrderStatus_NONE {
		return o.createTargetStatusTransitionError(targetStatus)
	}

	if targetStatus != OrderStatus_NONE {
		switch o.Status {
		case OrderStatus_NONE:
			if !noneStateValidTargetStates[targetStatus] {
				return o.createTargetStatusTransitionError(targetStatus)
			}
		case OrderStatus_LIVE:
			if !liveStateValidTargetStates[targetStatus] {
				return o.createTargetStatusTransitionError(targetStatus)
			}
		case OrderStatus_CANCELLED:
			if !cancelledStateValidTargetStates[targetStatus] {
				return o.createTargetStatusTransitionError(targetStatus)
			}
		case OrderStatus_FILLED:
			if !filledStateValidTargetStates[targetStatus] {
				return o.createTargetStatusTransitionError(targetStatus)
			}
		}
	}

	o.TargetStatus = targetStatus
	return nil
}

func (o *Order) IsTerminalState() bool {
	return o.Status == OrderStatus_FILLED || o.Status == OrderStatus_CANCELLED
}

func (o *Order) GetAvailableQty() *Decimal64 {
	quantity := *o.Quantity
	quantity.Sub(o.GetExposedQuantity())
	quantity.Sub(o.GetTradedQuantity())
	return &quantity
}

func (o *Order) createTargetStatusTransitionError(targetStatus OrderStatus) error {
	return fmt.Errorf("requested transition to target status %v is invalid for an order with status %v and target status %v",
		targetStatus.String(), o.Status.String(), o.TargetStatus.String())
}

func (o *Order) createStatusTransitionError(status OrderStatus) error {
	return fmt.Errorf("requested transition to status %v is invalid for an order with status %v and target status %v",
		status.String(), o.Status.String(), o.TargetStatus.String())
}

func NewOrder(id string, OrderSide Side, Quantity *Decimal64, Price *Decimal64, listingId int32,
	originatorId string, originatorRef string, rootOriginatorId string, rootOriginatorRef string) *Order {

	now := time.Now()

	return &Order{
		Id:                id,
		Side:              OrderSide,
		Quantity:          Quantity,
		Price:             Price,
		ListingId:         listingId,
		RemainingQuantity: Quantity,
		Status:            OrderStatus_NONE,
		TargetStatus:      OrderStatus_NONE,
		Created: &Timestamp{
			Seconds:     now.Unix(),
			Nanoseconds: int32(now.Nanosecond()),
		},
		OriginatorId:  originatorId,
		OriginatorRef: originatorRef,
		RootOriginatorId: rootOriginatorId,
		RootOriginatorRef: rootOriginatorRef,
	}

}

type Execution struct {
	Id    string
	Price Decimal64
	Qty   Decimal64
}

func (o *Order) AddExecution(execution Execution) error {

	if o.TargetStatus == OrderStatus_LIVE {
		err := o.SetStatus(OrderStatus_LIVE)
		if err != nil {
			return err
		}
	}

	o.LastExecId = execution.Id
	o.LastExecPrice = &execution.Price
	o.LastExecQuantity = &execution.Qty

	o.updateAveragePrice(execution.Price, execution.Qty)

	o.RemainingQuantity = ToDecimal64(o.RemainingQuantity.AsDecimal().Sub(execution.Qty.AsDecimal()))
	o.TradedQuantity = ToDecimal64(o.TradedQuantity.AsDecimal().Add(execution.Qty.AsDecimal()))

	if o.RemainingQuantity.AsDecimal().LessThanOrEqual(zero) {
		err := o.SetStatus(OrderStatus_FILLED)
		if err != nil {
			return err
		}
	}

	return nil
}

func (o *Order) updateAveragePrice(lastPx Decimal64, lastQty Decimal64) {
	totalTradeValue := o.AvgTradePrice.AsDecimal().Mul(o.TradedQuantity.AsDecimal()).Add(lastPx.AsDecimal().Mul(lastQty.AsDecimal()))
	totalTradedQnt := o.TradedQuantity.AsDecimal().Add(lastQty.AsDecimal())
	o.AvgTradePrice = ToDecimal64(totalTradeValue.Div(totalTradedQnt))
}
