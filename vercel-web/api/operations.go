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

type operationPayload struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Customer    string `json:"customer"`
	Email       string `json:"email"`
	Receipt     string `json:"receipt"`
	Subject     string `json:"subject"`
	Message     string `json:"message"`
	Status      string `json:"status"`
	Amount      int    `json:"amount"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Address     string `json:"address"`
	Plate       string `json:"plate"`
	Driver      string `json:"driver"`
	Origin      string `json:"origin"`
	Destination string `json:"destination"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	db, err := openOperationsDatabase()
	if err != nil {
		writeOperationsJSON(w, http.StatusInternalServerError, map[string]string{"message": "Database belum siap."})
		return
	}
	defer db.Close()

	if err := ensureOperationsSchema(db); err != nil {
		writeOperationsJSON(w, http.StatusInternalServerError, map[string]string{"message": "Tabel operasional belum siap."})
		return
	}

	switch r.Method {
	case http.MethodGet:
		handleOperationsGet(w, r, db)
	case http.MethodPost:
		handleOperationsPost(w, r, db)
	case http.MethodPatch:
		handleOperationsPatch(w, r, db)
	default:
		writeOperationsJSON(w, http.StatusMethodNotAllowed, map[string]string{"message": "method not allowed"})
	}
}

func handleOperationsGet(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.URL.Query().Get("action") {
	case "summary":
		writeOperationsJSON(w, http.StatusOK, fetchSummary(db))
	case "support":
		writeOperationsJSON(w, http.StatusOK, fetchRows(db, `select id::text, type, customer, email, coalesce(receipt,''), subject, message, status, created_at from public.support_cases order by created_at desc limit 50`, scanSupportCase))
	case "payments":
		writeOperationsJSON(w, http.StatusOK, fetchRows(db, `select id::text, customer, email, coalesce(receipt,''), amount, status, created_at from public.payments order by created_at desc limit 50`, scanPayment))
	case "branches":
		writeOperationsJSON(w, http.StatusOK, fetchRows(db, `select id::text, code, name, address, status, created_at from public.branches order by created_at desc limit 50`, scanBranch))
	case "vehicles":
		writeOperationsJSON(w, http.StatusOK, fetchRows(db, `select id::text, plate, driver, status, created_at from public.vehicles order by created_at desc limit 50`, scanVehicle))
	case "routes":
		writeOperationsJSON(w, http.StatusOK, fetchRows(db, `select id::text, receipt, driver, origin, destination, status, recommendation, updated_at from public.courier_routes order by updated_at desc limit 50`, scanRoute))
	default:
		writeOperationsJSON(w, http.StatusBadRequest, map[string]string{"message": "Aksi tidak dikenal."})
	}
}

func handleOperationsPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var payload operationPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeOperationsJSON(w, http.StatusBadRequest, map[string]string{"message": "invalid JSON body"})
		return
	}
	trimOperationPayload(&payload)

	switch payload.Type {
	case "chat", "claim":
		if payload.Customer == "" || payload.Email == "" || payload.Subject == "" || payload.Message == "" {
			writeOperationsJSON(w, http.StatusBadRequest, map[string]string{"message": "Data chat/klaim belum lengkap."})
			return
		}
		var id string
		err := db.QueryRow(`insert into public.support_cases (type, customer, email, receipt, subject, message, status) values ($1,$2,$3,nullif($4,''),$5,$6,'Open') returning id::text`, payload.Type, payload.Customer, payload.Email, payload.Receipt, payload.Subject, payload.Message).Scan(&id)
		writeOperationResult(w, err, id, "Tiket berhasil dibuat.")
	case "payment":
		if payload.Customer == "" || payload.Email == "" || payload.Amount <= 0 {
			writeOperationsJSON(w, http.StatusBadRequest, map[string]string{"message": "Data pembayaran belum lengkap."})
			return
		}
		var id string
		err := db.QueryRow(`insert into public.payments (customer, email, receipt, amount, status) values ($1,$2,nullif($3,''),$4,'Menunggu Pembayaran') returning id::text`, payload.Customer, payload.Email, payload.Receipt, payload.Amount).Scan(&id)
		writeOperationResult(w, err, id, "Invoice ongkir dibuat.")
	case "branch":
		if payload.Code == "" || payload.Name == "" || payload.Address == "" {
			writeOperationsJSON(w, http.StatusBadRequest, map[string]string{"message": "Data cabang belum lengkap."})
			return
		}
		var id string
		err := db.QueryRow(`insert into public.branches (code, name, address, status) values ($1,$2,$3,'Aktif') returning id::text`, payload.Code, payload.Name, payload.Address).Scan(&id)
		writeOperationResult(w, err, id, "Cabang berhasil disimpan.")
	case "vehicle":
		if payload.Plate == "" || payload.Driver == "" {
			writeOperationsJSON(w, http.StatusBadRequest, map[string]string{"message": "Data kendaraan belum lengkap."})
			return
		}
		var id string
		err := db.QueryRow(`insert into public.vehicles (plate, driver, status) values ($1,$2,'Siap Jalan') returning id::text`, payload.Plate, payload.Driver).Scan(&id)
		writeOperationResult(w, err, id, "Kendaraan berhasil disimpan.")
	case "route":
		if payload.Receipt == "" || payload.Driver == "" || payload.Origin == "" || payload.Destination == "" {
			writeOperationsJSON(w, http.StatusBadRequest, map[string]string{"message": "Data rute belum lengkap."})
			return
		}
		recommendation := routeRecommendation(payload.Origin, payload.Destination)
		var id string
		err := db.QueryRow(`insert into public.courier_routes (receipt, driver, origin, destination, status, recommendation) values ($1,$2,$3,$4,'Rute Diambil',$5) returning id::text`, payload.Receipt, payload.Driver, payload.Origin, payload.Destination, recommendation).Scan(&id)
		writeOperationResult(w, err, id, "Rute kurir berhasil diambil.")
	default:
		writeOperationsJSON(w, http.StatusBadRequest, map[string]string{"message": "Tipe data tidak dikenal."})
	}
}

