package logger

import (
	"log/slog"
	"os"
)

// New an-chior 初始化并返回一个配置好的 slog.Logger
func New(level slog.Level) *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})
	logger := slog.New(handler)
	return logger
}
