package repository

import (
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

func (r *ProjectRepository) FindByID(id uint) (*models.Project, error) {
	return nil, nil
}

func (r *ProjectRepository) FindAllByUserID(userID uint) (*repository.PaginatedEntities[models.Project], error) {
	return nil, nil
}

func (r *ProjectRepository) Create(project *models.Project) error {
	return nil
}

func (r *ProjectRepository) Update(project *models.Project) error {
	return nil
}

func (r *ProjectRepository) Delete(project *models.Project) error {
	return nil
}
