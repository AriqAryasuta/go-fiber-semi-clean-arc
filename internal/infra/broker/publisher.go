package broker

import (
	"context"

	"github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	channel *amqp091.Channel
}

func NewPublisher(conn *amqp091.Connection) (*Publisher, error) {
	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return &Publisher{channel: channel}, nil
}

func (p *Publisher) Publish(ctx context.Context, exchange string, routingKey string, body []byte) error {
	return p.channel.PublishWithContext(ctx, exchange, routingKey, false, false, amqp091.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
}

func (p *Publisher) Close() error {
	return p.channel.Close()
}
