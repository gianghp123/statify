package models

import (
	"github.com/gianghp/statify/internal/core/enums"
	"gorm.io/gorm"
)

type Project struct {
	gorm.Model
	Name                string `gorm:"not null"`
	Subdomain           string `gorm:"not null;unique"`
	UserID              uint
	User                User
	CurrentDeploymentID uint
	EffectiveStatus     enums.DeploymentStatus
	LatestDeploymentID  uint
	Deployments         []Deployment `gorm:"constraint:OnDelete:CASCADE;foreignKey:ProjectID"`
}
