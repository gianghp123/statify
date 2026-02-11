package enums

import (
	"database/sql/driver"
	"fmt"
)

type JobQueueType string

const (
	JobQueueTypeDeploymentDelete  JobQueueType = "deployment_delete"
	JobQueueTypeDeploymentProcess JobQueueType = "deployment_process"
	JobQueueTypeProjectDelete     JobQueueType = "project_delete"
)

func (p *JobQueueType) Scan(value interface{}) error {
	if value == nil {
		*p = ""
		return nil
	}
	switch v := value.(type) {
	case []byte:
		*p = JobQueueType(v)
	case string:
		*p = JobQueueType(v)
	default:
		return fmt.Errorf("unexpected type for JobQueueType: %T", value)
	}
	return nil
}

func (p JobQueueType) Value() (driver.Value, error) {
	return string(p), nil
}
