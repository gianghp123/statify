package enums

import (
	"database/sql/driver"
	"fmt"
)

type Role string

const (
	UserRoleAdmin Role = "ADMIN"
	UserRoleUser  Role = "USER"
)

func (p *Role) Scan(value interface{}) error {
	if value == nil {
		*p = ""
		return nil
	}
	switch v := value.(type) {
	case []byte:
		*p = Role(v)
	case string:
		*p = Role(v)
	default:
		return fmt.Errorf("unexpected type for Role: %T", value)
	}
	return nil
}

func (p Role) Value() (driver.Value, error) {
	return string(p), nil
}
