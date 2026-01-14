package repository

import (
	"context"

	"github.com/gianghp/statify/internal/core/repository"
	"github.com/gianghp/statify/internal/database/models"
	"gorm.io/gorm"
)

type IProjectRepository interface {
	FindBySubdomain(ctx context.Context, subdomain string) (*models.Project, error)
	FindByID(ctx context.Context, id uint) (*models.Project, error)
	FindAllByUserID(ctx context.Context, userID uint, page int, limit int) (*repository.PaginatedEntities[models.Project], error)
	Create(ctx context.Context, project *models.Project) error
	Update(ctx context.Context, project *models.Project) error
	Delete(ctx context.Context, project *models.Project) error
	Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error
	WithTx(tx *gorm.DB) IProjectRepository
}
