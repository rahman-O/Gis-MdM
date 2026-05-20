package logger

import (
	"log/slog"
	"os"
)

// New returns the default structured logger for the application.
func New() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}
