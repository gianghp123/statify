package models

import (
	"time"

	"gorm.io/gorm"
)

type DeploymentMetricsMinute struct {
	gorm.Model
	ProjectID uint    `gorm:"index;not null"`
	Project   Project `gorm:"foreignKey:ProjectID"`

	DeploymentID uint       `gorm:"index;not null"`
	Deployment   Deployment `gorm:"foreignKey:DeploymentID"`

	MinuteTs time.Time `gorm:"index;not null"`

	RequestCount  int64
	BytesServed   int64
	TotalDuration int64

	Status2xx int64 `gorm:"column:status_2xx"`
	Status3xx int64 `gorm:"column:status_3xx"`
	Status4xx int64 `gorm:"column:status_4xx"`
	Status5xx int64 `gorm:"column:status_5xx"`
}

func (DeploymentMetricsMinute) TableName() string {
	return "deployment_metrics_minute"
}
