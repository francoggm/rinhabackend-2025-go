package handlers

import (
	"encoding/json"
	"francoggm/rinhabackend-2025-go/internal/app/models"
	"net/http"
	"time"
)

func (h *Handlers) ProcessPayment(w http.ResponseWriter, r *http.Request) {
	var payment models.Payment
	err := json.NewDecoder(r.Body).Decode(&payment)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	event := &models.Event{
		CorrelationID: payment.CorrelationID,
		Amount:        payment.Amount,
		RequestedAt:   time.Now().UTC(),
	}

	h.events <- event

	w.WriteHeader(http.StatusAccepted)
}
