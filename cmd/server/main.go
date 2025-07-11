package main

import (
	"context"
	"fmt"
	"francoggm/rinhabackend-2025-go/internal/app/server"
	"francoggm/rinhabackend-2025-go/internal/app/workers"
	"francoggm/rinhabackend-2025-go/internal/app/workers/processors"
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
	paymentEventsCh := make(chan any, cfg.PaymentBufferSize)
	storageEventsCh := make(chan any, cfg.StorageBufferSize)

	// Worker processors
	paymentProcessor := processors.NewPaymentProcessor()
	storageProcessor := processors.NewStorageProcessor()

	// Worker orchestrators
	paymentOrchestrator := workers.NewOrchestrator(cfg.PaymentCount, paymentEventsCh, paymentProcessor)
	storageOrchestrator := workers.NewOrchestrator(cfg.StorageCount, storageEventsCh, storageProcessor)

	// Start workers in order of processing
	storageOrchestrator.StartWorkers(ctx)
	paymentOrchestrator.StartWorkers(ctx)

	server := server.NewServer(cfg, db, paymentEventsCh)
	if err := server.Run(); err != nil {
		panic(err)
	}

	close(paymentEventsCh)
	close(storageEventsCh)
}
