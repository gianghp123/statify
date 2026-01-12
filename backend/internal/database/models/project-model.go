package models

import "gorm.io/gorm"

type Project struct {
	gorm.Model
	Name                string `gorm:"not null"`
	Subdomain           string `gorm:"not null;unique"`
	UserID              uint
	User                User
	CurrentDeploymentID uint
	Deployments         []Deployment `gorm:"constraint:OnDelete:CASCADE;foreignKey:ProjectID"`
}
