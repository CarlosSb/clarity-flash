package api

import (
	"net/http"

	"github.com/aulaflash/backend/internal/handler"
)

// SetupRouter configura todas as rotas do servidor
func SetupRouter(
	sessionHandler *handler.SessionHandler,
	exportHandler *handler.ExportHandler,
) http.Handler {
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("GET /health", handler.HealthHandler)

	// Sessions
	mux.HandleFunc("POST /api/sessions/upload", sessionHandler.Upload)
	mux.HandleFunc("GET /api/sessions/{id}", sessionHandler.GetByID)
	mux.HandleFunc("GET /api/sessions", sessionHandler.ListByUser)
	mux.HandleFunc("DELETE /api/sessions/{id}", sessionHandler.Delete)

	// Export
	mux.HandleFunc("GET /api/export/{id}/csv", exportHandler.ExportCSV)
	mux.HandleFunc("GET /api/export/{id}/txt", exportHandler.ExportText)

	return mux
}
