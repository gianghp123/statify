package repository

import (
	"context"

	"github.com/gianghp/statify/internal/database/models"
	"gorm.io/gorm"
)

type DeploymentMetricsMinuteRepository struct {
	db *gorm.DB
}

func NewDeploymentMetricsMinuteRepository(db *gorm.DB) *DeploymentMetricsMinuteRepository {
	return &DeploymentMetricsMinuteRepository{db: db}
}

func (r *DeploymentMetricsMinuteRepository) Save(ctx context.Context, event *models.DeploymentMetricsMinute) error {
	return r.db.WithContext(ctx).Create(event).Error
}

func (r *DeploymentMetricsMinuteRepository) SaveMany(ctx context.Context, events []*models.DeploymentMetricsMinute) error {
	return r.db.WithContext(ctx).CreateInBatches(events, 100).Error
}
