package mqtt

import "github.com/RicardoCenci/iot-distributed-architecture/shared/logger"

type mockLogger struct {
	logs []string
}

func (m *mockLogger) Debug(msg string, args ...any) { m.logs = append(m.logs, "debug: "+msg) }
func (m *mockLogger) Info(msg string, args ...any)  { m.logs = append(m.logs, "info: "+msg) }
func (m *mockLogger) Warn(msg string, args ...any)  { m.logs = append(m.logs, "warn: "+msg) }
func (m *mockLogger) Error(msg string, args ...any) { m.logs = append(m.logs, "error: "+msg) }

var _ logger.Interface = (*mockLogger)(nil)
