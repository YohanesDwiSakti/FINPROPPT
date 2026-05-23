package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type manifestRequest struct {
	Receipt  string `json:"receipt"`
	Status   string `json:"status"`
	Location string `json:"location"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"message":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var payload manifestRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, `{"message":"invalid JSON body"}`, http.StatusBadRequest)
		return
	}

	payload.Receipt = strings.ToUpper(strings.TrimSpace(payload.Receipt))
	if payload.Receipt == "" || strings.TrimSpace(payload.Status) == "" {
		http.Error(w, `{"message":"receipt and status are required"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("Manifest %s berhasil diproses.", payload.Receipt),
	})
}
