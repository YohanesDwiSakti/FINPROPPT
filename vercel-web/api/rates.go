package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	origin := strings.TrimSpace(q.Get("origin"))
	destination := strings.TrimSpace(q.Get("destination"))
	weight, err := strconv.Atoi(strings.TrimSpace(q.Get("weight")))
	if origin == "" || destination == "" || err != nil || weight < 1 {
		http.Error(w, `{"message":"origin, destination, and weight are required"}`, http.StatusBadRequest)
		return
	}

	price := 18000 + ((weight - 1) * 9000)
	if !strings.EqualFold(origin, destination) {
		price += 7000
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"origin":      origin,
		"destination": destination,
		"weight_kg":   weight,
		"service":     "REG Bali",
		"price":       price,
		"estimate":    "1-2 hari kerja",
	})
}
