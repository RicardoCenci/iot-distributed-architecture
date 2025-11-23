module github.com/RicardoCenci/iot-distributed-architecture/workers/data

go 1.24.0

require (
	github.com/RicardoCenci/iot-distributed-architecture/shared v0.0.0
	github.com/golang-migrate/migrate/v4 v4.19.0
	github.com/lib/pq v1.10.9
	github.com/rabbitmq/amqp091-go v1.10.0
	google.golang.org/protobuf v1.36.10
)

require (
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
)

replace github.com/RicardoCenci/iot-distributed-architecture/shared => ../../shared
