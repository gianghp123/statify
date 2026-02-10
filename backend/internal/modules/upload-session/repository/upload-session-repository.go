package repository

import (
	"context"
	"time"

	"github.com/gianghp/statify/internal/database/models"
	"gorm.io/gorm"
)

type UploadSessionRepository struct {
	db *gorm.DB
}

func NewUploadSessionRepository(db *gorm.DB) *UploadSessionRepository {
	return &UploadSessionRepository{db: db}
}

func (r *UploadSessionRepository) Create(ctx context.Context, uploadSession *models.DeploymentUploadSession) error {
	return r.db.WithContext(ctx).Create(uploadSession).Error
}

func (r *UploadSessionRepository) FindUnexpiredByProjectID(ctx context.Context, projectID uint) (*models.DeploymentUploadSession, error) {
	var uploadSession models.DeploymentUploadSession
	if err := r.db.WithContext(ctx).Where("project_id = ? AND expires_at > ?", projectID, (time.Now().UTC().Add(time.Minute * 2))).First(&uploadSession).Error; err != nil {
		return nil, err
	}
	return &uploadSession, nil
}

func (r *UploadSessionRepository) Update(ctx context.Context, uploadSession *models.DeploymentUploadSession) error {
	return r.db.WithContext(ctx).Save(uploadSession).Error
}
