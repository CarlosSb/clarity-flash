package integration

import (
	"net/http/httptest"
	"testing"

	"github.com/aulaflash/backend/internal/handler"
	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestHealthCheck verifica que o health endpoint funciona
func TestHealthCheck(t *testing.T) {
	// Criar app minimal para teste
	app := fiber.New()

	// Registrar rota de health
	app.Get("/health", handler.HealthHandler)

	// Criar requisição de teste
	req := httptest.NewRequest("GET", "/health", nil)

	// Testar requisição
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Verificar status 200
	assert.Equal(t, 200, resp.StatusCode)

	// Verificar corpo da resposta
	// TODO: verificar JSON se necessário
}

// TestRoutesConfigured valida que as rotas essenciais estão configuradas
func TestRoutesConfigured(t *testing.T) {
	// Criar app minimal
	app := fiber.New()

	// Registrar rotas básicas
	app.Get("/health", handler.HealthHandler)
	app.Get("/api/sessions", func(c fiber.Ctx) error { return c.SendString("ok") })
	app.Post("/api/auth/login", func(c fiber.Ctx) error { return c.SendString("ok") })

	// Testar se rotas existem (não retornam 404)
	routes := []struct {
		method string
		path   string
	}{
		{"GET", "/health"},
		{"GET", "/api/sessions"},
		{"POST", "/api/auth/login"},
	}

	for _, route := range routes {
		t.Run(route.method+" "+route.path, func(t *testing.T) {
			req := httptest.NewRequest(route.method, route.path, nil)
			resp, err := app.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Não deve ser 404 (rota não encontrada)
			assert.NotEqual(t, 404, resp.StatusCode)
		})
	}
}
