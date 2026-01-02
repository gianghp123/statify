package repository

import (
	"github.com/gianghp/statify/internal/core/repository"
	"github.com/gianghp/statify/internal/database/models"
)

type IProjectRepository interface {
	FindByID(id uint) (*models.Project, error)
	FindAllByUserID(userID uint) (*repository.PaginatedEntities[models.Project], error)
	Create(project *models.Project) error
	Update(project *models.Project) error
	Delete(project *models.Project) error
}
