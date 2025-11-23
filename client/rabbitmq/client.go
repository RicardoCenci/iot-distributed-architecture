package rabbitmq

import (
	"fmt"
	"time"

	"github.com/RicardoCenci/iot-distributed-architecture/shared/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Client struct {
	conn   *amqp.Connection
	ch     *amqp.Channel
	logger logger.Interface
}

func NewClient(
	logger logger.Interface,
	url string,
) (*Client, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	logger.Debug("Connected to RabbitMQ", "url", url)

	return &Client{
		conn:   conn,
		ch:     ch,
		logger: logger,
	}, nil
}

func (c *Client) Publish(exchange string, routingKey string, body []byte) error {
	err := c.ch.Publish(
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

func (c *Client) Close() error {
	if c.ch != nil {
		c.ch.Close()
	}
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

