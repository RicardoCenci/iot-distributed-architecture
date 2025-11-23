package prometheus

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/RicardoCenci/iot-distributed-architecture/shared/logger"
	"github.com/RicardoCenci/iot-distributed-architecture/workers/metrics/parser"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type metricValue struct {
	value     float64
	timestamp time.Time
}

type timestampedCollector struct {
	mu sync.RWMutex

	cpuUsageDesc     *prometheus.Desc
	memoryUsageDesc  *prometheus.Desc
	diskUsageDesc    *prometheus.Desc
	networkUsageDesc *prometheus.Desc

	cpuUsage     map[string]metricValue
	memoryUsage  map[string]metricValue
	diskUsage    map[string]metricValue
	networkUsage map[string]metricValue
}

func newTimestampedCollector() *timestampedCollector {
	return &timestampedCollector{
		cpuUsageDesc: prometheus.NewDesc(
			"iot_device_cpu_usage_percent",
			"CPU usage percentage for IoT device",
			[]string{"device_id"},
			nil,
		),
		memoryUsageDesc: prometheus.NewDesc(
			"iot_device_memory_usage_percent",
			"Memory usage percentage for IoT device",
			[]string{"device_id"},
			nil,
		),
		diskUsageDesc: prometheus.NewDesc(
			"iot_device_disk_usage_percent",
			"Disk usage percentage for IoT device",
			[]string{"device_id"},
			nil,
		),
		networkUsageDesc: prometheus.NewDesc(
			"iot_device_network_usage_percent",
			"Network usage percentage for IoT device",
			[]string{"device_id"},
			nil,
		),
		cpuUsage:     make(map[string]metricValue),
		memoryUsage:  make(map[string]metricValue),
		diskUsage:    make(map[string]metricValue),
		networkUsage: make(map[string]metricValue),
	}
}

func (tc *timestampedCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- tc.cpuUsageDesc
	ch <- tc.memoryUsageDesc
	ch <- tc.diskUsageDesc
	ch <- tc.networkUsageDesc
}

func (tc *timestampedCollector) Collect(ch chan<- prometheus.Metric) {
	tc.mu.RLock()
	defer tc.mu.RUnlock()

	for deviceID, val := range tc.cpuUsage {
		metric := prometheus.MustNewConstMetric(
			tc.cpuUsageDesc,
			prometheus.GaugeValue,
			val.value,
			deviceID,
		)
		ch <- prometheus.NewMetricWithTimestamp(val.timestamp, metric)
	}

	for deviceID, val := range tc.memoryUsage {
		metric := prometheus.MustNewConstMetric(
			tc.memoryUsageDesc,
			prometheus.GaugeValue,
			val.value,
			deviceID,
		)
		ch <- prometheus.NewMetricWithTimestamp(val.timestamp, metric)
	}

	for deviceID, val := range tc.diskUsage {
		metric := prometheus.MustNewConstMetric(
			tc.diskUsageDesc,
			prometheus.GaugeValue,
			val.value,
			deviceID,
		)
		ch <- prometheus.NewMetricWithTimestamp(val.timestamp, metric)
	}

	for deviceID, val := range tc.networkUsage {
		metric := prometheus.MustNewConstMetric(
			tc.networkUsageDesc,
			prometheus.GaugeValue,
			val.value,
			deviceID,
		)
		ch <- prometheus.NewMetricWithTimestamp(val.timestamp, metric)
	}
}

func (tc *timestampedCollector) recordMetric(metric parser.MetricData) {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	deviceID := metric.DeviceID
	timestamp := metric.Timestamp

	tc.cpuUsage[deviceID] = metricValue{
		value:     float64(metric.CPUUsage),
		timestamp: timestamp,
	}
	tc.memoryUsage[deviceID] = metricValue{
		value:     float64(metric.MemoryUsage),
		timestamp: timestamp,
	}
	tc.diskUsage[deviceID] = metricValue{
		value:     float64(metric.DiskUsage),
		timestamp: timestamp,
	}
	tc.networkUsage[deviceID] = metricValue{
		value:     float64(metric.NetworkUsage),
		timestamp: timestamp,
	}
}

type Client struct {
	registry   *prometheus.Registry
	logger     logger.Interface
	httpServer *http.Server
	collector  *timestampedCollector
}

func NewClient(logger logger.Interface, listenAddress string) *Client {
	registry := prometheus.NewRegistry()
	collector := newTimestampedCollector()

	registry.MustRegister(collector)

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	server := &http.Server{
		Addr:    listenAddress,
		Handler: mux,
	}

	return &Client{
		registry:   registry,
		logger:     logger,
		httpServer: server,
		collector:  collector,
	}
}

func (c *Client) RecordMetric(metric parser.MetricData) error {
	c.collector.recordMetric(metric)

	c.logger.Debug("Recorded metrics",
		"device_id", metric.DeviceID,
		"cpu", metric.CPUUsage,
		"memory", metric.MemoryUsage,
		"disk", metric.DiskUsage,
		"network", metric.NetworkUsage,
		"timestamp", metric.Timestamp,
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
