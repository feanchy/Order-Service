package kafka

import (
	"context"

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

}
