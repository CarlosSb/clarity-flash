package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/aulaflash/backend/internal/auth"
)

// JWTOrFallbackAuth validates the Bearer token first. If no token,
// falls back to user_id from header or query param (for extension/legacy uploads).
// The resulting user_id is placed in the request context.
func JWTOrFallbackAuth(tokenService *auth.TokenService, fallbackHeader string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var userID string

		// Try JWT first
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenStr != authHeader {
				claims, err := tokenService.ValidateAccessToken(tokenStr)
				if err != nil {
					http.Error(w, `{"error":"token invalido"}`, http.StatusUnauthorized)
					return
				}
				userID = claims.UserID
			}
		}

		// Fallback to header or query param (for extension uploads)
		if userID == "" {
			userID = r.Header.Get(fallbackHeader)
		}
		if userID == "" {
			userID = r.URL.Query().Get("user_id")
		}
		if userID == "" {
			// For multipart form, we need to parse first
			if strings.Contains(r.Header.Get("Content-Type"), "multipart/") {
				err := r.ParseMultipartForm(50 << 20)
				if err == nil {
					userID = r.FormValue("user_id")
				}
			}
		}
		if userID == "" {
			http.Error(w, `{"error":"nao autorizado"}`, http.StatusUnauthorized)
			return
		}

		// Add user_id to request context
		ctx := contextWithUserID(r.Context(), userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

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
