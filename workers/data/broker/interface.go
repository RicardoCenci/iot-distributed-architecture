package broker

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageBroker interface {
	Connect() error
	Close() error
	Qos(prefetchCount, prefetchSize int, global bool) error
	ConsumeQueue(ctx context.Context, queue Queue, consumer string) (<-chan amqp.Delivery, error)
}

type Queue interface {
	GetName() string
}
