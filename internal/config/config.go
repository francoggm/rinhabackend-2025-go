package config

import "os"

type Config struct {
	Server
}

type Server struct {
	Port string
}

func NewConfig() *Config {
	return &Config{
		Server: Server{
			Port: os.Getenv("PORT"),
		},
	}
}
