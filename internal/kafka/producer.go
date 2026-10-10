package kafka

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/feanchy/Order-Service/internal/event"
	"github.com/segmentio/kafka-go"
)

type Producer interface {
	PublishOrderCreated(ctx context.Context, event event.OrderCreated) error
}

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(brokers []string, topic string) *KafkaProducer {
	return &KafkaProducer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (k *KafkaProducer) PublishOrderCreated(ctx context.Context, event event.OrderCreated) error {
	data, err := json.Marshal(event)

	if err != nil {
		return err

	}

	return k.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(strconv.Itoa(event.OrderID)),
		Value: data,
	})
}

func (k *KafkaProducer) Close() error {
	return k.writer.Close()
}
