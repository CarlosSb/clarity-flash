package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

// LocalStorage armazena arquivos localmente
type LocalStorage struct {
	UploadDir string
}

func NewLocalStorage(uploadDir string) (*LocalStorage, error) {
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, fmt.Errorf("criar diretorio upload: %w", err)
	}
	return &LocalStorage{UploadDir: uploadDir}, nil
}

// Save salva um arquivo enviado via multipart/form
func (s *LocalStorage) Save(file multipart.File, header *multipart.FileHeader, filename string) (string, error) {
	dst := filepath.Join(s.UploadDir, filename)

	out, err := os.Create(dst)
	if err != nil {
		return "", fmt.Errorf("criar arquivo: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		return "", fmt.Errorf("salvar arquivo: %w", err)
	}

	return dst, nil
}

// GetPath retorna o path absoluto de um arquivo por nome
func (s *LocalStorage) GetPath(filename string) string {
	return filepath.Join(s.UploadDir, filename)
}

// Delete remove um arquivo pelo path
func (s *LocalStorage) Delete(path string) error {
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("deletar arquivo: %w", err)
	}
	return nil
}
