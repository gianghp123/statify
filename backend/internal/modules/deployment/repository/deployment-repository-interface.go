package repository

import (
	"context"

	"github.com/gianghp/statify/internal/core/repository"
	"github.com/gianghp/statify/internal/database/models"
	"gorm.io/gorm"
)

type IDeploymentRepository interface {
	FindByID(ctx context.Context, id uint) (*models.Deployment, error)
	FindAllByUserID(ctx context.Context, userID uint, page int, limit int) (repository.PaginatedEntities[*models.Deployment], error)
	FindAllByProjectID(ctx context.Context, projectID uint, page int, limit int) (repository.PaginatedEntities[*models.Deployment], error)
	FindLatestByProjectID(ctx context.Context, projectID uint, status string) (*models.Deployment, error)
	Create(ctx context.Context, deployment *models.Deployment) error
	Update(ctx context.Context, deployment *models.Deployment) error
	Delete(ctx context.Context, deployment *models.Deployment) error
	Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error
	WithTx(tx *gorm.DB) IDeploymentRepository
}
