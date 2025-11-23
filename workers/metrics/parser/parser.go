package parser

import (
	"encoding/base64"
	"fmt"
	"time"

	protosensor "github.com/RicardoCenci/iot-distributed-architecture/shared/proto"
	"google.golang.org/protobuf/proto"
)

type MessageData struct {
	SensorID string `json:"sensor_id"`
	Data     struct {
		CPUUsage     float32 `json:"cpu_usage"`
		MemoryUsage  float32 `json:"memory_usage"`
		DiskUsage    float32 `json:"disk_usage"`
		NetworkUsage float32 `json:"network_usage"`
	} `json:"data"`
}

type MetricData struct {
	DeviceID     string
	CPUUsage     float32
	MemoryUsage  float32
	DiskUsage    float32
	NetworkUsage float32
	Timestamp    time.Time
}

func ParseMessage(body []byte) (MetricData, error) {
	decoded, err := base64.StdEncoding.DecodeString(string(body))
	if err != nil {
		return MetricData{}, fmt.Errorf("failed to decode base64 message: %w", err)
	}

	var metricsData protosensor.MetricsData

	if err := proto.Unmarshal(decoded, &metricsData); err != nil {
		return MetricData{}, fmt.Errorf("failed to parse protobuf message: %w", err)
	}

	var timestamp time.Time
	if metricsData.Timestamp > 0 {
		timestamp = time.Unix(metricsData.Timestamp, 0)
	} else {
		timestamp = time.Now()
	}

	return MetricData{
		DeviceID:     metricsData.SensorId,
		CPUUsage:     metricsData.CpuUsage,
		MemoryUsage:  metricsData.MemoryUsage,
		DiskUsage:    metricsData.DiskUsage,
		NetworkUsage: metricsData.NetworkUsage,
		Timestamp:    timestamp,
	}, nil
}
