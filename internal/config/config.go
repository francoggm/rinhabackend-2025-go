package config

import (
	"os"
	"strconv"
)

type Config struct {
	Database
	Workers
	Server
}

type Database struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type Workers struct {
	ProcessorCount      int
	ProcessorBufferSize int
	StorageCount        int
	StorageBufferSize   int
}

type Server struct {
	Port string
}

func NewConfig() *Config {
	return &Config{
		Database: Database{
			Host:     getEnvString("DB_HOST", "localhost"),
			Port:     getEnvString("DB_PORT", "5432"),
			User:     getEnvString("DB_USER", "postgres"),
			Password: getEnvString("DB_PASSWORD", "password"),
			Name:     getEnvString("DB_NAME", "payments"),
		},
		Workers: Workers{
			ProcessorCount:      getEnvInt("PROCESSOR_WORKERS_COUNT", 5),
			ProcessorBufferSize: getEnvInt("PROCESSOR_WORKERS_EVENTS_BUFFER_SIZE", 100),
			StorageCount:        getEnvInt("STORAGE_WORKERS_COUNT", 5),
			StorageBufferSize:   getEnvInt("STORAGE_WORKERS_EVENTS_BUFFER_SIZE", 100),
		},
		Server: Server{
			Port: getEnvString("SERVER_PORT", "8080"),
		},
	}
}

func getEnvString(key string, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	return value
}

func getEnvInt(key string, defaultValue int) int {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return intValue
}
