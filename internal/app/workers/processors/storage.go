package processors

import (
	"context"
	"francoggm/rinhabackend-2025-go/internal/app/services"
	"francoggm/rinhabackend-2025-go/internal/models"

	"go.uber.org/zap"
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
	payment := event.(*models.Payment)
	zap.L().Info("Processing storage event", zap.Any("payment", payment))

	return p.service.SavePayment(ctx, payment)
}
