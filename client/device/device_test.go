package device

import (
	"testing"

	"github.com/RicardoCenci/iot-distributed-architecture/client/drivers"
)

type mockDriver struct {
	sensorData      drivers.SensorData
	systemMetrics   drivers.SystemMetrics
	networkStatus   bool
	reconnectCalled bool
}

func (m *mockDriver) ProbeSensor() drivers.SensorData {
	return m.sensorData
}

func (m *mockDriver) ProbeSystemMetrics() drivers.SystemMetrics {
	return m.systemMetrics
}

func (m *mockDriver) CheckNetworkConnection() bool {
	return m.networkStatus
}

func (m *mockDriver) HandleReconnect() {
	m.reconnectCalled = true
}

func TestNewDevice(t *testing.T) {
	driver := &mockDriver{}
	device := NewDevice("test-device-123", driver)

	if device.DeviceID != "test-device-123" {
		t.Errorf("NewDevice() DeviceID = %v, want %v", device.DeviceID, "test-device-123")
	}
	if device.driver != driver {
		t.Error("NewDevice() driver not set correctly")
	}
}

func TestDevice_GetSensorData(t *testing.T) {
	expectedData := drivers.SensorData{
		Humidity:    45.5,
		Temperature: 23.7,
	}
	driver := &mockDriver{sensorData: expectedData}
	device := NewDevice("test-device", driver)

	data := device.GetSensorData()
	if data.Humidity != expectedData.Humidity {
		t.Errorf("GetSensorData() Humidity = %v, want %v", data.Humidity, expectedData.Humidity)
	}
	if data.Temperature != expectedData.Temperature {
		t.Errorf("GetSensorData() Temperature = %v, want %v", data.Temperature, expectedData.Temperature)
	}
}

func TestDevice_GetSystemMetrics(t *testing.T) {
	expectedMetrics := drivers.SystemMetrics{
		CPUUsage:     15.5,
		MemoryUsage:  30.2,
		DiskUsage:    45.8,
		NetworkUsage: 12.3,
	}
	driver := &mockDriver{systemMetrics: expectedMetrics}
	device := NewDevice("test-device", driver)

	metrics := device.GetSystemMetrics()
	if metrics.CPUUsage != expectedMetrics.CPUUsage {
		t.Errorf("GetSystemMetrics() CPUUsage = %v, want %v", metrics.CPUUsage, expectedMetrics.CPUUsage)
	}
	if metrics.MemoryUsage != expectedMetrics.MemoryUsage {
		t.Errorf("GetSystemMetrics() MemoryUsage = %v, want %v", metrics.MemoryUsage, expectedMetrics.MemoryUsage)
	}
	if metrics.DiskUsage != expectedMetrics.DiskUsage {
		t.Errorf("GetSystemMetrics() DiskUsage = %v, want %v", metrics.DiskUsage, expectedMetrics.DiskUsage)
	}
	if metrics.NetworkUsage != expectedMetrics.NetworkUsage {
		t.Errorf("GetSystemMetrics() NetworkUsage = %v, want %v", metrics.NetworkUsage, expectedMetrics.NetworkUsage)
	}
}

func TestDevice_IsConnectedToInternet(t *testing.T) {
	tests := []struct {
		name           string
		networkStatus  bool
		expectedResult bool
	}{
		{
			name:           "connected",
			networkStatus:  true,
			expectedResult: true,
		},
		{
			name:           "not connected",
			networkStatus:  false,
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			driver := &mockDriver{networkStatus: tt.networkStatus}
			device := NewDevice("test-device", driver)

			result := device.IsConnectedToInternet()
			if result != tt.expectedResult {
				t.Errorf("IsConnectedToInternet() = %v, want %v", result, tt.expectedResult)
			}
		})
	}
}
