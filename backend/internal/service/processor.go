package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"mime/multipart"

	"github.com/aulaflash/backend/internal/domain/model"
	"github.com/aulaflash/backend/internal/domain/repository"
	"github.com/aulaflash/backend/pkg/audio"
	"github.com/aulaflash/backend/pkg/llm"
	"github.com/aulaflash/backend/pkg/storage"
	"github.com/aulaflash/backend/pkg/stt"
)

// Processor orquestra o fluxo: upload -> transcricao -> resumo -> flashcards
type Processor struct {
	sessionRepo   repository.SessionRepository
	flashcardRepo repository.FlashcardRepository
	storage       *storage.LocalStorage
	audioProc     *audio.Processor
	sttClient     *stt.GroqClient
	llmClient     llm.LLMClient
}

func NewProcessor(
	sessionRepo repository.SessionRepository,
	flashcardRepo repository.FlashcardRepository,
	store *storage.LocalStorage,
	audioProc *audio.Processor,
	sttClient *stt.GroqClient,
	llmClient llm.LLMClient,
) *Processor {
	return &Processor{
		sessionRepo:   sessionRepo,
		flashcardRepo: flashcardRepo,
		storage:       store,
		audioProc:     audioProc,
		sttClient:     sttClient,
		llmClient:     llmClient,
	}
}

// Process recebe o audio, transcreve, gera resumo e flashcards
func (p *Processor) Process(ctx context.Context, session *repository.Session, file multipart.File, header *multipart.FileHeader) error {
	// Gera ID unico
	id, err := generateID()
	if err != nil {
		return fmt.Errorf("gerar id: %w", err)
	}
	session.ID = id

	// Salva audio original
	audioPath := id + "_" + header.Filename
	path, err := p.storage.Save(file, header, audioPath)
	if err != nil {
		return fmt.Errorf("salvar audio: %w", err)
	}

	// Cria sessao no banco
	if err := p.sessionRepo.Create(ctx, session); err != nil {
		return fmt.Errorf("criar sessao: %w", err)
	}

	// Processa de forma assincrona (em producao, isso seria um worker)
	go func() {
		if err := p.runPipeline(context.Background(), session, path); err != nil {
			_ = p.sessionRepo.UpdateStatus(context.Background(), session.ID, "failed")
		}
	}()

	return nil
}

// runPipeline executa: transcricao -> resumo -> flashcards
func (p *Processor) runPipeline(ctx context.Context, session *repository.Session, audioPath string) error {
	// Valida audio
	if err := p.audioProc.ValidateAudio(audioPath); err != nil {
		return err
	}

	// Converte para WAV
	wavPath, err := p.audioProc.ConvertToWAV(audioPath)
	if err != nil {
		return fmt.Errorf("converter audio: %w", err)
	}
	defer p.audioProc.Cleanup(wavPath)

	// Step 1: Transcricao via Groq Whisper
	transcript, err := p.sttClient.Transcribe(wavPath)
	if err != nil {
		return fmt.Errorf("transcricao: %w", err)
	}

	if err := p.sessionRepo.UpdateTranscript(ctx, session.ID, transcript); err != nil {
		return fmt.Errorf("salvar transcricao: %w", err)
	}

	// Step 2: Gera resumo via LLM
	summaryPrompt := model.SummaryPrompt(transcript)
	summaryJSON, err := p.llmClient.Generate(ctx, summaryPrompt)
	if err != nil {
		return fmt.Errorf("gerar resumo: %w", err)
	}

	var summary model.Summary
	if err := json.Unmarshal([]byte(extractJSON(summaryJSON)), &summary); err != nil {
		// Se falhar o parse, salva raw mesmo assim
		summary = model.Summary{
			Title:       session.Title,
			Description: "Resumo em processamento",
		}
	}

	summaryData, _ := json.Marshal(summary)
	if err := p.sessionRepo.UpdateSummary(ctx, session.ID, summaryData); err != nil {
		return fmt.Errorf("salvar resumo: %w", err)
	}

	// Step 3: Gera flashcards via LLM
	flashcardPrompt := model.FlashcardPrompt(transcript)
	flashcardJSON, err := p.llmClient.Generate(ctx, flashcardPrompt)
	if err != nil {
		return fmt.Errorf("gerar flashcards: %w", err)
	}

	var deck model.FlashcardDeck
	if err := json.Unmarshal([]byte(extractJSON(flashcardJSON)), &deck); err != nil {
		deck = model.FlashcardDeck{SessionID: session.ID}
	}

	// Salva flashcards no banco
	cards := make([]repository.Flashcard, len(deck.Cards))
	for i, card := range deck.Cards {
		cid, _ := generateID()
		cards[i] = repository.Flashcard{
			ID:         cid,
			SessionID:  session.ID,
			Front:      card.Front,
			Back:       card.Back,
			Difficulty: card.Difficulty,
		}
	}

	if err := p.flashcardRepo.BatchInsert(ctx, cards); err != nil {
		return fmt.Errorf("salvar flashcards: %w", err)
	}

	// Step 4: Limpa arquivo original de audio (privacidade)
	_ = p.storage.Delete(audioPath)

	// Step 5: Marca como completada
	return p.sessionRepo.UpdateStatus(ctx, session.ID, "completed")
}

// GetSession retorna uma sessao com resumo e flashcards
func (p *Processor) GetSession(ctx context.Context, id string) (*repository.Session, error) {
	return p.sessionRepo.GetByID(ctx, id)
}

// ListSessions retorna sessoes de um usuario
func (p *Processor) ListSessions(ctx context.Context, userID string) ([]repository.Session, error) {
	return p.sessionRepo.ListByUser(ctx, userID, 50, 0)
}

// DeleteSession remove uma sessao
func (p *Processor) DeleteSession(ctx context.Context, id string) error {
	// Busca path do audio para deletar
	session, err := p.sessionRepo.GetByID(ctx, id)
	if err == nil && session.AudioPath.Valid {
		_ = p.storage.Delete(session.AudioPath.String)
	}
	return p.sessionRepo.Delete(ctx, id)
}

// GetFlashcards retorna flashcards de uma sessao
func (p *Processor) GetFlashcards(ctx context.Context, sessionID string) ([]repository.Flashcard, error) {
	return p.flashcardRepo.GetBySession(ctx, sessionID)
}

// generateID gera um ID aleatorio de 16 bytes (32 chars hex)
func generateID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// extractJSON extrai o primeiro bloco JSON de uma resposta LLM
func extractJSON(s string) string {
	start := -1
	depth := 0
	for i, c := range s {
		if c == '{' {
			if depth == 0 {
				start = i
			}
			depth++
		} else if c == '}' {
			depth--
			if depth == 0 && start != -1 {
				return s[start : i+1]
			}
		}
	}
	return s
}
