package workers

import (
	"context"
	"fmt"
	"francoggm/rinhabackend-2025-go/internal/app/workers/processors"
)

type worker struct {
	id              int
	eventsCh        chan any
	eventsProcessor processors.Processor
}

func newWorker(id int, eventsCh chan any, eventsProcessor processors.Processor) *worker {
	return &worker{
		id:              id,
		eventsCh:        eventsCh,
		eventsProcessor: eventsProcessor,
	}
}

func (w *worker) start(ctx context.Context) {
	fmt.Printf("Worker %d started\n", w.id)
	defer fmt.Printf("Worker %d stopped\n", w.id)

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Worker %d stopping by done context\n", w.id)
			return
		case event, ok := <-w.eventsCh:
			if !ok {
				fmt.Printf("Worker %d stopping by closed channel\n", w.id)
				return
			}

			fmt.Printf("Worker %d received event: %v\n", w.id, event)
			if err := w.eventsProcessor.ProcessEvent(event); err != nil {
				fmt.Printf("Worker %d failed to process event: %v, error: %v\n", w.id, event, err)
			}
		}
	}
}
