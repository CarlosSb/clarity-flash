package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/aulaflash/backend/internal/api"
	"github.com/aulaflash/backend/internal/auth"
	"github.com/aulaflash/backend/internal/config"
	"github.com/aulaflash/backend/internal/handler"
	postgres "github.com/aulaflash/backend/internal/repository/postgres"
	"github.com/aulaflash/backend/internal/service"
	"github.com/aulaflash/backend/pkg/audio"
	"github.com/aulaflash/backend/pkg/llm"
	"github.com/aulaflash/backend/pkg/storage"
	"github.com/aulaflash/backend/pkg/stt"

	_ "github.com/lib/pq"
)

// main - Ponto de entrada da aplicação ClarityFlash
// Esta função inicializa TODOS os componentes do sistema em ordem:
// 1. Configurações → 2. Banco → 3. Serviços → 4. Handlers → 5. Servidor
func main() {
	// 📋 PASSO 1: Carregar configurações (variáveis de ambiente, arquivos)
	cfg := config.Load()
	log.Printf("Iniciando ClarityFlash server na porta %d", cfg.ServerPort)

	// 🗄️ PASSO 2: Conectar ao banco de dados PostgreSQL
	dsn, err := config.GetDSN(cfg) // Data Source Name (string de conexão)
	if err != nil {
		log.Fatalf("erro de configuracao do banco: %v", err)
	}
	db, err := sql.Open("postgres", dsn) // Abrir conexão com PostgreSQL
	if err != nil {
		log.Fatalf("erro ao conectar ao banco: %v", err)
	}
	defer db.Close() // Fechar conexão quando programa terminar

	// 🏥 Testar conexão com banco (ping)
	if err := db.Ping(); err != nil {
		log.Fatalf("erro ao pingar banco: %v", err)
	}
	log.Println("Banco de dados conectado")

	// 💾 PASSO 3: Inicializar sistema de armazenamento local
	store, err := storage.NewLocalStorage(cfg.UploadDir)
	if err != nil {
		log.Fatalf("erro ao inicializar storage: %v", err)
	}

	// 🎵 PASSO 4: Inicializar processador de áudio (conversão de formatos)
	audioProc, err := audio.NewProcessor("/tmp") // Pasta temporária para arquivos
	if err != nil {
		log.Fatalf("erro ao inicializar audio processor: %v", err)
	}

	// 🤖 PASSO 5: Configurar clientes de IA
	// STT (Speech-to-Text): Converte fala em texto
	sttClient := stt.NewGroqClient(cfg.GroqAPIKey, cfg.GroqModel)

	// LLM (Large Language Model): Gera resumos e flashcards
	var llmClient llm.LLMClient
	if cfg.UseOllama {
		// Modo alternativo: usar Ollama local (sem internet)
		llmClient = llm.NewOllamaClient(cfg.OllamaURL, cfg.LLMModel)
	} else {
		// Modo padrão: usar Groq (mais rápido, na nuvem)
		llmClient = llm.NewGroqLLMClient(cfg.GroqAPIKey, cfg.LLMModel)
	}

	// 📊 PASSO 6: Inicializar repositórios (acesso ao banco)
	sessionRepo := postgres.NewSessionRepository(db)     // CRUD de sessões
	flashcardRepo := postgres.NewFlashcardRepository(db) // CRUD de flashcards
	userRepo := postgres.NewUserRepository(db)           // CRUD de usuários

	// 🎼 PASSO 7: Criar o "maestro" - Processor que orquestra tudo
	proc := service.NewProcessor(sessionRepo, flashcardRepo, store, audioProc, sttClient, llmClient)

	// 🔐 PASSO 8: Configurar autenticação JWT
	jwtSecret := cfg.HuggingFaceToken // Usar token HF como fallback
	if jwtSecret == "" {
		jwtSecret = "dev_secret_key_change_in_production" // ⚠️ MUDAR EM PRODUÇÃO!
	}
	tokenService := auth.NewTokenService(jwtSecret)
	authSvc := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authSvc, tokenService)

	// 🎯 PASSO 9: Criar handlers HTTP (controladores das rotas)
	sessionHandler := handler.NewSessionHandler(proc) // /api/sessions/*
	exportHandler := handler.NewExportHandler(proc)   // /api/export/*

	// 🚀 PASSO 10: Montar aplicação Fiber com todas as rotas
	app := api.SetupFiberApp(sessionHandler, authHandler, tokenService, exportHandler)

	// 🌐 PASSO 11: Iniciar servidor HTTP
	addr := fmt.Sprintf(":%d", cfg.ServerPort) // Ex: ":8081"
	log.Printf("Servidor rodando em http://localhost%s", addr)
	log.Fatal(app.Listen(addr)) // Bloquear e servir para sempre
}
