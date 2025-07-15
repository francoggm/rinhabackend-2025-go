package services

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

type PaymentService struct {
	httpClient              *http.Client
	defaultURL              string
	fallbackURL             string
	isDefaultHealthy        atomic.Bool
	defaultMinResponseTime  atomic.Int32
	fallbackMinResponseTime atomic.Int32
}

func NewPaymentService(defaultURL, fallbackURL string) *PaymentService {
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     30 * time.Second,
		},
	}

	service := &PaymentService{
		httpClient:  httpClient,
		defaultURL:  defaultURL,
		fallbackURL: fallbackURL,
	}

	go service.startHealthChecker()
	return service
}

func (p *PaymentService) MakePayment(ctx context.Context, payment *models.Payment) error {
	payload, err := json.Marshal(payment)
	if err != nil {
		return err
	}

	var url, processingType string
	if p.isDefaultHealthy.Load() {
		url = p.defaultURL
		processingType = "default"
	} else {
		url = p.fallbackURL
		processingType = "fallback"
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url+"/payments", bytes.NewBuffer(payload))
	if err != nil {
		return nil
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request to %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 response from %s: %d", processingType, resp.StatusCode)
	}

	payment.ProcessingType = processingType

	return nil
}

func (p *PaymentService) startHealthChecker() {
	ticker := time.NewTicker(5100 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		ctx := context.Background()

		p.checkDefaultHealth(ctx)
		p.checkFallbackHealth(ctx)
	}
}

func (p *PaymentService) checkDefaultHealth(ctx context.Context) {
	defaultHealthCheck, err := p.checkHealth(ctx, p.defaultURL+"/payments/service-health")
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

func (p *PaymentService) checkFallbackHealth(ctx context.Context) {
	fallbackHealthCheck, err := p.checkHealth(ctx, p.fallbackURL+"/payments/service-health")
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

func (p *PaymentService) checkHealth(ctx context.Context, url string) (*models.HealthCheck, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := p.httpClient.Do(req)
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
