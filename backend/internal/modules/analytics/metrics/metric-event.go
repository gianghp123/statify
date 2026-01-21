package metrics

import "time"

type MetricEvent struct {
	DeploymentID uint
	ProjectID    uint
	Path         string
	StatusCode   int
	BytesServed  int64
	DurationMs   int64
	Timestamp    time.Time
}
