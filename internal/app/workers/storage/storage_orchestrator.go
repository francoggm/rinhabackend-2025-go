package storageworker

import (
	"context"
	"francoggm/rinhabackend-2025-go/internal/app/models"
	"francoggm/rinhabackend-2025-go/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

type StorageOrchestrator struct {
	events  chan *models.Event
	workers []*storageWorker
}

func NewOrchestrator(cfg *config.Config, events chan *models.Event, db *pgxpool.Pool) *StorageOrchestrator {
	var workers []*storageWorker
	for id := range cfg.StorageCount {
		worker := newWorker(id, events, db)
		workers = append(workers, worker)
	}

	return &StorageOrchestrator{
		events:  events,
		workers: workers,
	}
}

func (o *StorageOrchestrator) StartWorkers(ctx context.Context) {
	for _, worker := range o.workers {
		go worker.start(ctx)
	}
}
