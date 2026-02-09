package enums

import (
	"database/sql/driver"
	"fmt"
)

type JobQueueStatus string

const (
	JobQueueStatusPending    JobQueueStatus = "pending"
	JobQueueStatusProcessing JobQueueStatus = "processing"
	JobQueueStatusSuccess    JobQueueStatus = "success"
	JobQueueStatusFailed     JobQueueStatus = "failed"
)

func (p *JobQueueStatus) Scan(value interface{}) error {
	if value == nil {
		*p = ""
		return nil
	}
	switch v := value.(type) {
	case []byte:
		*p = JobQueueStatus(v)
	case string:
		*p = JobQueueStatus(v)
	default:
		return fmt.Errorf("unexpected type for JobQueueStatus: %T", value)
	}
	return nil
}

func (p JobQueueStatus) Value() (driver.Value, error) {
	return string(p), nil
}
