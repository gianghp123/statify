package projections

type ProjectOverviewMetrics struct {
	TotalRequests  int64 `json:"total_requests"`
	TotalBandwidth int64 `json:"total_bandwidth"`
	TotalDuration  int64 `json:"total_duration"`
	TotalErrors    int64 `json:"total_errors"`
}
