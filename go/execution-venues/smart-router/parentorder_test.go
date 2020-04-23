package main

import (
	"github.com/ettec/open-trading-platform/go/model"
	"testing"
)

func Test_parentOrder_onChildOrderUpdate(t *testing.T) {

	po := newParentOrder(*model.NewOrder("a", model.Side_BUY, IasD(20), IasD(50), 1, "oi", "or"))
	po.SetStatus(model.OrderStatus_LIVE)

	po.onChildOrderUpdate(&model.Order{Id: "a1", TargetStatus: model.OrderStatus_LIVE, Quantity: IasD(15), RemainingQuantity: IasD(15)})

	if !po.ExposedQuantity.Equal(IasD(15)) {
		t.FailNow()
	}

	if !po.GetAvailableQty().Equal(IasD(5)) {
		t.FailNow()
	}

	po.onChildOrderUpdate(&model.Order{Id: "a2", TargetStatus: model.OrderStatus_LIVE, Quantity: IasD(5), RemainingQuantity: IasD(5)})

	if !po.ExposedQuantity.Equal(IasD(20)) {
		t.FailNow()
	}

	if !po.GetAvailableQty().Equal(IasD(0)) {
		t.FailNow()
	}

	po.onChildOrderUpdate(&model.Order{Id: "a1", LastExecPrice: IasD(50), LastExecQuantity: IasD(5), LastExecSeqNo: 1,
		LastExecId: "e1", RemainingQuantity: IasD(10)})

	if !po.TradedQuantity.Equal(IasD(5)) {
		t.FailNow()
	}

	if !po.ExposedQuantity.Equal(IasD(15)) {
		t.FailNow()
	}

	po.onChildOrderUpdate(&model.Order{Id: "a2", TargetStatus: model.OrderStatus_LIVE, Quantity: IasD(5), RemainingQuantity: IasD(0),
		LastExecPrice: IasD(50), LastExecQuantity: IasD(5), LastExecSeqNo: 1,
		LastExecId: "e1"})

	if !po.TradedQuantity.Equal(IasD(10)) {
		t.FailNow()
	}

	if !po.ExposedQuantity.Equal(IasD(10)) {
		t.FailNow()
	}

	if !po.GetAvailableQty().Equal(IasD(0)) {
		t.FailNow()
	}

	po.onChildOrderUpdate(&model.Order{Id: "a1", LastExecPrice: IasD(50), LastExecQuantity: IasD(10), LastExecSeqNo: 2,
		LastExecId: "e2", RemainingQuantity: IasD(0)})

	if !po.TradedQuantity.Equal(IasD(20)) {
		t.FailNow()
	}

	if !po.ExposedQuantity.Equal(IasD(0)) {
		t.FailNow()
	}

	if !po.GetAvailableQty().Equal(IasD(0)) {
		t.FailNow()
	}

	if po.GetStatus() != model.OrderStatus_FILLED {
		t.FailNow()
	}

}
