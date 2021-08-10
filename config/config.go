package config

import "os"

type Config struct {
	SecretKey string
}

func New() *Config {
	return &Config{SecretKey: getEnv("SKERTKEY", "")}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
