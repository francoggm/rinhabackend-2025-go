package services

import (
	"context"
	"fmt"
	"francoggm/rinhabackend-2025-go/internal/models"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const insertPaymentQuery = `INSERT INTO payments (correlation_id, amount, processor_type, requested_at)
															VALUES ($1, $2, $3, $4)`

const getPaymentsSummaryBaseQuery = `SELECT processor_type, COUNT(*) AS total_requests, SUM(amount) AS total_amount FROM payments`

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
	sql, args := buildGetPaymentsSummaryQuery(from, to)

	rows, err := s.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := map[string]*models.ProcessorSummary{
		"default":  {},
		"fallback": {},
	}

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

func buildGetPaymentsSummaryQuery(from, to *time.Time) (string, []any) {
	var where []string
	var args []any

	if from != nil {
		where = append(where, fmt.Sprintf("requested_at >= $%d", len(where)+1))
		args = append(args, *from)
	}

	if to != nil {
		where = append(where, fmt.Sprintf("requested_at <= $%d", len(where)+1))
		args = append(args, *to)
	}

	sql := getPaymentsSummaryBaseQuery
	if len(where) > 0 {
		sql += " WHERE " + strings.Join(where, " AND ")
	}
	sql += " GROUP BY processor_type"

	return sql, args
}
