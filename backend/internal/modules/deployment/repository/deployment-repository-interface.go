package repository

import (
	"github.com/gianghp/statify/internal/core/repository"
	"github.com/gianghp/statify/internal/database/models"
)

type IDeploymentRepository interface {
	FindByID(id uint) (*models.Deployment, error)
	FindAllByProjectID(projectID uint) (*repository.PaginatedEntities[models.Deployment], error)
	Create(deployment *models.Deployment) error
	Update(deployment *models.Deployment) error
}
