package parser

import (
	"encoding/json"
	"fmt"
	"strings"
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
}

func ParseMessage(body []byte) (MetricData, error) {
	bodyStr := string(body)

	bodyStr = strings.ReplaceAll(bodyStr, "'", "\"")

	var msgData MessageData
	if err := json.Unmarshal([]byte(bodyStr), &msgData); err != nil {
		return MetricData{}, fmt.Errorf("failed to parse message: %w", err)
	}

	return MetricData{
		DeviceID:     msgData.SensorID,
		CPUUsage:     msgData.Data.CPUUsage,
		MemoryUsage:  msgData.Data.MemoryUsage,
		DiskUsage:    msgData.Data.DiskUsage,
		NetworkUsage: msgData.Data.NetworkUsage,
	}, nil
}
