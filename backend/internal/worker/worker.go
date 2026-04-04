package worker

import (
	"database/sql"
	"log"

	"github.com/aulaflash/backend/internal/config"
)

// Worker processa trabalhos em background
type Worker struct {
	db  *sql.DB
	cfg *config.Config
}

func NewWorker(db *sql.DB, cfg *config.Config) *Worker {
	return &Worker{db: db, cfg: cfg}
}

func (w *Worker) Start() error {
	log.Println("Worker iniciado")
	// No MVP o processamento e assincrono via goroutines no Handler.
	// Worker sera utilizado no futuro para retry, limpeza de arquivos, etc.
	return nil
}
