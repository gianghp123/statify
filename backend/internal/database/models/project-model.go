package models

import "gorm.io/gorm"

type Project struct {
	gorm.Model
	Name                string `gorm:"not null"`
	Subdomain           string `gorm:"not null"`
	UserID              uint
	User                User
	CurrentDeploymentID uint
}