func handleOperationsPatch(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var payload operationPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeOperationsJSON(w, http.StatusBadRequest, map[string]string{"message": "invalid JSON body"})
		return
	}
	trimOperationPayload(&payload)
	if payload.ID == "" || payload.Type == "" || payload.Status == "" {
		writeOperationsJSON(w, http.StatusBadRequest, map[string]string{"message": "ID, tipe, dan status wajib diisi."})
		return
	}

	var err error
	switch payload.Type {
	case "support":
		_, err = db.Exec(`update public.support_cases set status=$1 where id=$2`, payload.Status, payload.ID)
	case "payment":
		_, err = db.Exec(`update public.payments set status=$1 where id=$2`, payload.Status, payload.ID)
	case "vehicle":
		_, err = db.Exec(`update public.vehicles set status=$1 where id=$2`, payload.Status, payload.ID)
	case "route":
		_, err = db.Exec(`update public.courier_routes set status=$1, updated_at=now() where id=$2`, payload.Status, payload.ID)
	default:
		writeOperationsJSON(w, http.StatusBadRequest, map[string]string{"message": "Tipe update tidak dikenal."})
		return
	}
	writeOperationResult(w, err, payload.ID, "Status berhasil diperbarui.")
}

func ensureOperationsSchema(db *sql.DB) error {
	statements := []string{
		`create table if not exists public.support_cases (id uuid primary key default gen_random_uuid(), type text not null, customer text not null, email text not null, receipt text, subject text not null, message text not null, status text not null default 'Open', created_at timestamptz not null default now())`,
		`create table if not exists public.payments (id uuid primary key default gen_random_uuid(), customer text not null, email text not null, receipt text, amount integer not null, status text not null default 'Menunggu Pembayaran', created_at timestamptz not null default now())`,
		`create table if not exists public.branches (id uuid primary key default gen_random_uuid(), code text not null, name text not null, address text not null, status text not null default 'Aktif', created_at timestamptz not null default now())`,
		`create table if not exists public.vehicles (id uuid primary key default gen_random_uuid(), plate text not null, driver text not null, status text not null default 'Siap Jalan', created_at timestamptz not null default now())`,
		`create table if not exists public.courier_routes (id uuid primary key default gen_random_uuid(), receipt text not null, driver text not null, origin text not null, destination text not null, status text not null default 'Rute Diambil', recommendation text not null, updated_at timestamptz not null default now(), created_at timestamptz not null default now())`,
	}
	for _, statement := range statements {
		if _, err := db.Exec(statement); err != nil {
			return err
		}
	}
	return nil
}

func fetchSummary(db *sql.DB) map[string]int {
	queries := map[string]string{
		"packages":  `select count(*) from public.manifests`,
		"pickups":   `select count(*) from public.pickups`,
		"tickets":   `select count(*) from public.support_cases where status <> 'Selesai'`,
		"branches":  `select count(*) from public.branches`,
		"vehicles":  `select count(*) from public.vehicles`,
		"couriers":  `select count(distinct driver) from public.courier_routes`,
		"revenue":   `select coalesce(sum(amount),0) from public.payments where status='Lunas'`,
		"unpaid":    `select coalesce(sum(amount),0) from public.payments where status <> 'Lunas'`,
	}
	result := map[string]int{}
	for key, query := range queries {
		var value int
		_ = db.QueryRow(query).Scan(&value)
		result[key] = value
	}
	return result
}

