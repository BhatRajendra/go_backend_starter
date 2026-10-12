package logging

import (
	"context"
	"log/slog"
	"os"
	"runtime"
	"time"
)

type Logger struct {
	slog *slog.Logger
}

func New() *Logger {
	appEnv := os.Getenv("APP_ENV")

	// text in dev/test and json in prod
	var handler slog.Handler
	if appEnv == "PROD" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelInfo,
			AddSource: false,
		})
	} else {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		})
	}

	slogger := slog.New(handler).With( // With adds those key-value pairs to every log
		"service", "sandhi-backend",
		"environment", appEnv,
	)

	return &Logger{
		slog: slogger,
	}
}

// Add this helper to Logger.
func (l *Logger) log(level slog.Level, msg string, args ...any) {
	ctx := context.Background()

	if !l.slog.Enabled(ctx, level) {
		return
	}

	var pcs [1]uintptr
	runtime.Callers(3, pcs[:])

	record := slog.NewRecord(time.Now(), level, msg, pcs[0])
	record.Add(args...)

	_ = l.slog.Handler().Handle(ctx, record)
	if os.Getenv("APP_ENV") == "DEV" {
		println()
	}
}

func (l *Logger) Debug(msg string, args ...any) {
	l.log(slog.LevelDebug, msg, args...)
}

func (l *Logger) Info(msg string, args ...any) {
	l.log(slog.LevelInfo, msg, args...)
}

func (l *Logger) Warn(msg string, args ...any) {
	l.log(slog.LevelWarn, msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	l.log(slog.LevelError, msg, args...)
}

func (l *Logger) SecurityIssue(msg string, args ...any) {
	pairs := make([]any, 0, len(args)+2)
	pairs = append(pairs, "event_type", "security")
	pairs = append(pairs, args...)

	l.log(slog.LevelWarn, msg, pairs...)
}
