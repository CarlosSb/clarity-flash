package api

import (
	"net/http"

	"github.com/aulaflash/backend/internal/auth"
	"github.com/aulaflash/backend/internal/handler"
	"github.com/aulaflash/backend/internal/middleware"
	"github.com/aulaflash/backend/internal/websocket"
)

// Routes documenta todas as rotas disponiveis
var Routes = []Route{
	{Method: "GET", Path: "/health", Desc: "Health check"},
	{Method: "GET", Path: "/ws", Desc: "WebSocket para eventos em tempo real"},
	{Method: "GET", Path: "/ws-stream", Desc: "WebSocket stream de audio (frames binarios, extensao)"},
	{Method: "POST", Path: "/api/auth/register", Desc: "Registro de usuario"},
	{Method: "POST", Path: "/api/auth/login", Desc: "Login de usuario"},
	{Method: "POST", Path: "/api/sessions/upload", Desc: "Upload de audio para processamento", Auth: true},
	{Method: "POST", Path: "/api/sessions/stream-init", Desc: "Cria sessao vazia para streaming (Chrome extension)"},
	{Method: "GET", Path: "/api/sessions/{id}", Desc: "Detalhes de uma sessao", Auth: true},
	{Method: "GET", Path: "/api/sessions", Desc: "Lista sessoes do usuario", Auth: true},
	{Method: "DELETE", Path: "/api/sessions/{id}", Desc: "Deleta uma sessao", Auth: true},
	{Method: "PATCH", Path: "/api/sessions/{id}/audio-chunk", Desc: "Recebe chunk de audio em streaming (extensao, fallback HTTP)"},
	{Method: "POST", Path: "/api/sessions/{id}/audio-complete", Desc: "Sinaliza fim do streaming e inicia processamento"},
	{Method: "POST", Path: "/api/sessions/{id}/upload-complete", Desc: "Upload completo de audio (fallback offline da extensao)"},
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
	streamHandler *handler.StreamHandler,
	wsHub *websocket.Hub,
) http.Handler {
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("GET /health", handler.HealthHandler)

	// WebSocket (no auth needed — identified by user_id query param)
	mux.HandleFunc("GET /ws", wsHub.WSHandler())

	// Audio streaming WebSocket (raw binary frames from extension)
	mux.HandleFunc("GET /ws-stream", streamHandler.WSStream)

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

	// Stream-init route (creates empty session for chunk streaming)
	mux.HandleFunc("POST /api/sessions/stream-init", sessionHandler.StreamInit)

	// Streaming routes (Chrome extension real-time)
	mux.Handle("PATCH /api/sessions/{id}/audio-chunk", http.HandlerFunc(streamHandler.AudioChunk))
	mux.Handle("POST /api/sessions/{id}/audio-complete", http.HandlerFunc(streamHandler.AudioComplete))
	mux.Handle("POST /api/sessions/{id}/upload-complete", http.HandlerFunc(streamHandler.UploadComplete))

	// Export routes
	mux.Handle("GET /api/export/{id}/csv", authMW(http.HandlerFunc(exportHandler.ExportCSV)))
	mux.Handle("GET /api/export/{id}/txt", authMW(http.HandlerFunc(exportHandler.ExportText)))

	return mux
}
