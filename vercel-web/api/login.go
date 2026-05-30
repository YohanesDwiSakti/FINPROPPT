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

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type appUser struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeLoginJSON(w, http.StatusMethodNotAllowed, map[string]string{"message": "method not allowed"})
		return
	}

	var payload loginRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeLoginJSON(w, http.StatusBadRequest, map[string]string{"message": "invalid JSON body"})
		return
	}

	payload.Email = strings.ToLower(strings.TrimSpace(payload.Email))
	payload.Role = strings.TrimSpace(payload.Role)
	if payload.Email == "" || payload.Password == "" || (payload.Role != "customer" && payload.Role != "admin") {
		writeLoginJSON(w, http.StatusBadRequest, map[string]string{"message": "email, password, and role are required"})
		return
	}

	users, err := fetchUsersByEmail(payload.Email)
	if err != nil {
		writeLoginJSON(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}

	for _, user := range users {
		if user.Password == payload.Password && user.Role == payload.Role {
			writeLoginJSON(w, http.StatusOK, map[string]string{
				"id":    user.ID,
				"name":  user.Name,
				"email": user.Email,
				"role":  user.Role,
			})
			return
		}
	}

	writeLoginJSON(w, http.StatusUnauthorized, map[string]string{"message": "Email, password, atau role tidak sesuai."})
}

func fetchUsersByEmail(email string) ([]appUser, error) {
	url, key, err := supabaseLoginConfig()
	if err != nil {
		return demoLoginUsers(email), nil
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/rest/v1/app_users?email=eq.%s&select=id,name,email,password,role", url, email), nil)
	if err != nil {
		return nil, err
	}
	setSupabaseLoginHeaders(req, key)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("supabase login query failed: %s", string(body))
	}

	var users []appUser
	if err := json.Unmarshal(body, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func demoLoginUsers(email string) []appUser {
	users := map[string]appUser{
		"admin@tiki.test":    {ID: "demo-admin", Name: "Admin Hub Denpasar", Email: "admin@tiki.test", Password: "admin123", Role: "admin"},
		"customer@tiki.test": {ID: "demo-customer", Name: "Customer Demo", Email: "customer@tiki.test", Password: "customer123", Role: "customer"},
	}
	if user, ok := users[email]; ok {
		return []appUser{user}
	}
	return nil
}

func supabaseLoginConfig() (string, string, error) {
	url := strings.TrimRight(os.Getenv("SUPABASE_URL"), "/")
	key := os.Getenv("SUPABASE_SERVICE_ROLE_KEY")
	if url == "" || key == "" {
		return "", "", fmt.Errorf("supabase env is not configured")
	}
	return url, key, nil
}

func setSupabaseLoginHeaders(req *http.Request, key string) {
	req.Header.Set("apikey", key)
	req.Header.Set("Authorization", "Bearer "+key)
	req.Header.Set("Content-Type", "application/json")
}

func writeLoginJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func postLoginJSON(url string, key string, payload any) (*http.Response, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	setSupabaseLoginHeaders(req, key)
	req.Header.Set("Prefer", "return=representation")
	return http.DefaultClient.Do(req)
}
