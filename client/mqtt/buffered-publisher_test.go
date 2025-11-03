package mqtt

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/RicardoCenci/iot-distributed-architecture/client/queue"
	mqttProvider "github.com/eclipse/paho.mqtt.golang"
)

type mockMQTTClient struct {
	mqttProvider.Client
	publishError error
	publishCalls int
	publishFunc  func(topic string, qos byte, retained bool, payload interface{}) mqttProvider.Token
}

func (m *mockMQTTClient) Publish(topic string, qos byte, retained bool, payload interface{}) mqttProvider.Token {
	m.publishCalls++
	if m.publishFunc != nil {
		return m.publishFunc(topic, qos, retained, payload)
	}
	return &mockToken{err: m.publishError}
}

type mockToken struct {
	mqttProvider.Token
	err error
}

func (m *mockToken) Wait() bool {
	return true
}

func (m *mockToken) Error() error {
	return m.err
}

func TestBufferedPublisher_Run(t *testing.T) {
	tests := []struct {
		name            string
		publishError    error
		wantRequeue     bool
		validateMetrics func(*Metrics) bool
	}{
		{
			name:         "successful publish",
			publishError: nil,
			wantRequeue:  false,
			validateMetrics: func(m *Metrics) bool {
				return m.GetNumberOfMessages() == 1 && m.GetNumberOfErrors() == 0
			},
		},
		{
			name:         "failed publish",
			publishError: errors.New("publish failed"),
			wantRequeue:  true,
			validateMetrics: func(m *Metrics) bool {
				return m.GetNumberOfMessages() == 1 && m.GetNumberOfErrors() == 1
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := &mockLogger{}
			mockMQTT := &mockMQTTClient{publishError: tt.publishError}
			metrics := NewMetrics("test/topic")
			backoffConfig := queue.BackoffConfig{
				Base:       100 * time.Millisecond,
				Factor:     2,
				MaxDelay:   1 * time.Second,
				MaxRetries: 3,
			}
			q := queue.New[string](
				queue.WithCapacity[string](10),
				queue.WithBackoff[string](backoffConfig),
			)

			transformer := func(msg string) string {
				return "transformed: " + msg
			}

			publisher := BufferedPublisher[string]{
				Logger:             logger,
				Client:             &Client{client: mockMQTT, logger: logger},
				Metrics:            metrics,
				Queue:              q,
				MessageTransformer: transformer,
				QoS:                1,
				Topic:              "test/topic",
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			go publisher.Run(ctx)

			msg := queue.Message[string]{Data: "test message"}
			if err := q.Enqueue(msg); err != nil {
				t.Fatal(err)
			}

			maxWait := 2 * time.Second
			checkInterval := 50 * time.Millisecond
			elapsed := time.Duration(0)
			metricsValid := false

			for elapsed < maxWait {
				time.Sleep(checkInterval)
				elapsed += checkInterval
				if tt.validateMetrics(metrics) {
					metricsValid = true
					break
				}
			}

			if !metricsValid {
				t.Errorf("Metrics validation failed after %v. Messages: %d, Errors: %d",
					elapsed, metrics.GetNumberOfMessages(), metrics.GetNumberOfErrors())
			}

			if tt.wantRequeue {
				time.Sleep(500 * time.Millisecond)
				if mockMQTT.publishCalls < 2 {
					t.Errorf("Message should be requeued on error. Expected at least 2 publish calls, got %d",
						mockMQTT.publishCalls)
				}
			}

			cancel()
			q.Close()
		})
	}
}

func TestBufferedPublisher_MessageTransformer(t *testing.T) {
	logger := &mockLogger{}
	mockMQTT := &mockMQTTClient{}
	metrics := NewMetrics("test/topic")
	q := queue.New[int](queue.WithCapacity[int](10))

	transformer := func(msg int) string {
		return "number: " + string(rune(msg+'0'))
	}

	publisher := BufferedPublisher[int]{
		Logger:             logger,
		Client:             &Client{client: mockMQTT, logger: logger},
		Metrics:            metrics,
		Queue:              q,
		MessageTransformer: transformer,
		QoS:                1,
		Topic:              "test/topic",
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go publisher.Run(ctx)

	msg := queue.Message[int]{Data: 42}
	if err := q.Enqueue(msg); err != nil {
		t.Fatal(err)
	}

	time.Sleep(100 * time.Millisecond)

	if mockMQTT.publishCalls == 0 {
		t.Error("MessageTransformer should be called")
	}

	cancel()
	q.Close()
}

func TestBufferedPublisher_Close(t *testing.T) {
	logger := &mockLogger{}
	mockMQTT := &mockMQTTClient{}
	metrics := NewMetrics("test/topic")
	q := queue.New[string](queue.WithCapacity[string](10))

	publisher := BufferedPublisher[string]{
		Logger:             logger,
		Client:             &Client{client: mockMQTT, logger: logger},
		Metrics:            metrics,
		Queue:              q,
		MessageTransformer: func(msg string) string { return msg },
		QoS:                1,
		Topic:              "test/topic",
	}

	publisher.Close()

	if q.Len() != 0 && q.Cap() != 0 {
		err := q.Enqueue(queue.Message[string]{Data: "test"})
		if err == nil {
			t.Error("Close() should prevent enqueueing after close")
		}
	}
}

func TestBufferedPublisher_MultipleMessages(t *testing.T) {
	logger := &mockLogger{}
	mockMQTT := &mockMQTTClient{}
	metrics := NewMetrics("test/topic")
	q := queue.New[string](queue.WithCapacity[string](10))

	publisher := BufferedPublisher[string]{
		Logger:             logger,
		Client:             &Client{client: mockMQTT, logger: logger},
		Metrics:            metrics,
		Queue:              q,
		MessageTransformer: func(msg string) string { return msg },
		QoS:                1,
		Topic:              "test/topic",
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go publisher.Run(ctx)

	for i := 0; i < 5; i++ {
		msg := queue.Message[string]{Data: "message " + string(rune(i+'0'))}
		if err := q.Enqueue(msg); err != nil {
			t.Fatal(err)
		}
	}

	time.Sleep(200 * time.Millisecond)

	if metrics.GetNumberOfMessages() != 5 {
		t.Errorf("Multiple messages: numberOfMessages = %v, want %v", metrics.GetNumberOfMessages(), 5)
	}

	if mockMQTT.publishCalls != 5 {
		t.Errorf("Multiple messages: publishCalls = %v, want %v", mockMQTT.publishCalls, 5)
	}

	cancel()
	q.Close()
}
