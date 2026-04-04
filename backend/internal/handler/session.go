package handler

import (
	"encoding/json"
	"net/http"

	"github.com/aulaflash/backend/internal/domain/repository"
	"github.com/aulaflash/backend/internal/service"
)

// SessionHandler lida com HTTP requests de sessoes
type SessionHandler struct {
	processor *service.Processor
}

func NewSessionHandler(processor *service.Processor) *SessionHandler {
	return &SessionHandler{processor: processor}
}

// Upload recebe o audio e inicia o processamento
// POST /api/sessions/upload
func (h *SessionHandler) Upload(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(50 << 20); err != nil { // 50MB max no form
		http.Error(w, "erro ao ler formulario", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("audio")
	if err != nil {
		http.Error(w, "arquivo nao encontrado", http.StatusBadRequest)
		return
	}
	defer file.Close()

	userID := r.FormValue("user_id")
	if userID == "" {
		userID = "anonymous"
	}

	mode := r.FormValue("mode")
	if mode == "" {
		mode = "student"
	}

	// Cria sessao no banco
	session := &repository.Session{
		UserID: userID,
		Title:  header.Filename,
		Mode:   mode,
		Status: "processing",
	}

	if err := h.processor.Process(r.Context(), session, file, header); err != nil {
		http.Error(w, "erro ao processar audio: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message":    "audio recebido com sucesso",
		"session_id": session.ID,
		"status":     session.Status,
	})
}

// GetByID retorna uma sessao com seus dados processados
// GET /api/sessions/:id
func (h *SessionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "id obrigatorio", http.StatusBadRequest)
		return
	}

	session, err := h.processor.GetSession(r.Context(), id)
	if err != nil {
		http.Error(w, "sessao nao encontrada", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

// ListByUser retorna todas as sessoes de um usuario
// GET /api/sessions?user_id=xxx
func (h *SessionHandler) ListByUser(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user_id obrigatorio", http.StatusBadRequest)
		return
	}

	sessions, err := h.processor.ListSessions(r.Context(), userID)
	if err != nil {
		http.Error(w, "erro ao listar sessoes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sessions)
}

// Delete remove uma sessao e seus dados
// DELETE /api/sessions/:id
func (h *SessionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "id obrigatorio", http.StatusBadRequest)
		return
	}

	if err := h.processor.DeleteSession(r.Context(), id); err != nil {
		http.Error(w, "erro ao deletar sessao", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
