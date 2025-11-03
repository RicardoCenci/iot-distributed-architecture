package config

import (
	"testing"
)

func TestConfigOptions(t *testing.T) {
	t.Run("WithLog", func(t *testing.T) {
		cfg := NewConfig(WithLog(LogConfig{Level: "debug"}))
		if cfg.Log.Level != "debug" {
			t.Errorf("WithLog() Level = %v, want %v", cfg.Log.Level, "debug")
		}
	})

	t.Run("WithDevice", func(t *testing.T) {
		cfg := NewConfig(WithDevice(DeviceConfig{ID: "test-device"}))
		if cfg.Device.ID != "test-device" {
			t.Errorf("WithDevice() ID = %v, want %v", cfg.Device.ID, "test-device")
		}
	})

	t.Run("WithWiFi", func(t *testing.T) {
		wifi := &WiFiConfig{SSID: "MyWiFi"}
		cfg := NewConfig(WithWiFi(wifi))
		if cfg.WiFi == nil || cfg.WiFi.SSID != "MyWiFi" {
			t.Errorf("WithWiFi() SSID = %v, want %v", cfg.WiFi, wifi)
		}
	})

	t.Run("WithMQTT", func(t *testing.T) {
		mqtt := MQTTConfig{
			Broker: "tcp://localhost:1883",
			QoS:    1,
			Topics: map[Topic]TopicConfig{
				TopicDataJSON: {Topic: "test/topic"},
			},
		}
		cfg := NewConfig(WithMQTT(mqtt))
		if cfg.MQTT.Broker != "tcp://localhost:1883" {
			t.Errorf("WithMQTT() Broker = %v, want %v", cfg.MQTT.Broker, "tcp://localhost:1883")
		}
		if cfg.MQTT.QoS != 1 {
			t.Errorf("WithMQTT() QoS = %v, want %v", cfg.MQTT.QoS, 1)
		}
	})

	t.Run("WithBroker", func(t *testing.T) {
		cfg := NewConfig(WithBroker("tcp://test:1883"))
		if cfg.MQTT.Broker != "tcp://test:1883" {
			t.Errorf("WithBroker() Broker = %v, want %v", cfg.MQTT.Broker, "tcp://test:1883")
		}
	})

	t.Run("WithQoS", func(t *testing.T) {
		cfg := NewConfig(WithQoS(2))
		if cfg.MQTT.QoS != 2 {
			t.Errorf("WithQoS() QoS = %v, want %v", cfg.MQTT.QoS, 2)
		}
	})

	t.Run("WithTopics", func(t *testing.T) {
		topics := map[Topic]TopicConfig{
			TopicDataJSON: {Topic: "test/data"},
			TopicMetrics:  {Topic: "test/metrics"},
		}
		cfg := NewConfig(WithTopics(topics))
		if cfg.MQTT.Topics[TopicDataJSON].Topic != "test/data" {
			t.Errorf("WithTopics() TopicDataJSON = %v, want %v", cfg.MQTT.Topics[TopicDataJSON].Topic, "test/data")
		}
		if cfg.MQTT.Topics[TopicMetrics].Topic != "test/metrics" {
			t.Errorf("WithTopics() TopicMetrics = %v, want %v", cfg.MQTT.Topics[TopicMetrics].Topic, "test/metrics")
		}
	})

	t.Run("Merge multiple options", func(t *testing.T) {
		cfg := NewConfig(
			WithLog(LogConfig{Level: "debug"}),
			WithDevice(DeviceConfig{ID: "test-device"}),
			WithBroker("tcp://localhost:1883"),
			WithQoS(1),
		)
		if cfg.Log.Level != "debug" {
			t.Errorf("Merge() Log.Level = %v, want %v", cfg.Log.Level, "debug")
		}
		if cfg.Device.ID != "test-device" {
			t.Errorf("Merge() Device.ID = %v, want %v", cfg.Device.ID, "test-device")
		}
		if cfg.MQTT.Broker != "tcp://localhost:1883" {
			t.Errorf("Merge() MQTT.Broker = %v, want %v", cfg.MQTT.Broker, "tcp://localhost:1883")
		}
		if cfg.MQTT.QoS != 1 {
			t.Errorf("Merge() MQTT.QoS = %v, want %v", cfg.MQTT.QoS, 1)
		}
	})
}
