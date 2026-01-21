package repository

import (
	"context"
	"time"

	"github.com/gianghp/statify/internal/database/models"
	"github.com/gianghp/statify/internal/modules/analytics/repository/projections"
)

type DeploymentMetricsMinuteRepositoryInterface interface {
	Save(ctx context.Context, event *models.DeploymentMetricsMinute) error
	SaveMany(ctx context.Context, events []*models.DeploymentMetricsMinute) error
	GetProjectOverviewMetrics(ctx context.Context, projectID uint, startTime time.Time, endTime time.Time) (*projections.ProjectOverviewMetrics, error)
	GetPerformanceMetrics(ctx context.Context, projectID uint, startTime time.Time, endTime time.Time) (*projections.PerformanceMetrics, error)
	GetTimeSeriesPoints(ctx context.Context, projectID uint, startTime time.Time, endTime time.Time) ([]*projections.TimeSeriesPoint, error)
}
