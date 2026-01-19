package models

import (
	"time"

	"gorm.io/gorm"
)

type AnalyticsEvent struct {
	gorm.Model
	ProjectID   uint
	Project     Project
	Timestamp   time.Time
	Path        string
	VisitorHash string
	LoadTimeMs  int
	// Referrer    string
	SessionID string
}
