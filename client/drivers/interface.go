package drivers

type SensorData struct {
	Humidity    float32
	Temperature float32
}

type SystemMetrics struct {
	CPUUsage     float32
	MemoryUsage  float32
	DiskUsage    float32
	NetworkUsage float32
}

type DriverInterface interface {
	ProbeSensor() SensorData
	ProbeSystemMetrics() SystemMetrics
	CheckNetworkConnection() bool
	HandleReconnect()
}
