package orderstore

import (
	"context"
	"fmt"
	"github.com/ettec/open-trading-platform/go/execution-venue/internal/model"
	"github.com/golang/protobuf/proto"
	"github.com/segmentio/kafka-go"
)

type KafkaStore struct {
	writer *kafka.Writer
}

func NewKafkaStore(topic string, kafkaBrokerUrls []string) *KafkaStore {
	result := KafkaStore{}

	result.writer = kafka.NewWriter(kafka.WriterConfig{
		Brokers:  kafkaBrokerUrls,
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
		Async: true,
	})

	return &result
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
