package repository

import (
	"context"
	"time"

	"github.com/gianghp/statify/internal/core/sse"
	"github.com/gianghp/statify/internal/database/models"
	"github.com/gianghp/statify/internal/modules/analytics/repository/projections"
	"gorm.io/gorm"
)

type DeploymentMetricsMinuteRepository struct {
	db *gorm.DB
}

func NewDeploymentMetricsMinuteRepository(db *gorm.DB) *DeploymentMetricsMinuteRepository {
	return &DeploymentMetricsMinuteRepository{db: db}
}

func (r *DeploymentMetricsMinuteRepository) Save(ctx context.Context, event *models.DeploymentMetricsMinute) error {
	tx := r.db.WithContext(ctx)
	result := tx.Create(event)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected > 0 {
		notifyErr := tx.Exec(`
			SELECT pg_notify(?, json_build_object(
				'project_id', (?::bigint)
			)::text)`,
			sse.AnalyticsEvent,
			event.ProjectID,
		).Error

		if notifyErr != nil {
			return notifyErr
		}
	}

	return nil
}

func (r *DeploymentMetricsMinuteRepository) SaveMany(ctx context.Context, events []*models.DeploymentMetricsMinute) error {
	return r.db.WithContext(ctx).CreateInBatches(events, 100).Error
}

func (r *DeploymentMetricsMinuteRepository) GetProjectOverviewMetrics(ctx context.Context, projectID uint, startTime time.Time, endTime time.Time) (*projections.ProjectOverviewMetrics, error) {
	tx := r.db.WithContext(ctx)
	var metrics *projections.ProjectOverviewMetrics
	if err := tx.Raw(`
		SELECT
			COALESCE(SUM(request_count), 0)           AS total_requests,
			COALESCE(SUM(bytes_served), 0)            AS total_bandwidth,
			COALESCE(SUM(total_duration), 0)          AS total_duration,
			COALESCE(SUM(status_4xx + status_5xx), 0) AS total_errors
		FROM deployment_metrics_minute
		WHERE project_id = ?
		AND minute_ts BETWEEN ? AND ?;
	`, projectID, startTime, endTime).Scan(&metrics).Error; err != nil {
		return nil, err
	}
	return metrics, nil
}

func (r *DeploymentMetricsMinuteRepository) GetPerformanceMetrics(ctx context.Context, projectID uint, startTime time.Time, endTime time.Time) (*projections.PerformanceMetrics, error) {
	tx := r.db.WithContext(ctx)
	var metrics *projections.PerformanceMetrics
	if err := tx.Raw(`
	SELECT
			CASE
					WHEN SUM(request_count) = 0 THEN 0
					ELSE SUM(total_duration)::float / SUM(request_count)
			END AS avg_response_ms
	FROM deployment_metrics_minute
	WHERE project_id = ?
	AND minute_ts BETWEEN ? AND ?;
`, projectID, startTime, endTime).Scan(&metrics).Error; err != nil {
		return nil, err
	}
	return metrics, nil
}

func (r *DeploymentMetricsMinuteRepository) GetTimeSeriesPoints(ctx context.Context, projectID uint, startTime time.Time, endTime time.Time) ([]*projections.TimeSeriesPoint, error) {
	tx := r.db.WithContext(ctx)
	var metrics []*projections.TimeSeriesPoint
	if err := tx.Raw(`
	SELECT
		minute_ts        	AS timestamp,
		SUM(request_count) AS requests,
		SUM(bytes_served)  AS bandwidth
	FROM deployment_metrics_minute
	WHERE project_id = ?
	AND minute_ts BETWEEN ? AND ?
	GROUP BY minute_ts
	ORDER BY minute_ts;
`, projectID, startTime, endTime).Scan(&metrics).Error; err != nil {
		return nil, err
	}

	return metrics, nil
}
