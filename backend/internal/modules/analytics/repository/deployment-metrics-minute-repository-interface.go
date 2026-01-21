package repository

import (
	"context"

	"github.com/gianghp/statify/internal/database/models"
)

type DeploymentMetricsMinuteRepositoryInterface interface {
	Save(ctx context.Context, event *models.DeploymentMetricsMinute) error
	SaveMany(ctx context.Context, events []*models.DeploymentMetricsMinute) error
}
