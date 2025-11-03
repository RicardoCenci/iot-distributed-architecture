package drivers

import "math/rand/v2"

type RandomDataDriver struct{}

func NewRandomDataDriver() DriverInterface {
	return &RandomDataDriver{}
}

func (m *RandomDataDriver) ProbeSensor() SensorData {
	return SensorData{
		Humidity:    40.0 + (rand.Float32() * 20.0),
		Temperature: 20.0 + (rand.Float32() * 10.0),
	}
}

func (m *RandomDataDriver) ProbeSystemMetrics() SystemMetrics {
	return SystemMetrics{
		CPUUsage:     10.0 + (rand.Float32() * 20.0),
		MemoryUsage:  10.0 + (rand.Float32() * 20.0),
		DiskUsage:    10.0 + (rand.Float32() * 20.0),
		NetworkUsage: 10.0 + (rand.Float32() * 20.0),
	}
}

func (m *RandomDataDriver) CheckNetworkConnection() bool {
	return true // Maybe this could be random, to simulate instable networks
}

func (m *RandomDataDriver) HandleReconnect() {
}
