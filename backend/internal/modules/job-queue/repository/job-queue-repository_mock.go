package repository

import (
	"context"

	"github.com/gianghp/statify/internal/core/enums"
	"github.com/gianghp/statify/internal/database/models"
	"github.com/stretchr/testify/mock"
)

type JobQueueRepositoryMock struct {
	mock.Mock
}

func (m *JobQueueRepositoryMock) Create(ctx context.Context, job *models.JobQueue) error {
	args := m.Called(ctx, job)
	return args.Error(0)
}

func (m *JobQueueRepositoryMock) FindLatestByStatus(ctx context.Context, status enums.JobQueueStatus) (*models.JobQueue, error) {
	args := m.Called(ctx, status)
	return args.Get(0).(*models.JobQueue), args.Error(1)
}

func (m *JobQueueRepositoryMock) Update(ctx context.Context, job *models.JobQueue) error {
	args := m.Called(ctx, job)
	return args.Error(0)
}
