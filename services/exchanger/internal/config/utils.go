package config

import (
	"fmt"
	"os"
	"strconv"
)

func (db DBConfig) URL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		db.User,
		db.Password,
		db.Host,
		db.Port,
		db.Name,
		db.SSLMode,
	)
}

func getEnv(key, defaultValue string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if valStr, ok := os.LookupEnv(key); ok && valStr != "" {
		if v, err := strconv.Atoi(valStr); err == nil {
			return v
		}
	}
	return defaultValue
}
