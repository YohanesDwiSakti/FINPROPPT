package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type trackingStep struct {
	Date   string `json:"date"`
	Status string `json:"status"`
}

type trackingResponse struct {
	Receipt   string         `json:"receipt"`
	Status    string         `json:"status"`
	Estimate  string         `json:"estimate"`
	Location  string         `json:"location"`
	Timeline  []trackingStep `json:"timeline"`
	UpdatedAt string         `json:"updated_at"`
}

type rateResponse struct {
	Origin      string `json:"origin"`
	Destination string `json:"destination"`
	WeightKg    int    `json:"weight_kg"`
	Service     string `json:"service"`
	Price       int    `json:"price"`
	Estimate    string `json:"estimate"`
}

type manifestRequest struct {
	Receipt  string `json:"receipt"`
	Status   string `json:"status"`
	Location string `json:"location"`
}

type apiResponse struct {
	Message string `json:"message"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", withCORS(healthHandler))
	mux.HandleFunc("/api/tracking/", withCORS(trackingHandler))
	mux.HandleFunc("/api/rates", withCORS(rateHandler))
	mux.HandleFunc("/api/manifests", withCORS(manifestHandler))

	addr := ":5000"
	log.Printf("FINPROPPT backend listening on http://127.0.0.1%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

func withCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next(w, r)
	}
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, apiResponse{Message: "backend ready"})
}

func trackingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, apiResponse{Message: "method not allowed"})
		return
	}

	receipt := strings.TrimPrefix(r.URL.Path, "/api/tracking/")
	receipt = strings.ToUpper(strings.TrimSpace(receipt))
	if receipt == "" {
		writeJSON(w, http.StatusBadRequest, apiResponse{Message: "receipt is required"})
		return
	}

	writeJSON(w, http.StatusOK, trackingResponse{
		Receipt:   receipt,
		Status:    "SEDANG DIPROSES",
		Estimate:  "Besok, area Bali",
		Location:  "Hub Denpasar",
		UpdatedAt: time.Now().Format(time.RFC3339),
		Timeline: []trackingStep{
			{Date: "Hari ini, 14:30 WITA", Status: "Paket telah tiba di Hub Denpasar untuk proses sortir"},
			{Date: "Kemarin, 20:00 WIB", Status: "Paket dalam perjalanan menuju Denpasar"},
			{Date: "3 Mei 2026, 09:00 WIB", Status: "Paket diterima di TIKI Jakarta Gateway"},
		},
	})
}

func rateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, apiResponse{Message: "method not allowed"})
		return
	}

	q := r.URL.Query()
	origin := strings.TrimSpace(q.Get("origin"))
	destination := strings.TrimSpace(q.Get("destination"))
	weightRaw := strings.TrimSpace(q.Get("weight"))
	if origin == "" || destination == "" || weightRaw == "" {
		writeJSON(w, http.StatusBadRequest, apiResponse{Message: "origin, destination, and weight are required"})
		return
	}

	weight, err := strconv.Atoi(weightRaw)
	if err != nil || weight < 1 {
		writeJSON(w, http.StatusBadRequest, apiResponse{Message: "weight must be a positive number"})
		return
	}

	price := 18000 + ((weight - 1) * 9000)
	if !strings.EqualFold(origin, destination) {
		price += 7000
	}

	writeJSON(w, http.StatusOK, rateResponse{
		Origin:      origin,
		Destination: destination,
		WeightKg:    weight,
		Service:     "REG Bali",
		Price:       price,
		Estimate:    "1-2 hari kerja",
	})
}

func manifestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, apiResponse{Message: "method not allowed"})
		return
	}

	var payload manifestRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Message: "invalid JSON body"})
		return
	}

	payload.Receipt = strings.ToUpper(strings.TrimSpace(payload.Receipt))
	payload.Status = strings.TrimSpace(payload.Status)
	payload.Location = strings.TrimSpace(payload.Location)
	if payload.Receipt == "" || payload.Status == "" {
		writeJSON(w, http.StatusBadRequest, apiResponse{Message: "receipt and status are required"})
		return
	}

	writeJSON(w, http.StatusOK, apiResponse{
		Message: fmt.Sprintf("Paket %s berhasil diupdate: %s - %s", payload.Receipt, payload.Status, payload.Location),
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("write response error: %v", err)
	}
}
