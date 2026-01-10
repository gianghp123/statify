package configs

import "github.com/gianghp/statify/internal/utils"

func LoadMinioConfig() (string, string) {
	accessKeyID := utils.GetEnv("MINIO_USER", "")
	secretAccessKey := utils.GetEnv("MINIO_PASSWORD", "")

	return accessKeyID, secretAccessKey
}
