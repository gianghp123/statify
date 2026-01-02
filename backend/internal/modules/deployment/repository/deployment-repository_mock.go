package repository

import (
	"github.com/gianghp/statify/internal/core/repository"
	"github.com/gianghp/statify/internal/database/models"
	"github.com/stretchr/testify/mock"
)

type DeploymentRepositoryMock struct {
	mock.Mock
}

func (m *DeploymentRepositoryMock) FindByID(id uint) (*models.Deployment, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Deployment), args.Error(1)
}

func (m *DeploymentRepositoryMock) FindAllByProjectID(projectID uint) (*repository.PaginatedEntities[models.Deployment], error) {
	args := m.Called(projectID)
	return args.Get(0).(*repository.PaginatedEntities[models.Deployment]), args.Error(1)
}

func (m *DeploymentRepositoryMock) Create(deployment *models.Deployment) error {
	args := m.Called(deployment)
	return args.Error(0)
}

func (m *DeploymentRepositoryMock) Update(deployment *models.Deployment) error {
	args := m.Called(deployment)
	return args.Error(0)
}
