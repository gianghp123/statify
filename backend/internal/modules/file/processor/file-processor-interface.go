package processor

import (
	"context"

	"github.com/gianghp/statify/internal/database/models"
)

type IFileProcessor interface {
	ProcessDeploymentFiles(ctx context.Context, deployment *models.Deployment) error
	DeleteMinioFolder(ctx context.Context, prefix string) error
}
