package models

import (
	"github.com/gianghp/statify/internal/core/enums"
	"gorm.io/gorm"
)

type JobQueue struct {
	gorm.Model
	Type         enums.JobQueueType
	DeploymentID uint
	Deployment
	Payload    string
	Status     enums.JobQueueStatus
	RetryCount int
	Error      string
}
