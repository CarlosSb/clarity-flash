package repository

import (
	"context"
	"database/sql"
)

// Session representa uma gravacao de aula/reuniao
type Session struct {
	ID           string         `json:"id"`
	UserID       string         `json:"user_id"`
	Title        string         `json:"title"`
	Description  string         `json:"description"`
	Duration     int            `json:"duration"` // segundos
	Status       string         `json:"status"`   // processing, completed, failed
	Mode         string         `json:"mode"`     // student, professional
	CreatedAt    sql.NullTime   `json:"created_at"`
	UpdatedAt    sql.NullTime   `json:"updated_at"`
	Transcript   sql.NullString `json:"-"`
	AudioPath    sql.NullString `json:"-"`
	SummaryData  sql.NullString `json:"-"`
	FlashcardData sql.NullString `json:"-"`
}

// SessionRepository define as operacoes de sessao no banco
type SessionRepository interface {
	Create(ctx context.Context, s *Session) error
	GetByID(ctx context.Context, id string) (*Session, error)
	ListByUser(ctx context.Context, userID string, limit, offset int) ([]Session, error)
	UpdateStatus(ctx context.Context, id, status string) error
	UpdateTranscript(ctx context.Context, id, transcript string) error
	UpdateSummary(ctx context.Context, id string, summaryData []byte) error
	UpdateFlashcards(ctx context.Context, id string, flashcardData []byte) error
	Delete(ctx context.Context, id string) error
}
