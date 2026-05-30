package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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

	if err := saveManifestToSupabase(payload); err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": fmt.Sprintf("Manifest %s berhasil diproses. Supabase belum aktif: %s", payload.Receipt, err.Error()),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("Manifest %s berhasil disimpan ke Supabase.", payload.Receipt),
	})
}

func saveManifestToSupabase(payload manifestRequest) error {
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
