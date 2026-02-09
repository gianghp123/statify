package repository

import (
	"context"

	"github.com/gianghp/statify/internal/core/enums"
	"github.com/gianghp/statify/internal/database/models"
	"gorm.io/gorm"
)

type JobQueueRepository struct {
	db *gorm.DB
}

func NewJobQueueRepository(db *gorm.DB) *JobQueueRepository {
	return &JobQueueRepository{db: db}
}

func (r *JobQueueRepository) Create(ctx context.Context, job *models.JobQueue) error {
	return nil
}

func (r *JobQueueRepository) Update(ctx context.Context, job *models.JobQueue) error {
	return nil
}

func (r *JobQueueRepository) FindLatestByStatus(ctx context.Context, status enums.JobQueueStatus) (*models.JobQueue, error) {
	return nil, nil
}
