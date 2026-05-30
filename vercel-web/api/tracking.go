package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type trackingStep struct {
	Date   string `json:"date"`
	Status string `json:"status"`
}

type shipmentTracking struct {
	Receipt   string
	Status    string
	Location  sql.NullString
	UpdatedAt time.Time
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

	shipment, err := fetchShipment(receipt)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		if err != sql.ErrNoRows {
			fmt.Println("tracking database error:", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "Resi belum bisa dicek. Coba lagi."})
			return
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "Nomor resi belum terdaftar."})
		return
	}

	location := strings.TrimSpace(shipment.Location.String)
	if location == "" {
		location = "Hub Denpasar"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"receipt":    shipment.Receipt,
		"status":     shipment.Status,
		"estimate":   estimateForStatus(shipment.Status),
		"location":   location,
		"updated_at": shipment.UpdatedAt.Format(time.RFC3339),
		"timeline": []trackingStep{
			{Date: shipment.UpdatedAt.Format("02 Jan 2006, 15:04 WITA"), Status: statusTimelineText(shipment.Status, location)},
		},
	})
}

func fetchShipment(receipt string) (shipmentTracking, error) {
	db, err := openTrackingDatabase()
	if err != nil {
		return shipmentTracking{}, err
	}
	defer db.Close()

	var shipment shipmentTracking
	err = db.QueryRow(
		`select receipt, status, location, updated_at
		 from public.manifests
		 where receipt = $1`,
		receipt,
	).Scan(&shipment.Receipt, &shipment.Status, &shipment.Location, &shipment.UpdatedAt)
	return shipment, err
}

func openTrackingDatabase() (*sql.DB, error) {
	url := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if url == "" {
		return nil, fmt.Errorf("DATABASE_URL belum diset")
	}
	if !strings.Contains(url, "sslmode=") {
		if strings.Contains(url, "?") {
			url += "&sslmode=require"
		} else {
			url += "?sslmode=require"
		}
	}
	if !strings.Contains(url, "default_query_exec_mode=") {
		url += "&default_query_exec_mode=simple_protocol"
	}
	return sql.Open("pgx", url)
}

func estimateForStatus(status string) string {
	switch strings.ToLower(status) {
	case "sudah diterima":
		return "Selesai"
	case "cancel":
		return "Pengiriman dibatalkan"
	case "dalam perjalanan":
		return "1-2 hari kerja"
	default:
		return "Sedang diproses"
	}
}

func statusTimelineText(status string, location string) string {
	switch strings.ToLower(status) {
	case "sudah diterima":
		return "Paket sudah diterima di " + location
	case "cancel":
		return "Pengiriman dibatalkan oleh admin"
	case "dalam perjalanan":
		return "Paket sedang dalam perjalanan. Posisi terakhir: " + location
	default:
		return "Paket sedang diproses di " + location
	}
}
