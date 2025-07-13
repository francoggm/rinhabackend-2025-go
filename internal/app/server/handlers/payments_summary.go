package handlers

import (
	"encoding/json"
	"net/http"
	"time"
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

	result, err := h.storageService.GetPaymentsSummary(ctx, from, to)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(result); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
