package repository

import (
	"context"

	"github.com/gianghp/statify/internal/database/models"
)

type IUploadSessionRepository interface {
	Create(ctx context.Context, uploadSession *models.DeploymentUploadSession) error
	FindByID(ctx context.Context, id uint) (*models.DeploymentUploadSession, error)
	FindUnexpiredByProjectID(ctx context.Context, projectID uint) (*models.DeploymentUploadSession, error)
	Update(ctx context.Context, uploadSession *models.DeploymentUploadSession) error
}
