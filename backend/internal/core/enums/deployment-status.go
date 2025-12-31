package enums

import "database/sql/driver"

type DeploymentStatus string

const (
	DeploymentStatusFailed   DeploymentStatus = "FAILED"
	DeploymentStatusUploaded DeploymentStatus = "UPLOADED"
	DeploymentStatusReady    DeploymentStatus = "READY"
	DeploymentStatusDeleted  DeploymentStatus = "DELETED"
)

func (p *DeploymentStatus) Scan(value interface{}) error {
	*p = DeploymentStatus(value.([]byte))
	return nil
}

func (p DeploymentStatus) Value() (driver.Value, error) {
	return string(p), nil
}
