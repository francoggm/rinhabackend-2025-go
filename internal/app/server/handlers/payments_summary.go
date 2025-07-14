package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"go.uber.org/zap"
)

func (h *Handlers) GetPaymentsSummary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	query := r.URL.Query()

	var from, to *time.Time
	layout := time.RFC3339

	if fromStr := query.Get("from"); fromStr != "" {
		if t, err := time.Parse(layout, fromStr); err == nil {
			from = &t
		}
	}

	if toStr := query.Get("to"); toStr != "" {
		if t, err := time.Parse(layout, toStr); err == nil {
			to = &t
		}
	}

	summary, err := h.storageService.GetPaymentsSummary(ctx, from, to)
	if err != nil {
		zap.L().Error("failed to get payments summary", zap.Error(err))
		http.Error(w, "failed to get payments summary", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(summary); err != nil {
		zap.L().Error("failed to encode payments summary", zap.Error(err))
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
}
