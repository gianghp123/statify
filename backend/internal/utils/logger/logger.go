package logger

import (
	"log/slog"
	"os"
	"strings"
	"sync"

	"github.com/gianghp/statify/internal/utils"
)

var (
	logger *slog.Logger
	once   sync.Once
)

// Init initializes the global logger
func Init() {
	once.Do(func() {
		logFormat := strings.ToLower(utils.GetEnv("LOG_FORMAT", ""))
		var handler slog.Handler

		if logFormat == "json" {
			handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelInfo,
			})
		} else {
			handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelInfo,
			})
		}

		logger = slog.New(handler)
		slog.SetDefault(logger)
	})
}

// Get returns the global logger instance
func Get() *slog.Logger {
	if logger == nil {
		Init()
	}
	return logger
}
