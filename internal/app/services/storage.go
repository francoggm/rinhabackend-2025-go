package services

import (
	"context"
	"francoggm/rinhabackend-2025-go/internal/models"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

const insertPaymentQuery = ` INSERT INTO payments (correlation_id, amount, processor_type, requested_at)
															VALUES ($1, $2, $3, $4)`

type StorageService struct {
	db *pgxpool.Pool
}

func NewStorageService(db *pgxpool.Pool) *StorageService {
	return &StorageService{
		db: db,
	}
}

func (s *StorageService) SavePayment(ctx context.Context, payment *models.Payment) error {
	_, err := s.db.Exec(
		ctx,
		insertPaymentQuery,
		payment.CorrelationID,
		payment.Amount,
		payment.ProcessingType,
		payment.RequestedAt,
	)
	return err
}

func (s *StorageService) GetPaymentsSummary(ctx context.Context, from, to *time.Time) (map[string]*models.ProcessorSummary, error) {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	queryBuilder := psql.Select(
		"processing_type",
		"COUNT(*) AS total_requests",
		"SUM(amount) AS total_amount",
	).
		From("payments").
		GroupBy("processing_type")

	if from != nil {
		queryBuilder = queryBuilder.Where(squirrel.GtOrEq{"requested_at": *from})
	}

	if to != nil {
		queryBuilder = queryBuilder.Where(squirrel.Lt{"requested_at": *to})
	}

	sql, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := s.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]*models.ProcessorSummary)
	for rows.Next() {
		var processingType string
		var processorSummary models.ProcessorSummary

		if err := rows.Scan(&processingType, &processorSummary.TotalRequests, &processorSummary.TotalAmount); err != nil {
			return nil, err
		}

		result[processingType] = &processorSummary
	}

	return result, nil
}
