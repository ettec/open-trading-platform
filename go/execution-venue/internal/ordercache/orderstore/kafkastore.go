package orderstore

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

type KafkaStore struct {
	writer          *kafka.Writer
	log             *log.Logger
	topic           string
	kafkaBrokerUrls []string
	execVenueId     string
}

func NewKafkaStore(topic string, kafkaBrokerUrls []string, execVenueId string) (*KafkaStore, error) {



	result := KafkaStore{
		log:             log.New(os.Stdout, "Topic: "+topic+" ", log.Lshortfile|log.Ltime),
		topic:           topic,
		kafkaBrokerUrls: kafkaBrokerUrls,
		execVenueId:     execVenueId,
	}

	result.writer = kafka.NewWriter(kafka.WriterConfig{
		Brokers:      kafkaBrokerUrls,
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		Async:        true,
		BatchTimeout: 10 * time.Millisecond,
	})

	return &result, nil
}

func (ks *KafkaStore) RecoverInitialCache() (map[string]*model.Order, error) {

	ks.log.Println("restoring order state from topic:")
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        ks.kafkaBrokerUrls,
		Topic:          ks.topic,
		ReadBackoffMin: 100 * time.Millisecond,
		ReadBackoffMax: 200 * time.Millisecond,
		MaxWait:        150 * time.Millisecond,
	})
	defer reader.Close()

	result := map[string]*model.Order{}
	now := time.Now()

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

		if msg.Time.After(now) {
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

	return result, nil
}

func (ks *KafkaStore) Close() {
	ks.writer.Close()
}

func (ks *KafkaStore) Write(order *model.Order) error {

	orderBytes, err := proto.Marshal(order)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Key:   []byte(order.Id),
		Value: orderBytes,
	}

	err = ks.writer.WriteMessages(context.Background(), msg)

	if err != nil {
		return fmt.Errorf("failed to write order to kafka order store: %w", err)
	}

	return nil
}
