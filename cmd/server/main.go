package main

import (
	"context"
	"fmt"
	"francoggm/rinhabackend-2025-go/internal/app/server"
	"francoggm/rinhabackend-2025-go/internal/app/services"
	"francoggm/rinhabackend-2025-go/internal/app/workers"
	"francoggm/rinhabackend-2025-go/internal/app/workers/processors"
	"francoggm/rinhabackend-2025-go/internal/config"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.NewConfig()

	uri := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)

	ctx := context.Background()

	dbCfg, err := pgxpool.ParseConfig(uri)
	if err != nil {
		panic(fmt.Errorf("failed to parse database URI: %w", err))
	}

	dbCfg.MaxConns = 100
	dbCfg.MinConns = int32(cfg.Workers.StorageCount)
	dbCfg.HealthCheckPeriod = 30 * time.Second
	dbCfg.MaxConnIdleTime = 5 * time.Minute
	dbCfg.MaxConnLifetime = 30 * time.Minute

	db, err := pgxpool.NewWithConfig(ctx, dbCfg)
	if err != nil {
		panic(err)
	}

	if err := db.Ping(ctx); err != nil {
		panic(fmt.Errorf("failed to connect to database: %w", err))
	}

	// Worker queues
	paymentEventsCh := make(chan any, cfg.PaymentBufferSize)
	storageEventsCh := make(chan any, cfg.StorageBufferSize)

	// Services
	paymentService := services.NewPaymentService(cfg.PaymentProcessorConfig.DefaultURL, cfg.PaymentProcessorConfig.FallbackURL)
	storageService := services.NewStorageService(db)

	// Worker processors
	paymentProcessor := processors.NewPaymentProcessor(paymentService, storageEventsCh)
	storageProcessor := processors.NewStorageProcessor(storageService)

	// Worker orchestrators
	paymentOrchestrator := workers.NewOrchestrator(cfg.PaymentCount, paymentEventsCh, paymentProcessor)
	storageOrchestrator := workers.NewOrchestrator(cfg.StorageCount, storageEventsCh, storageProcessor)

	// Start workers in order of processing
	storageOrchestrator.StartWorkers(ctx)
	paymentOrchestrator.StartWorkers(ctx)

	server := server.NewServer(cfg, storageService, paymentEventsCh)
	if err := server.Run(); err != nil {
		panic(err)
	}

	close(paymentEventsCh)
	close(storageEventsCh)
}
