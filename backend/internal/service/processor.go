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

// Processor - O "Maestro" do sistema de IA
// Esta struct é como um maestro de orquestra que coordena todo o processamento:
// 1. Recebe áudio → 2. Salva arquivo → 3. Transcreve fala → 4. Gera resumo → 5. Cria flashcards
// Tudo roda em background para não bloquear o usuário
type Processor struct {
	sessionRepo   repository.SessionRepository // Acesso ao banco (sessões)
	flashcardRepo repository.FlashcardRepository // Acesso ao banco (flashcards)
	storage       *storage.LocalStorage         // Sistema de arquivos local
	audioProc     *audio.Processor              // Processamento de áudio (conversão)
	sttClient     *stt.GroqClient               // IA de fala para texto (Whisper)
	llmClient     llm.LLMClient                 // IA de texto (Llama) para resumos/flashcards
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

// Process - Ponto de entrada para processamento de áudio
// Esta função é chamada imediatamente após o upload. Ela:
// 1. Prepara a sessão, 2. Salva o arquivo, 3. Inicia processamento assíncrono
func (p *Processor) Process(ctx context.Context, session *repository.Session, file multipart.File, header *multipart.FileHeader) error {
	// 🆔 PASSO 1: Gerar ID único para a sessão
	// Cada sessão precisa de um ID único (como CPF, mas aleatório)
	id, err := generateID()
	if err != nil {
		return fmt.Errorf("gerar id: %w", err)
	}
	session.ID = id // Salvar ID na sessão

	// 💾 PASSO 2: Salvar arquivo de áudio no disco
	// Combina ID + nome original para evitar conflitos
	audioPath := id + "_" + header.Filename // Ex: "abc123_audio.mp3"
	path, err := p.storage.Save(file, header, audioPath)
	if err != nil {
		return fmt.Errorf("salvar audio: %w", err)
	}

	// 🗃️ PASSO 3: Criar registro da sessão no banco de dados
	// Agora o frontend já pode ver que a sessão existe (status: processing)
	if err := p.sessionRepo.Create(ctx, session); err != nil {
		return fmt.Errorf("criar sessao: %w", err)
	}

	// 🚀 PASSO 4: Iniciar processamento ASSÍNCRONO (não bloqueia usuário!)
	// Em produção, isso seria feito por workers dedicados (como filas)
	go func() {
		// Rodar pipeline completo em background
		if err := p.runPipeline(context.Background(), session, path); err != nil {
			// Se der erro, marcar sessão como falhada
			_ = p.sessionRepo.UpdateStatus(context.Background(), session.ID, "failed")
		}
	}()

	return nil // Retorna imediatamente - processamento continua em background
}

// runPipeline - O "coração" do processamento de IA
// Esta função executa todo o pipeline mágico em 5 etapas:
// 1. Validar áudio → 2. Converter formato → 3. Transcrever → 4. Gerar resumo → 5. Criar flashcards
func (p *Processor) runPipeline(ctx context.Context, session *repository.Session, audioPath string) error {
	// 🔍 ETAPA 1: Validar arquivo de áudio
	// Verificar se o arquivo existe e está em formato válido
	if err := p.audioProc.ValidateAudio(audioPath); err != nil {
		return err // Arquivo inválido ou corrompido
	}

	// 🔄 ETAPA 2: Converter para formato WAV
	// A IA de transcrição funciona melhor com WAV (formato padrão)
	wavPath, err := p.audioProc.ConvertToWAV(audioPath)
	if err != nil {
		return fmt.Errorf("converter audio: %w", err)
	}
	defer p.audioProc.Cleanup(wavPath) // Limpar arquivo temporário no final

	// 🎤 ETAPA 3: Transcrição - Converter fala em texto
	// Usa Groq Whisper (modelo de IA especializado em fala)
	transcript, err := p.sttClient.Transcribe(wavPath)
	if err != nil {
		return fmt.Errorf("transcricao: %w", err)
	}

	// 💾 Salvar transcrição no banco (usuário já pode ver!)
	if err := p.sessionRepo.UpdateTranscript(ctx, session.ID, transcript); err != nil {
		return fmt.Errorf("salvar transcricao: %w", err)
	}

	// 📝 ETAPA 4: Gerar resumo inteligente
	// Cria um prompt personalizado com a transcrição
	summaryPrompt := model.SummaryPrompt(transcript)
	// Envia para IA de linguagem (Llama) e recebe resposta em JSON
	summaryJSON, err := p.llmClient.Generate(ctx, summaryPrompt)
	if err != nil {
		return fmt.Errorf("gerar resumo: %w", err)
	}

	// 🔧 Processar resposta da IA
	var summary model.Summary
	if err := json.Unmarshal([]byte(extractJSON(summaryJSON)), &summary); err != nil {
		// Se a IA não retornou JSON válido, criar resumo básico
		summary = model.Summary{
			Title:       session.Title,
			Description: "Resumo em processamento",
		}
	}

	// 💾 Salvar resumo no banco
	summaryData, _ := json.Marshal(summary)
	if err := p.sessionRepo.UpdateSummary(ctx, session.ID, summaryData); err != nil {
		return fmt.Errorf("salvar resumo: %w", err)
	}

	// 🃏 ETAPA 5: Gerar flashcards de estudo
	// Mesmo processo: criar prompt → enviar para IA → processar resposta
	flashcardPrompt := model.FlashcardPrompt(transcript)
	flashcardJSON, err := p.llmClient.Generate(ctx, flashcardPrompt)
	if err != nil {
		return fmt.Errorf("gerar flashcards: %w", err)
	}

	// 🔧 Processar resposta da IA
	var deck model.FlashcardDeck
	if err := json.Unmarshal([]byte(extractJSON(flashcardJSON)), &deck); err != nil {
		// Se falhar, criar deck vazio (melhor que erro)
		deck = model.FlashcardDeck{SessionID: session.ID}
	}

	// 💾 Salvar flashcards no banco (bulk insert para performance)
	cards := make([]repository.Flashcard, len(deck.Cards))
	for i, card := range deck.Cards {
		cid, _ := generateID() // ID único para cada flashcard
		cards[i] = repository.Flashcard{
			ID:         cid,
			SessionID:  session.ID,      // Vincular à sessão
			Front:      card.Front,      // Pergunta
			Back:       card.Back,       // Resposta
			Difficulty: card.Difficulty, // Dificuldade (easy/medium/hard)
		}
	}

	// Inserir todas as flashcards de uma vez (mais eficiente)
	if err := p.flashcardRepo.BatchInsert(ctx, cards); err != nil {
		return fmt.Errorf("salvar flashcards: %w", err)
	}

	// 🧹 ETAPA 6: Limpeza de privacidade
	// Deletar arquivo de áudio original (LGPD/GDPR compliance)
	_ = p.storage.Delete(audioPath) // Ignorar erro se arquivo não existir

	// ✅ ETAPA 7: Finalizar sessão
	// Marcar como completa - usuário pode ver resultado final
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

// extractJSON - Extrator inteligente de JSON de respostas de IA
// As IAs de linguagem nem sempre retornam JSON "limpo". Elas podem adicionar
// texto extra antes/depois. Esta função encontra o primeiro bloco JSON válido.
//
// Exemplo: "Aqui está seu JSON: {"nome": "João"} Espero que ajude!"
// Retorna: {"nome": "João"}
func extractJSON(s string) string {
	start := -1  // Posição onde começa o JSON
	depth := 0   // Profundidade das chaves (para JSONs aninhados)

	// Percorrer cada caractere da string
	for i, c := range s {
		if c == '{' {
			if depth == 0 {
				start = i // Marcar início do JSON quando encontra primeira '{'
			}
			depth++ // Entrar em um nível mais profundo
		} else if c == '}' {
			depth-- // Sair de um nível
			if depth == 0 && start != -1 {
				// Achamos o fim do JSON! Retornar do start até aqui
				return s[start : i+1]
			}
		}
	}

	// Se não encontrou JSON válido, retornar string original
	return s
}
