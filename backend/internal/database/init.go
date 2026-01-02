package database

import (
	"github.com/gianghp/statify/internal/configs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDatabase() (*gorm.DB, error) {
	config := configs.LoadDatabaseConfig()
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  config,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		TranslateError: true,
	})
	return db, err
}
