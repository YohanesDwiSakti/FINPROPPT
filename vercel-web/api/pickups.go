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

type pickupRequest struct {
	CustomerName  string `json:"customer_name"`
	CustomerEmail string `json:"customer_email"`
	Phone         string `json:"phone"`
	Address       string `json:"address"`
	PickupDate    string `json:"pickup_date"`
	PickupTime    string `json:"pickup_time"`
	Note          string `json:"note"`
}

type pickupConfirmRequest struct {
	ID                string `json:"id"`
	Status            string `json:"status"`
	ConfirmationPhoto string `json:"confirmation_photo"`
	ConfirmedBy       string `json:"confirmed_by"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleListPickups(w)
	case http.MethodPost:
		handleCreatePickup(w, r)
	case http.MethodPatch:
		handleConfirmPickup(w, r)
	default:
		writePickupJSON(w, http.StatusMethodNotAllowed, map[string]string{"message": "method not allowed"})
	}
}

func handleListPickups(w http.ResponseWriter) {
	db, err := openPickupDatabase()
	if err != nil {
		writePickupJSON(w, http.StatusInternalServerError, map[string]string{"message": "Data pickup belum bisa dimuat."})
		return
	}
	defer db.Close()

	rows, err := db.Query(
		`select id::text, customer_name, customer_email, phone, address, pickup_date::text,
		        pickup_time, coalesce(note, ''), status, coalesce(confirmation_photo, ''),
		        coalesce(confirmed_by, ''), updated_at
		   from public.pickups
		  order by created_at desc
		  limit 40`,
	)
	if err != nil {
		writePickupJSON(w, http.StatusInternalServerError, map[string]string{"message": "Data pickup belum bisa dimuat."})
		return
	}
	defer rows.Close()

	items := []map[string]any{}
	for rows.Next() {
		var id, name, email, phone, address, date, pickupTime, note, status, photo, confirmedBy string
		var updatedAt time.Time
		if err := rows.Scan(&id, &name, &email, &phone, &address, &date, &pickupTime, &note, &status, &photo, &confirmedBy, &updatedAt); err != nil {
			writePickupJSON(w, http.StatusInternalServerError, map[string]string{"message": "Data pickup belum bisa dimuat."})
			return
		}
		items = append(items, map[string]any{
			"id":                 id,
			"customer_name":      name,
			"customer_email":     email,
			"phone":              phone,
			"address":            address,
			"pickup_date":        date,
			"pickup_time":        pickupTime,
			"note":               note,
			"status":             status,
			"confirmation_photo": photo,
			"confirmed_by":       confirmedBy,
			"updated_at":         updatedAt.Format(time.RFC3339),
		})
	}
	writePickupJSON(w, http.StatusOK, items)
}

func handleCreatePickup(w http.ResponseWriter, r *http.Request) {
	var payload pickupRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writePickupJSON(w, http.StatusBadRequest, map[string]string{"message": "invalid JSON body"})
		return
	}

	payload.CustomerName = strings.TrimSpace(payload.CustomerName)
	payload.CustomerEmail = strings.ToLower(strings.TrimSpace(payload.CustomerEmail))
	payload.Phone = strings.TrimSpace(payload.Phone)
	payload.Address = strings.TrimSpace(payload.Address)
	payload.PickupDate = strings.TrimSpace(payload.PickupDate)
	payload.PickupTime = strings.TrimSpace(payload.PickupTime)
	payload.Note = strings.TrimSpace(payload.Note)

	if payload.CustomerName == "" || payload.CustomerEmail == "" || payload.Phone == "" || payload.Address == "" || payload.PickupDate == "" || payload.PickupTime == "" {
		writePickupJSON(w, http.StatusBadRequest, map[string]string{"message": "Nama, nomor HP, alamat, tanggal, dan jam pickup wajib diisi."})
		return
	}

	db, err := openPickupDatabase()
	if err != nil {
		writePickupJSON(w, http.StatusInternalServerError, map[string]string{"message": "Jadwal pickup belum bisa dibuat."})
		return
	}
	defer db.Close()

	var id string
	err = db.QueryRow(
		`insert into public.pickups (customer_name, customer_email, phone, address, pickup_date, pickup_time, note)
		 values ($1, $2, $3, $4, $5, $6, $7)
		 returning id::text`,
		payload.CustomerName,
		payload.CustomerEmail,
		payload.Phone,
		payload.Address,
		payload.PickupDate,
		payload.PickupTime,
		payload.Note,
	).Scan(&id)
	if err != nil {
		writePickupJSON(w, http.StatusInternalServerError, map[string]string{"message": "Jadwal pickup belum bisa dibuat."})
		return
	}

	writePickupJSON(w, http.StatusOK, map[string]string{
		"id":      id,
		"message": "Jadwal pickup berhasil dibuat.",
	})
}

func handleConfirmPickup(w http.ResponseWriter, r *http.Request) {
	var payload pickupConfirmRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writePickupJSON(w, http.StatusBadRequest, map[string]string{"message": "invalid JSON body"})
		return
	}

	payload.ID = strings.TrimSpace(payload.ID)
	payload.Status = strings.TrimSpace(payload.Status)
	payload.ConfirmationPhoto = strings.TrimSpace(payload.ConfirmationPhoto)
	payload.ConfirmedBy = strings.TrimSpace(payload.ConfirmedBy)

	if payload.ID == "" || payload.Status == "" {
		writePickupJSON(w, http.StatusBadRequest, map[string]string{"message": "Data pickup dan status wajib diisi."})
		return
	}
	if payload.Status == "Dijemput" && payload.ConfirmationPhoto == "" {
		writePickupJSON(w, http.StatusBadRequest, map[string]string{"message": "Foto konfirmasi wajib diupload sebelum pickup dikonfirmasi."})
		return
	}
	if payload.Status != "Dijemput" && payload.Status != "Cancel" {
		writePickupJSON(w, http.StatusBadRequest, map[string]string{"message": "Status pickup tidak valid."})
		return
	}

	db, err := openPickupDatabase()
	if err != nil {
		writePickupJSON(w, http.StatusInternalServerError, map[string]string{"message": "Pickup belum bisa dikonfirmasi."})
		return
	}
	defer db.Close()

	result, err := db.Exec(
		`update public.pickups
		    set status = $1,
		        confirmation_photo = nullif($2, ''),
		        confirmed_by = nullif($3, ''),
		        updated_at = now()
		  where id = $4`,
		payload.Status,
		payload.ConfirmationPhoto,
		payload.ConfirmedBy,
		payload.ID,
	)
	if err != nil {
		writePickupJSON(w, http.StatusInternalServerError, map[string]string{"message": "Pickup belum bisa dikonfirmasi."})
		return
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		writePickupJSON(w, http.StatusNotFound, map[string]string{"message": "Data pickup tidak ditemukan."})
		return
	}

	writePickupJSON(w, http.StatusOK, map[string]string{"message": "Status pickup berhasil diperbarui."})
}

func openPickupDatabase() (*sql.DB, error) {
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

func writePickupJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}
