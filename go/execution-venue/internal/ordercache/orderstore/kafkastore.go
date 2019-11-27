package orderstore

import (
	"context"
	"fmt"
	"github.com/coronationstreet/open-trading-platform/execution-venue/pb"
	"github.com/segmentio/kafka-go"
	"os"
)

type KafkaStore struct {

}

func NewKafkaStore(topic string, kafkaUrl string, partition int) (*KafkaStore, error) {
	result := KafkaStore{}

	conn, _ := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)



	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666 )
	if err != nil {
		return nil, fmt.Errorf("Unable to create file store: %w", err)
	}

	result.file = file

	return &result, nil
}

func (ks *KafkaStore) Close() {
	.file.Close()
}

func (fs *KafkaStore) Write(order *pb.Order) error {
	bytes, err := proto.Marshal(order)
	if err != nil {
		return fmt.Errorf("unable to convert order %v to bytes: %w", order, err)
	}
	_, err = fs.file.Write(bytes)
	if err != nil {
		return fmt.Errorf("unable to write order %v bytes: %w", order, err)
	}
	return nil
}
