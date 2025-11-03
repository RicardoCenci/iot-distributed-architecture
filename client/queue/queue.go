package queue

import (
	"context"
	"errors"
	"math"
	"sync"
	"time"
)

var (
	ErrClosed = errors.New("queue closed")
	ErrFull   = errors.New("queue full")
)

type Queue[T any] struct {
	channel chan Message[T]
	mutex   sync.Mutex
	closed  bool
	done    chan struct{}
	backoff BackoffConfig
}

type Option[T any] func(*Queue[T])

type Message[T any] struct {
	NumberOfRetries int
	Data            T
}

func WithCapacity[T any](capacity int) Option[T] {
	return func(q *Queue[T]) {
		q.channel = make(chan Message[T], capacity)
	}
}

func WithBackoff[T any](backoff BackoffConfig) Option[T] {
	return func(q *Queue[T]) {
		q.backoff = backoff
	}
}

func New[T any](options ...Option[T]) *Queue[T] {
	q := &Queue[T]{
		channel: make(chan Message[T], 100),
		done:    make(chan struct{}),
		backoff: BackoffConfig{
			Base:       2 * time.Second,
			Factor:     2,
			MaxDelay:   time.Second * 10,
			MaxRetries: 3,
		},
	}

	for _, option := range options {
		option(q)
	}

	return q
}

func (q *Queue[T]) Enqueue(item Message[T]) error {
	q.mutex.Lock()
	if q.closed {
		q.mutex.Unlock()
		return ErrClosed
	}
	ch := q.channel
	q.mutex.Unlock()

	select {
	case ch <- item:
		return nil
	default:
		return ErrFull
	}
}

func (q *Queue[T]) Items() <-chan Message[T] {
	q.mutex.Lock()
	ch := q.channel
	q.mutex.Unlock()
	return ch
}

func (q *Queue[T]) Dequeue(ctx context.Context) (Message[T], bool, error) {
	q.mutex.Lock()
	ch := q.channel
	q.mutex.Unlock()

	select {
	case v, ok := <-ch:
		if !ok {
			return Message[T]{}, false, ErrClosed
		}
		return v, true, nil
	case <-ctx.Done():
		return Message[T]{}, false, ctx.Err()
	}
}

func (q *Queue[T]) Close() {
	q.mutex.Lock()
	if !q.closed {
		q.closed = true
		close(q.done)
		close(q.channel)
	}
	q.mutex.Unlock()
}

func (q *Queue[T]) Len() int {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return len(q.channel)
}

func (q *Queue[T]) Cap() int {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return cap(q.channel)
}

type BackoffConfig struct {
	Base       time.Duration
	Factor     float64
	MaxDelay   time.Duration
	MaxRetries int // 0 means unlimited retries
}

func (c BackoffConfig) delayForAttempt(attempt int) time.Duration {
	if attempt < 1 {
		attempt = 1
	}
	base := c.Base
	if base <= 0 {
		base = 100 * time.Millisecond
	}
	factor := c.Factor
	if factor <= 0 {
		factor = 2
	}
	d := float64(base) * math.Pow(factor, float64(attempt-1))
	dur := time.Duration(d)
	if c.MaxDelay > 0 && dur > c.MaxDelay {
		return c.MaxDelay
	}
	return dur
}

func (q *Queue[T]) RequeueAfter(ctx context.Context, item Message[T], delay time.Duration) {
	q.mutex.Lock()
	if q.closed {
		q.mutex.Unlock()
		return
	}
	q.mutex.Unlock()

	go func() {
		if delay <= 0 {
			_ = q.Enqueue(item)
			return
		}
		t := time.NewTimer(delay)
		defer t.Stop()
		select {
		case <-t.C:
			_ = q.Enqueue(item)
		case <-ctx.Done():
			return
		case <-q.done:
			return
		}
	}()
}

func (q *Queue[T]) RequeueWithBackoff(ctx context.Context, item Message[T]) {
	q.mutex.Lock()
	if q.closed {
		q.mutex.Unlock()
		return
	}
	q.mutex.Unlock()

	go func() {
		if item.NumberOfRetries < 1 {
			item.NumberOfRetries = 1
		}

		for {
			delay := q.backoff.delayForAttempt(item.NumberOfRetries)
			t := time.NewTimer(delay)
			select {
			case <-t.C:
				if err := q.Enqueue(item); err == nil {
					return
				} else if err == ErrClosed {
					return
				}
				// queue full: consider retrying
				if q.backoff.MaxRetries > 0 && item.NumberOfRetries >= q.backoff.MaxRetries {
					return
				}
				item.NumberOfRetries++
			case <-ctx.Done():
				t.Stop()
				return
			case <-q.done:
				t.Stop()
				return
			}
		}
	}()
}
