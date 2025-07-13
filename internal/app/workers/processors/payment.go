package processors

import (
	"context"
	"fmt"
	paymentservice "francoggm/rinhabackend-2025-go/internal/app/services"
	"francoggm/rinhabackend-2025-go/internal/models"
)

type PaymentProcessor struct {
	service         *paymentservice.PaymentService
	storageEventsCh chan any
}

func NewPaymentProcessor(service *paymentservice.PaymentService, storageEventsCh chan any) *PaymentProcessor {
	return &PaymentProcessor{
		service:         service,
		storageEventsCh: storageEventsCh,
	}
}

func (p *PaymentProcessor) ProcessEvent(ctx context.Context, event any) error {
	payment, ok := event.(*models.Payment)
	if !ok {
		return fmt.Errorf("invalid event type: %T, expected models.Payment", event)
	}

	if err := p.service.MakePayment(ctx, payment); err != nil {
		return err
	}

	// Request suceeded, send an event to the storage workers
	p.storageEventsCh <- payment

	return nil
}
