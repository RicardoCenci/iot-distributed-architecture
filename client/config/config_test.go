package config

import (
	"os"
	"testing"
	"time"

	"github.com/RicardoCenci/iot-distributed-architecture/shared/logger"
)

func TestNewConfig(t *testing.T) {
	cfg := NewConfig()
	if cfg.Log.Level != "info" {
		t.Errorf("NewConfig() Log.Level = %v, want %v", cfg.Log.Level, "info")
	}
}

func TestConfig_LoadFromTomlFile(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		wantErr  bool
		validate func(*Config) bool
	}{
		{
			name: "valid config",
			content: `[log]
level=debug

[log.source]
enabled=true
relative=false
as_json=true

[device]
id=test-device-123

[mqtt]
broker=tcp://localhost:1883
qos=1

[mqtt.topics.data_json]
topic=iot.device.data.json

[mqtt.topics.metrics]
topic=iot/device/metrics`,
			wantErr: false,
			validate: func(c *Config) bool {
				return c.Log.Level == "debug" &&
					c.Log.Source.Enabled == true &&
					c.Device.ID == "test-device-123" &&
					c.MQTT.Broker == "tcp://localhost:1883" &&
					c.MQTT.QoS == 1 &&
					c.MQTT.Topics[TopicDataJSON].Topic == "iot.device.data.json" &&
					c.MQTT.Topics[TopicMetrics].Topic == "iot/device/metrics"
			},
		},
		{
			name: "config with buffer settings",
			content: `[device]
id=test-device

[mqtt]
broker=tcp://localhost:1883
qos=1

[mqtt.topics.data_json]
topic=iot.device.data.json

[mqtt.topics.data_json.buffer]
capacity=5

[mqtt.topics.data_json.buffer.backoff]
baseInSeconds=3
factor=2
maxDelayInSeconds=15
maxRetries=5

[mqtt.topics.metrics]
topic=iot/device/metrics`,
			wantErr: false,
			validate: func(c *Config) bool {
				return c.MQTT.Topics[TopicDataJSON].Buffer.Capacity == 5 &&
					c.MQTT.Topics[TopicDataJSON].Buffer.Backoff.Base == 3*time.Second &&
					c.MQTT.Topics[TopicDataJSON].Buffer.Backoff.Factor == 2 &&
					c.MQTT.Topics[TopicDataJSON].Buffer.Backoff.MaxDelay == 15*time.Second &&
					c.MQTT.Topics[TopicDataJSON].Buffer.Backoff.MaxRetries == 5
			},
		},
		{
			name: "config with wifi",
			content: `[device]
id=test-device

[mqtt]
broker=tcp://localhost:1883
qos=1

[mqtt.topics.data_json]
topic=iot.device.data.json

[mqtt.topics.metrics]
topic=iot/device/metrics

[wifi]
ssid=MyWiFi`,
			wantErr: false,
			validate: func(c *Config) bool {
				return c.WiFi != nil && c.WiFi.SSID == "MyWiFi"
			},
		},
		{
			name:     "non-existent file",
			content:  "",
			wantErr:  true,
			validate: func(*Config) bool { return true },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fileName string
			if tt.content != "" {
				tmpfile, err := os.CreateTemp("", "test_config_*.toml")
				if err != nil {
					t.Fatal(err)
				}
				fileName = tmpfile.Name()
				if err := os.WriteFile(fileName, []byte(tt.content), 0644); err != nil {
					t.Fatal(err)
				}
				defer os.Remove(fileName)
			} else {
				fileName = "non_existent_file.toml"
			}

			cfg := NewConfig()
			err := cfg.LoadFromTomlFile(fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadFromTomlFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !tt.validate(cfg) {
				t.Error("LoadFromTomlFile() validation failed")
			}
		})
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &Config{
				Device: DeviceConfig{ID: "test-device"},
				MQTT: MQTTConfig{
					Broker: "tcp://localhost:1883",
					QoS:    1,
					Topics: map[Topic]TopicConfig{
						TopicDataJSON: {Topic: "iot.device.data.json"},
						TopicMetrics:  {Topic: "iot/device/metrics"},
					},
				},
				Log: logger.Config{Level: "info"},
			},
			wantErr: false,
		},
		{
			name: "missing device id",
			config: &Config{
				MQTT: MQTTConfig{
					Broker: "tcp://localhost:1883",
					QoS:    1,
					Topics: map[Topic]TopicConfig{
						TopicDataJSON: {Topic: "iot.device.data.json"},
						TopicMetrics:  {Topic: "iot/device/metrics"},
					},
				},
				Log: logger.Config{Level: "info"},
			},
			wantErr: true,
		},
		{
			name: "missing mqtt broker",
			config: &Config{
				Device: DeviceConfig{ID: "test-device"},
				MQTT: MQTTConfig{
					QoS: 1,
					Topics: map[Topic]TopicConfig{
						TopicDataJSON: {Topic: "iot.device.data.json"},
						TopicMetrics:  {Topic: "iot/device/metrics"},
					},
				},
				Log: logger.Config{Level: "info"},
			},
			wantErr: true,
		},
		{
			name: "invalid qos",
			config: &Config{
				Device: DeviceConfig{ID: "test-device"},
				MQTT: MQTTConfig{
					Broker: "tcp://localhost:1883",
					QoS:    3,
					Topics: map[Topic]TopicConfig{
						TopicDataJSON: {Topic: "iot.device.data.json"},
						TopicMetrics:  {Topic: "iot/device/metrics"},
					},
				},
				Log: logger.Config{Level: "info"},
			},
			wantErr: true,
		},
		{
			name: "missing topics",
			config: &Config{
				Device: DeviceConfig{ID: "test-device"},
				MQTT: MQTTConfig{
					Broker: "tcp://localhost:1883",
					QoS:    1,
				},
				Log: logger.Config{Level: "info"},
			},
			wantErr: true,
		},
		{
			name: "missing topic data_json",
			config: &Config{
				Device: DeviceConfig{ID: "test-device"},
				MQTT: MQTTConfig{
					Broker: "tcp://localhost:1883",
					QoS:    1,
					Topics: map[Topic]TopicConfig{
						TopicMetrics: {Topic: "iot/device/metrics"},
					},
				},
				Log: logger.Config{Level: "info"},
			},
			wantErr: true,
		},
		{
			name: "empty topic",
			config: &Config{
				Device: DeviceConfig{ID: "test-device"},
				MQTT: MQTTConfig{
					Broker: "tcp://localhost:1883",
					QoS:    1,
					Topics: map[Topic]TopicConfig{
						TopicDataJSON: {Topic: ""},
						TopicMetrics:  {Topic: "iot/device/metrics"},
					},
				},
				Log: logger.Config{Level: "info"},
			},
			wantErr: true,
		},
		{
			name: "missing wifi ssid",
			config: &Config{
				Device: DeviceConfig{ID: "test-device"},
				MQTT: MQTTConfig{
					Broker: "tcp://localhost:1883",
					QoS:    1,
					Topics: map[Topic]TopicConfig{
						TopicDataJSON: {Topic: "iot.device.data.json"},
						TopicMetrics:  {Topic: "iot/device/metrics"},
					},
				},
				WiFi: &WiFiConfig{},
				Log:  logger.Config{Level: "info"},
			},
			wantErr: true,
		},
		{
			name: "missing log level",
			config: &Config{
				Device: DeviceConfig{ID: "test-device"},
				MQTT: MQTTConfig{
					Broker: "tcp://localhost:1883",
					QoS:    1,
					Topics: map[Topic]TopicConfig{
						TopicDataJSON: {Topic: "iot.device.data.json"},
						TopicMetrics:  {Topic: "iot/device/metrics"},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid log level",
			config: &Config{
				Device: DeviceConfig{ID: "test-device"},
				MQTT: MQTTConfig{
					Broker: "tcp://localhost:1883",
					QoS:    1,
					Topics: map[Topic]TopicConfig{
						TopicDataJSON: {Topic: "iot.device.data.json"},
						TopicMetrics:  {Topic: "iot/device/metrics"},
					},
				},
				Log: logger.Config{Level: "invalid"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
