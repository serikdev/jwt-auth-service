package handler

import (
	"encoding/json"
	"fmt"
	"jwt-service/internal/entities"
	"net/http"

	"github.com/google/uuid"
)

// --- ВСПОМОГАТЕЛЬНЫЕ МОДЕЛИ ---------------------------------------------

// TokenPair описывает пару выдаваемых токенов.
// swagger:model
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// RefreshRequest тело запроса на обновление токенов.
// swagger:model
type RefreshRequest struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// ------------------------------------------------------------------------

type AuthService interface {
	GenerateTokens(userGUID, ip string) (*entities.TokenPair, error)
	RefreshTokens(oldAccessToken, refreshToken, clientIP string) (*entities.TokenPair, error)
}

type AuthHandler struct{ authUC AuthService }

func NewAuthHandler(authUC AuthService) *AuthHandler { return &AuthHandler{authUC: authUC} }

// HandleAuth выдаёт пару JWT.
//
// @Summary      Login / первичная авторизация
// @Description  Принимает GUID (или генерирует новый) и возвращает Access‑/Refresh‑пару.
// @Tags         auth
// @Produce      json
// @Param        guid  query  string  false  "Клиентский GUID"
// @Success      200   {object}  TokenPair
// @Failure      400   {string}  string  "invalid guid format"
// @Failure      500   {string}  string  "internal error"
// @Router       /auth [get]
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
	_ = json.NewEncoder(w).Encode(tokens)
}

// HandleRefresh обновляет пару JWT.
//
// @Summary      Refresh tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        tokens  body  RefreshRequest  true  "Старая пара токенов"
// @Success      200     {object}  TokenPair
// @Failure      400     {string}  string  "invalid request body"
// @Failure      401     {string}  string  "unauthorized"
// @Router       /refresh [post]
func (h *AuthHandler) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
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
	_ = json.NewEncoder(w).Encode(tokens)
}

// ------------------------------------------------------------------------

func getClientIP(r *http.Request) string {
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		return fwd
	}
	return r.RemoteAddr
}
