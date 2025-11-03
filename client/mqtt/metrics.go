package mqtt

import (
	"sync/atomic"
	"time"

	"github.com/RicardoCenci/iot-distributed-architecture/client/logger"
)

type Metrics struct {
	sumPublishingTimes int64
	numberOfMessages   int64
	numberOfErrors     int64
	topic              string
}

func NewMetrics(topic string) *Metrics {
	return &Metrics{
		sumPublishingTimes: 0,
		numberOfMessages:   0,
		numberOfErrors:     0,
		topic:              topic,
	}
}
func (m *Metrics) Update(timeToPublish time.Duration, err error) {
	atomic.AddInt64(&m.numberOfMessages, 1)
	atomic.AddInt64(&m.sumPublishingTimes, int64(timeToPublish))

	if err != nil {
		atomic.AddInt64(&m.numberOfErrors, 1)
	}
}

func (m *Metrics) GetAvgPublishingTime() time.Duration {
	numberOfMessages := atomic.LoadInt64(&m.numberOfMessages)
	if numberOfMessages == 0 {
		return 0
	}

	return time.Duration(atomic.LoadInt64(&m.sumPublishingTimes)) / time.Duration(numberOfMessages)
}

func (m *Metrics) GetNumberOfErrors() int64 {
	return m.numberOfErrors
}

func (m *Metrics) GetNumberOfMessages() int64 {
	return m.numberOfMessages
}

func (m *Metrics) Print(logger logger.Interface) {
	logger.Info(
		"Topic Metrics",
		"topic", m.topic,
		"number_of_messages", m.GetNumberOfMessages(),
		"number_of_errors", m.GetNumberOfErrors(),
		"avg_publishing_time_ms", m.GetAvgPublishingTime().Milliseconds(),
	)
}
