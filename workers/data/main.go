package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/yourorg/workers-data/broker/rabbitmq"
	"github.com/yourorg/workers-data/consumer"
)

const USER = "data-worker"
const PASSWORD = "RxYy9aiobBDBc0v6W6Ab4RnOuaBgClwx"
const DOMAIN = "localhost"
const PORT = "5672"
const QUEUE_NAME = "data-queue"

func main() {
	// TODO: Substituir logger por um igual o do client
	// TODO: Testar e2e com um producer

	log.Println("Connecting to RabbitMQ...")
	rabbitMQ := rabbitmq.NewBroker(
		fmt.Sprintf(
			"amqp://%s:%s@%s:%s",
			USER,
			PASSWORD,
			DOMAIN,
			PORT,
		),
	)

	if err := rabbitMQ.Connect(); err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	log.Println("Connected to RabbitMQ")

	defer rabbitMQ.Close()

	consumer := consumer.NewDataConsumer(rabbitMQ)

	queue := rabbitmq.NewQueue(QUEUE_NAME)

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
