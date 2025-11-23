package rabbitmq

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
	MessageTransformer func(T) []byte
	Exchange           string
	RoutingKey         string
}

func (bp *BufferedPublisher[T]) Run(ctx context.Context) {
	for msg := range bp.Queue.Items() {
		payload := bp.MessageTransformer(msg.Data)

		startTime := time.Now()

		err := bp.Client.Publish(
			bp.Exchange,
			bp.RoutingKey,
			payload,
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

		bp.Logger.Debug("Published message", "exchange", bp.Exchange, "routing_key", bp.RoutingKey, "payload_size", len(payload))
	}
}

func (bp *BufferedPublisher[T]) Close() {
	bp.Queue.Close()
}

