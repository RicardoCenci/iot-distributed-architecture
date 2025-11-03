package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/RicardoCenci/iot-distributed-architecture/workers/data/broker/rabbitmq"
	"github.com/RicardoCenci/iot-distributed-architecture/workers/data/config"
	"github.com/RicardoCenci/iot-distributed-architecture/workers/data/consumer"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// TODO: Substituir logger por um igual o do client
	// TODO: Testar e2e com um producer

	config := config.NewConfig()

	log.Println("Connecting to RabbitMQ...")

	rabbitMQ := rabbitmq.NewBroker(
		fmt.Sprintf(
			"amqp://%s:%s@%s:%s",
			config.User,
			config.Password,
			config.Domain,
			config.Port,
		),
	)

	if err := rabbitMQ.Connect(); err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	log.Println("Connected to RabbitMQ")

	defer rabbitMQ.Close()

	consumer := consumer.NewDataConsumer(rabbitMQ)

	queue := rabbitmq.NewQueue(config.QueueName)

	if err := rabbitMQ.SetupQueueChannel(queue); err != nil {
		log.Fatalf("Failed to setup queue channel: %v", err)
	}

	log.Println("Starting consumer...")

	if err := consumer.Start(context.Background(), queue, func(delivery amqp.Delivery) error {
		log.Printf("Processing message: %v", delivery.Body)
		return nil
	}); err != nil {
		log.Fatalf("Failed to start consumer: %v", err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	log.Println("Data worker is running. Press Ctrl+C to stop.")
	<-c
	log.Println("Shutting down data worker...")
}
