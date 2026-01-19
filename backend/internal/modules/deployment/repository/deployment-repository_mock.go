package repository

import (
	"context"

	"github.com/gianghp/statify/internal/core/repository"
	"github.com/gianghp/statify/internal/database/models"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type DeploymentRepositoryMock struct {
	mock.Mock
}

func (m *DeploymentRepositoryMock) FindByID(ctx context.Context, id uint) (*models.Deployment, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Deployment), args.Error(1)
}

func (m *DeploymentRepositoryMock) FindAllByUserID(ctx context.Context, userID uint, page int, limit int) (repository.PaginatedEntities[*models.Deployment], error) {
	args := m.Called(ctx, userID, page, limit)
	return args.Get(0).(repository.PaginatedEntities[*models.Deployment]), args.Error(1)
}

func (m *DeploymentRepositoryMock) FindAllByProjectID(ctx context.Context, projectID uint, page int, limit int) (repository.PaginatedEntities[*models.Deployment], error) {
	args := m.Called(ctx, projectID, page, limit)
	return args.Get(0).(repository.PaginatedEntities[*models.Deployment]), args.Error(1)
}

func (m *DeploymentRepositoryMock) FindLatestByProjectID(ctx context.Context, projectID uint, status string) (*models.Deployment, error) {
	args := m.Called(ctx, projectID, status)
	return args.Get(0).(*models.Deployment), args.Error(1)
}

func (m *DeploymentRepositoryMock) Create(ctx context.Context, deployment *models.Deployment) error {
	args := m.Called(ctx, deployment)
	return args.Error(0)
}

func (m *DeploymentRepositoryMock) Update(ctx context.Context, deployment *models.Deployment) error {
	args := m.Called(ctx, deployment)
	return args.Error(0)
}

func (m *DeploymentRepositoryMock) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	args := m.Called(ctx, fn)
	return args.Error(0)
}

func (m *DeploymentRepositoryMock) WithTx(tx *gorm.DB) IDeploymentRepository {
	args := m.Called(tx)
	return args.Get(0).(IDeploymentRepository)
}

func (m *DeploymentRepositoryMock) Delete(ctx context.Context, deployment *models.Deployment) error {
	args := m.Called(ctx, deployment)
	return args.Error(0)
}

func (m *DeploymentRepositoryMock) ClaimNextQueued(ctx context.Context) (*models.Deployment, error) {
	args := m.Called(ctx)
	return args.Get(0).(*models.Deployment), args.Error(1)
}

func (m *DeploymentRepositoryMock) MarkFailed(ctx context.Context, id uint, error string) error {
	args := m.Called(ctx, id, error)
	return args.Error(0)
}

func (m *DeploymentRepositoryMock) MarkReady(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
