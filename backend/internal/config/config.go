package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort       int
	DatabaseURL      string
	GroqAPIKey       string
	HuggingFaceToken string
	OllamaURL        string
	GroqModel        string
	LLMModel         string
	UploadDir        string
	UseOllama        bool
}

func Load() *Config {
	// Tenta carregar .env da raiz do projeto (2 níveis acima do backend/)
	rootDir, err := os.Getwd()
	if err == nil {
		for i := 0; i < 3; i++ {
			envPath := filepath.Join(rootDir, ".env")
			if _, err := os.Stat(envPath); err == nil {
				_ = godotenv.Load(envPath)
				break
			}
			rootDir = filepath.Dir(rootDir)
		}
	}

	port := 8081
	if p := os.Getenv("SERVER_PORT"); p != "" {
		if val, err := strconv.Atoi(p); err == nil {
			port = val
		}
	}

	useOllama := os.Getenv("USE_OLLAMA") == "true"

	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "/tmp/aulaflash-uploads"
	}

	return &Config{
		ServerPort:       port,
		DatabaseURL:      getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/aulaflash?sslmode=disable"),
		GroqAPIKey:       os.Getenv("GROQ_API_KEY"),
		HuggingFaceToken: os.Getenv("HUGGING_FACE_TOKEN"),
		OllamaURL:        getEnv("OLLAMA_URL", "http://localhost:11434"),
		GroqModel:        getEnv("GROQ_MODEL", "whisper-large-v3"),
		LLMModel:         getEnv("LLM_MODEL", "neuralmagic/Qwen2.5-72B-Instruct-quantized.w4a16"),
		UploadDir:        uploadDir,
		UseOllama:        useOllama,
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func GetDSN(cfg *Config) (string, error) {
	if cfg.DatabaseURL == "" {
		return "", fmt.Errorf("DATABASE_URL nao configurado")
	}
	return cfg.DatabaseURL, nil
}
