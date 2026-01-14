package repository

import (
	"context"

	"github.com/gianghp/statify/internal/core/repository"
	"github.com/gianghp/statify/internal/database/models"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type ProjectRepositoryMock struct {
	mock.Mock
}

func (m *ProjectRepositoryMock) FindByID(ctx context.Context, id uint) (*models.Project, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Project), args.Error(1)
}

func (m *ProjectRepositoryMock) FindAllByUserID(ctx context.Context, userID uint, page int, limit int) (repository.PaginatedEntities[*models.Project], error) {
	args := m.Called(ctx, userID, page, limit)
	return args.Get(0).(repository.PaginatedEntities[*models.Project]), args.Error(1)
}

func (m *ProjectRepositoryMock) FindBySubdomain(ctx context.Context, subdomain string) (*models.Project, error) {
	args := m.Called(ctx, subdomain)
	return args.Get(0).(*models.Project), args.Error(1)
}

func (m *ProjectRepositoryMock) Create(ctx context.Context, project *models.Project) error {
	args := m.Called(ctx, project)
	return args.Error(0)
}

func (m *ProjectRepositoryMock) Update(ctx context.Context, project *models.Project) error {
	args := m.Called(ctx, project)
	return args.Error(0)
}

func (m *ProjectRepositoryMock) Delete(ctx context.Context, project *models.Project) error {
	args := m.Called(ctx, project)
	return args.Error(0)
}

func (m *ProjectRepositoryMock) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	args := m.Called(ctx, fn)
	return args.Error(0)
}

func (m *ProjectRepositoryMock) WithTx(tx *gorm.DB) IProjectRepository {
	args := m.Called(tx)
	return args.Get(0).(IProjectRepository)
}
