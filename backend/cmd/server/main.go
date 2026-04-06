package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/aulaflash/backend/internal/api"
	"github.com/aulaflash/backend/internal/auth"
	"github.com/aulaflash/backend/internal/config"
	"github.com/aulaflash/backend/internal/handler"
	"github.com/aulaflash/backend/internal/middleware"
	postgres "github.com/aulaflash/backend/internal/repository/postgres"
	"github.com/aulaflash/backend/internal/service"
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
	audioProc, err := audio.NewProcessor("/tmp")
	if err != nil {
		log.Fatalf("erro ao inicializar audio processor: %v", err)
	}

	// STT e LLM clients — ambos usando Groq (unificado)
	sttClient := stt.NewGroqClient(cfg.GroqAPIKey, cfg.GroqModel)

	var llmClient llm.LLMClient
	if cfg.UseOllama {
		// Fallback: Ollama local (ainda disponivel se quiser alternar)
		llmClient = llm.NewOllamaClient(cfg.OllamaURL, cfg.LLMModel)
	} else {
		llmClient = llm.NewGroqLLMClient(cfg.GroqAPIKey, cfg.LLMModel)
	}

	// Repositories
	sessionRepo := postgres.NewSessionRepository(db)
	flashcardRepo := postgres.NewFlashcardRepository(db)
	userRepo := postgres.NewUserRepository(db)

	// Processor (orquestrador)
	proc := service.NewProcessor(sessionRepo, flashcardRepo, store, audioProc, sttClient, llmClient)

	// JWT & Auth
	jwtSecret := cfg.HuggingFaceToken // Usando como fallback se JWT_SECRET não estiver setado
	if jwtSecret == "" {
		jwtSecret = "dev_secret_key_change_in_production"
	}
	tokenService := auth.NewTokenService(jwtSecret)
	authSvc := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authSvc, tokenService)

	// Handlers
	sessionHandler := handler.NewSessionHandler(proc)
	exportHandler := handler.NewExportHandler(proc)

	// Router com middleware
	mux := middleware.CORS(api.SetupRouter(sessionHandler, authHandler, tokenService, exportHandler))

	addr := fmt.Sprintf(":%d", cfg.ServerPort)
	log.Printf("Servidor rodando em http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
