package processors

import (
	"context"
	"fmt"
	"francoggm/rinhabackend-2025-go/internal/app/services"
	"francoggm/rinhabackend-2025-go/internal/models"
)

type StorageProcessor struct {
	service *services.StorageService
}

func NewStorageProcessor(service *services.StorageService) *StorageProcessor {
	return &StorageProcessor{
		service: service,
	}
}

func (p *StorageProcessor) ProcessEvent(ctx context.Context, event any) error {
	payment, ok := event.(*models.Payment)
	if !ok {
		return fmt.Errorf("invalid event type: %T, expected models.Payment", event)
	}

	return p.service.SavePayment(ctx, payment)
}
