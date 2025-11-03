package drivers

import (
	"testing"
)

func TestNewRandomDataDriver(t *testing.T) {
	driver := NewRandomDataDriver()
	if driver == nil {
		t.Error("NewRandomDataDriver() returned nil")
	}

	_, ok := driver.(DriverInterface)
	if !ok {
		t.Error("NewRandomDataDriver() does not implement DriverInterface")
	}
}

func TestRandomDataDriver_ProbeSensor(t *testing.T) {
	driver := NewRandomDataDriver()

	data := driver.ProbeSensor()

	if data.Humidity < 40.0 || data.Humidity > 60.0 {
		t.Errorf("ProbeSensor() Humidity = %v, want between 40.0 and 60.0", data.Humidity)
	}

	if data.Temperature < 20.0 || data.Temperature > 30.0 {
		t.Errorf("ProbeSensor() Temperature = %v, want between 20.0 and 30.0", data.Temperature)
	}

	data2 := driver.ProbeSensor()
	if data.Humidity == data2.Humidity && data.Temperature == data2.Temperature {
		t.Error("ProbeSensor() should return different values on each call")
	}
}

func TestRandomDataDriver_ProbeSystemMetrics(t *testing.T) {
	driver := NewRandomDataDriver()

	metrics := driver.ProbeSystemMetrics()

	if metrics.CPUUsage < 10.0 || metrics.CPUUsage > 30.0 {
		t.Errorf("ProbeSystemMetrics() CPUUsage = %v, want between 10.0 and 30.0", metrics.CPUUsage)
	}

	if metrics.MemoryUsage < 10.0 || metrics.MemoryUsage > 30.0 {
		t.Errorf("ProbeSystemMetrics() MemoryUsage = %v, want between 10.0 and 30.0", metrics.MemoryUsage)
	}

	if metrics.DiskUsage < 10.0 || metrics.DiskUsage > 30.0 {
		t.Errorf("ProbeSystemMetrics() DiskUsage = %v, want between 10.0 and 30.0", metrics.DiskUsage)
	}

	if metrics.NetworkUsage < 10.0 || metrics.NetworkUsage > 30.0 {
		t.Errorf("ProbeSystemMetrics() NetworkUsage = %v, want between 10.0 and 30.0", metrics.NetworkUsage)
	}

	metrics2 := driver.ProbeSystemMetrics()
	if metrics.CPUUsage == metrics2.CPUUsage &&
		metrics.MemoryUsage == metrics2.MemoryUsage &&
		metrics.DiskUsage == metrics2.DiskUsage &&
		metrics.NetworkUsage == metrics2.NetworkUsage {
		t.Error("ProbeSystemMetrics() should return different values on each call")
	}
}

func TestRandomDataDriver_CheckNetworkConnection(t *testing.T) {
	driver := NewRandomDataDriver()

	result := driver.CheckNetworkConnection()
	if !result {
		t.Error("CheckNetworkConnection() = false, want true")
	}

	result2 := driver.CheckNetworkConnection()
	if result != result2 {
		t.Error("CheckNetworkConnection() should return consistent values")
	}
}

func TestRandomDataDriver_HandleReconnect(t *testing.T) {
	driver := NewRandomDataDriver()

	driver.HandleReconnect()

	if driver == nil {
		t.Error("HandleReconnect() should not panic")
	}
}
