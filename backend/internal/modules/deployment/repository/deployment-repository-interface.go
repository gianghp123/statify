package repository

import (
	"context"

	"github.com/gianghp/statify/internal/core/repository"
	"github.com/gianghp/statify/internal/database/models"
)

type IDeploymentRepository interface {
	FindByID(ctx context.Context, id uint) (*models.Deployment, error)
	FindAllByProjectID(ctx context.Context, projectID uint) (*repository.PaginatedEntities[models.Deployment], error)
	FindLatestByProjectID(ctx context.Context, projectID uint) (*models.Deployment, error)
	Create(ctx context.Context, deployment *models.Deployment) error
	Update(ctx context.Context, deployment *models.Deployment) error
}
