package models

import "time"

type Event struct {
	CorrelationID string
	Amount        float32
	RequestedAt   time.Time
}
