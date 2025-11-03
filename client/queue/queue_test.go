package queue

import (
	"context"
	"testing"
	"time"
)

func TestNewQueue(t *testing.T) {
	q := New[int]()
	if q == nil {
		t.Error("New() returned nil")
	}
	if q.channel == nil {
		t.Error("New() channel is nil")
	}
	if q.done == nil {
		t.Error("New() done channel is nil")
	}
}

func TestQueue_WithCapacity(t *testing.T) {
	q := New[int](WithCapacity[int](5))
	if cap(q.channel) != 5 {
		t.Errorf("WithCapacity() capacity = %v, want %v", cap(q.channel), 5)
	}
}

func TestQueue_WithBackoff(t *testing.T) {
	backoff := BackoffConfig{
		Base:       3 * time.Second,
		Factor:     3,
		MaxDelay:   15 * time.Second,
		MaxRetries: 5,
	}
	q := New[int](WithBackoff[int](backoff))
	if q.backoff.Base != backoff.Base {
		t.Errorf("WithBackoff() Base = %v, want %v", q.backoff.Base, backoff.Base)
	}
	if q.backoff.Factor != backoff.Factor {
		t.Errorf("WithBackoff() Factor = %v, want %v", q.backoff.Factor, backoff.Factor)
	}
	if q.backoff.MaxDelay != backoff.MaxDelay {
		t.Errorf("WithBackoff() MaxDelay = %v, want %v", q.backoff.MaxDelay, backoff.MaxDelay)
	}
	if q.backoff.MaxRetries != backoff.MaxRetries {
		t.Errorf("WithBackoff() MaxRetries = %v, want %v", q.backoff.MaxRetries, backoff.MaxRetries)
	}
}

func TestQueue_Enqueue(t *testing.T) {
	q := New[int](WithCapacity[int](2))

	msg := Message[int]{Data: 42}
	err := q.Enqueue(msg)
	if err != nil {
		t.Errorf("Enqueue() error = %v, want nil", err)
	}
	if q.Len() != 1 {
		t.Errorf("Enqueue() length = %v, want %v", q.Len(), 1)
	}

	err = q.Enqueue(Message[int]{Data: 43})
	if err != nil {
		t.Errorf("Enqueue() error = %v, want nil", err)
	}

	err = q.Enqueue(Message[int]{Data: 44})
	if err != ErrFull {
		t.Errorf("Enqueue() error = %v, want %v", err, ErrFull)
	}
}

func TestQueue_Enqueue_Closed(t *testing.T) {
	q := New[int]()
	q.Close()

	err := q.Enqueue(Message[int]{Data: 42})
	if err != ErrClosed {
		t.Errorf("Enqueue() on closed queue error = %v, want %v", err, ErrClosed)
	}
}

func TestQueue_Dequeue(t *testing.T) {
	q := New[int]()
	ctx := context.Background()

	msg := Message[int]{Data: 42}
	if err := q.Enqueue(msg); err != nil {
		t.Fatal(err)
	}

	dequeued, ok, err := q.Dequeue(ctx)
	if err != nil {
		t.Errorf("Dequeue() error = %v, want nil", err)
	}
	if !ok {
		t.Error("Dequeue() ok = false, want true")
	}
	if dequeued.Data != 42 {
		t.Errorf("Dequeue() Data = %v, want %v", dequeued.Data, 42)
	}
}

func TestQueue_Dequeue_ContextCancel(t *testing.T) {
	q := New[int]()
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	_, ok, err := q.Dequeue(ctx)
	if err == nil {
		t.Error("Dequeue() should return error when context is cancelled")
	}
	if ok {
		t.Error("Dequeue() ok = true, want false when context cancelled")
	}
}

func TestQueue_Dequeue_Closed(t *testing.T) {
	q := New[int]()
	ctx := context.Background()

	q.Close()

	_, ok, err := q.Dequeue(ctx)
	if err != ErrClosed {
		t.Errorf("Dequeue() on closed queue error = %v, want %v", err, ErrClosed)
	}
	if ok {
		t.Error("Dequeue() on closed queue ok = true, want false")
	}
}

