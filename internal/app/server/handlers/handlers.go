package handlers

import (
	"francoggm/rinhabackend-2025-go/internal/app/models"
	"francoggm/rinhabackend-2025-go/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Handlers struct {
	cfg    *config.Config
	db     *pgxpool.Pool
	events chan *models.Event
}

func NewHandlers(cfg *config.Config, db *pgxpool.Pool, events chan *models.Event) *Handlers {
	return &Handlers{
		cfg:    cfg,
		db:     db,
		events: events,
	}
}
