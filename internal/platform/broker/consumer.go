package broker

import (
	"context"

	"github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	channel *amqp091.Channel
}

func NewConsumer(conn *amqp091.Connection) (*Consumer, error) {
	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return &Consumer{channel: channel}, nil
}

func (c *Consumer) Consume(queue string) (<-chan amqp091.Delivery, error) {
	return c.channel.Consume(
		queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
}

func (c *Consumer) Start(ctx context.Context, queue string, handler func(context.Context, amqp091.Delivery) error) error {
	deliveries, err := c.Consume(queue)
	if err != nil {
		return err
	}
	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-deliveries:
			if !ok {
				return nil
			}
			if err := handler(ctx, msg); err != nil {
				_ = msg.Nack(false, true)
				continue
			}
			_ = msg.Ack(false)
		}
	}
}

func (c *Consumer) Close() error {
	return c.channel.Close()
}
