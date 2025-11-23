package rabbitmq

import (
	"context"
	"fmt"
	"time"

	"github.com/RicardoCenci/iot-distributed-architecture/shared/logger"
	"github.com/RicardoCenci/iot-distributed-architecture/shared/workers/broker"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Broker struct {
	conn   *amqp.Connection
	ch     *amqp.Channel
	url    string
	logger logger.Interface
}

func NewBroker(url string, logger logger.Interface) *Broker {
	return &Broker{
		conn:   nil,
		ch:     nil,
		url:    url,
		logger: logger,
	}
}

func (r *Broker) Connect() error {
	maxRetries := 30

	for i := 0; i < maxRetries; i++ {
		conn, err := amqp.Dial(r.url)
		if err == nil {
			r.conn = conn
			r.logger.Info("Connected to RabbitMQ")
			return nil
		}

		if i < maxRetries-1 {
			r.logger.Warn("Failed to connect to RabbitMQ, retrying", "retry_number", i+1, "max_retries", maxRetries)
			time.Sleep(2 * time.Second)
		}
	}

	return fmt.Errorf("failed to connect to RabbitMQ after %d attempts", maxRetries)
}

func (r *Broker) SetupQueueChannel(q *Queue) error {
	ch, err := r.conn.Channel()
	if err != nil {
		return err
	}

	r.ch = ch

	_, err = ch.QueueDeclare(
		q.queueName,
		q.options.durable,
		q.options.deleteWhenUnused,
		q.options.exclusive,
		q.options.noWait,
		q.options.arguments,
	)

	if err != nil {
		ch.Close()
		r.conn.Close()
		return err
	}

	return nil
}

func (r *Broker) ConsumeQueue(ctx context.Context, queue broker.Queue, consumer string) (<-chan amqp.Delivery, error) {
	q, ok := queue.(*Queue)
	if !ok {
		return nil, fmt.Errorf("rabbitmq broker expects *rabbitmq.Queue got %T", queue)
	}

	return r.ch.ConsumeWithContext(
		ctx,
		q.queueName,
		consumer,
		q.options.exclusive,
		q.options.noWait,
		false,
		false,
		nil,
	)
}

func (r *Broker) Close() error {
	if r.ch != nil {
		return r.ch.Close()
	}

	if r.conn != nil {
		return r.conn.Close()
	}

	return nil
}

func (r *Broker) Qos(prefetchCount, prefetchSize int, global bool) error {
	return r.ch.Qos(prefetchCount, prefetchSize, global)
}
