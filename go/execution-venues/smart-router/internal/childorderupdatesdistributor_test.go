package internal

import (
	"github.com/ettec/open-trading-platform/go/execution-venues/common/orderstore"
	"github.com/ettec/open-trading-platform/go/model"
	"testing"
)

func Test_childOrderUpdatesDistributor(t *testing.T) {

	updates := make(chan orderstore.ChildOrder, 100)

	d := NewChildOrderUpdatesDistributor(updates)

	astream := d.NewOrderStream("a", 200)
	bstream := d.NewOrderStream("b", 200)

	updates <- orderstore.ChildOrder{
		ParentOrderId: "a",
		Child:         &model.Order{Id: "a1"},
	}

	updates <- orderstore.ChildOrder{
		ParentOrderId: "b",
		Child:         &model.Order{Id: "b1"},
	}

	d.Start()

	u := <-astream.GetStream()
	if u.Id != "a1" {
		t.FailNow()
	}

	u = <-bstream.GetStream()
	if u.Id != "b1" {
		t.FailNow()
	}

	cstream := d.NewOrderStream("c", 200)

	updates <- orderstore.ChildOrder{
		ParentOrderId: "a",
		Child:         &model.Order{Id: "a1"},
	}

	updates <- orderstore.ChildOrder{
		ParentOrderId: "c",
		Child:         &model.Order{Id: "c1"},
	}

	u = <-astream.GetStream()
	if u.Id != "a1" {
		t.FailNow()
	}

	u = <-cstream.GetStream()
	if u.Id != "c1" {
		t.FailNow()
	}

}

func Test_closingChildOrderStream(t *testing.T) {

	updates := make(chan orderstore.ChildOrder, 100)

	d := NewChildOrderUpdatesDistributor(updates)

	astream := d.NewOrderStream("a", 200)
	bstream := d.NewOrderStream("b", 200)

	updates <- orderstore.ChildOrder{
		ParentOrderId: "a",
		Child:         &model.Order{Id: "a1"},
	}

	updates <- orderstore.ChildOrder{
		ParentOrderId: "b",
		Child:         &model.Order{Id: "b1"},
	}

	d.Start()

	u := <-astream.GetStream()
	if u.Id != "a1" {
		t.FailNow()
	}

	u = <-bstream.GetStream()
	if u.Id != "b1" {
		t.FailNow()
	}

	astream.Close()

	updates <- orderstore.ChildOrder{
		ParentOrderId: "a",
		Child:         &model.Order{Id: "a1"},
	}

	updates <- orderstore.ChildOrder{
		ParentOrderId: "b",
		Child:         &model.Order{Id: "b1"},
	}

	_, ok := <-astream.GetStream()
	if ok {
		t.FailNow()
	}

	u = <-bstream.GetStream()
	if u.Id != "b1" {
		t.FailNow()
	}

}

func Test_blockedStreamDoesNotStopOtherStreamEvents(t *testing.T) {

	updates := make(chan orderstore.ChildOrder, 100)

	d := NewChildOrderUpdatesDistributor(updates)

	d.NewOrderStream("a", 1)
	bstream := d.NewOrderStream("b", 1)

	updates <- orderstore.ChildOrder{
		ParentOrderId: "a",
		Child:         &model.Order{Id: "a1"},
	}

	updates <- orderstore.ChildOrder{
		ParentOrderId: "b",
		Child:         &model.Order{Id: "b1"},
	}

	d.Start()

	u := <-bstream.GetStream()
	if u.Id != "b1" {
		t.FailNow()
	}

	updates <- orderstore.ChildOrder{
		ParentOrderId: "a",
		Child:         &model.Order{Id: "a1"},
	}

	updates <- orderstore.ChildOrder{
		ParentOrderId: "b",
		Child:         &model.Order{Id: "b1"},
	}

	u = <-bstream.GetStream()
	if u.Id != "b1" {
		t.FailNow()
	}

}
