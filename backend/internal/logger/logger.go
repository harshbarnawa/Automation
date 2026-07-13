package logger

import (
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/harshbarnawa/mintok/backend/internal/config"
)

func New(cfg config.Config) *slog.Logger {
	return NewWithWriter(cfg, os.Stdout)
}

func NewWithWriter(cfg config.Config, writer io.Writer) *slog.Logger {
	options := &slog.HandlerOptions{
		Level: parseLevel(cfg.LogLevel),
	}

	if cfg.Environment == "production" {
		return slog.New(slog.NewJSONHandler(writer, options)).With(
			"service", cfg.ServiceName,
			"environment", cfg.Environment,
		)
	}

	return slog.New(slog.NewTextHandler(writer, options)).With(
		"service", cfg.ServiceName,
		"environment", cfg.Environment,
	)
}

func parseLevel(value string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(value)) {
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
