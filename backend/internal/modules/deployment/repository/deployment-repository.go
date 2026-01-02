package repository

import (
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

func (r *DeploymentRepository) FindByID(id uint) (*models.Deployment, error) {
	return nil, nil
}

func (r *DeploymentRepository) FindAllByProjectID(projectID uint) (*repository.PaginatedEntities[models.Deployment], error) {
	return nil, nil
}

func (r *DeploymentRepository) Create(deployment *models.Deployment) error {
	return nil
}

func (r *DeploymentRepository) Update(deployment *models.Deployment) error {
	return nil
}
