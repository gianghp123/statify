package projections

import "time"

type TimeSeriesPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Requests  int64     `json:"requests"`
	Bandwidth int64     `json:"bandwidth"`
}
