package processorworker

import (
	"context"
	"francoggm/rinhabackend-2025-go/internal/app/models"
	"francoggm/rinhabackend-2025-go/internal/config"
)

type ProcessorOrchestrator struct {
	events  chan *models.Event
	workers []*processorWorker
}

func NewOrchestrator(cfg *config.Config, events, storageEvents chan *models.Event) *ProcessorOrchestrator {
	var workers []*processorWorker
	for id := range cfg.ProcessorCount {
		worker := newWorker(id, events, storageEvents)
		workers = append(workers, worker)
	}

	return &ProcessorOrchestrator{
		events:  events,
		workers: workers,
	}
}

func (o *ProcessorOrchestrator) StartWorkers(ctx context.Context) {
	for _, worker := range o.workers {
		go worker.start(ctx)
	}
}
