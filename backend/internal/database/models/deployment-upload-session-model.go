package models

import (
	"time"

	"github.com/gianghp/statify/internal/core/enums"
	"gorm.io/gorm"
)

type DeploymentUploadSession struct {
	gorm.Model
	ProjectID    uint
	UploadKey    string
	OutputPrefix string
	PresignedUrl string
	Status       enums.UploadStatus
	ExpiredAt    time.Time
}
