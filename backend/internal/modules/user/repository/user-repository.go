package repository

import (
	"context"

	"github.com/gianghp/statify/internal/core/repository"
	"github.com/gianghp/statify/internal/database/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindAll(ctx context.Context, page int, limit int) (repository.PaginatedEntities[*models.User], error) {
	result := repository.PaginatedEntities[*models.User]{
		Entities: []*models.User{},
		Pagination: repository.Pagination{
			TotalCount: 0,
			Page:       page,
			Limit:      limit,
		},
	}

	db := r.db.WithContext(ctx).Model(&models.User{})

	if err := db.Count(&result.Pagination.TotalCount).Error; err != nil {
		return result, err
	}
	if err := db.Offset((page - 1) * limit).Limit(limit).Find(&result.Entities).Error; err != nil {
		return result, err
	}

	return result, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *UserRepository) Delete(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Delete(user).Error
}
