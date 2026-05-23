package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type trackingStep struct {
	Date   string `json:"date"`
	Status string `json:"status"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	receipt := strings.ToUpper(strings.TrimSpace(r.URL.Query().Get("receipt")))
	if receipt == "" {
		receipt = strings.ToUpper(strings.TrimPrefix(r.URL.Path, "/api/tracking/"))
	}
	if receipt == "" {
		http.Error(w, `{"message":"receipt is required"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"receipt":    receipt,
		"status":     "SEDANG DIPROSES",
		"estimate":   "Besok, area Bali",
		"location":   "Hub Denpasar",
		"updated_at": time.Now().Format(time.RFC3339),
		"timeline": []trackingStep{
			{Date: "Hari ini, 14:30 WITA", Status: "Paket telah tiba di Hub Denpasar untuk proses sortir"},
			{Date: "Kemarin, 20:00 WIB", Status: "Paket dalam perjalanan menuju Denpasar"},
			{Date: "3 Mei 2026, 09:00 WIB", Status: "Paket diterima di TIKI Jakarta Gateway"},
		},
	})
}
