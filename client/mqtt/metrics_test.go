package mqtt

import (
	"testing"
	"time"
)

func TestNewMetrics(t *testing.T) {
	topic := "test/topic"
	metrics := NewMetrics(topic)

	if metrics.topic != topic {
		t.Errorf("NewMetrics() topic = %v, want %v", metrics.topic, topic)
	}
	if metrics.GetNumberOfMessages() != 0 {
		t.Errorf("NewMetrics() numberOfMessages = %v, want %v", metrics.GetNumberOfMessages(), 0)
	}
	if metrics.GetNumberOfErrors() != 0 {
		t.Errorf("NewMetrics() numberOfErrors = %v, want %v", metrics.GetNumberOfErrors(), 0)
	}
}

func TestMetrics_Update(t *testing.T) {
	metrics := NewMetrics("test/topic")

	tests := []struct {
		name     string
		duration time.Duration
		err      error
		wantMsgs int64
		wantErrs int64
	}{
		{
			name:     "successful update",
			duration: 100 * time.Millisecond,
			err:      nil,
			wantMsgs: 1,
			wantErrs: 0,
		},
		{
			name:     "failed update",
			duration: 200 * time.Millisecond,
			err:      &testError{msg: "test error"},
			wantMsgs: 2,
			wantErrs: 1,
		},
		{
			name:     "multiple updates",
			duration: 50 * time.Millisecond,
			err:      nil,
			wantMsgs: 3,
			wantErrs: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metrics.Update(tt.duration, tt.err)
			if metrics.GetNumberOfMessages() != tt.wantMsgs {
				t.Errorf("Update() numberOfMessages = %v, want %v", metrics.GetNumberOfMessages(), tt.wantMsgs)
			}
			if metrics.GetNumberOfErrors() != tt.wantErrs {
				t.Errorf("Update() numberOfErrors = %v, want %v", metrics.GetNumberOfErrors(), tt.wantErrs)
			}
		})
	}
}

func TestMetrics_GetAvgPublishingTime(t *testing.T) {
	metrics := NewMetrics("test/topic")

	if metrics.GetAvgPublishingTime() != 0 {
		t.Errorf("GetAvgPublishingTime() with no messages = %v, want %v", metrics.GetAvgPublishingTime(), 0)
	}

	metrics.Update(100*time.Millisecond, nil)
	metrics.Update(200*time.Millisecond, nil)
	metrics.Update(300*time.Millisecond, nil)

	avg := metrics.GetAvgPublishingTime()
	expected := 200 * time.Millisecond

	if avg != expected {
		t.Errorf("GetAvgPublishingTime() = %v, want %v", avg, expected)
	}
}

func TestMetrics_GetNumberOfMessages(t *testing.T) {
	metrics := NewMetrics("test/topic")

	if metrics.GetNumberOfMessages() != 0 {
		t.Errorf("GetNumberOfMessages() = %v, want %v", metrics.GetNumberOfMessages(), 0)
	}

	metrics.Update(100*time.Millisecond, nil)
	if metrics.GetNumberOfMessages() != 1 {
		t.Errorf("GetNumberOfMessages() = %v, want %v", metrics.GetNumberOfMessages(), 1)
	}

	metrics.Update(100*time.Millisecond, nil)
	if metrics.GetNumberOfMessages() != 2 {
		t.Errorf("GetNumberOfMessages() = %v, want %v", metrics.GetNumberOfMessages(), 2)
	}
}

func TestMetrics_GetNumberOfErrors(t *testing.T) {
	metrics := NewMetrics("test/topic")

	if metrics.GetNumberOfErrors() != 0 {
		t.Errorf("GetNumberOfErrors() = %v, want %v", metrics.GetNumberOfErrors(), 0)
	}

	metrics.Update(100*time.Millisecond, &testError{msg: "error1"})
	if metrics.GetNumberOfErrors() != 1 {
		t.Errorf("GetNumberOfErrors() = %v, want %v", metrics.GetNumberOfErrors(), 1)
	}

	metrics.Update(100*time.Millisecond, nil)
	if metrics.GetNumberOfErrors() != 1 {
		t.Errorf("GetNumberOfErrors() = %v, want %v", metrics.GetNumberOfErrors(), 1)
	}

	metrics.Update(100*time.Millisecond, &testError{msg: "error2"})
	if metrics.GetNumberOfErrors() != 2 {
		t.Errorf("GetNumberOfErrors() = %v, want %v", metrics.GetNumberOfErrors(), 2)
	}
}

func TestMetrics_Print(t *testing.T) {
	metrics := NewMetrics("test/topic")
	logger := &mockLogger{}

	metrics.Update(100*time.Millisecond, nil)
	metrics.Update(200*time.Millisecond, &testError{msg: "test error"})

	metrics.Print(logger)

	if len(logger.logs) == 0 {
		t.Error("Print() should log metrics")
	}
}

type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}
