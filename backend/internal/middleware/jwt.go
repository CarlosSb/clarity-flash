package middleware

import (
	"context"
	"strings"

	"github.com/aulaflash/backend/internal/auth"
	"github.com/gofiber/fiber/v3"
)



// contextKey is an unexported type for context keys
type contextKey string

const userIDKey contextKey = "user_id"

func contextWithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// GetUserIDFromContext extracts the user_id from context
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDKey).(string)
	return userID, ok
}

// FiberJWTOrFallbackAuth is the Fiber version of JWTOrFallbackAuth
func FiberJWTOrFallbackAuth(tokenService *auth.TokenService, fallbackHeader string) fiber.Handler {
	return func(c fiber.Ctx) error {
		var userID string

		// Try JWT first
		authHeader := string(c.Get("Authorization"))
		if authHeader != "" {
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenStr != authHeader {
				claims, err := tokenService.ValidateAccessToken(tokenStr)
				if err != nil {
					return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
						"error": "token invalido",
					})
				}
				userID = claims.UserID
			}
		}

		// Fallback to header or query param (for extension uploads)
		if userID == "" {
			userID = string(c.Get(fallbackHeader))
		}
		if userID == "" {
			userID = string(c.Query("user_id"))
		}
		if userID == "" {
			// For multipart form, we need to parse first
			if strings.Contains(string(c.Get("Content-Type")), "multipart/") {
				form, err := c.MultipartForm()
				if err == nil && form.Value != nil {
					if vals, exists := form.Value["user_id"]; exists && len(vals) > 0 {
						userID = vals[0]
					}
				}
			}
		}
		if userID == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "nao autorizado",
			})
		}

		// Add user_id to Fiber context locals
		c.Locals("user_id", userID)
		return c.Next()
	}
}
