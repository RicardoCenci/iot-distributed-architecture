package app

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/RicardoCenci/iot-distributed-architecture/client/config"
	"github.com/RicardoCenci/iot-distributed-architecture/client/device"
	"github.com/RicardoCenci/iot-distributed-architecture/client/logger"
	"github.com/RicardoCenci/iot-distributed-architecture/client/mqtt"
	"github.com/RicardoCenci/iot-distributed-architecture/client/queue"
)

type App struct {
	device *device.Device
	logger logger.Interface
	config *config.Config
}

type MetricMessage struct {
	DeviceID     string
	Timestamp    time.Time
	CPUUsage     float32
	MemoryUsage  float32
	DiskUsage    float32
	NetworkUsage float32
}

type DataMessage struct {
	DeviceID    string
	Timestamp   time.Time
	Humidity    float32
	Temperature float32
}

func NewApp(config *config.Config, device *device.Device, logger logger.Interface) *App {
	return &App{
		config: config,
		device: device,
		logger: logger,
	}
}

func (a *App) Run(ctx context.Context) {

	a.logger.Info("Starting application")

	client, err := mqtt.NewClient(a.logger, a.config.MQTT.Broker, a.config.Device.ID)

	if err != nil {
		a.logger.Error("Failed to create MQTT client", "error", err)
		return
	}

	dataMetrics := mqtt.NewMetrics(a.config.MQTT.Topics[config.TopicDataJSON].Topic)

	dataPublisher := mqtt.BufferedPublisher[DataMessage]{
		Logger:  a.logger,
		Client:  client,
		Metrics: dataMetrics,
		Queue: queue.New(
			queue.WithCapacity[DataMessage](5),
			queue.WithBackoff[DataMessage](queue.BackoffConfig{
				Base:       2 * time.Second,
				Factor:     2,
				MaxDelay:   time.Second * 10,
				MaxRetries: 3,
			}),
		),
		MessageTransformer: func(msg DataMessage) string {
			return fmt.Sprintf("{'sensor_id': '%s', 'data': {'humidity': %f, 'temperature': %f}}", msg.DeviceID, msg.Humidity, msg.Temperature)
		},
		QoS:   a.config.MQTT.QoS,
		Topic: a.config.MQTT.Topics[config.TopicDataJSON].Topic,
	}

	metricMetrics := mqtt.NewMetrics(a.config.MQTT.Topics[config.TopicMetrics].Topic)

	metricPublisher := mqtt.BufferedPublisher[MetricMessage]{
		Logger:  a.logger,
		Client:  client,
		Metrics: metricMetrics,
		Queue: queue.New(
			queue.WithCapacity[MetricMessage](5),
			queue.WithBackoff[MetricMessage](queue.BackoffConfig{
				Base:       2 * time.Second,
				Factor:     2,
				MaxDelay:   time.Second * 10,
				MaxRetries: 3,
			}),
		),
		MessageTransformer: func(msg MetricMessage) string {
			return fmt.Sprintf("{'sensor_id': '%s', 'data': {'cpu_usage': %f, 'memory_usage': %f, 'disk_usage': %f, 'network_usage': %f}}", msg.DeviceID, msg.CPUUsage, msg.MemoryUsage, msg.DiskUsage, msg.NetworkUsage)
		},
		QoS:   a.config.MQTT.QoS,
		Topic: a.config.MQTT.Topics[config.TopicMetrics].Topic,
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		dataPublisher.Run(ctx)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		metricPublisher.Run(ctx)
	}()

	ticker := time.NewTicker(time.Second)
	logPublishMetricsTick := time.NewTicker(time.Second * 5)

	defer ticker.Stop()
	defer logPublishMetricsTick.Stop()

	for {
		select {
		case <-ctx.Done():
			a.logger.Debug("Flushing buffers", "data_queue_len", dataPublisher.Queue.Len(), "metric_queue_len", metricPublisher.Queue.Len())

			dataPublisher.Queue.Close()
			metricPublisher.Queue.Close()

			wg.Wait()

			dataMetrics.Print(a.logger)
			metricMetrics.Print(a.logger)

			a.logger.Debug("Closing MQTT client")

			client.Close()
			return
		case <-ticker.C:
			timestamp := time.Now()

			sensorData := a.device.GetSensorData()

			dataPublisher.Queue.Enqueue(queue.Message[DataMessage]{
				Data: DataMessage{
					DeviceID:    a.device.DeviceID,
					Timestamp:   timestamp,
					Humidity:    sensorData.Humidity,
					Temperature: sensorData.Temperature,
				},
			})

			metricData := a.device.GetSystemMetrics()

			metricPublisher.Queue.Enqueue(queue.Message[MetricMessage]{
				Data: MetricMessage{
					DeviceID:     a.device.DeviceID,
					Timestamp:    timestamp,
					CPUUsage:     metricData.CPUUsage,
					MemoryUsage:  metricData.MemoryUsage,
					DiskUsage:    metricData.DiskUsage,
					NetworkUsage: metricData.NetworkUsage,
				},
			})
		case <-logPublishMetricsTick.C:
			dataMetrics.Print(a.logger)
			metricMetrics.Print(a.logger)
		}
	}
}
