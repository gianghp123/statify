package repository

import (
	"github.com/gianghp/statify/internal/core/repository"
	"github.com/gianghp/statify/internal/database/models"
	"github.com/stretchr/testify/mock"
)

type ProjectRepositoryMock struct {
	mock.Mock
}

func (m *ProjectRepositoryMock) FindByID(id uint) (*models.Project, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Project), args.Error(1)
}

func (m *ProjectRepositoryMock) FindAllByUserID(userID uint) (*repository.PaginatedEntities[models.Project], error) {
	args := m.Called(userID)
	return args.Get(0).(*repository.PaginatedEntities[models.Project]), args.Error(1)
}

func (m *ProjectRepositoryMock) Create(project *models.Project) error {
	args := m.Called(project)
	return args.Error(0)
}

func (m *ProjectRepositoryMock) Update(project *models.Project) error {
	args := m.Called(project)
	return args.Error(0)
}

func (m *ProjectRepositoryMock) Delete(project *models.Project) error {
	args := m.Called(project)
	return args.Error(0)
}
