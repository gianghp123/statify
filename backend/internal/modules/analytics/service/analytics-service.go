package service

import (
	"context"
	"time"

	"github.com/gianghp/statify/internal/core"
	"github.com/gianghp/statify/internal/modules/analytics/dtos/response"
	"github.com/gianghp/statify/internal/modules/analytics/repository"
	"github.com/gianghp/statify/internal/utils"
	"github.com/jinzhu/copier"
)

type AnalyticsService struct {
	deploymentMetricsMinuteRepository repository.DeploymentMetricsMinuteRepositoryInterface
}

func NewAnalyticsService(deploymentMetricsMinuteRepository repository.DeploymentMetricsMinuteRepositoryInterface) *AnalyticsService {
	return &AnalyticsService{
		deploymentMetricsMinuteRepository: deploymentMetricsMinuteRepository,
	}
}

func (s *AnalyticsService) GetComprehensiveMetrics(ctx context.Context, projectID uint, startTime time.Time, endTime time.Time) (*response.ComprehensiveMetricsDTO, error) {
	performanceMetrics, err := s.deploymentMetricsMinuteRepository.GetPerformanceMetrics(ctx, projectID, startTime, endTime)
	if err != nil {
		return nil, core.ParseDatabaseError(err)
	}

	projectOverviewMetrics, err := s.deploymentMetricsMinuteRepository.GetProjectOverviewMetrics(ctx, projectID, startTime, endTime)
	if err != nil {
		return nil, core.ParseDatabaseError(err)
	}

	timeSeriesPoints, err := s.deploymentMetricsMinuteRepository.GetTimeSeriesPoints(ctx, projectID, startTime, endTime)
	if err != nil {
		return nil, core.ParseDatabaseError(err)
	}

	performanceMetricsDto, err := utils.EntityToDto[response.PerformanceMetricsDTO](performanceMetrics, copier.Option{IgnoreEmpty: false})
	if err != nil {
		return nil, core.InternalError(err.Error())
	}

	projectOverviewMetricsDto, err := utils.EntityToDto[response.ProjectOverviewMetricsDTO](projectOverviewMetrics, copier.Option{IgnoreEmpty: false})
	if err != nil {
		return nil, core.InternalError(err.Error())
	}
	projectOverviewMetricsDto.AvgResponseMs = float64(projectOverviewMetrics.TotalDuration) / float64(projectOverviewMetrics.TotalRequests)
	projectOverviewMetricsDto.ErrorRatePercent = float64(projectOverviewMetrics.TotalErrors) / float64(projectOverviewMetrics.TotalRequests)
	projectOverviewMetricsDto.TotalBandwidth = float64(projectOverviewMetrics.TotalBandwidth) / (1024 * 1024)

	timeSeriesPointsDto, err := utils.EntitiesToDto[response.TimeSeriesPointDTO](timeSeriesPoints, copier.Option{IgnoreEmpty: false})
	if err != nil {
		return nil, core.InternalError(err.Error())
	}
	for i := range timeSeriesPointsDto {
		timeSeriesPointsDto[i].Bandwidth = float64(timeSeriesPointsDto[i].Bandwidth) / (1024 * 1024)
	}
	return &response.ComprehensiveMetricsDTO{
		PerformanceMetrics:     performanceMetricsDto,
		ProjectOverviewMetrics: projectOverviewMetricsDto,
		TimeSeriesPoints:       timeSeriesPointsDto,
	}, nil
}
