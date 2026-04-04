package repository

import (
	"context"
	"database/sql"
)

// Flashcard representa um cartao individual no banco
type Flashcard struct {
	ID         string       `json:"id"`
	SessionID  string       `json:"session_id"`
	Front      string       `json:"front"`
	Back       string       `json:"back"`
	Difficulty int          `json:"difficulty"`
	Known      bool         `json:"known"`
	CreatedAt  sql.NullTime `json:"created_at"`
}

type FlashcardRepository interface {
	BatchInsert(ctx context.Context, cards []Flashcard) error
	GetBySession(ctx context.Context, sessionID string) ([]Flashcard, error)
	MarkKnown(ctx context.Context, id string, known bool) error
}
