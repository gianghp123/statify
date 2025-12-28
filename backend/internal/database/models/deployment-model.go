package models

import "gorm.io/gorm"

type Deployment struct {
	gorm.Model
	ProjectID          uint
	Project            Project
	Status             string `gorm:"default:'PENDING'"`
	OutputPrefix       string `gorm:"not null"`
	SourceZipObjectKey string `gorm:"not null"`
	ValidationError    string
}
