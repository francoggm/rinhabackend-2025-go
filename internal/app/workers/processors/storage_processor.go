package processors

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type StorageProcessor struct {
	db *pgxpool.Pool
}

func NewStorageProcessor() *StorageProcessor {
	return &StorageProcessor{}
}

func (p *StorageProcessor) ProcessEvent(ctx context.Context, event any) error { return nil }
