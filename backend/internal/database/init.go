package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDatabase(connectStr string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  connectStr,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		TranslateError: true,
	})
	return db, err
}
