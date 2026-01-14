package repository

import (
	"context"

	"github.com/gianghp/statify/internal/core/repository"
	"github.com/gianghp/statify/internal/database/models"
)

type IUserRepository interface {
	FindByID(ctx context.Context, id uint) (*models.User, error)
	FindAll(ctx context.Context, page int, limit int) (repository.PaginatedEntities[*models.User], error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByUsername(ctx context.Context, username string) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, user *models.User) error
}
