package workers

import (
	"fmt"
	"francoggm/rinhabackend-2025-go/internal/app/models"
)

type Worker struct {
	id     int
	events chan *models.WorkerEvent
}

func NewWorker(id int, events chan *models.WorkerEvent) *Worker {
	return &Worker{
		id:     id,
		events: events,
	}
}

func (w *Worker) Start() {
	fmt.Printf("Worker %d started\n", w.id)
}
