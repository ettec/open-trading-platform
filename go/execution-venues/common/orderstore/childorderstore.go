package orderstore

import (
	"context"
	"github.com/ettec/open-trading-platform/go/model"
	"github.com/gogo/protobuf/proto"
	"github.com/segmentio/kafka-go"
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

func GetChildOrders(id string, kafkaBrokerUrls []string) (initialState map[string][]*model.Order, updates <-chan ChildOrder, err error) {

	topic := "orders"

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        kafkaBrokerUrls,
		Topic:          topic,
		ReadBackoffMin: 100 * time.Millisecond,
		ReadBackoffMax: 200 * time.Millisecond,
		MaxWait:        150 * time.Millisecond,
	})

	initialState, updates, err = getChildOrdersFromReader(id, reader)
	if err != nil {
		return nil, nil, err
	}

	return initialState, updates, nil
}

func getChildOrdersFromReader(id string, reader orderReader) (map[string][]*model.Order, <-chan ChildOrder, error) {
	isChildOrder := func(order *model.Order) bool {
		return id == order.GetOriginatorId()
	}

	getParentOrderId := func(order *model.Order) string {
		return order.OriginatorRef
	}

	childOrders, err := getInitialState(reader, isChildOrder)

	initialState := map[string][]*model.Order{}
	for _, order := range childOrders {
		initialState[order.OriginatorRef] = append(initialState[order.OriginatorRef], order)
	}

	if err != nil {
		return nil, nil, err
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
	return initialState, updates, nil
}
