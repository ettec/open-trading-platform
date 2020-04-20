package internal

import (
	"context"
	"fmt"
	"github.com/ettec/open-trading-platform/go/model"
	"github.com/golang/protobuf/proto"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"time"
)

type ChildOrder struct {
	parentOrderId string
	child *model.Order
}

func New(kafkaBrokerUrls []string, execVenueId string) (<-chan *ChildOrder , error) {

	topic := "orders"

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        kafkaBrokerUrls,
		Topic:          topic,
		ReadBackoffMin: 100 * time.Millisecond,
		ReadBackoffMax: 200 * time.Millisecond,
		MaxWait:        150 * time.Millisecond,
	})
	defer reader.Close()

	result := map[string]*model.Order{}
	start := time.Now()

	readMsgCnt := 0
	for {

		deadline, _ := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))

		msg, err := reader.ReadMessage(deadline)

		if err != nil {
			if err != context.DeadlineExceeded {
				return nil, fmt.Errorf("failed to restore order state from store: %w", err)
			} else {
				break
			}
		}

		if msg.Time.After(start) {
			ks.log.Printf("child order state restored, %v orders for execution venue %v reloaded from %v messages", len(result), ks.execVenueId, readMsgCnt)
			break
		}

		readMsgCnt++
		order := model.Order{}
		err = proto.Unmarshal(msg.Value, &order)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal order whilst recovering state:%w", err)
		}

		if ks.execVenueId == order.GetPlacedWithExecVenueId() {
			result[order.Id] = &order
		}

	}

	ks.log.Printf("order state restored, %v orders for execution venue %v reloaded from %v messages", len(result), ks.execVenueId, readMsgCnt)






	return &result, nil
}

func (ks *ChildOrderStream) recoverInitialState() (map[string]*model.Order, error) {


}

func (ks *ChildOrderStream) Close() {
	ks.writer.Close()
}

