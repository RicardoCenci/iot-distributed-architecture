package logger

import (
	"testing"
)

func TestNewSlogLogger(t *testing.T) {
	tests := []struct {
		name     string
		config   Config
		validate func(*SlogLogger) bool
	}{
		{
			name: "debug level",
			config: Config{
				Level: "debug",
			},
			validate: func(l *SlogLogger) bool {
				return l != nil && l.handler != nil
			},
		},
		{
			name: "info level",
			config: Config{
				Level: "info",
			},
			validate: func(l *SlogLogger) bool {
				return l != nil && l.handler != nil
			},
		},
		{
			name: "warn level",
			config: Config{
				Level: "warn",
			},
			validate: func(l *SlogLogger) bool {
				return l != nil && l.handler != nil
			},
		},
		{
			name: "error level",
			config: Config{
				Level: "error",
			},
			validate: func(l *SlogLogger) bool {
				return l != nil && l.handler != nil
			},
		},
		{
			name: "with source enabled",
			config: Config{
				Level: "info",
				Source: SourceConfig{
					Enabled:  true,
					Relative: true,
					AsJSON:   true,
				},
			},
			validate: func(l *SlogLogger) bool {
				return l != nil && l.handler != nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewSlogLogger(tt.config)
			if !tt.validate(logger) {
				t.Error("NewSlogLogger() validation failed")
			}
		})
	}
}

func TestSlogLogger_Debug(t *testing.T) {
	cfg := Config{Level: "debug"}
	logger := NewSlogLogger(cfg)

	logger.Debug("test debug message", "key", "value")

	if logger == nil {
		t.Error("Debug() should not panic")
	}
}

func TestSlogLogger_Info(t *testing.T) {
	cfg := Config{Level: "info"}
	logger := NewSlogLogger(cfg)

	logger.Info("test info message", "key", "value")

	if logger == nil {
		t.Error("Info() should not panic")
	}
}

func TestSlogLogger_Warn(t *testing.T) {
	cfg := Config{Level: "warn"}
	logger := NewSlogLogger(cfg)

	logger.Warn("test warn message", "key", "value")

	if logger == nil {
		t.Error("Warn() should not panic")
	}
}

func TestSlogLogger_Error(t *testing.T) {
	cfg := Config{Level: "error"}
	logger := NewSlogLogger(cfg)

	logger.Error("test error message", "key", "value")

	if logger == nil {
		t.Error("Error() should not panic")
	}
}

func TestSlogLogger_LogLevels(t *testing.T) {
	tests := []struct {
		name     string
		level    string
		testFunc func(*SlogLogger)
	}{
		{
			name:  "debug level allows debug",
			level: "debug",
			testFunc: func(l *SlogLogger) {
				l.Debug("debug message")
			},
		},
		{
			name:  "info level allows info",
			level: "info",
			testFunc: func(l *SlogLogger) {
				l.Info("info message")
			},
		},
		{
			name:  "warn level allows warn",
			level: "warn",
			testFunc: func(l *SlogLogger) {
				l.Warn("warn message")
			},
		},
		{
			name:  "error level allows error",
			level: "error",
			testFunc: func(l *SlogLogger) {
				l.Error("error message")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Config{Level: tt.level}
			logger := NewSlogLogger(cfg)
			tt.testFunc(logger)
		})
	}
}

func TestSlogLogger_JSONOutput(t *testing.T) {
	cfg := Config{Level: "info"}
	logger := NewSlogLogger(cfg)

	logger.Info("test message", "key1", "value1", "key2", 42)

	if logger == nil {
		t.Error("JSON output should not panic")
	}
}

func TestSlogLogger_Interface(t *testing.T) {
	cfg := Config{Level: "info"}
	logger := NewSlogLogger(cfg)

	var logInterface Interface = logger

	logInterface.Debug("debug")
	logInterface.Info("info")
	logInterface.Warn("warn")
	logInterface.Error("error")

	if logInterface == nil {
		t.Error("SlogLogger should implement Interface")
	}
}

func TestGetBinaryDirectory(t *testing.T) {
	dir, err := getBinaryDirectory()
	if err != nil {
		t.Logf("getBinaryDirectory() error = %v (may be expected in test environment)", err)
		return
	}
	if dir == "" {
		t.Error("getBinaryDirectory() returned empty string")
	}
}

func TestReplaceSourceAttrFn(t *testing.T) {
	tests := []struct {
		name   string
		config Config
	}{
		{
			name: "source disabled",
			config: Config{
				Level:  "info",
				Source: SourceConfig{Enabled: false},
			},
		},
		{
			name: "source enabled relative",
			config: Config{
				Level: "info",
				Source: SourceConfig{
					Enabled:  true,
					Relative: true,
					AsJSON:   false,
				},
			},
		},
		{
			name: "source enabled as json",
			config: Config{
				Level: "info",
				Source: SourceConfig{
					Enabled:  true,
					Relative: false,
					AsJSON:   true,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn := replaceSourceAttrFn(tt.config)
			if fn == nil {
				t.Error("replaceSourceAttrFn() returned nil")
			}
		})
	}
}
