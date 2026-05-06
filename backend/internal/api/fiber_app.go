package api

import (
	"github.com/aulaflash/backend/internal/auth"
	"github.com/aulaflash/backend/internal/handler"
	"github.com/aulaflash/backend/internal/middleware"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

// SetupFiberApp configura a aplicação Fiber com todas as rotas e middlewares
func SetupFiberApp(
	sessionHandler *handler.SessionHandler,
	authHandler *handler.AuthHandler,
	tokenService *auth.TokenService,
	exportHandler *handler.ExportHandler,
) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error",
			})
		},
	})

	// CORS middleware
	app.Use(cors.New())

	// Health check
	app.Get("/health", handler.HealthHandler)

	// Public auth routes
	app.Post("/api/auth/register", authHandler.FiberRegister)
	app.Post("/api/auth/login", authHandler.FiberLogin)

	// Protected routes group with JWT auth middleware
	protected := app.Group("", JWTOrFallbackAuth(tokenService, "X-User-ID"))

	// Session routes
	protected.Post("/api/sessions/upload", sessionHandler.FiberUpload)
	protected.Get("/api/sessions/:id", sessionHandler.FiberGetByID)
	protected.Get("/api/sessions", sessionHandler.FiberListByUser)
	protected.Delete("/api/sessions/:id", sessionHandler.FiberDelete)

	// Export routes
	protected.Get("/api/export/:id/csv", exportHandler.FiberExportCSV)
	protected.Get("/api/export/:id/txt", exportHandler.FiberExportText)

	return app
}

// JWTOrFallbackAuth converts the existing JWT middleware to Fiber format
func JWTOrFallbackAuth(tokenService *auth.TokenService, userIDHeader string) fiber.Handler {
	return middleware.FiberJWTOrFallbackAuth(tokenService, userIDHeader)
}