package main

import (
	"github.com/ettec/otp-model"
	"testing"
	"time"
)

func Test_sendUpdatesConflationConflatesWhenPauseBeforeReceivingInitialOrder(t *testing.T) {

	before := time.Now()

	out := make(chan orderAndWriteTime)
	sent := make(chan *model.Order, 100)
	send := func(order *model.Order) error {
		sent <- order
		return nil
	}

	go sendUpdates(out, send)

	time.Sleep(maxInitialOrderConflationInterval)

	out <- orderAndWriteTime{order: &model.Order{Id: "1", Version: 0}, writeTime: before}
	out <- orderAndWriteTime{order: &model.Order{Id: "1", Version: 1}, writeTime: before}

	o := <-sent

	if o.Id == "1" && o.Version != 1 {
		t.FailNow()
	}

	out <- orderAndWriteTime{order: &model.Order{Id: "1", Version: 2}, writeTime: before}

	o = <-sent

	if o.Id == "1" && o.Version != 2 {
		t.FailNow()
	}

}

func Test_sendUpdatesConflationWhenNoOrderIsSendThenOneOrderIsSend(t *testing.T) {

	before := time.Now()

	out := make(chan orderAndWriteTime)
	sent := make(chan *model.Order, 100)
	send := func(order *model.Order) error {
		sent <- order
		return nil
	}

	go sendUpdates(out, send)

	time.Sleep(maxInitialOrderConflationInterval)

	out <- orderAndWriteTime{order: &model.Order{Id: "1", Version: 0}, writeTime: before}

	o := <-sent

	if o.Id == "1" && o.Version != 0 {
		t.FailNow()
	}

	out <- orderAndWriteTime{order: &model.Order{Id: "1", Version: 1}, writeTime: before}

	o = <-sent

	if o.Id == "1" && o.Version != 1 {
		t.FailNow()
	}

}

func Test_sendUpdatesConflationWhenNoNewOrderIsSent(t *testing.T) {

	before := time.Now()

	out := make(chan orderAndWriteTime)
	sent := make(chan *model.Order, 100)
	send := func(order *model.Order) error {
		sent <- order
		return nil
	}

	go sendUpdates(out, send)

	out <- orderAndWriteTime{order: &model.Order{Id: "1", Version: 0}, writeTime: before}
	out <- orderAndWriteTime{order: &model.Order{Id: "1", Version: 1}, writeTime: before}

	out <- orderAndWriteTime{order: &model.Order{Id: "2", Version: 0}, writeTime: before}
	out <- orderAndWriteTime{order: &model.Order{Id: "2", Version: 1}, writeTime: before}

	time.Sleep(maxInitialOrderConflationInterval)

	o := <-sent

	if o.Id == "1" && o.Version != 1 {
		t.FailNow()
	}

	if o.Id == "2" && o.Version != 1 {
		t.FailNow()
	}

	o = <-sent

	if o.Id == "1" && o.Version != 1 {
		t.FailNow()
	}

	if o.Id == "2" && o.Version != 1 {
		t.FailNow()
	}

	out <- orderAndWriteTime{order: &model.Order{Id: "2", Version: 2}, writeTime: before}
	out <- orderAndWriteTime{order: &model.Order{Id: "2", Version: 3}, writeTime: before}
	o = <-sent

	if o.Id == "2" && o.Version != 2 {
		t.FailNow()
	}

	o = <-sent

	if o.Id == "2" && o.Version != 3 {
		t.FailNow()
	}

}

func Test_sendUpdatesConflationWhenNewOrderIsSent(t *testing.T) {

	before := time.Now()

	out := make(chan orderAndWriteTime)
	sent := make(chan *model.Order, 100)
	send := func(order *model.Order) error {
		sent <- order
		return nil
	}

	go sendUpdates(out, send)

	out <- orderAndWriteTime{order: &model.Order{Id: "1", Version: 0}, writeTime: before}
	out <- orderAndWriteTime{order: &model.Order{Id: "1", Version: 1}, writeTime: before}

	out <- orderAndWriteTime{order: &model.Order{Id: "2", Version: 0}, writeTime: before}

	after := time.Now()

	out <- orderAndWriteTime{order: &model.Order{Id: "2", Version: 1}, writeTime: after}

	o := <-sent

	if o.Id == "1" && o.Version != 1 {
		t.FailNow()
	}

	if o.Id == "2" && o.Version != 1 {
		t.FailNow()
	}

	o = <-sent

	if o.Id == "1" && o.Version != 1 {
		t.FailNow()
	}

	if o.Id == "2" && o.Version != 1 {
		t.FailNow()
	}

	out <- orderAndWriteTime{order: &model.Order{Id: "2", Version: 2}, writeTime: before}
	out <- orderAndWriteTime{order: &model.Order{Id: "2", Version: 3}, writeTime: before}
	o = <-sent

	if o.Id == "2" && o.Version != 2 {
		t.FailNow()
	}

	o = <-sent

	if o.Id == "2" && o.Version != 3 {
		t.FailNow()
	}

}
