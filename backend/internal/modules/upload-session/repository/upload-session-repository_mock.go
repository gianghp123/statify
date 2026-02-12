package repository

import (
	"context"

	"github.com/gianghp/statify/internal/database/models"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type UploadSessionRepositoryMock struct {
	mock.Mock
}

func (m *UploadSessionRepositoryMock) Create(ctx context.Context, uploadSession *models.DeploymentUploadSession) error {
	args := m.Called(ctx, uploadSession)
	return args.Error(0)
}

func (m *UploadSessionRepositoryMock) FindByID(ctx context.Context, id uint) (*models.DeploymentUploadSession, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.DeploymentUploadSession), args.Error(1)
}

func (m *UploadSessionRepositoryMock) FindUnexpiredByProjectID(ctx context.Context, projectID uint) (*models.DeploymentUploadSession, error) {
	args := m.Called(ctx, projectID)
	return args.Get(0).(*models.DeploymentUploadSession), args.Error(1)
}

func (m *UploadSessionRepositoryMock) Update(ctx context.Context, uploadSession *models.DeploymentUploadSession) error {
	args := m.Called(ctx, uploadSession)
	return args.Error(0)
}

func (m *UploadSessionRepositoryMock) WithTx(tx *gorm.DB) IUploadSessionRepository {
	args := m.Called(tx)
	return args.Get(0).(IUploadSessionRepository)
}
