package models

import "time"

type Payment struct {
	CorrelationID string    `json:"correlationId"`
	Amount        float32   `json:"amount"`
	RequestedAt   time.Time `json:"requestedAt"`
}
