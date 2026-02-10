package enums

import (
	"database/sql/driver"
	"fmt"
)

type UploadStatus string

const (
	UploadStatusUploaded UploadStatus = "UPLOADED"
	UploadStatusWaiting  UploadStatus = "WAITING"
	UploadStatusExpired  UploadStatus = "EXPIRED"
)

func (p *UploadStatus) Scan(value interface{}) error {
	if value == nil {
		*p = ""
		return nil
	}
	switch v := value.(type) {
	case []byte:
		*p = UploadStatus(v)
	case string:
		*p = UploadStatus(v)
	default:
		return fmt.Errorf("unexpected type for UploadStatus: %T", value)
	}
	return nil
}

func (p UploadStatus) Value() (driver.Value, error) {
	return string(p), nil
}
