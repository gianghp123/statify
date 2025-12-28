package configs

import (
	"fmt"

	"github.com/gianghp/statify/internal/utils"
	"gorm.io/gorm"
)

var DB *gorm.DB

func LoadDatabaseConfig() string {
	// Load environment variables or use defaults
	host := utils.GetEnv("POSTGRES_HOST", "localhost")
	port := utils.GetEnv("POSTGRES_PORT", "5432")
	user := utils.GetEnv("POSTGRES_USER", "myuser")
	password := utils.GetSecret("POSTGRES_PASSWORD", "")
	dbname := utils.GetEnv("POSTGRES_DB", "mydb")
	sslmode := utils.GetEnv("POSTGRES_SSLMODE", "disable") // usually disable for local dev

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)
	return dsn
}
