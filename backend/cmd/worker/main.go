package main

import (
	"database/sql"
	"log"

	"github.com/aulaflash/backend/internal/config"
	"github.com/aulaflash/backend/internal/worker"

	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Load()
	log.Println("Iniciando worker do AulaFlash")

	dsn, err := config.GetDSN(cfg)
	if err != nil {
		log.Fatalf("erro de configuracao do banco: %v", err)
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("erro ao conectar ao banco: %v", err)
	}
	defer db.Close()

	w := worker.NewWorker(db, cfg)
	if err := w.Start(); err != nil {
		log.Fatalf("erro ao iniciar worker: %v", err)
	}

	// Mantem worker rodando
	select {}
}
