package processor

import (
	"context"

	"github.com/gianghp/statify/internal/database/models"
	"github.com/stretchr/testify/mock"
)

type FileProcessorMock struct {
	mock.Mock
}

func (m *FileProcessorMock) ProcessDeploymentFiles(ctx context.Context, deployment *models.Deployment) error {
	args := m.Called(ctx, deployment)
	return args.Error(0)
}

func (m *FileProcessorMock) DeleteMinioFolder(ctx context.Context, prefix string) error {
	args := m.Called(ctx, prefix)
	return args.Error(0)
}
