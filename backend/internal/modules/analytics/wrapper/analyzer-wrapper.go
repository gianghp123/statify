package wrapper

import (
	"time"

	"github.com/gianghp/statify/internal/modules/analytics/metrics"
	"github.com/gin-gonic/gin"
)

type AnalyzerWrapper struct {
	metricCollector *metrics.MetricCollector
}

func NewAnalyzerWrapper(metricCollector *metrics.MetricCollector) *AnalyzerWrapper {
	return &AnalyzerWrapper{
		metricCollector: metricCollector,
	}
}

func (w *AnalyzerWrapper) Wrap(handler func(c *gin.Context) (*StaticServeResult, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		result, err := handler(c)
		duration := time.Since(start)
		if result != nil {
			w.metricCollector.EmitMetric(metrics.MetricEvent{
				DeploymentID: result.DeploymentID,
				ProjectID:    result.ProjectID,
				Path:         c.Request.URL.Path,
				StatusCode:   result.StatusCode,
				BytesServed:  result.BytesServed,
				DurationMs:   duration.Milliseconds(),
				Timestamp:    time.Now(),
			})
		}

		_ = err
	}
}
