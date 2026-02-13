package middlewares

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggingMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		method := c.Request.Method
		path := c.Request.URL.Path

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				logger.Error(
					"request failed",
					slog.String("method", method),
					slog.String("path", path),
					slog.Int("status", status),
					slog.Duration("latency", latency),
					slog.Any("error", e.Err),
				)
			}
			return
		}

		logger.Info(
			"request completed",
			slog.String("method", method),
			slog.String("path", path),
			slog.Int("status", status),
			slog.Duration("latency", latency),
		)
	}
}
