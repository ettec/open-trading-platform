package internal

import (
	"context"
	"github.com/ettec/open-trading-platform/go/model"
	"github.com/golang/protobuf/proto"
	"github.com/segmentio/kafka-go"
	logger "log"
	"os"
	"time"
)

type ChildOrder struct {
	ParentOrderId string
	Child         *model.Order
}

type orderReader interface {
	Close() error
	ReadMessage(ctx context.Context) (kafka.Message, error)
}

var errLog = logger.New(os.Stderr, "", logger.Ltime|logger.Lshortfile)

func GetChildOrders(id string, kafkaBrokerUrls []string, bufferSize int) (<-chan ChildOrder, error) {

	topic := "orders"

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        kafkaBrokerUrls,
		Topic:          topic,
		ReadBackoffMin: 100 * time.Millisecond,
		ReadBackoffMax: 200 * time.Millisecond,
		MaxWait:        150 * time.Millisecond,
	})

	isChildOrder := func(order *model.Order) bool {
		return id == order.GetOriginatorId()
	}

	getParentOrderId := func(order *model.Order) string {
		return order.OriginatorRef
	}

	updates := make(chan ChildOrder, bufferSize)

	go func() {
		defer reader.Close()

		for {

			msg, err := reader.ReadMessage(context.Background())

			if err != nil {
				errLog.Printf("exiting read loop as error occurred whilst streaming Child orders:%v", err)
				break
			}

			order := &model.Order{}
			err = proto.Unmarshal(msg.Value, order)
			if err != nil {
				errLog.Printf("exiting read loop, failed to unmarshal order:%v", err)
				break
			}

			if isChildOrder(order) {
				updates <- ChildOrder{
					ParentOrderId: getParentOrderId(order),
					Child:         order,
				}
			}

		}

	}()

	return updates, nil
}

func getChildOrdersFromReader(id string, reader orderReader) (<-chan ChildOrder, error) {
	isChildOrder := func(order *model.Order) bool {
		return id == order.GetOriginatorId()
	}

	getParentOrderId := func(order *model.Order) string {
		return order.OriginatorRef
	}

	updates := make(chan ChildOrder, 1000)

	go func() {
		defer reader.Close()

		for {

			msg, err := reader.ReadMessage(context.Background())

			if err != nil {
				errLog.Printf("exiting read loop as error occurred whilst streaming Child orders:%v", err)
				break
			}

			order := &model.Order{}
			err = proto.Unmarshal(msg.Value, order)
			if err != nil {
				errLog.Printf("exiting read loop, failed to unmarshal order:%v", err)
				break
			}

			if isChildOrder(order) {
				updates <- ChildOrder{
					ParentOrderId: getParentOrderId(order),
					Child:         order,
				}
			}

		}

	}()
	return updates, nil
}
