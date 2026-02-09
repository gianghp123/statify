package repository

import (
	"context"

	"github.com/gianghp/statify/internal/core/enums"
	"github.com/gianghp/statify/internal/database/models"
)

type IJobQueueRepository interface {
	Create(ctx context.Context, job *models.JobQueue) error
	Update(ctx context.Context, job *models.JobQueue) error
	FindLatestByStatus(ctx context.Context, status enums.JobQueueStatus) (*models.JobQueue, error)
}
