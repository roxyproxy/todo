package config

import "os"

type Config struct {
	SecretKey string
	DBUrl     string
	GrpcPort  string
	HttpPort  string
}

//TODO create DB struct

func New() *Config {
	return &Config{
		SecretKey: getEnv("SECRETKEY", ""),
		DBUrl:     getEnv("DATABASE_URL", ""),
		GrpcPort:  getEnv("GRPCPORT", ":5000"),
		HttpPort:  getEnv("HTTPPORT", ":5001"),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
