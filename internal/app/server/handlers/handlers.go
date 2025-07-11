package handlers

import (
	"francoggm/rinhabackend-2025-go/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Handlers struct {
	cfg             *config.Config
	db              *pgxpool.Pool
	paymentEventsCh chan any
}

func NewHandlers(cfg *config.Config, db *pgxpool.Pool, paymentEventsCh chan any) *Handlers {
	return &Handlers{
		cfg:             cfg,
		db:              db,
		paymentEventsCh: paymentEventsCh,
	}
}
