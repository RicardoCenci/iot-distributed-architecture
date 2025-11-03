package config

import (
	"fmt"
	"time"

	"github.com/RicardoCenci/iot-distributed-architecture/shared/logger"
)

var defaultBackoffConfig = BackoffConfig{
	Base:       2 * time.Second,
	Factor:     2,
	MaxDelay:   10 * time.Second,
	MaxRetries: 3,
}

var defaultBufferConfig = BufferConfig{
	Capacity: 10,
	Backoff:  defaultBackoffConfig,
}

func NewConfig(options ...Option) *Config {
	config := &Config{
		Log: logger.Config{
			Level: "info",
		},
	}

	config.Merge(options...)

	return config
}

func (c *Config) LoadFromTomlFile(filename string) error {
	configMap, err := parseConfig(filename)

	if err != nil {
		return err
	}

	if v := configMap.Get("log"); v != nil {
		if m, ok := v.(map[string]interface{}); ok {
			if s, ok := m["level"].(string); ok {
				c.Log.Level = s
			}

			if v := m["source"]; v != nil {
				if m, ok := v.(map[string]interface{}); ok {
					c.Log.Source.Enabled = m["enabled"].(bool)
					c.Log.Source.Relative = m["relative"].(bool)
					c.Log.Source.AsJSON = m["as_json"].(bool)
				}
			}

		}
	}

	if v := configMap.Get("device"); v != nil {
		if s, ok := v.(map[string]interface{}); ok {
			c.Device.ID = s["id"].(string)
		}
	}

	if v := configMap.Get("wifi"); v != nil {
		if c.WiFi == nil {
			c.WiFi = &WiFiConfig{}
		}

		if s, ok := v.(map[string]interface{}); ok {

			for k, v := range s {
				if k == "ssid" {
					c.WiFi.SSID = v.(string)
				}
			}
		}
	}

	if v := configMap.Get("mqtt.broker"); v != nil {
		if s, ok := v.(string); ok {
			c.MQTT.Broker = s
		}
	}

	if v := configMap.Get("mqtt.qos"); v != nil {
		switch n := v.(type) {
		case int:
			c.MQTT.QoS = n
		case int64:
			c.MQTT.QoS = int(n)
		case float64:
			c.MQTT.QoS = int(n)
		}
	}

	for _, topic := range TOPICS {
		key := fmt.Sprintf("mqtt.topics.%s", string(topic))

		if v := configMap.Get(key); v != nil {
			if s, ok := v.(map[string]interface{}); ok {
				if c.MQTT.Topics == nil {
					c.MQTT.Topics = make(map[Topic]TopicConfig)
				}

				cfg := TopicConfig{
					Topic: s["topic"].(string),
				}

				cfg.Buffer = defaultBufferConfig

				if v := s["buffer"]; v != nil {
					if m, ok := v.(map[string]interface{}); ok {

						if v := m["capacity"]; v != nil {
							if i, ok := v.(int); ok {
								cfg.Buffer.Capacity = i
							}
						}

						backoff := defaultBackoffConfig

						if v := m["backoff"]; v != nil {
							if m, ok := v.(map[string]interface{}); ok {

								if v := m["baseInSeconds"]; v != nil {
									if i, ok := v.(int); ok {
										backoff.Base = time.Duration(i) * time.Second
									}
								}

								if v := m["factor"]; v != nil {
									if i, ok := v.(int); ok {
										backoff.Factor = i
									}
								}

								if v := m["maxDelayInSeconds"]; v != nil {
									if i, ok := v.(int); ok {
										backoff.MaxDelay = time.Duration(i) * time.Second
									}
								}

								if v := m["maxRetries"]; v != nil {
									if i, ok := v.(int); ok {
										backoff.MaxRetries = i
									}
								}

							}
						}

						cfg.Buffer.Backoff = backoff
					}

				}

				c.MQTT.Topics[topic] = cfg

			}
		}
	}

	return nil
}

func (c *Config) Validate() error {
	if c.Device.ID == "" {
		return fmt.Errorf("device id is required")
	}

	if c.MQTT.Broker == "" {
		return fmt.Errorf("mqtt broker is required")
	}

	if c.MQTT.QoS < 0 || c.MQTT.QoS > 2 {
		return fmt.Errorf("mqtt qos must be between 0 and 2")
	}

	if c.MQTT.Topics == nil {
		return fmt.Errorf("mqtt topics are required")
	}

	for _, topic := range TOPICS {
		if v, ok := c.MQTT.Topics[topic]; !ok {
			return fmt.Errorf("mqtt topic %s is required", topic)
		} else if v.Topic == "" {
			return fmt.Errorf("mqtt topic %s is empty", topic)
		}
	}

	if c.WiFi != nil {
		if c.WiFi.SSID == "" {
			return fmt.Errorf("wifi ssid is required")
		}
	}

	if c.Log.Level == "" {
		return fmt.Errorf("log level is required")
	}

	if c.Log.Level != "debug" && c.Log.Level != "info" && c.Log.Level != "warn" && c.Log.Level != "error" {
		return fmt.Errorf("log level must be debug, info, warn or error")
	}

	return nil
}
