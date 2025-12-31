package enums

import "database/sql/driver"

type Role string

const (
	UserRoleAdmin Role = "ADMIN"
	UserRoleUser  Role = "USER"
)

func (p *Role) Scan(value interface{}) error {
	*p = Role(value.([]byte))
	return nil
}

func (p Role) Value() (driver.Value, error) {
	return string(p), nil
}
