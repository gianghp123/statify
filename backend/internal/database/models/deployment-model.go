package models

import (
	"github.com/gianghp/statify/internal/core/enums"
	"gorm.io/gorm"
)

type Deployment struct {
	gorm.Model
	ProjectID       uint
	Project         Project
	Status          enums.DeploymentStatus `gorm:"default:'UPLOADED';type:deployment_status"`
	OutputPrefix    string                 `gorm:"not null"`
	ValidationError string
}
