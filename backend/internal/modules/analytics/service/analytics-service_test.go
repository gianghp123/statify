package service

import (
	"context"
	"testing"
	"time"

	"github.com/gianghp/statify/internal/modules/analytics/dtos/response"
	"github.com/gianghp/statify/internal/modules/analytics/repository"
	"github.com/gianghp/statify/internal/modules/analytics/repository/projections"
	"github.com/stretchr/testify/assert"
)

func TestAnalyticsService_GetComprehensiveMetrics(t *testing.T) {
	start := time.Now().Add(-time.Hour)
	end := time.Now()

	tests := []struct {
		name         string
		projectID    uint
		startTime    time.Time
		endTime      time.Time
		setupMocks   func(repo *repository.DeploymentMetricsMinuteRepositoryMock)
		expectedFunc func(t *testing.T, metrics response.ComprehensiveMetricsDTO, err error)
	}{
		{
			name:      "Get comprehensive metrics successfully",
			projectID: 1,
			startTime: start,
			endTime:   end,
			setupMocks: func(repo *repository.DeploymentMetricsMinuteRepositoryMock) {
				repo.
					On("GetProjectOverviewMetrics", context.TODO(), uint(1), start, end).
					Return(&projections.ProjectOverviewMetrics{
						TotalRequests:  2,
						TotalBandwidth: 100,
						TotalDuration:  10,
						TotalErrors:    0,
					}, nil)

				repo.
					On("GetPerformanceMetrics", context.TODO(), uint(1), start, end).
					Return(&projections.PerformanceMetrics{
						AvgResponseMs: 10,
					}, nil)

				repo.
					On("GetTimeSeriesPoints", context.TODO(), uint(1), start, end).
					Return([]*projections.TimeSeriesPoint{
						{Timestamp: start, Requests: 2, Bandwidth: 100},
						{Timestamp: start.Add(-45 * time.Minute), Requests: 3, Bandwidth: 150},
						{Timestamp: start.Add(-30 * time.Minute), Requests: 1, Bandwidth: 50},
						{Timestamp: start.Add(-15 * time.Minute), Requests: 2, Bandwidth: 100},
						{Timestamp: end, Requests: 2, Bandwidth: 100},
					}, nil)
			},
			expectedFunc: func(t *testing.T, metrics response.ComprehensiveMetricsDTO, err error) {
				assert.NoError(t, err)

				expected := response.ComprehensiveMetricsDTO{
					PerformanceMetrics: response.PerformanceMetricsDTO{
						AvgResponseMs: 10,
					},
					ProjectOverviewMetrics: response.ProjectOverviewMetricsDTO{
						TotalRequests:    2,
						TotalBandwidth:   float64(100) / (1024 * 1024),
						AvgResponseMs:    float64(10) / 2,
						ErrorRatePercent: float64(0) / 2,
					},
					TimeSeriesPoints: []response.TimeSeriesPointDTO{
						{Timestamp: start, Requests: 2, Bandwidth: float64(100) / (1024 * 1024)},
						{Timestamp: start.Add(-45 * time.Minute), Requests: 3, Bandwidth: float64(150) / (1024 * 1024)},
						{Timestamp: start.Add(-30 * time.Minute), Requests: 1, Bandwidth: float64(50) / (1024 * 1024)},
						{Timestamp: start.Add(-15 * time.Minute), Requests: 2, Bandwidth: float64(100) / (1024 * 1024)},
						{Timestamp: end, Requests: 2, Bandwidth: float64(100) / (1024 * 1024)},
					},
				}

				assert.Equal(t, expected, metrics)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := new(repository.DeploymentMetricsMinuteRepositoryMock)
			test.setupMocks(repo)

			service := NewAnalyticsService(repo)
			metrics, err := service.GetComprehensiveMetrics(
				context.TODO(),
				test.projectID,
				test.startTime,
				test.endTime,
			)

			test.expectedFunc(t, *metrics, err)
			repo.AssertExpectations(t)
		})
	}
}
