package handler

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	origin := strings.TrimSpace(q.Get("origin"))
	destination := strings.TrimSpace(q.Get("destination"))
	weightInput, err := strconv.ParseFloat(strings.TrimSpace(q.Get("weight")), 64)
	if origin == "" || destination == "" || err != nil || weightInput <= 0 {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"message":"Kota asal, kota tujuan, dan berat wajib diisi."}`, http.StatusBadRequest)
		return
	}

	weight := int(math.Ceil(weightInput))
	price := 18000 + ((weight - 1) * 9000)
	if !strings.EqualFold(origin, destination) {
		price += 7000
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"origin":      origin,
		"destination": destination,
		"weight_kg":   weight,
		"input_kg":    weightInput,
		"service":     "REG Bali",
		"price":       price,
		"estimate":    "1-2 hari kerja",
	})
}
