package storageworker

import (
	"context"
	"fmt"
	"francoggm/rinhabackend-2025-go/internal/app/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type storageWorker struct {
	id     int
	events chan *models.Event
	db     *pgxpool.Pool
}

func newWorker(id int, events chan *models.Event, db *pgxpool.Pool) *storageWorker {
	return &storageWorker{
		id:     id,
		events: events,
		db:     db,
	}
}

func (w *storageWorker) start(ctx context.Context) {
	fmt.Printf("Storage Worker %d started\n", w.id)

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Storage Worker %d stopping by done context\n", w.id)
		case event, ok := <-w.events:
			if !ok {
				fmt.Printf("Storage Worker %d stopping by closed channel\n", w.id)
				return
			}

			fmt.Printf("Storage Worker %d received event: %v\n", w.id, event)
			w.processEvent(event)
		}
	}
}

func (w *storageWorker) processEvent(event *models.Event) {}
