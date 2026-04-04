package postgres

import (
	"context"
	"database/sql"

	"github.com/aulaflash/backend/internal/domain/repository"
)

type FlashcardRepository struct {
	db *sql.DB
}

func NewFlashcardRepository(db *sql.DB) *FlashcardRepository {
	return &FlashcardRepository{db: db}
}

func (r *FlashcardRepository) BatchInsert(ctx context.Context, cards []repository.Flashcard) error {
	if len(cards) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx,
		`INSERT INTO flashcards (id, session_id, front, back, difficulty, known)
		 VALUES ($1, $2, $3, $4, $5, $6)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, c := range cards {
		if _, err := stmt.ExecContext(ctx, c.ID, c.SessionID, c.Front, c.Back, c.Difficulty, c.Known); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *FlashcardRepository) GetBySession(ctx context.Context, sessionID string) ([]repository.Flashcard, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, session_id, front, back, difficulty, known FROM flashcards
		 WHERE session_id = $1 ORDER BY difficulty ASC`, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []repository.Flashcard
	for rows.Next() {
		var c repository.Flashcard
		if err := rows.Scan(&c.ID, &c.SessionID, &c.Front, &c.Back, &c.Difficulty, &c.Known); err != nil {
			return nil, err
		}
		cards = append(cards, c)
	}
	return cards, nil
}

func (r *FlashcardRepository) MarkKnown(ctx context.Context, id string, known bool) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE flashcards SET known = $1 WHERE id = $2`, known, id)
	return err
}
