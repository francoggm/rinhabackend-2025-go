package models

type Payment struct {
	CorrelationID string  `json:"correlationId"`
	Amount        float32 `json:"amount"`
}
