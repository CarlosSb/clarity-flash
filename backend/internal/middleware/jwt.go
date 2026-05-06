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

// FiberJWTOrFallbackAuth - Middleware de autenticação inteligente
// Esta função é como um "porteiro" que verifica se a pessoa pode entrar
// no sistema. Ela é "inteligente" porque tenta múltiplas formas de autenticação.
func FiberJWTOrFallbackAuth(tokenService *auth.TokenService, fallbackHeader string) fiber.Handler {
	return func(c fiber.Ctx) error {
		var userID string

		// 🥇 TENTATIVA 1: Token JWT (padrão web moderno)
		// Verifica se existe header "Authorization: Bearer <token>"
		authHeader := string(c.Get("Authorization"))
		if authHeader != "" {
			// Remove o prefixo "Bearer " para pegar só o token
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenStr != authHeader { // Se conseguiu remover "Bearer "
				// Valida o token e extrai as informações (claims)
				claims, err := tokenService.ValidateAccessToken(tokenStr)
				if err != nil {
					return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
						"error": "token invalido",
					})
				}
				userID = claims.UserID // Sucesso! Pegamos o userID do token
			}
		}

		// 🥈 TENTATIVA 2: Header customizado (para extensões Chrome)
		// Algumas requisições (como da extensão) podem não ter JWT,
		// então aceitamos um header simples com o userID
		if userID == "" {
			userID = string(c.Get(fallbackHeader)) // Ex: X-User-ID
		}

		// 🥉 TENTATIVA 3: Query parameter (fallback adicional)
		// Para casos extremos onde nem header conseguimos passar
		if userID == "" {
			userID = string(c.Query("user_id")) // Ex: ?user_id=123
		}

		// 🔍 TENTATIVA 4: Campo no formulário multipart (para uploads)
		// Quando é upload de arquivo, o userID pode vir dentro do form
		if userID == "" {
			contentType := string(c.Get("Content-Type"))
			if strings.Contains(contentType, "multipart/") {
				// Parse do formulário multipart
				form, err := c.MultipartForm()
				if err == nil && form.Value != nil {
					// Procura pelo campo "user_id" no form
					if vals, exists := form.Value["user_id"]; exists && len(vals) > 0 {
						userID = vals[0] // Primeiro valor encontrado
					}
				}
			}
		}

		// ❌ Se nenhuma tentativa funcionou, bloquear acesso
		if userID == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "nao autorizado",
			})
		}

		// ✅ Sucesso! Guardar userID no contexto da requisição
		// Agora todos os handlers podem acessar c.Locals("user_id")
		c.Locals("user_id", userID)
		return c.Next() // Permitir que a requisição continue
	}
}
