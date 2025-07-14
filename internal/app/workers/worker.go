package workers

import (
	"context"
	"francoggm/rinhabackend-2025-go/internal/app/workers/processors"

	"go.uber.org/zap"
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
	zap.L().Info("Starting worker", zap.Int("id", w.id))
	defer zap.L().Info("Worker stopped", zap.Int("id", w.id))

	for {
		select {
		case <-ctx.Done():
			zap.L().Info("Worker context done", zap.Int("id", w.id))
			return
		case event, ok := <-w.eventsCh:
			if !ok {
				zap.L().Info("Worker channel closed", zap.Int("id", w.id))
				return
			}

			if err := w.eventsProcessor.ProcessEvent(ctx, event); err != nil {
				zap.L().Error("Failed to process event", zap.Int("worker_id", w.id), zap.Any("event", event), zap.Error(err))
			}
		}
	}
}
