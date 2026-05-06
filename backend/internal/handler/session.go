package handler

import (
	"github.com/aulaflash/backend/internal/domain/repository"
	"github.com/aulaflash/backend/internal/service"
	"github.com/gofiber/fiber/v3"
)

// SessionHandler lida com HTTP requests de sessoes
type SessionHandler struct {
	processor *service.Processor
}

func NewSessionHandler(processor *service.Processor) *SessionHandler {
	return &SessionHandler{processor: processor}
}



// FiberUpload is the Fiber version of Upload
func (h *SessionHandler) FiberUpload(c fiber.Ctx) error {
	file, err := c.FormFile("audio")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "arquivo nao encontrado",
		})
	}

	src, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "erro ao abrir arquivo",
		})
	}
	defer src.Close()

	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		userID = "anonymous"
	}

	mode := string(c.FormValue("mode"))
	if mode == "" {
		mode = "student"
	}

	// Cria sessao no banco
	session := &repository.Session{
		UserID: userID,
		Title:  file.Filename,
		Mode:   mode,
		Status: "processing",
	}

  if err := h.processor.Process(c.Context(), session, src, file); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "erro ao processar audio: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":    "audio recebido com sucesso",
		"session_id": session.ID,
		"status":     session.Status,
	})
}



// FiberGetByID is the Fiber version of GetByID
func (h *SessionHandler) FiberGetByID(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "id obrigatorio",
		})
	}

	session, err := h.processor.GetSession(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "sessao nao encontrada",
		})
	}

	return c.JSON(session)
}



// FiberListByUser is the Fiber version of ListByUser
func (h *SessionHandler) FiberListByUser(c fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user_id obrigatorio",
		})
	}

	sessions, err := h.processor.ListSessions(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "erro ao listar sessoes",
		})
	}

	return c.JSON(sessions)
}



// FiberDelete is the Fiber version of Delete
func (h *SessionHandler) FiberDelete(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "id obrigatorio",
		})
	}

	if err := h.processor.DeleteSession(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "erro ao deletar sessao",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
