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

func (r *ProjectRepository) FindAllByUserID(ctx context.Context, userID uint) (*repository.PaginatedEntities[models.Project], error) {
	var projects []models.Project
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&projects).Error; err != nil {
		return nil, err
	}
	return &repository.PaginatedEntities[models.Project]{
		Entities: projects,
		Pagination: repository.Pagination{
			TotalCount: len(projects),
			Page:       1,
			Limit:      10,
		},
	}, nil
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
