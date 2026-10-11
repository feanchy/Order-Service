package kafka

import (
	"context"
	"encoding/json"
	"log"
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
	log.Printf("publishing event: event_id=%q order_id=%d",
		event.EventID, event.OrderID)
	data, err := json.Marshal(event)

	log.Printf("kafka payload: %s", data)

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
