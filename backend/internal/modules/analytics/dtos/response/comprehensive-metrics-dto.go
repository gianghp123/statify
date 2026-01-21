package response

import (
	"time"
)

type ComprehensiveMetricsDTO struct {
	PerformanceMetrics     PerformanceMetricsDTO     `json:"performance_metrics" `
	ProjectOverviewMetrics ProjectOverviewMetricsDTO `json:"project_overview_metrics"`
	TimeSeriesPoints       []TimeSeriesPointDTO      `json:"time_series_points"`
}

type PerformanceMetricsDTO struct {
	AvgResponseMs float64 `json:"avg_response_ms"`
}

type ProjectOverviewMetricsDTO struct {
	TotalRequests    int64   `json:"total_requests"`
	TotalBandwidth   float64 `json:"total_bandwidth"`
	AvgResponseMs    float64 `json:"avg_response_ms"`
	ErrorRatePercent float64 `json:"error_rate_percent"`
}

type TimeSeriesPointDTO struct {
	Timestamp time.Time `json:"timestamp"`
	Requests  int64     `json:"requests"`
	Bandwidth float64   `json:"bandwidth"`
}
