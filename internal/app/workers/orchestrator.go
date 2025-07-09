package workers

import "francoggm/rinhabackend-2025-go/internal/app/models"

type Orchestrator struct {
	events  chan *models.WorkerEvent
	workers []*Worker
}

func NewOrchestrator(workerCount int, events chan *models.WorkerEvent) *Orchestrator {
	var workers []*Worker
	for i := range workerCount {
		worker := NewWorker(i, events)
		workers = append(workers, worker)
	}

	return &Orchestrator{
		events:  events,
		workers: workers,
	}
}

func (o *Orchestrator) StartWorkers() {
	for _, worker := range o.workers {
		go worker.Start()
	}
}
