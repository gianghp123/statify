package repository

import (
	"context"

	"github.com/gianghp/statify/internal/core/enums"
	"github.com/gianghp/statify/internal/database/models"
	"gorm.io/gorm"
)

type IJobQueueRepository interface {
	Create(ctx context.Context, job *models.JobQueue) error
	Update(ctx context.Context, job *models.JobQueue) error
	ClaimNextQueueByType(ctx context.Context, jobType enums.JobQueueType) (*models.JobQueue, error)
	WithTx(tx *gorm.DB) IJobQueueRepository
}
