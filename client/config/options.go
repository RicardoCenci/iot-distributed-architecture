package config

import "time"

type DeviceConfig struct {
	ID string `json:"id"`
}

type WiFiConfig struct {
	SSID string `json:"ssid"`
}

type Topic string

const (
	TopicDataJSON Topic = "data_json"
	TopicMetrics  Topic = "metrics"
)

var TOPICS = []Topic{TopicDataJSON, TopicMetrics}

type BackoffConfig struct {
	Base       time.Duration `json:"baseInSeconds"`
	Factor     int           `json:"factor"`
	MaxDelay   time.Duration `json:"maxDelayInSeconds"`
	MaxRetries int           `json:"maxRetries"`
}

type BufferConfig struct {
	Capacity int           `json:"capacity"`
	Backoff  BackoffConfig `json:"backoff"`
}

type TopicConfig struct {
	Topic  string       `json:"topic"`
	Buffer BufferConfig `json:"buffer"`
}

type MQTTConfig struct {
	Broker string                `json:"broker"`
	Topics map[Topic]TopicConfig `json:"topics"`
	QoS    int                   `json:"qos"`
}

type LogConfig struct {
	Level  string          `json:"level"`
	Source LogSourceConfig `json:"source,omitempty"`
}

type LogSourceConfig struct {
	Enabled  bool `json:"enabled,omitempty"`
	Relative bool `json:"relative,omitempty"`
	AsJSON   bool `json:"as_json,omitempty"`
}

type Config struct {
	Log    LogConfig    `json:"log"`
	Device DeviceConfig `json:"device"`
	WiFi   *WiFiConfig  `json:"wifi,omitempty"`
	MQTT   MQTTConfig   `json:"mqtt"`
}

type Option func(*Config)

func WithLog(log LogConfig) Option {
	return func(c *Config) {
		c.Log = log
	}
}

func WithDevice(device DeviceConfig) Option {
	return func(c *Config) {
		c.Device = device
	}
}

func WithWiFi(wifi *WiFiConfig) Option {
	return func(c *Config) {
		c.WiFi = wifi
	}
}

func WithMQTT(mqtt MQTTConfig) Option {
	return func(c *Config) {
		c.MQTT = mqtt
	}
}

func WithTopics(topics map[Topic]TopicConfig) Option {
	return func(c *Config) {
		c.MQTT.Topics = topics
	}
}

func WithQoS(qoS int) Option {
	return func(c *Config) {
		c.MQTT.QoS = qoS
	}
}

func WithBroker(broker string) Option {
	return func(c *Config) {
		c.MQTT.Broker = broker
	}
}

func (c *Config) Merge(options ...Option) *Config {
	for _, option := range options {
		option(c)
	}

	return c
}
