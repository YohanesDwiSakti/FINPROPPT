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
	Destination   string `json:"destination"`
	Weight        int    `json:"weight"`
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
	if err := ensurePickupAutomationSchema(db); err != nil {
		writePickupJSON(w, http.StatusInternalServerError, map[string]string{"message": "Data pickup belum bisa dimuat."})
		return
	}

	rows, err := db.Query(
		`select id::text, customer_name, customer_email, phone, address,
		        coalesce(destination, ''), coalesce(weight_kg, 0), coalesce(receipt, ''),
		        pickup_date::text,
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
		var id, name, email, phone, address, destination, receipt, date, pickupTime, note, status, photo, confirmedBy string
		var weight int
		var updatedAt time.Time
		if err := rows.Scan(&id, &name, &email, &phone, &address, &destination, &weight, &receipt, &date, &pickupTime, &note, &status, &photo, &confirmedBy, &updatedAt); err != nil {
			writePickupJSON(w, http.StatusInternalServerError, map[string]string{"message": "Data pickup belum bisa dimuat."})
			return
		}
		items = append(items, map[string]any{
			"id":                 id,
			"customer_name":      name,
			"customer_email":     email,
			"phone":              phone,
			"address":            address,
			"destination":        destination,
			"weight":             weight,
			"receipt":            receipt,
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
	payload.Destination = strings.TrimSpace(payload.Destination)
	payload.PickupDate = strings.TrimSpace(payload.PickupDate)
	payload.PickupTime = strings.TrimSpace(payload.PickupTime)
	payload.Note = strings.TrimSpace(payload.Note)

	if payload.CustomerName == "" || payload.CustomerEmail == "" || payload.Phone == "" || payload.Address == "" || payload.Destination == "" || payload.PickupDate == "" || payload.PickupTime == "" || payload.Weight <= 0 {
		writePickupJSON(w, http.StatusBadRequest, map[string]string{"message": "Nama, nomor HP, alamat, tujuan, berat, tanggal, dan jam pickup wajib diisi."})
		return
	}

	db, err := openPickupDatabase()
	if err != nil {
		writePickupJSON(w, http.StatusInternalServerError, map[string]string{"message": "Jadwal pickup belum bisa dibuat."})
		return
	}
	defer db.Close()
	if err := ensurePickupAutomationSchema(db); err != nil {
		writePickupJSON(w, http.StatusInternalServerError, map[string]string{"message": "Jadwal pickup belum bisa dibuat."})
		return
	}

	var id string
	receipt := generatePickupReceipt()
	price := automaticOngkir(payload.Weight)
	err = db.QueryRow(
		`insert into public.pickups (customer_name, customer_email, phone, address, destination, weight_kg, receipt, pickup_date, pickup_time, note)
		 values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		 returning id::text`,
		payload.CustomerName,
		payload.CustomerEmail,
		payload.Phone,
		payload.Address,
		payload.Destination,
		payload.Weight,
		receipt,
		payload.PickupDate,
		payload.PickupTime,
		payload.Note,
	).Scan(&id)
	if err != nil {
		writePickupJSON(w, http.StatusInternalServerError, map[string]string{"message": "Jadwal pickup belum bisa dibuat."})
		return
	}
	if err := createAutomaticShipment(db, receipt, payload, price); err != nil {
		writePickupJSON(w, http.StatusInternalServerError, map[string]string{"message": "Pickup dibuat, tapi otomatisasi resi belum berhasil."})
		return
	}

	writePickupJSON(w, http.StatusOK, map[string]string{
		"id":      id,
		"receipt": receipt,
		"message": fmt.Sprintf("Pickup dibuat otomatis. Resi %s, invoice Rp %s, dan rute driver sudah disiapkan.", receipt, formatRupiah(price)),
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
	if err := ensurePickupAutomationSchema(db); err != nil {
		writePickupJSON(w, http.StatusInternalServerError, map[string]string{"message": "Pickup belum bisa dikonfirmasi."})
		return
	}

	var receipt string
	if err := db.QueryRow(`select coalesce(receipt, '') from public.pickups where id = $1`, payload.ID).Scan(&receipt); err != nil {
		writePickupJSON(w, http.StatusNotFound, map[string]string{"message": "Data pickup tidak ditemukan."})
		return
	}

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
	syncShipmentAfterPickup(db, receipt, payload.Status, payload.ConfirmedBy)

	writePickupJSON(w, http.StatusOK, map[string]string{"message": "Status pickup dan resi otomatis diperbarui."})
}

func ensurePickupAutomationSchema(db *sql.DB) error {
	statements := []string{
		`alter table public.pickups add column if not exists destination text`,
		`alter table public.pickups add column if not exists weight_kg integer`,
		`alter table public.pickups add column if not exists receipt text`,
		`create unique index if not exists pickups_receipt_unique on public.pickups(receipt) where receipt is not null`,
		`create table if not exists public.manifests (id uuid primary key default gen_random_uuid(), receipt text unique not null, status text not null, location text, updated_by text, created_at timestamptz not null default now(), updated_at timestamptz not null default now())`,
		`create table if not exists public.payments (id uuid primary key default gen_random_uuid(), customer text not null, email text not null, receipt text, amount integer not null, status text not null default 'Menunggu Pembayaran', created_at timestamptz not null default now())`,
		`create table if not exists public.courier_routes (id uuid primary key default gen_random_uuid(), receipt text not null, driver text not null, origin text not null, destination text not null, status text not null default 'Menunggu Pickup', recommendation text not null, updated_at timestamptz not null default now(), created_at timestamptz not null default now())`,
	}
	for _, statement := range statements {
		if _, err := db.Exec(statement); err != nil {
			return err
		}
	}
	return nil
}

func generatePickupReceipt() string {
	return "TKI-DEN-" + time.Now().Format("060102150405")
}

func automaticOngkir(weight int) int {
	if weight < 1 {
		weight = 1
	}
	return 25000 + ((weight - 1) * 9000)
}

func createAutomaticShipment(db *sql.DB, receipt string, payload pickupRequest, price int) error {
	_, err := db.Exec(
		`insert into public.manifests (receipt, status, location, updated_by)
		 values ($1, 'Pickup Dijadwalkan', $2, 'system-auto')
		 on conflict (receipt) do update set status=excluded.status, location=excluded.location, updated_by='system-auto', updated_at=now()`,
		receipt,
		"Menunggu pickup di "+payload.Address,
	)
	if err != nil {
		return err
	}
	_, err = db.Exec(
		`insert into public.courier_routes (receipt, driver, origin, destination, status, recommendation)
		 values ($1, 'Petugas Pickup Denpasar', 'Hub Denpasar', $2, 'Menunggu Pickup', $3)`,
		receipt,
		payload.Destination,
		fmt.Sprintf("Auto route: jemput paket di %s, bawa ke Hub Denpasar, lalu teruskan menuju %s. Prioritaskan pickup sesuai jadwal %s %s.", payload.Address, payload.Destination, payload.PickupDate, payload.PickupTime),
	)
	if err != nil {
		return err
	}
	_, err = db.Exec(
		`insert into public.payments (customer, email, receipt, amount, status)
		 values ($1, $2, $3, $4, 'Menunggu Pembayaran')`,
		payload.CustomerName,
		payload.CustomerEmail,
		receipt,
		price,
	)
	return err
}

func syncShipmentAfterPickup(db *sql.DB, receipt string, status string, confirmedBy string) {
	if receipt == "" {
		return
	}
	shipmentStatus := "Dalam Perjalanan"
	location := "Paket sudah dijemput kurir"
	routeStatus := "Dalam Perjalanan"
	if status == "Cancel" {
		shipmentStatus = "Cancel"
		location = "Pickup dibatalkan"
		routeStatus = "Cancel"
	}
	_, _ = db.Exec(`update public.manifests set status=$1, location=$2, updated_by=nullif($3,''), updated_at=now() where receipt=$4`, shipmentStatus, location, confirmedBy, receipt)
	_, _ = db.Exec(`update public.courier_routes set status=$1, updated_at=now() where receipt=$2`, routeStatus, receipt)
}

func formatRupiah(amount int) string {
	if amount <= 0 {
		return "0"
	}
	text := fmt.Sprintf("%d", amount)
	chunks := []string{}
	for len(text) > 3 {
		chunks = append([]string{text[len(text)-3:]}, chunks...)
		text = text[:len(text)-3]
	}
	chunks = append([]string{text}, chunks...)
	return strings.Join(chunks, ".")
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
