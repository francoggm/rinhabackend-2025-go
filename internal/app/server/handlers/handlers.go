package handlers

import (
	"francoggm/rinhabackend-2025-go/internal/app/services"
	"francoggm/rinhabackend-2025-go/internal/config"
)

type Handlers struct {
	cfg             *config.Config
	storageService  *services.StorageService
	paymentEventsCh chan any
}

func NewHandlers(cfg *config.Config, storageService *services.StorageService, paymentEventsCh chan any) *Handlers {
	return &Handlers{
		cfg:             cfg,
		storageService:  storageService,
		paymentEventsCh: paymentEventsCh,
	}
}
