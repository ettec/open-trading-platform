package orderstore

import (
	"context"
	"github.com/ettec/open-trading-platform/go/model"
	"github.com/golang/protobuf/proto"
	"github.com/segmentio/kafka-go"
	"testing"
	"time"
)

type testReader struct {
	closeChan chan bool
	msgs      chan kafka.Message
}

type orderAndTime struct {
	order *model.Order
	time  time.Time
}

func newTestReader(orders []orderAndTime) (*testReader, error) {

	t := &testReader{
		closeChan: make(chan bool),
		msgs:      make(chan kafka.Message, 100),
	}

	for _, o := range orders {

		orderBytes, err := proto.Marshal(o.order)
		if err != nil {
			return nil, err
		}

		msg := kafka.Message{
			Key:   []byte(o.order.Id),
			Value: orderBytes,
			Time:  o.time,
		}

		t.msgs <- msg
	}

	return t, nil

}

func (t testReader) Close() error {
	t.closeChan <- true
	return nil
}

func (t testReader) ReadMessage(ctx context.Context) (kafka.Message, error) {
	return <-t.msgs, nil
}

func Test_getChildOrdersFromReader(t *testing.T) {

	now := time.Now()
	afterNow := now.Add(20 * time.Minute)

	id := "testid"

	args := []orderAndTime{
		{&model.Order{

			Id:            "aa",
			Quantity:      model.IasD(1),
			OriginatorId:  id,
			OriginatorRef: "a",
		}, now},
		{&model.Order{

			Id:            "aa",
			Quantity:      model.IasD(2),
			OriginatorId:  id,
			OriginatorRef: "a",
		}, now},
		{&model.Order{

			Id:            "aa",
			Quantity:      model.IasD(5),
			OriginatorId:  "other",
			OriginatorRef: "a",
		}, now},
		{&model.Order{

			Id:            "ab",
			Quantity:      model.IasD(8),
			OriginatorId:  id,
			OriginatorRef: "a",
		}, now},

		{&model.Order{

			Id:            "bb",
			Quantity:      model.IasD(1),
			OriginatorId:  id,
			OriginatorRef: "b",
		}, now},
		{&model.Order{

			Id:            "bb",
			Quantity:      model.IasD(2),
			OriginatorId:  id,
			OriginatorRef: "b",
		}, afterNow},
		{&model.Order{

			Id:            "bc",
			Quantity:      model.IasD(10),
			OriginatorId:  id,
			OriginatorRef: "b",
		}, afterNow},
	}

	tr, err := newTestReader(args)
	if err != nil {
		panic(err)
	}

	initial, updates, err := getChildOrdersFromReader(id, tr)
	if err != nil {
		panic(err)
	}

	if len(initial) != 2 {
		t.FailNow()
	}

	if !initial["a"][0].Quantity.Equal(model.IasD(2)) {
		t.FailNow()
	}

	if !initial["a"][1].Quantity.Equal(model.IasD(8)) {
		t.FailNow()
	}

	if !initial["b"][0].Quantity.Equal(model.IasD(2)) {
		t.FailNow()
	}

	update := <-updates

	if update.parentOrderId != "b" || !update.child.Quantity.Equal(model.IasD(10)) {
		t.FailNow()
	}

}
