package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer() *Consumer {
}

func (c *Consumer) Run(ctx context.Context) error {
	for {
		msg, err := c.reader.FetchMessage(ctx)
	}
}

func (c *Consumer) Close() error {
	err := c.reader.Close()
	if err != nil {
		return err
	}
	return nil
}
