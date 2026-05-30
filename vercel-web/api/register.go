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

type registerRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeRegisterJSON(w, http.StatusMethodNotAllowed, map[string]string{"message": "method not allowed"})
		return
	}

	var payload registerRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeRegisterJSON(w, http.StatusBadRequest, map[string]string{"message": "invalid JSON body"})
		return
	}

	payload.Name = strings.TrimSpace(payload.Name)
	payload.Email = strings.ToLower(strings.TrimSpace(payload.Email))
	if payload.Name == "" || payload.Email == "" || len(payload.Password) < 6 {
		writeRegisterJSON(w, http.StatusBadRequest, map[string]string{"message": "Nama, email, dan password minimal 6 karakter wajib diisi."})
		return
	}

	url, key, err := supabaseRegisterConfig()
	if err != nil {
		writeRegisterJSON(w, http.StatusOK, map[string]string{
			"id":    "local-customer",
			"name":  payload.Name,
			"email": payload.Email,
			"role":  "customer",
		})
		return
	}

	body, _ := json.Marshal(map[string]string{
		"name":     payload.Name,
		"email":    payload.Email,
		"password": payload.Password,
		"role":     "customer",
	})
	req, err := http.NewRequest(http.MethodPost, url+"/rest/v1/app_users", bytes.NewReader(body))
	if err != nil {
		writeRegisterJSON(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}
	setSupabaseRegisterHeaders(req, key)
	req.Header.Set("Prefer", "return=representation")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		writeRegisterJSON(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		writeRegisterJSON(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}
	if resp.StatusCode >= 300 {
		writeRegisterJSON(w, http.StatusBadRequest, map[string]string{"message": "Email sudah terdaftar atau Supabase menolak request."})
		return
	}

	var users []map[string]any
	if err := json.Unmarshal(respBody, &users); err != nil || len(users) == 0 {
		writeRegisterJSON(w, http.StatusOK, map[string]string{"name": payload.Name, "email": payload.Email, "role": "customer"})
		return
	}
	writeRegisterJSON(w, http.StatusOK, users[0])
}

func supabaseRegisterConfig() (string, string, error) {
	url := strings.TrimRight(os.Getenv("SUPABASE_URL"), "/")
	key := os.Getenv("SUPABASE_SERVICE_ROLE_KEY")
	if url == "" || key == "" {
		return "", "", fmt.Errorf("supabase env is not configured")
	}
	return url, key, nil
}

func setSupabaseRegisterHeaders(req *http.Request, key string) {
	req.Header.Set("apikey", key)
	req.Header.Set("Authorization", "Bearer "+key)
	req.Header.Set("Content-Type", "application/json")
}

func writeRegisterJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}
