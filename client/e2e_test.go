package main

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/RicardoCenci/iot-distributed-architecture/client/app"
	"github.com/RicardoCenci/iot-distributed-architecture/client/config"
	"github.com/RicardoCenci/iot-distributed-architecture/client/device"
	"github.com/RicardoCenci/iot-distributed-architecture/client/drivers"
	"github.com/RicardoCenci/iot-distributed-architecture/client/logger"
)

func TestE2E_AppLifecycle(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	cfg := config.NewConfig(
		config.WithDevice(config.DeviceConfig{ID: "e2e-test-device"}),
		config.WithBroker("tcp://test.mosquitto.org:1883"),
		config.WithQoS(1),
		config.WithTopics(map[config.Topic]config.TopicConfig{
			config.TopicDataJSON: {Topic: "iot/e2e/test/data"},
			config.TopicMetrics:  {Topic: "iot/e2e/test/metrics"},
		}),
		config.WithLog(config.LogConfig{Level: "info"}),
	)

	if err := cfg.Validate(); err != nil {
		t.Fatalf("Config validation failed: %v", err)
	}

	log := logger.NewSlogLogger(cfg)
	dev := device.NewDevice(cfg.Device.ID, drivers.NewRandomDataDriver())
	appInstance := app.NewApp(cfg, dev, log)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	done := make(chan bool)
	go func() {
		appInstance.Run(ctx)
		done <- true
	}()

	select {
	case <-done:
		t.Log("E2E test completed successfully")
	case <-time.After(6 * time.Second):
		t.Error("E2E test timed out")
	}
}

func TestE2E_ConfigLoadAndValidate(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "e2e_config_*.toml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	configContent := `[log]
level=info

[device]
id=e2e-test-device-123

[mqtt]
broker=tcp://test.mosquitto.org:1883
qos=1

[mqtt.topics.data_json]
topic=iot/e2e/test/data/json

[mqtt.topics.metrics]
topic=iot/e2e/test/metrics
`

	if err := os.WriteFile(tmpfile.Name(), []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := config.NewConfig()
	if err := cfg.LoadFromTomlFile(tmpfile.Name()); err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if err := cfg.Validate(); err != nil {
		t.Fatalf("Config validation failed: %v", err)
	}

	if cfg.Device.ID != "e2e-test-device-123" {
		t.Errorf("Device ID = %v, want %v", cfg.Device.ID, "e2e-test-device-123")
	}

	if cfg.MQTT.Broker != "tcp://test.mosquitto.org:1883" {
		t.Errorf("MQTT Broker = %v, want %v", cfg.MQTT.Broker, "tcp://test.mosquitto.org:1883")
	}
}

func TestE2E_DeviceDataFlow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	driver := drivers.NewRandomDataDriver()
	dev := device.NewDevice("e2e-device", driver)

	sensorData := dev.GetSensorData()
	if sensorData.Humidity < 40.0 || sensorData.Humidity > 60.0 {
		t.Errorf("Sensor Humidity = %v, want between 40.0 and 60.0", sensorData.Humidity)
	}
	if sensorData.Temperature < 20.0 || sensorData.Temperature > 30.0 {
		t.Errorf("Sensor Temperature = %v, want between 20.0 and 30.0", sensorData.Temperature)
	}

	systemMetrics := dev.GetSystemMetrics()
	if systemMetrics.CPUUsage < 10.0 || systemMetrics.CPUUsage > 30.0 {
		t.Errorf("CPU Usage = %v, want between 10.0 and 30.0", systemMetrics.CPUUsage)
	}
	if systemMetrics.MemoryUsage < 10.0 || systemMetrics.MemoryUsage > 30.0 {
		t.Errorf("Memory Usage = %v, want between 10.0 and 30.0", systemMetrics.MemoryUsage)
	}

	if !dev.IsConnectedToInternet() {
		t.Error("Device should report connected to internet")
	}
}

func TestE2E_FullIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	tmpfile, err := os.CreateTemp("", "e2e_config_*.toml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	configContent := `[log]
level=info

[device]
id=e2e-integration-test

[mqtt]
broker=tcp://test.mosquitto.org:1883
qos=1

[mqtt.topics.data_json]
topic=iot/e2e/integration/data

[mqtt.topics.metrics]
topic=iot/e2e/integration/metrics
`

	if err := os.WriteFile(tmpfile.Name(), []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := config.NewConfig()
	if err := cfg.LoadFromTomlFile(tmpfile.Name()); err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if err := cfg.Validate(); err != nil {
		t.Fatalf("Config validation failed: %v", err)
	}

	log := logger.NewSlogLogger(cfg)
	driver := drivers.NewRandomDataDriver()
	dev := device.NewDevice(cfg.Device.ID, driver)
	appInstance := app.NewApp(cfg, dev, log)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	done := make(chan bool)
	go func() {
		appInstance.Run(ctx)
		done <- true
	}()

	select {
	case <-done:
		t.Log("Full integration test completed successfully")
	case <-time.After(4 * time.Second):
		t.Error("Full integration test timed out")
	}
}
