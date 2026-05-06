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



// FiberUpload - Handler para upload de áudio
// Esta é a função mais importante! Recebe o arquivo de áudio e inicia
// todo o processo de IA (transcrição, resumo, flashcards)
func (h *SessionHandler) FiberUpload(c fiber.Ctx) error {
	// 📁 PASSO 1: Receber o arquivo enviado pelo usuário
	// O arquivo vem como "multipart/form-data" (formulário com arquivo)
	file, err := c.FormFile("audio") // Campo chamado "audio" no form
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "arquivo nao encontrado", // Arquivo não foi enviado
		})
	}

	// 📖 PASSO 2: Abrir o arquivo para leitura
	// O Fiber nos dá o arquivo, mas precisamos abrir para poder ler
	src, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "erro ao abrir arquivo",
		})
	}
	defer src.Close() // IMPORTANTE: fechar arquivo quando terminar

	// 👤 PASSO 3: Identificar o usuário (vem do middleware de auth)
	// O middleware já verificou autenticação e colocou userID no contexto
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		userID = "anonymous" // Fallback para usuários não logados
	}

	// ⚙️ PASSO 4: Pegar configurações opcionais
	// Modo de uso: "student" (aprendizado) ou "professional" (trabalho)
	mode := string(c.FormValue("mode")) // Campo opcional do formulário
	if mode == "" {
		mode = "student" // Valor padrão
	}

	// 💾 PASSO 5: Criar registro da sessão no banco
	// Antes mesmo de processar, salvamos que uma sessão começou
	session := &repository.Session{
		UserID: userID,           // Quem enviou
		Title:  file.Filename,    // Nome do arquivo como título
		Mode:   mode,             // Modo de uso
		Status: "processing",     // Status inicial: processando
	}

	// 🤖 PASSO 6: Iniciar processamento com IA
	// Esta é a parte "mágica" - envia para o Processor que vai:
	// 1. Salvar arquivo, 2. Transcrever, 3. Gerar resumo, 4. Criar flashcards
	if err := h.processor.Process(c.Context(), session, src, file); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "erro ao processar audio: " + err.Error(),
		})
	}

	// ✅ PASSO 7: Responder sucesso
	// Retornar informações da sessão criada
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":    "audio recebido com sucesso",
		"session_id": session.ID,    // ID único da sessão (gerado automaticamente)
		"status":     session.Status, // "processing" inicialmente
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
