package logger

import (
	"log/slog"
	"os"
	"strings"
)

// New creates a configured *slog.Logger.
//   - In production it emits JSON to stdout with source location.
//   - In other environments it emits human-readable text.
func New(level, env string) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level:     parseLevel(level),
		AddSource: strings.ToLower(env) == "production",
	}

	var handler slog.Handler
	if strings.ToLower(env) == "production" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}

func parseLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
