package mqtt

import (
	"fmt"

	"github.com/RicardoCenci/iot-distributed-architecture/shared/logger"
	mqttProvider "github.com/eclipse/paho.mqtt.golang"
)

type Client struct {
	client mqttProvider.Client
	logger logger.Interface
}

func NewClient(
	logger logger.Interface,
	broker string,
	clientID string,
	user string,
	password string,
) (*Client, error) {
	options := mqttProvider.NewClientOptions()
	options.AddBroker(broker)
	options.SetClientID(clientID)
	options.SetUsername(user)
	options.SetPassword(password)

	options.OnConnect = func(client mqttProvider.Client) {
		logger.Debug("Connected to MQTT broker")
	}

	options.OnConnectionLost = func(client mqttProvider.Client, err error) {
		logger.Error("Connection Lost", "error", err.Error())
	}

	client := mqttProvider.NewClient(options)

	logger.Debug("Connecting to MQTT broker",
		"broker", broker,
		"clientID", clientID,
	)
	token := client.Connect()

	if token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("failed to connect to MQTT broker: %s", token.Error())
	}

	return &Client{
		client: client,
		logger: logger,
	}, nil
}

func (c *Client) Publish(topic string, payload []byte, qos int, retained bool) error {
	token := c.client.Publish(
		topic,
		byte(qos),
		retained,
		payload,
	)

	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

func (c *Client) Close() error {
	c.client.Disconnect(1000)
	return nil
}
