package config

import "os"

// Config represents a config info used in application.
type Config struct {
	SecretKey string
	DBUrl     string
	GrpcPort  string
	HTTPPort  string
}

// TODO create DB struct

// New returns config object.
func New() *Config {
	return &Config{
		SecretKey: getEnv("SECRETKEY", ""),
		DBUrl:     getEnv("DATABASE_URL", ""),
		GrpcPort:  getEnv("GRPCPORT", ":5000"),
		HTTPPort:  getEnv("HTTPPORT", ":5001"),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
