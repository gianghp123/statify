package repository

import (
	"context"

	"github.com/gianghp/statify/internal/core/repository"
	"github.com/gianghp/statify/internal/database/models"
	"gorm.io/gorm"
)

type DeploymentRepository struct {
	db *gorm.DB
}

func NewDeploymentRepository(db *gorm.DB) *DeploymentRepository {
	return &DeploymentRepository{db: db}
}

func (r *DeploymentRepository) FindByID(ctx context.Context, id uint) (*models.Deployment, error) {
	var deployment models.Deployment
	if err := r.db.WithContext(ctx).Preload("Project").First(&deployment, id).Error; err != nil {
		return nil, err
	}
	return &deployment, nil
}

func (r *DeploymentRepository) FindAllByUserID(ctx context.Context, userID uint, page int, limit int) (repository.PaginatedEntities[*models.Deployment], error) {
	result := repository.PaginatedEntities[*models.Deployment]{
		Entities: []*models.Deployment{},
		Pagination: repository.Pagination{
			TotalCount: 0,
			Page:       page,
			Limit:      limit,
		},
	}
	db := r.db.WithContext(ctx).Debug().Model(&models.Deployment{}).Preload("Project").Joins("JOIN projects ON projects.id = deployments.project_id").Where("projects.user_id = ?", userID)

	if err := db.Count(&result.Pagination.TotalCount).Error; err != nil {
		return result, err
	}

	if err := db.Order("created_at DESC").Offset((page - 1) * limit).Limit(limit).Find(&result.Entities).Error; err != nil {
		return result, err
	}
	return result, nil
}

func (r *DeploymentRepository) FindAllByProjectID(ctx context.Context, projectID uint, page int, limit int) (repository.PaginatedEntities[*models.Deployment], error) {
	result := repository.PaginatedEntities[*models.Deployment]{
		Entities: []*models.Deployment{},
		Pagination: repository.Pagination{
			TotalCount: 0,
			Page:       page,
			Limit:      limit,
		},
	}
	db := r.db.WithContext(ctx).Model(&models.Deployment{}).Where("project_id = ?", projectID)

	if err := db.Count(&result.Pagination.TotalCount).Error; err != nil {
		return result, err
	}

	if err := db.Order("created_at DESC").Offset((page - 1) * limit).Limit(limit).Find(&result.Entities).Error; err != nil {
		return result, err
	}

	return result, nil
}

func (r *DeploymentRepository) FindLatestByProjectID(ctx context.Context, projectID uint, status string) (*models.Deployment, error) {
	var deployment models.Deployment
	db := r.db.WithContext(ctx).Where("project_id = ?", projectID)
	if status != "" {
		db = db.Where("status = ?", status)
	}
	err := db.Order("created_at DESC").First(&deployment).Error

	if err != nil {
		return nil, err
	}
	return &deployment, nil
}

func (r *DeploymentRepository) Create(ctx context.Context, deployment *models.Deployment) error {
	return r.db.WithContext(ctx).Create(deployment).Error
}

func (r *DeploymentRepository) Update(ctx context.Context, deployment *models.Deployment) error {
	return r.db.WithContext(ctx).Save(deployment).Error
}

func (r *DeploymentRepository) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return r.db.WithContext(ctx).Transaction(fn)
}

func (r *DeploymentRepository) WithTx(tx *gorm.DB) IDeploymentRepository {
	if tx == nil {
		return r
	}
	return &DeploymentRepository{db: tx}
}
