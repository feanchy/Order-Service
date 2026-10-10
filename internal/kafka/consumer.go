package kafka

import (
	"context"
	"encoding/json"

	"github.com/feanchy/Order-Service/internal/event"
	"github.com/feanchy/Order-Service/internal/service"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader  *kafka.Reader
	handler service.OrderEventHandler
}

func NewConsumer(
	topic,
	groupID string,
	broker []string,
	handler service.OrderEventHandler,
) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: broker,
		Topic:   topic,
		GroupID: groupID,
	})

	return &Consumer{
		reader:  reader,
		handler: handler,
	}

}

func (c *Consumer) Run(ctx context.Context) error {
	for {
		message, err := c.reader.FetchMessage(ctx)

		if err != nil {
			return err
		}

		var e event.OrderCreated

		err = json.Unmarshal(message.Value, &e)
		if err != nil {
			return err
		}

		err = c.handler.HandleOrderCreated(ctx, e)
		if err != nil {
			return err
		}

		err = c.reader.CommitMessages(ctx, message)
		if err != nil {
			return err
		}

	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
