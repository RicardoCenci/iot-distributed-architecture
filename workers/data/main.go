package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/RicardoCenci/iot-distributed-architecture/shared/logger"
	"github.com/RicardoCenci/iot-distributed-architecture/workers/data/broker/rabbitmq"
	"github.com/RicardoCenci/iot-distributed-architecture/workers/data/config"
	"github.com/RicardoCenci/iot-distributed-architecture/workers/data/consumer"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// TODO: Testar e2e com um producer

	config := config.NewConfig()

	loggerConfig := logger.Config{
		Level: config.Log.Level,
		Source: logger.SourceConfig{
			Enabled:  config.Log.Source.Enabled,
			Relative: config.Log.Source.Relative,
			AsJSON:   config.Log.Source.AsJSON,
		},
	}

	logger := logger.NewSlogLogger(loggerConfig)

	url := fmt.Sprintf(
		"amqp://%s:%s@%s:%s",
		config.User,
		config.Password,
		config.Domain,
		config.Port,
	)

	logger.Debug("Connecting to RabbitMQ with URL", "url", url)

	rabbitMQ := rabbitmq.NewBroker(
		url,
		logger,
	)

	if err := rabbitMQ.Connect(); err != nil {
		logger.Error("Failed to connect to RabbitMQ", "error", err)
		os.Exit(1)
	}

	defer rabbitMQ.Close()

	consumer := consumer.NewDataConsumer(rabbitMQ, logger)

	queue := rabbitmq.NewQueue(config.QueueName)

	logger.Debug("Setting up queue channel", "queue", queue.GetName())

	if err := rabbitMQ.SetupQueueChannel(queue); err != nil {
		logger.Error("Failed to setup queue channel", "error", err)
		os.Exit(1)
	}

	logger.Info("Starting consumer")

	if err := consumer.Start(context.Background(), queue, func(delivery amqp.Delivery) error {
		logger.Info("Processing message", "message", string(delivery.Body))
		return nil
	}); err != nil {
		logger.Error("Failed to start consumer", "error", err)
		os.Exit(1)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	logger.Info("Data worker is running. Press Ctrl+C to stop.")
	<-c
	logger.Info("Shutting down data worker")
}
