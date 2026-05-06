package handler

import (
	"encoding/csv"
	"strings"

	"github.com/aulaflash/backend/internal/service"
	"github.com/gofiber/fiber/v3"
)

// ExportHandler lida com exportacao de flashcards
type ExportHandler struct {
	processor *service.Processor
}

func NewExportHandler(processor *service.Processor) *ExportHandler {
	return &ExportHandler{processor: processor}
}

// FiberExportCSV is the Fiber version of ExportCSV
func (h *ExportHandler) FiberExportCSV(c fiber.Ctx) error {
	sessionID := c.Params("id")
	if sessionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "id obrigatorio",
		})
	}

	cards, err := h.processor.GetFlashcards(c.Context(), sessionID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "flashcards nao encontrados",
		})
	}

	c.Set("Content-Type", "text/csv; charset=utf-8")
	c.Set("Content-Disposition", "attachment; filename=flashcards.csv")

	var csvContent strings.Builder
	writer := csv.NewWriter(&csvContent)
	writer.Write([]string{"Front", "Back", "Difficulty"})

	for _, card := range cards {
		writer.Write([]string{card.Front, card.Back, difficultyLabel(card.Difficulty)})
	}

	writer.Flush()
	return c.SendString(csvContent.String())
}

// FiberExportText is the Fiber version of ExportText
func (h *ExportHandler) FiberExportText(c fiber.Ctx) error {
	sessionID := c.Params("id")
	if sessionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "id obrigatorio",
		})
	}

	cards, err := h.processor.GetFlashcards(c.Context(), sessionID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "flashcards nao encontrados",
		})
	}

	c.Set("Content-Type", "text/plain; charset=utf-8")
	c.Set("Content-Disposition", "attachment; filename=flashcards.txt")

	var txtContent strings.Builder
	for i, card := range cards {
		txtContent.WriteString(formatCard(i, card))
	}

	return c.SendString(txtContent.String())
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
