package prometheus

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/RicardoCenci/iot-distributed-architecture/shared/logger"
	"github.com/RicardoCenci/iot-distributed-architecture/workers/metrics/parser"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Client struct {
	registry   *prometheus.Registry
	logger     logger.Interface
	httpServer *http.Server
	mu         sync.RWMutex

	cpuUsageGauge     *prometheus.GaugeVec
	memoryUsageGauge  *prometheus.GaugeVec
	diskUsageGauge    *prometheus.GaugeVec
	networkUsageGauge *prometheus.GaugeVec
}

func NewClient(logger logger.Interface, listenAddress string) *Client {
	registry := prometheus.NewRegistry()

	cpuUsageGauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "iot_device_cpu_usage_percent",
			Help: "CPU usage percentage for IoT device",
		},
		[]string{"device_id"},
	)

	memoryUsageGauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "iot_device_memory_usage_percent",
			Help: "Memory usage percentage for IoT device",
		},
		[]string{"device_id"},
	)

	diskUsageGauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "iot_device_disk_usage_percent",
			Help: "Disk usage percentage for IoT device",
		},
		[]string{"device_id"},
	)

	networkUsageGauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "iot_device_network_usage_percent",
			Help: "Network usage percentage for IoT device",
		},
		[]string{"device_id"},
	)

	registry.MustRegister(cpuUsageGauge)
	registry.MustRegister(memoryUsageGauge)
	registry.MustRegister(diskUsageGauge)
	registry.MustRegister(networkUsageGauge)

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	server := &http.Server{
		Addr:    listenAddress,
		Handler: mux,
	}

	return &Client{
		registry:          registry,
		logger:            logger,
		httpServer:        server,
		cpuUsageGauge:     cpuUsageGauge,
		memoryUsageGauge:  memoryUsageGauge,
		diskUsageGauge:    diskUsageGauge,
		networkUsageGauge: networkUsageGauge,
	}
}

func (c *Client) RecordMetric(metric parser.MetricData) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	deviceID := metric.DeviceID

	c.cpuUsageGauge.WithLabelValues(deviceID).Set(float64(metric.CPUUsage))
	c.memoryUsageGauge.WithLabelValues(deviceID).Set(float64(metric.MemoryUsage))
	c.diskUsageGauge.WithLabelValues(deviceID).Set(float64(metric.DiskUsage))
	c.networkUsageGauge.WithLabelValues(deviceID).Set(float64(metric.NetworkUsage))

	c.logger.Debug("Recorded metrics", "device_id", deviceID,
		"cpu", metric.CPUUsage,
		"memory", metric.MemoryUsage,
		"disk", metric.DiskUsage,
		"network", metric.NetworkUsage,
	)

	return nil
}

func (c *Client) Start() error {
	go func() {
		if err := c.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			c.logger.Error("Failed to start Prometheus HTTP server", "error", err)
		}
	}()

	c.logger.Info("Prometheus metrics server started", "address", c.httpServer.Addr)
	return nil
}

func (c *Client) Close() error {
	if c.httpServer != nil {
		return c.httpServer.Close()
	}
	return nil
}

func (c *Client) GetMetricsEndpoint() string {
	return fmt.Sprintf("http://%s/metrics", c.httpServer.Addr)
}
