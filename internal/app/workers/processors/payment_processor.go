package processors

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"francoggm/rinhabackend-2025-go/internal/models"
	"net/http"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
)

type PaymentProcessor struct {
	storageEventsCh         chan any
	httpClient              *http.Client
	defaultURL              string
	fallbackURL             string
	isDefaultHealthy        atomic.Bool
	defaultMinResponseTime  atomic.Int32
	fallbackMinResponseTime atomic.Int32
}

func NewPaymentProcessor(defaultURL, fallbackURL string) *PaymentProcessor {
	httpClient := &http.Client{
		Timeout: 2 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     30 * time.Second,
		},
	}

	p := &PaymentProcessor{
		httpClient:  httpClient,
		defaultURL:  defaultURL,
		fallbackURL: fallbackURL,
	}

	go p.startHealthChecker()

	return p
}

func (p *PaymentProcessor) ProcessEvent(ctx context.Context, event any) error {
	payment, ok := event.(*models.Payment)
	if !ok {
		return fmt.Errorf("invalid event type: %T, expected models.Payment", event)
	}

	payload, err := json.Marshal(payment)
	if err != nil {
		return err
	}

	var url string
	if p.isDefaultHealthy.Load() {
		url = p.defaultURL
	} else {
		url = p.fallbackURL
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request to %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 response from %s: %d", url, resp.StatusCode)
	}

	// Request suceeded, send an event to the storage workers
	p.storageEventsCh <- payment

	return nil
}

func (p *PaymentProcessor) startHealthChecker() {
	ticker := time.NewTicker(51000 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		ctx := context.Background()

		p.checkDefaultHealth(ctx)
		p.checkFallbackHealth(ctx)
	}
}

func (p *PaymentProcessor) checkDefaultHealth(ctx context.Context) {
	defaultHealthCheck, err := p.checkHealth(ctx, p.defaultURL+"/health")
	if err != nil {
		zap.L().Error("default health check failed", zap.Error(err))

		p.isDefaultHealthy.Store(false)
		p.defaultMinResponseTime.Store(1000) // Default to a high value if health check fails
		return
	}

	zap.L().Info("default health check result",
		zap.Bool("isFailing", defaultHealthCheck.IsFailing),
		zap.Int("minResponseTime", defaultHealthCheck.MinResponseTime))

	p.isDefaultHealthy.Store(!defaultHealthCheck.IsFailing)
	p.defaultMinResponseTime.Store(int32(defaultHealthCheck.MinResponseTime))
}

func (p *PaymentProcessor) checkFallbackHealth(ctx context.Context) {
	fallbackHealthCheck, err := p.checkHealth(ctx, p.fallbackURL+"/health")
	if err != nil {
		zap.L().Error("fallback health check failed", zap.Error(err))

		p.fallbackMinResponseTime.Store(1000) // Default to a high value if health check fails
		return
	}

	zap.L().Info("fallback health check result",
		zap.Bool("isFailing", fallbackHealthCheck.IsFailing),
		zap.Int("minResponseTime", fallbackHealthCheck.MinResponseTime))

	p.fallbackMinResponseTime.Store(int32(fallbackHealthCheck.MinResponseTime))
}

func (p *PaymentProcessor) checkHealth(ctx context.Context, url string) (*models.HealthCheck, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("health check request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("health check failed with status code: %d", resp.StatusCode)
	}

	var healthCheck models.HealthCheck
	if err := json.NewDecoder(resp.Body).Decode(&healthCheck); err != nil {
		return nil, fmt.Errorf("failed to decode health check response: %w", err)
	}

	return &healthCheck, nil
}
