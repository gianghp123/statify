package repository

import (
	"context"
	"time"

	"github.com/gianghp/statify/internal/database/models"
	"github.com/gianghp/statify/internal/modules/analytics/repository/projections"
	"github.com/stretchr/testify/mock"
)

type DeploymentMetricsMinuteRepositoryMock struct {
	mock.Mock
}

func (m *DeploymentMetricsMinuteRepositoryMock) Save(ctx context.Context, event *models.DeploymentMetricsMinute) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *DeploymentMetricsMinuteRepositoryMock) SaveMany(ctx context.Context, events []*models.DeploymentMetricsMinute) error {
	args := m.Called(ctx, events)
	return args.Error(0)
}

func (m *DeploymentMetricsMinuteRepositoryMock) GetProjectOverviewMetrics(ctx context.Context, projectID uint, startTime time.Time, endTime time.Time) (*projections.ProjectOverviewMetrics, error) {
	args := m.Called(ctx, projectID, startTime, endTime)
	return args.Get(0).(*projections.ProjectOverviewMetrics), args.Error(1)
}

func (m *DeploymentMetricsMinuteRepositoryMock) GetPerformanceMetrics(ctx context.Context, projectID uint, startTime time.Time, endTime time.Time) (*projections.PerformanceMetrics, error) {
	args := m.Called(ctx, projectID, startTime, endTime)
	return args.Get(0).(*projections.PerformanceMetrics), args.Error(1)
}

func (m *DeploymentMetricsMinuteRepositoryMock) GetTimeSeriesPoints(ctx context.Context, projectID uint, startTime time.Time, endTime time.Time) ([]*projections.TimeSeriesPoint, error) {
	args := m.Called(ctx, projectID, startTime, endTime)
	return args.Get(0).([]*projections.TimeSeriesPoint), args.Error(1)
}
