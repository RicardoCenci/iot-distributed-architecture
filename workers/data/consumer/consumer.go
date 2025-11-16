package consumer

import (
	"context"
	"fmt"

	"github.com/RicardoCenci/iot-distributed-architecture/shared/logger"
	"github.com/RicardoCenci/iot-distributed-architecture/workers/data/broker"
	amqp "github.com/rabbitmq/amqp091-go"
)

type DataConsumer struct {
	broker broker.MessageBroker
	logger logger.Interface
}

func NewDataConsumer(broker broker.MessageBroker, logger logger.Interface) *DataConsumer {
	return &DataConsumer{
		broker: broker,
		logger: logger,
	}
}

func (dp *DataConsumer) Connect() error {
	if err := dp.broker.Connect(); err != nil {
		return fmt.Errorf("failed to connect to broker: %w", err)
	}

	return nil
}

func (dp *DataConsumer) Start(
	ctx context.Context, queue broker.Queue,
	handler func(delivery amqp.Delivery) error,
) error {
	if err := dp.broker.Qos(1, 0, false); err != nil {
		return fmt.Errorf("failed to set QoS: %w", err)
	}

	msgChan, err := dp.broker.ConsumeQueue(
		ctx,
		queue,
		"data-consumer",
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

func (dp *DataConsumer) Close() error {
	return dp.broker.Close()
}
