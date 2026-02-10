package models

import (
	"time"

	"gorm.io/gorm"
)

type DeploymentUploadSession struct {
	gorm.Model
	ProjectID    uint
	UploadKey    string
	OutputPrefix string
	PresignedUrl string
	ExpiredAt    time.Time
}
