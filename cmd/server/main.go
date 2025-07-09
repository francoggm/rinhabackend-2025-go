package main

import (
	"francoggm/rinhabackend-2025-go/internal/app/models"
	"francoggm/rinhabackend-2025-go/internal/app/server"
	"francoggm/rinhabackend-2025-go/internal/app/workers"
	"francoggm/rinhabackend-2025-go/internal/config"
)

func main() {
	cfg := config.NewConfig()

	events := make(chan *models.WorkerEvent, cfg.WorkerEventsBufferSize)

	orchestrator := workers.NewOrchestrator(cfg.WorkerCount, events)
	orchestrator.StartWorkers()

	server := server.NewServer(cfg)
	if err := server.Run(); err != nil {
		panic(err)
	}
}
