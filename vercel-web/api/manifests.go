package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type manifestRequest struct {
	Receipt  string `json:"receipt"`
	Status   string `json:"status"`
	Location string `json:"location"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		manifests, err := listManifestsFromDB()
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "Data status kiriman belum bisa dimuat."})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(manifests)
		return
	}

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

	if err := saveManifestToSupabase(payload); err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": fmt.Sprintf("Manifest %s berhasil diproses. Supabase belum aktif: %s", payload.Receipt, err.Error()),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("Data resi %s berhasil disimpan.", payload.Receipt),
	})
}

func saveManifestToSupabase(payload manifestRequest) error {
	if err := saveManifestToDB(payload); err == nil {
		return nil
	}

	url := strings.TrimRight(os.Getenv("SUPABASE_URL"), "/")
	key := os.Getenv("SUPABASE_SERVICE_ROLE_KEY")
	if url == "" || key == "" {
		return fmt.Errorf("environment variables belum diset")
	}

	body, err := json.Marshal(map[string]string{
		"receipt":  payload.Receipt,
		"status":   payload.Status,
		"location": payload.Location,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url+"/rest/v1/manifests?on_conflict=receipt", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("apikey", key)
	req.Header.Set("Authorization", "Bearer "+key)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "resolution=merge-duplicates")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request ditolak: %s", string(respBody))
	}
	return nil
}

func saveManifestToDB(payload manifestRequest) error {
	db, err := openManifestDatabase()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(
		`insert into public.manifests (receipt, status, location)
		 values ($1, $2, $3)
		 on conflict (receipt)
		 do update set status = excluded.status, location = excluded.location, updated_at = now()`,
		payload.Receipt,
		payload.Status,
		payload.Location,
	)
	return err
}

func listManifestsFromDB() ([]map[string]any, error) {
	db, err := openManifestDatabase()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(
		`select receipt, status, coalesce(location, ''), updated_at
		 from public.manifests
		 order by updated_at desc
		 limit 20`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	manifests := []map[string]any{}
	for rows.Next() {
		var receipt, status, location string
		var updatedAt time.Time
		if err := rows.Scan(&receipt, &status, &location, &updatedAt); err != nil {
			return nil, err
		}
		manifests = append(manifests, map[string]any{
			"receipt":    receipt,
			"status":     status,
			"location":   location,
			"updated_at": updatedAt.Format(time.RFC3339),
		})
	}
	return manifests, rows.Err()
}

func openManifestDatabase() (*sql.DB, error) {
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
