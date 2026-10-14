// Package logging configures the process-wide slog logger.
package logging

import (
	"log/slog"
	"os"
	"strings"
)

// Setup installs a slog default logger (text or JSON) at the given level and
// returns it. slog.SetDefault also routes the standard `log` package through it,
// so existing log.Printf calls share the same format and destination.
func Setup(format, level string) *slog.Logger {
	opts := &slog.HandlerOptions{Level: parseLevel(level)}

	var handler slog.Handler
	if strings.EqualFold(format, "json") {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
	return logger
}

func parseLevel(level string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(level)) {
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
