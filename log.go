package log

import (
	"context"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
)

var (
	ReleaseMode    bool
	defaultLogger  atomic.Pointer[slog.Logger]
	logLoggerLevel slog.LevelVar
)

func init() {
	ReleaseMode, _ = strconv.ParseBool(os.Getenv("LOG_PURE"))

	opts := &slog.HandlerOptions{
		Level: &logLoggerLevel,
	}
	level, found := os.LookupEnv("LOG_LEVEL")
	if !found {
		SetLogLoggerLevel(slog.LevelInfo)
		goto INIT
	}

	switch strings.ToLower(level) {
	case "debug":
		SetLogLoggerLevel(slog.LevelDebug)
	case "info":
		SetLogLoggerLevel(slog.LevelInfo)
	case "warn":
		SetLogLoggerLevel(slog.LevelWarn)
	case "error":
		SetLogLoggerLevel(slog.LevelError)
	}

INIT:
	defaultLogger.Store(slog.New(NewHandler(opts, os.Stderr)))
}

func SetLogLoggerLevel(level slog.Level) (oldLevel slog.Level) {
	oldLevel = logLoggerLevel.Level()
	logLoggerLevel.Set(level)
	return
}

func Default() *slog.Logger {
	return defaultLogger.Load()
}

func SetDefault(l *slog.Logger) {
	defaultLogger.Store(l)
}

func Debug(msg string, args ...any) {
	Default().Debug(msg, args...)
}

func DebugContext(ctx context.Context, msg string, args ...any) {
	Default().DebugContext(ctx, msg, args...)
}

func Info(msg string, args ...any) {
	Default().Info(msg, args...)
}

func InfoContext(ctx context.Context, msg string, args ...any) {
	Default().InfoContext(ctx, msg, args...)
}

func Warn(msg string, args ...any) {
	Default().Warn(msg, args...)
}

func WarnContext(ctx context.Context, msg string, args ...any) {
	Default().WarnContext(ctx, msg, args...)
}

func Error(msg string, args ...any) {
	Default().Error(msg, args...)
}

func ErrorContext(ctx context.Context, msg string, args ...any) {
	Default().ErrorContext(ctx, msg, args...)
}

func Log(ctx context.Context, level slog.Level, msg string, args ...any) {
	Default().Log(ctx, level, msg, args...)
}

func LogAttrs(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	Default().LogAttrs(ctx, level, msg, attrs...)
}

func With(args ...any) *slog.Logger {
	return Default().With(args...)
}
