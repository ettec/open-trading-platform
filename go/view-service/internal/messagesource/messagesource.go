package messagesource

import (
	"context"
	"github.com/segmentio/kafka-go"
)

type Source interface {
	ReadMessage(ctx context.Context) ([]byte, error)
	Close() error
}

type KafkaMessageSource struct {
	reader *kafka.Reader
}

func NewKafkaMessageSource(topic string, brokerUrls []string) *KafkaMessageSource {
	ks := &KafkaMessageSource{}
	ks.reader =  kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokerUrls,
		Topic:   topic,
	})

	return ks
}

func (k *KafkaMessageSource) ReadMessage(ctx context.Context) ([]byte, error) {
	m, err := k.reader.ReadMessage(ctx)
	if err != nil {
		return nil, err
	}

	return m.Value, nil
}

func (k *KafkaMessageSource) Close() error {
	return k.reader.Close()
}




