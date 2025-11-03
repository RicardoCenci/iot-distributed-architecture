package device

import "github.com/RicardoCenci/iot-distributed-architecture/client/drivers"

type Device struct {
	DeviceID string
	driver   drivers.DriverInterface
}

func NewDevice(deviceID string, driver drivers.DriverInterface) *Device {
	return &Device{
		DeviceID: deviceID,
		driver:   driver,
	}
}

func (d *Device) GetSensorData() drivers.SensorData {
	return d.driver.ProbeSensor()
}

func (d *Device) GetSystemMetrics() drivers.SystemMetrics {
	return d.driver.ProbeSystemMetrics()
}

func (d *Device) IsConnectedToInternet() bool {
	return d.driver.CheckNetworkConnection()
}
