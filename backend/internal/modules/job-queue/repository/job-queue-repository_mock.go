package repository

import (
	"context"

	"github.com/gianghp/statify/internal/core/enums"
	"github.com/gianghp/statify/internal/database/models"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type JobQueueRepositoryMock struct {
	mock.Mock
}

func (m *JobQueueRepositoryMock) Create(ctx context.Context, job *models.JobQueue) error {
	args := m.Called(ctx, job)
	return args.Error(0)
}

func (m *JobQueueRepositoryMock) ClaimNextQueueByType(ctx context.Context, jobType enums.JobQueueType) (*models.JobQueue, error) {
	args := m.Called(ctx, jobType)
	return args.Get(0).(*models.JobQueue), args.Error(1)
}

func (m *JobQueueRepositoryMock) Update(ctx context.Context, job *models.JobQueue) error {
	args := m.Called(ctx, job)
	return args.Error(0)
}

func (m *JobQueueRepositoryMock) WithTx(tx *gorm.DB) IJobQueueRepository {
	args := m.Called(tx)
	return args.Get(0).(IJobQueueRepository)
}
