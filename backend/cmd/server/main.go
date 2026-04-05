package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aulaflash/backend/internal/api"
	"github.com/aulaflash/backend/internal/auth"
	"github.com/aulaflash/backend/internal/config"
	"github.com/aulaflash/backend/internal/handler"
	"github.com/aulaflash/backend/internal/middleware"
	postgres "github.com/aulaflash/backend/internal/repository/postgres"
	"github.com/aulaflash/backend/internal/service"
	"github.com/aulaflash/backend/internal/websocket"
	"github.com/aulaflash/backend/pkg/audio"
	"github.com/aulaflash/backend/pkg/llm"
	"github.com/aulaflash/backend/pkg/storage"
	"github.com/aulaflash/backend/pkg/stt"

	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Load()
	log.Printf("Iniciando ClarityFlash server na porta %d", cfg.ServerPort)

	// Banco de dados
	dsn, err := config.GetDSN(cfg)
	if err != nil {
		log.Fatalf("erro de configuracao do banco: %v", err)
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("erro ao conectar ao banco: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("erro ao pingar banco: %v", err)
	}
	log.Println("Banco de dados conectado")

	// Storage local
	store, err := storage.NewLocalStorage(cfg.UploadDir)
	if err != nil {
		log.Fatalf("erro ao inicializar storage: %v", err)
	}

	// Audio processor
	audioProc, err := audio.NewProcessor(os.TempDir())
	if err != nil {
		log.Fatalf("erro ao inicializar audio processor: %v", err)
	}

	// STT e LLM clients — ambos usando Groq (unificado)
	sttClient := stt.NewGroqClient(cfg.GroqAPIKey, cfg.GroqSTTModel)

	var llmClient llm.LLMClient
	if cfg.UseOllama {
		// Fallback: Ollama local (ainda disponivel se quiser alternar)
		ollamaURL := cfg.OllamaURL
		ollamaModel := cfg.OllamaModel
		llmClient = llm.NewOllamaClient(ollamaURL, ollamaModel)
	} else {
		llmClient = llm.NewGroqLLMClient(cfg.GroqAPIKey, cfg.GroqLLMModel)
	}

	// WebSocket hub
	hub := websocket.NewHub()

	// Repositories
	sessionRepo := postgres.NewSessionRepository(db)
	flashcardRepo := postgres.NewFlashcardRepository(db)
	userRepo := postgres.NewUserRepository(db)

	// Processor (orquestrador)
	proc := service.NewProcessor(sessionRepo, flashcardRepo, store, audioProc, sttClient, llmClient, hub, cfg.UploadDir)

	// Stream handler (recebe chunks do frontend em tempo real)
	streamHandler := handler.NewStreamHandler(cfg.UploadDir, proc)
	// Enable real-time transcription: Whisper -> Hub -> frontend
	streamHandler.SetTranscriber(sttClient)
	streamHandler.SetPusher(hub)

	// JWT & Auth
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}
	tokenService := auth.NewTokenService(jwtSecret)
	authSvc := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authSvc, tokenService)

	// Handlers
	sessionHandler := handler.NewSessionHandler(proc, userRepo)
	exportHandler := handler.NewExportHandler(proc)

	// Router com middleware
	mux := middleware.CORS(api.SetupRouter(sessionHandler, authHandler, tokenService, exportHandler, streamHandler, hub))

	addr := fmt.Sprintf(":%d", cfg.ServerPort)
	log.Printf("Servidor rodando em http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
