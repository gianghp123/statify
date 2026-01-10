package middlewares

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func ErrorLoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		method := c.Request.Method
		path := c.Request.URL.Path

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				log.Printf(
					"[ERROR] method=%s path=%s status=%d latency=%s error=%v",
					method,
					path,
					status,
					latency,
					e.Err,
				)
			}
			return
		}

		log.Printf(
			"[INFO] method=%s path=%s status=%d latency=%s",
			method,
			path,
			status,
			latency,
		)
	}
}
