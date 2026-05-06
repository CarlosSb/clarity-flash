package api

import (
	"github.com/aulaflash/backend/internal/auth"
	"github.com/aulaflash/backend/internal/handler"
	"github.com/aulaflash/backend/internal/middleware"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

// SetupFiberApp configura a aplicação Fiber com todas as rotas e middlewares
// Esta função é como o "mapa da casa" - define todas as portas de entrada (rotas)
// e os guardas de segurança (middlewares) da nossa aplicação
func SetupFiberApp(
	sessionHandler *handler.SessionHandler, // Gerencia sessões de áudio
	authHandler *handler.AuthHandler, // Gerencia autenticação (login/register)
	tokenService *auth.TokenService, // Serviço para gerar/verificar tokens JWT
	exportHandler *handler.ExportHandler, // Gerencia exportação de flashcards
) *fiber.App {
	// Criamos uma nova aplicação Fiber com configuração personalizada
	app := fiber.New(fiber.Config{
		// Handler de erro global - se algo der errado em qualquer rota,
		// esta função será chamada automaticamente
		ErrorHandler: func(c fiber.Ctx, err error) error {
			// Retorna erro 500 com mensagem genérica (não expõe detalhes internos)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error",
			})
		},
	})

	// 🔒 CORS middleware - permite que o frontend (que roda em porta diferente)
	// possa fazer requisições para o backend sem ser bloqueado pelo navegador
	// Ex: frontend em localhost:5173 → backend em localhost:8081
	app.Use(cors.New())

	// 💚 Health check - rota simples para verificar se o servidor está vivo
	// Muito útil para monitoramento e load balancers
	app.Get("/health", handler.HealthHandler)

	// 🔓 Rotas públicas - não precisam de autenticação
	// Qualquer pessoa pode se registrar ou fazer login
	app.Post("/api/auth/register", authHandler.FiberRegister) // Criar nova conta
	app.Post("/api/auth/login", authHandler.FiberLogin)       // Fazer login

	// 🔐 Grupo de rotas protegidas - TODAS precisam de autenticação JWT
	// O middleware JWTOrFallbackAuth verifica o token antes de executar o handler
	protected := app.Group("", JWTOrFallbackAuth(tokenService, "X-User-ID"))

	// 📤 Rotas de sessão - gerenciar gravações de áudio
	protected.Post("/api/sessions/upload", sessionHandler.FiberUpload) // Upload de áudio
	protected.Get("/api/sessions/:id", sessionHandler.FiberGetByID)    // Ver detalhes de uma sessão
	protected.Get("/api/sessions", sessionHandler.FiberListByUser)     // Listar sessões do usuário
	protected.Delete("/api/sessions/:id", sessionHandler.FiberDelete)  // Deletar sessão

	// 📊 Rotas de exportação - baixar flashcards em diferentes formatos
	protected.Get("/api/export/:id/csv", exportHandler.FiberExportCSV)  // CSV para Anki
	protected.Get("/api/export/:id/txt", exportHandler.FiberExportText) // Texto simples

	return app
}

// JWTOrFallbackAuth converts the existing JWT middleware to Fiber format
func JWTOrFallbackAuth(tokenService *auth.TokenService, userIDHeader string) fiber.Handler {
	return middleware.FiberJWTOrFallbackAuth(tokenService, userIDHeader)
}
