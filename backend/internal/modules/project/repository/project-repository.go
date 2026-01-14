package repository

import (
	"context"

	"github.com/gianghp/statify/internal/core/repository"
	"github.com/gianghp/statify/internal/database/models"
	"gorm.io/gorm"
)

type ProjectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) FindBySubdomain(ctx context.Context, subdomain string) (*models.Project, error) {
	var project models.Project
	if err := r.db.WithContext(ctx).First(&project, "subdomain = ?", subdomain).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *ProjectRepository) FindByID(ctx context.Context, id uint) (*models.Project, error) {
	var project models.Project
	if err := r.db.WithContext(ctx).First(&project, id).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *ProjectRepository) FindAllByUserID(ctx context.Context, userID uint, page int, limit int) (repository.PaginatedEntities[*models.Project], error) {
	result := repository.PaginatedEntities[*models.Project]{
		Entities: []*models.Project{},
		Pagination: repository.Pagination{
			TotalCount: 0,
			Page:       page,
			Limit:      limit,
		},
	}

	db := r.db.WithContext(ctx).Model(&models.Project{}).Where("user_id = ?", userID)

	if err := db.Count(&result.Pagination.TotalCount).Error; err != nil {
		return result, err
	}

	if err := db.Offset((page - 1) * limit).Limit(limit).Find(&result.Entities).Error; err != nil {
		return result, err
	}

	return result, nil
}

func (r *ProjectRepository) Create(ctx context.Context, project *models.Project) error {
	return r.db.WithContext(ctx).Create(project).Error
}

func (r *ProjectRepository) Update(ctx context.Context, project *models.Project) error {
	return r.db.WithContext(ctx).Save(project).Error
}

func (r *ProjectRepository) Delete(ctx context.Context, project *models.Project) error {
	return r.db.WithContext(ctx).Unscoped().Delete(project).Error
}

func (r *ProjectRepository) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return r.db.WithContext(ctx).Transaction(fn)
}

func (r *ProjectRepository) WithTx(tx *gorm.DB) IProjectRepository {
	if tx == nil {
		return r
	}
	return &ProjectRepository{db: tx}
}
