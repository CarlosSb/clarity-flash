package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/aulaflash/backend/internal/config"

	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Load()
	log.Println("Rodando migracoes do banco de dados")

	// Determinar o diretorio base (onde o binario esta ou CWD)
	baseDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("erro ao obter diretorio atual: %v", err)
	}
	migrationsDir := filepath.Join(baseDir, "migrations")

	// Se nao existir, tenta relativo ao source
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		migrationsDir = filepath.Join(baseDir, "..", "migrations")
	}

	dsn, err := config.GetDSN(cfg)
	if err != nil {
		log.Fatalf("erro de configuracao do banco: %v", err)
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("erro ao conectar ao banco: %v", err)
	}
	defer db.Close()

	if err := runMigrations(db, migrationsDir); err != nil {
		log.Fatalf("erro ao executar migracoes: %v", err)
	}

	log.Println("Migracoes executadas com sucesso")
}

func runMigrations(db *sql.DB, dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("ler diretorio de migracoes (%s): %w", dir, err)
	}

	var files []os.DirEntry
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			files = append(files, e)
		}
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	if len(files) == 0 {
		return fmt.Errorf("nenhuma migracao encontrada em %s", dir)
	}

	for _, f := range files {
		path := filepath.Join(dir, f.Name())
		sqlBytes, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("ler %s: %w", path, err)
		}

		log.Printf("Aplicando %s...", f.Name())
		if _, err := db.Exec(string(sqlBytes)); err != nil {
			return fmt.Errorf("executar %s: %w", path, err)
		}
	}

	return nil
}
