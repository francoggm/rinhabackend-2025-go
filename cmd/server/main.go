package main

import (
	"context"
	"fmt"
	"francoggm/rinhabackend-2025-go/internal/app/models"
	"francoggm/rinhabackend-2025-go/internal/app/server"
	processorworker "francoggm/rinhabackend-2025-go/internal/app/workers/processor"
	storageworker "francoggm/rinhabackend-2025-go/internal/app/workers/storage"
	"francoggm/rinhabackend-2025-go/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.NewConfig()

	uri := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Database.Port,
		cfg.Name,
	)

	ctx := context.Background()

	db, err := pgxpool.New(ctx, uri)
	if err != nil {
		panic(err)
	}

	// Worker queues
	processorEvents := make(chan *models.Event, cfg.ProcessorBufferSize)
	storageEvents := make(chan *models.Event, cfg.StorageBufferSize)

	// Worker orchestrators
	processorOrchestrator := processorworker.NewOrchestrator(cfg, processorEvents, storageEvents)
	storageOrchestrator := storageworker.NewOrchestrator(cfg, storageEvents, db)

	storageOrchestrator.StartWorkers(ctx)
	processorOrchestrator.StartWorkers(ctx)

	server := server.NewServer(cfg, db, processorEvents)
	if err := server.Run(); err != nil {
		panic(err)
	}

	close(processorEvents)
	close(storageEvents)
}
