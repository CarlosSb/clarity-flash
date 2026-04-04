package middleware

import (
	"net/http"
	"os"
	"strings"
)

// SimpleAuth middleware basico com API key (MVP)
func SimpleAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := os.Getenv("API_KEY")
		if apiKey == "" {
			// Sem API key configurada = skip auth
			next.ServeHTTP(w, r)
			return
		}

		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") || strings.TrimPrefix(auth, "Bearer ") != apiKey {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error":"nao autorizado"}`))
			return
		}

		next.ServeHTTP(w, r)
	})
}
