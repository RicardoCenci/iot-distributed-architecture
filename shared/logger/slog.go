package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type SlogLogger struct {
	handler slog.Handler
	context []any
}

func NewSlogLogger(config Config, context ...any) *SlogLogger {
	replaceAttr := func(groups []string, a slog.Attr) slog.Attr { return a }

	if config.Source.Enabled && (config.Source.Relative || !config.Source.AsJSON) {
		replaceAttr = replaceSourceAttrFn(config)
	}

	handlerOption := &slog.HandlerOptions{
		Level:       slog.LevelInfo,
		AddSource:   config.Source.Enabled,
		ReplaceAttr: replaceAttr,
	}

	switch config.Level {
	case "debug":
		handlerOption.Level = slog.LevelDebug
	case "info":
		handlerOption.Level = slog.LevelInfo
	case "warn":
		handlerOption.Level = slog.LevelWarn
	case "error":
		handlerOption.Level = slog.LevelError
	}

	return &SlogLogger{
		handler: slog.NewJSONHandler(os.Stdout, handlerOption),
		context: context,
	}
}

func replaceSourceAttrFn(config Config) func(groups []string, a slog.Attr) slog.Attr {
	binDir, err := getBinaryDirectory()

	if err != nil {
		return func(groups []string, a slog.Attr) slog.Attr {
			return a
		}
	}

	return func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.SourceKey && binDir != "" {
			if src, ok := a.Value.Any().(*slog.Source); ok && src != nil {

				if config.Source.Relative {
					if rel, err := filepath.Rel(binDir, src.File); err == nil {
						src.File = rel
					}
				}

				if !config.Source.AsJSON {
					a.Value = slog.StringValue(
						fmt.Sprintf("%s@%s:%d", src.File, src.Function, src.Line),
					)
				}

			}
		}
		return a
	}
}

func getBinaryDirectory() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}

	exe, err = filepath.EvalSymlinks(exe)

	if err != nil {
		return "", err
	}

	dir := filepath.Dir(exe)

	if strings.Contains(dir, string(os.PathSeparator)+"go-build") || strings.HasPrefix(dir, os.TempDir()) {
		if wd, err := os.Getwd(); err == nil {
			return wd, nil
		}
	}
	return dir, nil
}

func (l *SlogLogger) log(level slog.Level, msg string, args ...any) {
	ctx := context.Background()
	if !l.handler.Enabled(ctx, level) {
		return
	}
	pc, _, _, ok := runtime.Caller(2)
	if !ok {
		pc = 0
	}
	rec := slog.NewRecord(time.Now(), level, msg, pc)
	rec.Add(args...)
	_ = l.handler.Handle(ctx, rec)
}

func (l *SlogLogger) Debug(msg string, args ...any) {
	l.log(slog.LevelDebug, msg, append(l.context, args...)...)
}

func (l *SlogLogger) Info(msg string, args ...any) {
	l.log(slog.LevelInfo, msg, append(l.context, args...)...)
}

func (l *SlogLogger) Warn(msg string, args ...any) {
	l.log(slog.LevelWarn, msg, append(l.context, args...)...)
}

func (l *SlogLogger) Error(msg string, args ...any) {
	l.log(slog.LevelError, msg, append(l.context, args...)...)
}

func (l *SlogLogger) WithContext(args ...any) *SlogLogger {
	return &SlogLogger{
		handler: l.handler,
		context: append(l.context, args...),
	}
}
