package enums

import (
	"database/sql/driver"
	"fmt"
)

type DeploymentStatus string

const (
	DeploymentStatusQueued     DeploymentStatus = "QUEUED"
	DeploymentStatusFailed     DeploymentStatus = "FAILED"
	DeploymentStatusUploaded   DeploymentStatus = "UPLOADED"
	DeploymentStatusReady      DeploymentStatus = "READY"
	DeploymentStatusDeleted    DeploymentStatus = "DELETED"
	DeploymentStatusProcessing DeploymentStatus = "PROCESSING"
	DeploymentStatusLive       DeploymentStatus = "LIVE"
)

func (p *DeploymentStatus) Scan(value interface{}) error {
	if value == nil {
		*p = ""
		return nil
	}
	switch v := value.(type) {
	case []byte:
		*p = DeploymentStatus(v)
	case string:
		*p = DeploymentStatus(v)
	default:
		return fmt.Errorf("unexpected type for DeploymentStatus: %T", value)
	}
	return nil
}

func (p DeploymentStatus) Value() (driver.Value, error) {
	return string(p), nil
}
