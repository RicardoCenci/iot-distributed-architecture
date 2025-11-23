package consumer

import (
	"context"
	"fmt"

	"github.com/RicardoCenci/iot-distributed-architecture/shared/logger"
	"github.com/RicardoCenci/iot-distributed-architecture/shared/workers/broker"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	broker       broker.MessageBroker
	logger       logger.Interface
	consumerName string
}

func NewConsumer(broker broker.MessageBroker, logger logger.Interface, consumerName string) *Consumer {
	return &Consumer{
		broker:       broker,
		logger:       logger,
		consumerName: consumerName,
	}
}

func (c *Consumer) Connect() error {
	if err := c.broker.Connect(); err != nil {
		return fmt.Errorf("failed to connect to broker: %w", err)
	}

	return nil
}

func (c *Consumer) Start(
	ctx context.Context, queue broker.Queue,
	handler func(delivery amqp.Delivery) error,
) error {
	if err := c.broker.Qos(1, 0, false); err != nil {
		return fmt.Errorf("failed to set QoS: %w", err)
	}

	msgChan, err := c.broker.ConsumeQueue(
		ctx,
		queue,
		c.consumerName,
	)

	if err != nil {
		return fmt.Errorf("failed to consume queue: %w", err)
	}

	go func() {
		for delivery := range msgChan {
			if err := handler(delivery); err != nil {
				delivery.Nack(false, true)
			} else {
				delivery.Ack(true)
			}
		}
	}()

	return nil
}

func (c *Consumer) Close() error {
	return c.broker.Close()
}
