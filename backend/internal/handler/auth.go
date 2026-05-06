package handler

import (
	"context"

	"github.com/aulaflash/backend/internal/auth"
	"github.com/aulaflash/backend/internal/domain/repository"
	"github.com/gofiber/fiber/v3"
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



// FiberRegister is the Fiber version of Register
func (h *AuthHandler) FiberRegister(c fiber.Ctx) error {
	var req RegisterRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "corpo da requisicao invalido",
		})
	}

	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "email e senha obrigatorios",
		})
	}

	if req.Mode == "" {
		req.Mode = "student"
	}

	user, err := h.authService.Register(c.Context(), req.Name, req.Email, req.Password, req.Mode)
	if err != nil {
		if err.Error() == "user already exists" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "usuario ja existe com este email",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "erro ao registrar usuario",
		})
	}

	token, err := h.tokenService.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "erro ao gerar token",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"token": token,
		"user": fiber.Map{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"mode":  user.Mode,
		},
	})
}



// FiberLogin is the Fiber version of Login
func (h *AuthHandler) FiberLogin(c fiber.Ctx) error {
	var req LoginRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "corpo da requisicao invalido",
		})
	}

	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "email e senha obrigatorios",
		})
	}

	user, err := h.authService.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "email ou senha incorretos",
		})
	}

	token, err := h.tokenService.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "erro ao gerar token",
		})
	}

	return c.JSON(fiber.Map{
		"token": token,
		"user": fiber.Map{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"mode":  user.Mode,
		},
	})
}
