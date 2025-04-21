package handler

import (
	"encoding/json"
	"fmt"
	"jwt-service/internal/usecases"
	"net/http"

	"github.com/google/uuid"
)

type AuthHandler struct {
	authUC *usecases.AuthUseCase
}

func NewAuthHandler(authUC *usecases.AuthUseCase) *AuthHandler {
	return &AuthHandler{authUC: authUC}
}

func (h *AuthHandler) HandleAuth(w http.ResponseWriter, r *http.Request) {
	guid := r.URL.Query().Get("guid")

	if guid == "" {
		guid = uuid.New().String()
		fmt.Println("Generated new GUID", guid)
	}
	if _, err := uuid.Parse(guid); err != nil {
		http.Error(w, "invalid guid format", http.StatusBadRequest)
		return
	}

	clientIP := getClientIP(r)
	tokens, err := h.authUC.GenerateTokens(guid, clientIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokens)
}

func (h *AuthHandler) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	clientIP := getClientIP(r)
	tokens, err := h.authUC.RefreshTokens(req.AccessToken, req.RefreshToken, clientIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokens)
}

func getClientIP(r *http.Request) string {
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}
