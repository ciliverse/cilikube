package logger

import (
	"log/slog"
	"os"
)

// New initializes and returns a configured slog.Logger
func New(level slog.Level) *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})
	logger := slog.New(handler)
	return logger
}
