package rabbitmq

import (
	"fmt"
	"sync"
	"time"

	"github.com/RicardoCenci/iot-distributed-architecture/shared/logger"
)

type Metrics struct {
	mu              sync.RWMutex
	topic           string
	totalPublished  int64
	totalFailed     int64
	totalLatency    time.Duration
	minLatency      time.Duration
	maxLatency      time.Duration
}

func NewMetrics(topic string) *Metrics {
	return &Metrics{
		topic:      topic,
		minLatency: time.Hour,
	}
}

func (m *Metrics) Update(latency time.Duration, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err != nil {
		m.totalFailed++
	} else {
		m.totalPublished++
		m.totalLatency += latency

		if latency < m.minLatency {
			m.minLatency = latency
		}
		if latency > m.maxLatency {
			m.maxLatency = latency
		}
	}
}

func (m *Metrics) Print(logger logger.Interface) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var avgLatency time.Duration
	if m.totalPublished > 0 {
		avgLatency = m.totalLatency / time.Duration(m.totalPublished)
	}

	logger.Info("RabbitMQ Publisher Metrics",
		"topic", m.topic,
		"total_published", m.totalPublished,
		"total_failed", m.totalFailed,
		"avg_latency_ms", avgLatency.Milliseconds(),
		"min_latency_ms", m.minLatency.Milliseconds(),
		"max_latency_ms", m.maxLatency.Milliseconds(),
	)
}

func (m *Metrics) String() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var avgLatency time.Duration
	if m.totalPublished > 0 {
		avgLatency = m.totalLatency / time.Duration(m.totalPublished)
	}

	return fmt.Sprintf("Topic: %s, Published: %d, Failed: %d, Avg Latency: %v",
		m.topic, m.totalPublished, m.totalFailed, avgLatency)
}

