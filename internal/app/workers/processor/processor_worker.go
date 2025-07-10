package processorworker

import (
	"context"
	"fmt"
	"francoggm/rinhabackend-2025-go/internal/app/models"
)

type processorWorker struct {
	id            int
	events        chan *models.Event
	storageEvents chan *models.Event
}

func newWorker(id int, events, storageEvents chan *models.Event) *processorWorker {
	return &processorWorker{
		id:            id,
		events:        events,
		storageEvents: storageEvents,
	}
}

func (w *processorWorker) start(ctx context.Context) {
	fmt.Printf("Processor Worker %d started\n", w.id)

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Processor Worker %d stopping by done context\n", w.id)
		case event, ok := <-w.events:
			if !ok {
				fmt.Printf("Processor Worker %d stopping by closed channel\n", w.id)
				return
			}

			fmt.Printf("Processor Worker %d received event: %v\n", w.id, event)
			w.processEvent(event)
		}
	}
}

func (w *processorWorker) processEvent(event *models.Event) {
	w.storageEvents <- event // Forward the event to storage worker
}
