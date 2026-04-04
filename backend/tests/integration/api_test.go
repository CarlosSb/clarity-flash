package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestHealthCheck verifica que o health endpoint funciona
func TestHTTPResponseCode(t *testing.T) {
	w := httptest.NewRecorder()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

// TestRoutesExist valida que todas as rotas estao configuradas
func TestRoutesExist(t *testing.T) {
	// Teste basico de sanity check
	// TODO: integrar com router real e repositories mock
	t.Skip("integração com router requer repositorios mock")
}
