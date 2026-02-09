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
	return r.db.WithContext(ctx).Create(job).Error
}

func (r *JobQueueRepository) Update(ctx context.Context, job *models.JobQueue) error {
	return r.db.WithContext(ctx).Save(job).Error
}

func (r *JobQueueRepository) ClaimNextQueueByType(ctx context.Context, jobType enums.JobQueueType) (*models.JobQueue, error) {
	var job *models.JobQueue

	err := r.db.WithContext(ctx).Raw(`
		UPDATE job_queues
		SET status = ?, updated_at = NOW()
		WHERE id = (
				SELECT id
				FROM job_queues
				WHERE type = ?
				AND status = ?
				ORDER BY created_at ASC
				FOR UPDATE SKIP LOCKED
				LIMIT 1
		)
		RETURNING *
	`, enums.JobQueueStatusProcessing,
		jobType,
		enums.JobQueueStatusPending,
	).Scan(&job).Error

	if err != nil {
		return nil, err
	}

	if err := r.db.WithContext(ctx).
		Preload("Deployment").
		First(&job, job.ID).Error; err != nil {
		return nil, err
	}

	return job, nil
}