func TestQueue_Items(t *testing.T) {
	q := New[int]()

	items := q.Items()
	if items == nil {
		t.Error("Items() returned nil channel")
	}

	msg := Message[int]{Data: 42}
	if err := q.Enqueue(msg); err != nil {
		t.Fatal(err)
	}

	select {
	case item := <-items:
		if item.Data != 42 {
			t.Errorf("Items() Data = %v, want %v", item.Data, 42)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Items() channel did not receive message")
	}
}

func TestQueue_Close(t *testing.T) {
	q := New[int]()

	if q.closed {
		t.Error("Queue should not be closed initially")
	}

	q.Close()

	if !q.closed {
		t.Error("Queue should be closed after Close()")
	}

	select {
	case _, ok := <-q.channel:
		if ok {
			t.Error("Queue channel should be closed")
		}
	default:
		t.Error("Queue channel should be closed")
	}
}

func TestQueue_Len(t *testing.T) {
	q := New[int]()

	if q.Len() != 0 {
		t.Errorf("Len() = %v, want %v", q.Len(), 0)
	}

	q.Enqueue(Message[int]{Data: 1})
	if q.Len() != 1 {
		t.Errorf("Len() = %v, want %v", q.Len(), 1)
	}

	q.Enqueue(Message[int]{Data: 2})
	if q.Len() != 2 {
		t.Errorf("Len() = %v, want %v", q.Len(), 2)
	}
}

func TestQueue_Cap(t *testing.T) {
	q := New[int](WithCapacity[int](10))

	if q.Cap() != 10 {
		t.Errorf("Cap() = %v, want %v", q.Cap(), 10)
	}
}

func TestBackoffConfig_delayForAttempt(t *testing.T) {
	tests := []struct {
		name    string
		config  BackoffConfig
		attempt int
		wantMin time.Duration
		wantMax time.Duration
	}{
		{
			name: "first attempt",
			config: BackoffConfig{
				Base:   2 * time.Second,
				Factor: 2,
			},
			attempt: 1,
			wantMin: 2 * time.Second,
			wantMax: 2 * time.Second,
		},
		{
			name: "second attempt",
			config: BackoffConfig{
				Base:   2 * time.Second,
				Factor: 2,
			},
			attempt: 2,
			wantMin: 4 * time.Second,
			wantMax: 4 * time.Second,
		},
		{
			name: "with max delay",
			config: BackoffConfig{
				Base:     2 * time.Second,
				Factor:   2,
				MaxDelay: 5 * time.Second,
			},
			attempt: 5,
			wantMin: 0,
			wantMax: 5 * time.Second,
		},
		{
			name: "zero attempt",
			config: BackoffConfig{
				Base:   2 * time.Second,
				Factor: 2,
			},
			attempt: 0,
			wantMin: 2 * time.Second,
			wantMax: 2 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.config.delayForAttempt(tt.attempt)
			if got < tt.wantMin || got > tt.wantMax {
				t.Errorf("delayForAttempt() = %v, want between %v and %v", got, tt.wantMin, tt.wantMax)
			}
		})
	}
}

func TestQueue_RequeueAfter(t *testing.T) {
	q := New[int]()
	ctx := context.Background()

	msg := Message[int]{Data: 42}
	q.RequeueAfter(ctx, msg, 50*time.Millisecond)

	time.Sleep(100 * time.Millisecond)

	dequeued, ok, err := q.Dequeue(ctx)
	if err != nil {
		t.Errorf("RequeueAfter() error = %v, want nil", err)
	}
	if !ok {
		t.Error("RequeueAfter() message not requeued")
	}
	if dequeued.Data != 42 {
		t.Errorf("RequeueAfter() Data = %v, want %v", dequeued.Data, 42)
	}
}

func TestQueue_RequeueAfter_Closed(t *testing.T) {
	q := New[int]()
	ctx := context.Background()

	q.Close()

	msg := Message[int]{Data: 42}
	q.RequeueAfter(ctx, msg, 50*time.Millisecond)

	time.Sleep(100 * time.Millisecond)

	if q.Len() != 0 {
		t.Error("RequeueAfter() should not requeue when queue is closed")
	}
}

func TestQueue_RequeueWithBackoff(t *testing.T) {
	q := New[int](WithBackoff[int](BackoffConfig{
		Base:       50 * time.Millisecond,
		Factor:     2,
		MaxDelay:   200 * time.Millisecond,
		MaxRetries: 3,
	}))
	ctx := context.Background()

	msg := Message[int]{Data: 42, NumberOfRetries: 0}
	q.RequeueWithBackoff(ctx, msg)

	time.Sleep(150 * time.Millisecond)

	dequeued, ok, err := q.Dequeue(ctx)
	if err != nil {
		t.Errorf("RequeueWithBackoff() error = %v, want nil", err)
	}
	if !ok {
		t.Error("RequeueWithBackoff() message not requeued")
	}
	if dequeued.Data != 42 {
		t.Errorf("RequeueWithBackoff() Data = %v, want %v", dequeued.Data, 42)
	}
	if dequeued.NumberOfRetries != 1 {
		t.Errorf("RequeueWithBackoff() NumberOfRetries = %v, want %v", dequeued.NumberOfRetries, 1)
	}
}

func TestQueue_RequeueWithBackoff_MaxRetries(t *testing.T) {
	q := New[int](WithBackoff[int](BackoffConfig{
		Base:       10 * time.Millisecond,
		Factor:     2,
		MaxDelay:   100 * time.Millisecond,
		MaxRetries: 2,
	}), WithCapacity[int](1))
	ctx := context.Background()

	msg := Message[int]{Data: 42, NumberOfRetries: 2}
	q.Enqueue(Message[int]{Data: 1})

	q.RequeueWithBackoff(ctx, msg)

	time.Sleep(200 * time.Millisecond)

	if q.Len() > 1 {
		t.Error("RequeueWithBackoff() should respect MaxRetries")
	}
}
