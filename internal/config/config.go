package config

import (
	"os"
	"strconv"
)

type Config struct {
	Workers
	Server
}

type Workers struct {
	WorkerCount            int
	WorkerEventsBufferSize int
}

type Server struct {
	Port string
}

func NewConfig() *Config {
	return &Config{
		Workers: Workers{
			WorkerCount:            getEnvInt("WORKERS_COUNT", 5),
			WorkerEventsBufferSize: getEnvInt("WORKERS_EVENTS_BUFFER_SIZE", 100),
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
