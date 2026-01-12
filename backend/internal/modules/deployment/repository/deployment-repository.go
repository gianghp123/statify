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
	if err := r.db.WithContext(ctx).First(&deployment, id).Error; err != nil {
		return nil, err
	}
	return &deployment, nil
}

func (r *DeploymentRepository) FindAllByProjectID(ctx context.Context, projectID uint) (*repository.PaginatedEntities[models.Deployment], error) {
	var deployments []models.Deployment
	if err := r.db.WithContext(ctx).Where("project_id = ?", projectID).Find(&deployments).Error; err != nil {
		return nil, err
	}
	return &repository.PaginatedEntities[models.Deployment]{
		Entities: deployments,
		Pagination: repository.Pagination{
			TotalCount: len(deployments),
			Page:       1,
			Limit:      10,
		},
	}, nil
}

func (r *DeploymentRepository) FindLatestByProjectID(ctx context.Context, projectID uint) (*models.Deployment, error) {
	var deployment models.Deployment
	err := r.db.WithContext(ctx).Where("project_id = ?", projectID).Order("created_at DESC").First(&deployment).Error

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
