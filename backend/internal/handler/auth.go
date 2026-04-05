package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aulaflash/backend/internal/auth"
	"github.com/aulaflash/backend/internal/domain/repository"
)

type AuthService interface {
	Register(ctx context.Context, name, email, password, mode string) (*repository.User, error)
	Login(ctx context.Context, email, password string) (*repository.User, error)
}

type AuthHandler struct {
	authService  AuthService
	tokenService *auth.TokenService
}

func NewAuthHandler(authService AuthService, tokenService *auth.TokenService) *AuthHandler {
	return &AuthHandler{
		authService:  authService,
		tokenService: tokenService,
	}
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Mode     string `json:"mode"` // optional, defaults to "student"
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"corpo da requisicao invalido"}`, http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, `{"error":"email e senha obrigatorios"}`, http.StatusBadRequest)
		return
	}

	if req.Mode == "" {
		req.Mode = "student"
	}

	user, err := h.authService.Register(r.Context(), req.Name, req.Email, req.Password, req.Mode)
	if err != nil {
		if err.Error() == "user already exists" {
			http.Error(w, `{"error":"usuario ja existe com este email"}`, http.StatusConflict)
			return
		}
		http.Error(w, `{"error":"erro ao registrar usuario"}`, http.StatusInternalServerError)
		return
	}

	token, err := h.tokenService.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		http.Error(w, `{"error":"erro ao gerar token"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"token": token,
		"user": map[string]any{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"mode":  user.Mode,
		},
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"corpo da requisicao invalido"}`, http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, `{"error":"email e senha obrigatorios"}`, http.StatusBadRequest)
		return
	}

	user, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, `{"error":"email ou senha incorretos"}`, http.StatusUnauthorized)
		return
	}

	token, err := h.tokenService.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		http.Error(w, `{"error":"erro ao gerar token"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"token": token,
		"user": map[string]any{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"mode":  user.Mode,
		},
	})
}
