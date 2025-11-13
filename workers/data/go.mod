module github.com/RicardoCenci/iot-distributed-architecture/workers/data

go 1.24.0

require (
	github.com/RicardoCenci/iot-distributed-architecture/shared v0.0.0
	github.com/rabbitmq/amqp091-go v1.10.0
)

replace github.com/RicardoCenci/iot-distributed-architecture/shared => ../../shared
