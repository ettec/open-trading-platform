package main

import (
	"context"
	"github.com/ettec/otp-common/model"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestConflatesWhenPauseBeforeReceivingInitialOrder(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	before := time.Now()

	out := make(chan orderAndWriteTime)
	sent := make(chan *model.Order, 100)
	send := func(order *model.Order) error {
		sent <- order
		return nil
	}

	go func() {
		err := sendOrderUpdates(ctx, out, send)
		assert.NoError(t, err)
	}()

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

func TestConflationWhenNoOrderIsSentThenOneOrderIsSent(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	before := time.Now()

	out := make(chan orderAndWriteTime)
	sent := make(chan *model.Order, 100)
	send := func(order *model.Order) error {
		sent <- order
		return nil
	}

	go func() {
		err := sendOrderUpdates(ctx, out, send)
		assert.NoError(t, err)
	}()

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

func TestConflationWhenNoNewOrderIsSent(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	before := time.Now()

	out := make(chan orderAndWriteTime)
	sent := make(chan *model.Order, 100)
	send := func(order *model.Order) error {
		sent <- order
		return nil
	}

	go func() {
		err := sendOrderUpdates(ctx, out, send)
		assert.NoError(t, err)
	}()

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

func TestConflationWhenNewOrderIsSent(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	before := time.Now()

	out := make(chan orderAndWriteTime)
	sent := make(chan *model.Order, 100)
	send := func(order *model.Order) error {
		sent <- order
		return nil
	}

	go func() {
		err := sendOrderUpdates(ctx, out, send)
		assert.NoError(t, err)
	}()

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
