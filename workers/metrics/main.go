package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/RicardoCenci/iot-distributed-architecture/shared/logger"
	"github.com/RicardoCenci/iot-distributed-architecture/shared/workers/broker/rabbitmq"
	"github.com/RicardoCenci/iot-distributed-architecture/shared/workers/consumer"
	"github.com/RicardoCenci/iot-distributed-architecture/workers/metrics/config"
	"github.com/RicardoCenci/iot-distributed-architecture/workers/metrics/parser"
	"github.com/RicardoCenci/iot-distributed-architecture/workers/metrics/prometheus"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// TODO: test sending metrics to prometheus
	cfg := config.NewConfig()

	loggerConfig := logger.Config{
		Level: cfg.Log.Level,
		Source: logger.SourceConfig{
			Enabled:  cfg.Log.Source.Enabled,
			Relative: cfg.Log.Source.Relative,
			AsJSON:   cfg.Log.Source.AsJSON,
		},
	}

	log := logger.NewSlogLogger(loggerConfig)

	log.Debug("Starting metrics worker", "config", cfg)

	url := fmt.Sprintf(
		"amqp://%s:%s@%s:%s",
		cfg.User,
		cfg.Password,
		cfg.Domain,
		cfg.Port,
	)

	log.Debug("Connecting to RabbitMQ with URL", "url", url)

	rabbitMQ := rabbitmq.NewBroker(
		url,
		log,
	)

	if err := rabbitMQ.Connect(); err != nil {
		log.Error("Failed to connect to RabbitMQ", "error", err)
		os.Exit(1)
	}

	defer rabbitMQ.Close()

	prometheusClient := prometheus.NewClient(log, ":2112")

	if err := prometheusClient.Start(); err != nil {
		log.Error("Failed to start Prometheus client", "error", err)
		os.Exit(1)
	}

	defer prometheusClient.Close()

	log.Info("Prometheus metrics endpoint", "endpoint", prometheusClient.GetMetricsEndpoint())

	metricsConsumer := consumer.NewConsumer(rabbitMQ, log, "metrics-consumer")

	queue := rabbitmq.NewQueue(cfg.QueueName)

	log.Debug("Setting up queue channel", "queue", queue.GetName())

	if err := rabbitMQ.SetupQueueChannel(queue); err != nil {
		log.Error("Failed to setup queue channel", "error", err)
		os.Exit(1)
	}

	log.Info("Starting consumer")

	if err := metricsConsumer.Start(context.Background(), queue, func(delivery amqp.Delivery) error {
		log.Debug("Received message", "message", string(delivery.Body))

		metricData, err := parser.ParseMessage(delivery.Body)
		if err != nil {
			log.Error("Failed to parse message", "error", err, "message", string(delivery.Body))
			return err
		}

		if err := prometheusClient.RecordMetric(metricData); err != nil {
			log.Error("Failed to record metric", "error", err)
			return err
		}

		return nil
	}); err != nil {
		log.Error("Failed to start consumer", "error", err)
		os.Exit(1)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	log.Info("Metrics worker is running. Press Ctrl+C to stop.")
	<-c
	log.Info("Shutting down metrics worker")
}
