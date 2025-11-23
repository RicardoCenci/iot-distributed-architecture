module github.com/RicardoCenci/iot-distributed-architecture/client

go 1.24.0

require (
	github.com/RicardoCenci/iot-distributed-architecture/shared v0.0.0
	github.com/eclipse/paho.mqtt.golang v1.5.1
	github.com/rabbitmq/amqp091-go v1.10.0
)

require (
	github.com/gorilla/websocket v1.5.3 // indirect
	golang.org/x/net v0.44.0 // indirect
	golang.org/x/sync v0.17.0 // indirect
)

replace github.com/RicardoCenci/iot-distributed-architecture/shared => ../shared
