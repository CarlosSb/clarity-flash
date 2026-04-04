package handler

import (
	"encoding/csv"
	"net/http"
	"strings"

	"github.com/aulaflash/backend/internal/service"
)

// ExportHandler lida com exportacao de flashcards
type ExportHandler struct {
	processor *service.Processor
}

func NewExportHandler(processor *service.Processor) *ExportHandler {
	return &ExportHandler{processor: processor}
}

// ExportCSV exporta flashcards em formato CSV (compativel com Anki)
// GET /api/export/:id/csv
func (h *ExportHandler) ExportCSV(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("id")
	if sessionID == "" {
		http.Error(w, "id obrigatorio", http.StatusBadRequest)
		return
	}

	cards, err := h.processor.GetFlashcards(r.Context(), sessionID)
	if err != nil {
		http.Error(w, "flashcards nao encontrados", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=flashcards.csv")

	writer := csv.NewWriter(w)
	writer.Write([]string{"Front", "Back", "Difficulty"})

	for _, card := range cards {
		writer.Write([]string{card.Front, card.Back, difficultyLabel(card.Difficulty)})
	}

	writer.Flush()
}

// ExportText exporta flashcards em texto simples
// GET /api/export/:id/txt
func (h *ExportHandler) ExportText(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("id")
	if sessionID == "" {
		http.Error(w, "id obrigatorio", http.StatusBadRequest)
		return
	}

	cards, err := h.processor.GetFlashcards(r.Context(), sessionID)
	if err != nil {
		http.Error(w, "flashcards nao encontrados", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=flashcards.txt")

	for i, card := range cards {
		w.Write([]byte(formatCard(i, card)))
	}
}

func difficultyLabel(d int) string {
	switch d {
	case 1:
		return "easy"
	case 3:
		return "hard"
	default:
		return "medium"
	}
}

func formatCard(i int, card any) string {
	type Card struct {
		Front, Back string
	}
	c := card.(Card)
	var sb strings.Builder
	sb.WriteString("Card ")
	sb.WriteRune(rune(i + 1))
	sb.WriteString("\n")
	sb.WriteString("Q: ")
	sb.WriteString(c.Front)
	sb.WriteString("\n")
	sb.WriteString("A: ")
	sb.WriteString(c.Back)
	sb.WriteString("\n\n")
	return sb.String()
}
