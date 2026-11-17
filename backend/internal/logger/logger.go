package logger

import (
	"context"
	"log/slog"
	"os"
)

type Logger struct {
	*slog.Logger
}

func New(level string) *Logger {
	var lvl slog.Level
	switch level {
	case "debug":
		lvl = slog.LevelDebug
	case "info":
		lvl = slog.LevelInfo
	case "warn":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: lvl,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{Key: "timestamp", Value: a.Value}
			}
			return a
		},
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	return &Logger{slog.New(handler)}
}

func (l *Logger) WithRequestID(requestID string) *Logger {
	return &Logger{l.Logger.With("request_id", requestID)}
}

func (l *Logger) WithUser(ctx context.Context) *Logger {
	if ctx == nil {
		return l
	}
	var args []any
	if uid, ok := ctx.Value("user_id").(string); ok && uid != "" {
		args = append(args, "user_id", uid)
	}
	if tid, ok := ctx.Value("tenant_id").(string); ok && tid != "" {
		args = append(args, "tenant_id", tid)
	}
	if len(args) == 0 {
		return l
	}
	return &Logger{l.Logger.With(args...)}
}

func (l *Logger) WithError(err error) *Logger {
	return &Logger{l.Logger.With("error", err.Error())}
}