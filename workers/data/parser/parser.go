package parser

import (
	"encoding/base64"
	"fmt"
	"time"

	protosensor "github.com/RicardoCenci/iot-distributed-architecture/shared/proto"
	"github.com/RicardoCenci/iot-distributed-architecture/workers/data/database"
	"google.golang.org/protobuf/proto"
)

func ParseMessage(body []byte) (database.SensorData, error) {
	decoded, err := base64.StdEncoding.DecodeString(string(body))
	if err != nil {
		return database.SensorData{}, fmt.Errorf("failed to decode base64 message: %w", err)
	}

	var sensorData protosensor.SensorData
	if err := proto.Unmarshal(decoded, &sensorData); err != nil {
		return database.SensorData{}, fmt.Errorf("failed to parse protobuf message: %w", err)
	}

	var timestamp time.Time
	if sensorData.Timestamp > 0 {
		timestamp = time.Unix(sensorData.Timestamp, 0)
	} else {
		timestamp = time.Now()
	}

	return database.SensorData{
		DeviceID:    sensorData.SensorId,
		Timestamp:   timestamp,
		Humidity:    sensorData.Humidity,
		Temperature: sensorData.Temperature,
	}, nil
}
