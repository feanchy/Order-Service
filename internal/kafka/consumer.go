package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"log"

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
		Brokers:     broker,
		Topic:       topic,
		GroupID:     groupID,
		StartOffset: kafka.FirstOffset,
	})

	return &Consumer{
		reader:  reader,
		handler: handler,
	}

}

func (c *Consumer) Run(ctx context.Context) error {
	for {
		log.Println("run consumer")
		message, err := c.reader.FetchMessage(ctx)

		if err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}
			return err
		}

		var e event.OrderCreated

		log.Printf(
			"kafka message: topic=%s partition=%d offset=%d value=%s",
			message.Topic,
			message.Partition,
			message.Offset,
			string(message.Value),
		)

		if len(message.Value) == 0 {
			return errors.New("empty kafka message")
		}

		err = json.Unmarshal(message.Value, &e)
		if err != nil {
			return err
		}

		log.Printf("decoded event: event_id=%q order_id=%d offset=%d",
			e.EventID, e.OrderID, message.Offset)

		if e.OrderID <= 0 {
			return errors.New("invalid order id")
		}

		if e.EventID == "" {
			return errors.New("empty event id")
		}

		log.Println(e)

		err = c.handler.HandleOrderCreated(ctx, e)
		if err != nil {
			return err
		}

		log.Println(e)

		err = c.reader.CommitMessages(ctx, message)
		if err != nil {
			return err
		}

	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
