package handlers

import "francoggm/rinhabackend-2025-go/internal/config"

type Handlers struct {
	cfg *config.Config
}

func NewHandlers(cfg *config.Config) *Handlers {
	return &Handlers{cfg: cfg}
}
