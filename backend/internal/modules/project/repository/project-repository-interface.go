package repository

import (
	"context"

	"github.com/gianghp/statify/internal/core/repository"
	"github.com/gianghp/statify/internal/database/models"
)

type IProjectRepository interface {
	FindByID(ctx context.Context, id uint) (*models.Project, error)
	FindAllByUserID(ctx context.Context, userID uint) (*repository.PaginatedEntities[models.Project], error)
	Create(ctx context.Context, project *models.Project) error
	Update(ctx context.Context, project *models.Project) error
	Delete(ctx context.Context, project *models.Project) error
}
