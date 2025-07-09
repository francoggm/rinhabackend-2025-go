package main

import (
	"francoggm/rinhabackend-2025-go/internal/app/server"
	"francoggm/rinhabackend-2025-go/internal/config"
)

func main() {
	cfg := config.NewConfig()

	server := server.NewServer(cfg)
	if err := server.Run(); err != nil {
		panic(err)
	}
}
