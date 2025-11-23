package mqtt

import (
	"context"
	"time"

	"github.com/RicardoCenci/iot-distributed-architecture/client/queue"
	"github.com/RicardoCenci/iot-distributed-architecture/shared/logger"
)

type BufferedPublisher[T any] struct {
	Logger             logger.Interface
	Client             *Client
	Metrics            *Metrics
	Queue              *queue.Queue[T]
	MessageTransformer func(T) ([]byte, error)
	QoS                int
	Topic              string
}

func (bp *BufferedPublisher[T]) Run(ctx context.Context) {

	for msg := range bp.Queue.Items() {

		payload, err := bp.MessageTransformer(msg.Data)
		if err != nil {
			bp.Logger.Error("Failed to transform message", "error", err)
			bp.Queue.RequeueWithBackoff(ctx, msg)
			continue
		}

		startTime := time.Now()

		err = bp.Client.Publish(
			bp.Topic,
			payload,
			bp.QoS,
			false,
		)

		bp.Metrics.Update(
			time.Since(startTime),
			err,
		)

		if err != nil {
			bp.Logger.Error("Failed to publish message", "error", err, "retries", msg.NumberOfRetries)
			bp.Queue.RequeueWithBackoff(ctx, msg)
			continue
		}

		bp.Logger.Debug("Published message", "topic", bp.Topic, "payload_size", len(payload))
	}
}

func (bp *BufferedPublisher[T]) Close() {
	bp.Queue.Close()
}
