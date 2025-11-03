package app

import (
	"context"
	"testing"
	"time"

	"github.com/RicardoCenci/iot-distributed-architecture/client/config"
	"github.com/RicardoCenci/iot-distributed-architecture/client/device"
	"github.com/RicardoCenci/iot-distributed-architecture/client/drivers"
	"github.com/RicardoCenci/iot-distributed-architecture/client/logger"
)

type mockLogger struct {
	logs []string
}

func (m *mockLogger) Debug(msg string, args ...any) { m.logs = append(m.logs, "debug: "+msg) }
func (m *mockLogger) Info(msg string, args ...any)  { m.logs = append(m.logs, "info: "+msg) }
func (m *mockLogger) Warn(msg string, args ...any)  { m.logs = append(m.logs, "warn: "+msg) }
func (m *mockLogger) Error(msg string, args ...any) { m.logs = append(m.logs, "error: "+msg) }

var _ logger.Interface = (*mockLogger)(nil)

func TestNewApp(t *testing.T) {
	cfg := &config.Config{
		Device: config.DeviceConfig{ID: "test-device"},
		MQTT: config.MQTTConfig{
			Broker: "tcp://localhost:1883",
			QoS:    1,
			Topics: map[config.Topic]config.TopicConfig{
				config.TopicDataJSON: {Topic: "test/data"},
				config.TopicMetrics:  {Topic: "test/metrics"},
			},
		},
		Log: config.LogConfig{Level: "info"},
	}
	dev := device.NewDevice("test-device", drivers.NewRandomDataDriver())
	log := &mockLogger{}

	app := NewApp(cfg, dev, log)

	if app.config != cfg {
		t.Error("NewApp() config not set correctly")
	}
	if app.device != dev {
		t.Error("NewApp() device not set correctly")
	}
	if app.logger != log {
		t.Error("NewApp() logger not set correctly")
	}
}

func TestApp_Run_ContextCancel(t *testing.T) {
	cfg := &config.Config{
		Device: config.DeviceConfig{ID: "test-device"},
		MQTT: config.MQTTConfig{
			Broker: "tcp://test.mosquitto.org:1883",
			QoS:    1,
			Topics: map[config.Topic]config.TopicConfig{
				config.TopicDataJSON: {Topic: "test/data"},
				config.TopicMetrics:  {Topic: "test/metrics"},
			},
		},
		Log: config.LogConfig{Level: "info"},
	}
	dev := device.NewDevice("test-device", drivers.NewRandomDataDriver())
	log := &mockLogger{}

	app := NewApp(cfg, dev, log)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	done := make(chan bool)
	go func() {
		app.Run(ctx)
		done <- true
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Error("App.Run() should exit when context is cancelled")
	}
}

func TestApp_Run_MessageEnqueuing(t *testing.T) {
	cfg := &config.Config{
		Device: config.DeviceConfig{ID: "test-device"},
		MQTT: config.MQTTConfig{
			Broker: "tcp://test.mosquitto.org:1883",
			QoS:    1,
			Topics: map[config.Topic]config.TopicConfig{
				config.TopicDataJSON: {Topic: "test/data"},
				config.TopicMetrics:  {Topic: "test/metrics"},
			},
		},
		Log: config.LogConfig{Level: "info"},
	}
	dev := device.NewDevice("test-device", drivers.NewRandomDataDriver())
	log := &mockLogger{}

	app := NewApp(cfg, dev, log)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		time.Sleep(500 * time.Millisecond)
		cancel()
	}()

	done := make(chan bool)
	go func() {
		app.Run(ctx)
		done <- true
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Error("App.Run() should exit when context is cancelled")
	}
}

func TestApp_DataMessage(t *testing.T) {
	msg := DataMessage{
		DeviceID:    "test-device",
		Timestamp:   time.Now(),
		Humidity:    45.5,
		Temperature: 23.7,
	}

	if msg.DeviceID != "test-device" {
		t.Errorf("DataMessage DeviceID = %v, want %v", msg.DeviceID, "test-device")
	}
	if msg.Humidity != 45.5 {
		t.Errorf("DataMessage Humidity = %v, want %v", msg.Humidity, 45.5)
	}
	if msg.Temperature != 23.7 {
		t.Errorf("DataMessage Temperature = %v, want %v", msg.Temperature, 23.7)
	}
}

func TestApp_MetricMessage(t *testing.T) {
	msg := MetricMessage{
		DeviceID:     "test-device",
		Timestamp:    time.Now(),
		CPUUsage:     15.5,
		MemoryUsage:  30.2,
		DiskUsage:    45.8,
		NetworkUsage: 12.3,
	}

	if msg.DeviceID != "test-device" {
		t.Errorf("MetricMessage DeviceID = %v, want %v", msg.DeviceID, "test-device")
	}
	if msg.CPUUsage != 15.5 {
		t.Errorf("MetricMessage CPUUsage = %v, want %v", msg.CPUUsage, 15.5)
	}
	if msg.MemoryUsage != 30.2 {
		t.Errorf("MetricMessage MemoryUsage = %v, want %v", msg.MemoryUsage, 30.2)
	}
	if msg.DiskUsage != 45.8 {
		t.Errorf("MetricMessage DiskUsage = %v, want %v", msg.DiskUsage, 45.8)
	}
	if msg.NetworkUsage != 12.3 {
		t.Errorf("MetricMessage NetworkUsage = %v, want %v", msg.NetworkUsage, 12.3)
	}
}
