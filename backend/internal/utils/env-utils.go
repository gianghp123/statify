package utils

import (
	"log/slog"
	"os"
	"strings"
)

func GetEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func GetSecret(envName string, defaultValue string) string {
	fileEnv := envName + "_FILE"
	filePath := GetEnv(fileEnv, "")

	if filePath != "" {
		content, err := os.ReadFile(filePath)
		if err == nil {
			return strings.TrimSpace(string(content))
		}
		slog.Warn("Could not read secret file", "path", filePath, "error", err)
	}

	return GetEnv(envName, defaultValue)
}