func fetchRows(db *sql.DB, query string, scanner func(*sql.Rows) (map[string]any, error)) any {
	rows, err := db.Query(query)
	if err != nil {
		return []map[string]any{}
	}
	defer rows.Close()
	items := []map[string]any{}
	for rows.Next() {
		item, err := scanner(rows)
		if err == nil {
			items = append(items, item)
		}
	}
	return items
}

func scanSupportCase(rows *sql.Rows) (map[string]any, error) {
	var id, typ, customer, email, receipt, subject, message, status string
	var createdAt time.Time
	err := rows.Scan(&id, &typ, &customer, &email, &receipt, &subject, &message, &status, &createdAt)
	return map[string]any{"id": id, "type": typ, "customer": customer, "email": email, "receipt": receipt, "subject": subject, "message": message, "status": status, "created_at": createdAt.Format(time.RFC3339)}, err
}

func scanPayment(rows *sql.Rows) (map[string]any, error) {
	var id, customer, email, receipt, status string
	var amount int
	var createdAt time.Time
	err := rows.Scan(&id, &customer, &email, &receipt, &amount, &status, &createdAt)
	return map[string]any{"id": id, "customer": customer, "email": email, "receipt": receipt, "amount": amount, "status": status, "created_at": createdAt.Format(time.RFC3339)}, err
}

func scanBranch(rows *sql.Rows) (map[string]any, error) {
	var id, code, name, address, status string
	var createdAt time.Time
	err := rows.Scan(&id, &code, &name, &address, &status, &createdAt)
	return map[string]any{"id": id, "code": code, "name": name, "address": address, "status": status, "created_at": createdAt.Format(time.RFC3339)}, err
}

func scanVehicle(rows *sql.Rows) (map[string]any, error) {
	var id, plate, driver, status string
	var createdAt time.Time
	err := rows.Scan(&id, &plate, &driver, &status, &createdAt)
	return map[string]any{"id": id, "plate": plate, "driver": driver, "status": status, "created_at": createdAt.Format(time.RFC3339)}, err
}

func scanRoute(rows *sql.Rows) (map[string]any, error) {
	var id, receipt, driver, origin, destination, status, recommendation string
	var updatedAt time.Time
	err := rows.Scan(&id, &receipt, &driver, &origin, &destination, &status, &recommendation, &updatedAt)
	return map[string]any{"id": id, "receipt": receipt, "driver": driver, "origin": origin, "destination": destination, "status": status, "recommendation": recommendation, "updated_at": updatedAt.Format(time.RFC3339)}, err
}

func trimOperationPayload(payload *operationPayload) {
	payload.ID = strings.TrimSpace(payload.ID)
	payload.Type = strings.TrimSpace(payload.Type)
	payload.Customer = strings.TrimSpace(payload.Customer)
	payload.Email = strings.ToLower(strings.TrimSpace(payload.Email))
	payload.Receipt = strings.ToUpper(strings.TrimSpace(payload.Receipt))
	payload.Subject = strings.TrimSpace(payload.Subject)
	payload.Message = strings.TrimSpace(payload.Message)
	payload.Status = strings.TrimSpace(payload.Status)
	payload.Name = strings.TrimSpace(payload.Name)
	payload.Code = strings.ToUpper(strings.TrimSpace(payload.Code))
	payload.Address = strings.TrimSpace(payload.Address)
	payload.Plate = strings.ToUpper(strings.TrimSpace(payload.Plate))
	payload.Driver = strings.TrimSpace(payload.Driver)
	payload.Origin = strings.TrimSpace(payload.Origin)
	payload.Destination = strings.TrimSpace(payload.Destination)
}

func routeRecommendation(origin string, destination string) string {
	return fmt.Sprintf("AI rekomendasi: mulai dari %s, prioritaskan jalur utama menuju %s, hindari jam padat 16.00-19.00, lalu update status setiap titik transit.", origin, destination)
}

func openOperationsDatabase() (*sql.DB, error) {
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

func writeOperationResult(w http.ResponseWriter, err error, id string, message string) {
	if err != nil {
		writeOperationsJSON(w, http.StatusInternalServerError, map[string]string{"message": "Data belum bisa disimpan."})
		return
	}
	writeOperationsJSON(w, http.StatusOK, map[string]string{"id": id, "message": message})
}

func writeOperationsJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}
