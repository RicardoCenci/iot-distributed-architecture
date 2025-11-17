package parser

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/RicardoCenci/iot-distributed-architecture/workers/data/database"
)

type MessageData struct {
	SensorID string `json:"sensor_id"`
	Data     struct {
		Humidity    float32 `json:"humidity"`
		Temperature float32 `json:"temperature"`
	} `json:"data"`
}

func ParseMessage(body []byte) (database.SensorData, error) {
	bodyStr := string(body)

	bodyStr = strings.ReplaceAll(bodyStr, "'", "\"")

	var msgData MessageData
	if err := json.Unmarshal([]byte(bodyStr), &msgData); err != nil {
		return database.SensorData{}, fmt.Errorf("failed to parse message: %w", err)
	}

	return database.SensorData{
		DeviceID:    msgData.SensorID,
		Timestamp:   time.Now(),
		Humidity:    msgData.Data.Humidity,
		Temperature: msgData.Data.Temperature,
	}, nil
}
