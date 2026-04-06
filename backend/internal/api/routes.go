package api

import (
	"net/http"

	"github.com/aulaflash/backend/internal/auth"
	"github.com/aulaflash/backend/internal/handler"
	"github.com/aulaflash/backend/internal/middleware"
)

// Routes documenta todas as rotas disponiveis
var Routes = []Route{
	{Method: "GET", Path: "/health", Desc: "Health check"},
	{Method: "POST", Path: "/api/auth/register", Desc: "Registro de usuario"},
	{Method: "POST", Path: "/api/auth/login", Desc: "Login de usuario"},
	{Method: "POST", Path: "/api/sessions/upload", Desc: "Upload de audio para processamento", Auth: true},
	{Method: "GET", Path: "/api/sessions/{id}", Desc: "Detalhes de uma sessao", Auth: true},
	{Method: "GET", Path: "/api/sessions", Desc: "Lista sessoes do usuario", Auth: true},
	{Method: "DELETE", Path: "/api/sessions/{id}", Desc: "Deleta uma sessao", Auth: true},
	{Method: "GET", Path: "/api/export/{id}/csv", Desc: "Exporta flashcards em CSV (Anki)", Auth: true},
	{Method: "GET", Path: "/api/export/{id}/txt", Desc: "Exporta flashcards em texto", Auth: true},
}

type Route struct {
	Method string
	Path   string
	Desc   string
	Auth   bool // requires JWT auth
}

// SetupRouter configura todas as rotas do servidor
func SetupRouter(
	sessionHandler *handler.SessionHandler,
	authHandler *handler.AuthHandler,
	tokenService *auth.TokenService,
	exportHandler *handler.ExportHandler,
) http.Handler {
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("GET /health", handler.HealthHandler)

	// Public auth routes
	mux.HandleFunc("POST /api/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/auth/login", authHandler.Login)

	// Protected routes — JWT com fallback para user_id (extensao)
	authMW := func(next http.Handler) http.Handler {
		return middleware.JWTOrFallbackAuth(tokenService, "X-User-ID", next)
	}

	// Upload routes (legacy full-file upload)
	mux.Handle("POST /api/sessions/upload", authMW(http.HandlerFunc(sessionHandler.Upload)))
	mux.Handle("GET /api/sessions/{id}", authMW(http.HandlerFunc(sessionHandler.GetByID)))
	mux.Handle("GET /api/sessions", authMW(http.HandlerFunc(sessionHandler.ListByUser)))
	mux.Handle("DELETE /api/sessions/{id}", authMW(http.HandlerFunc(sessionHandler.Delete)))

	// Export routes
	mux.Handle("GET /api/export/{id}/csv", authMW(http.HandlerFunc(exportHandler.ExportCSV)))
	mux.Handle("GET /api/export/{id}/txt", authMW(http.HandlerFunc(exportHandler.ExportText)))

	return mux
}
