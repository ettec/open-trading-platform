package main

import (
	"context"
	"github.com/segmentio/kafka-go"
	"time"
)


type KafkaMessageSource struct {
	reader *kafka.Reader
}

func NewKafkaMessageSource(readerConfig kafka.ReaderConfig) *KafkaMessageSource {
	ks := &KafkaMessageSource{}
	ks.reader = kafka.NewReader(readerConfig)

	return ks
}

func (k *KafkaMessageSource) ReadMessage(ctx context.Context) (key []byte, value []byte, writeTime time.Time, err error) {
	m, err := k.reader.ReadMessage(ctx)
	if err != nil {
		return nil, nil, time.Time{}, err
	}

	return m.Key, m.Value, m.Time, nil
}

func (k *KafkaMessageSource) Close() error {
	return k.reader.Close()
}
