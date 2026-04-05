package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/aulaflash/backend/internal/domain/repository"
)

type SessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Create(ctx context.Context, s *repository.Session) error {
	query := `INSERT INTO sessions (id, user_id, title, status, mode, created_at, updated_at)
				  VALUES ($1, $2, $3, $4, $5, $6, $7)`
	now := time.Now()
	_, err := r.db.ExecContext(ctx, query,
		s.ID, s.UserID, s.Title, s.Status, s.Mode, now, now)
	return err
}

func (r *SessionRepository) GetByID(ctx context.Context, id string) (*repository.Session, error) {
	s := &repository.Session{}
	var desc, audioPath, transcript, summaryData sql.NullString
	query := `SELECT id, user_id, title, description, duration, status, mode,
				  created_at, updated_at, transcript, audio_path, summary_data
				  FROM sessions WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&s.ID, &s.UserID, &s.Title, &desc, &s.Duration,
		&s.Status, &s.Mode, &s.CreatedAt, &s.UpdatedAt,
		&transcript, &audioPath, &summaryData)
	s.Description = nullString(desc)
	s.Transcript = transcript
	s.AudioPath = audioPath
	s.SummaryData = summaryData
	return s, err
}

func (r *SessionRepository) ListByUser(ctx context.Context, userID string, limit, offset int) ([]repository.Session, error) {
	query := `SELECT s.id, s.user_id, s.title, s.description, s.duration, s.status, s.mode,
				  s.created_at, s.updated_at, COUNT(f.id)::int as flashcard_count
				  FROM sessions s
				  LEFT JOIN flashcards f ON f.session_id = s.id
				  WHERE s.user_id = $1
				  GROUP BY s.id
				  ORDER BY s.created_at DESC LIMIT $2 OFFSET $3`
	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []repository.Session
	for rows.Next() {
		var s repository.Session
		var desc sql.NullString
		if err := rows.Scan(
			&s.ID, &s.UserID, &s.Title, &desc,
			&s.Duration, &s.Status, &s.Mode, &s.CreatedAt, &s.UpdatedAt, &s.FlashcardCount,
		); err != nil {
			return nil, err
		}
		s.Description = nullString(desc)
		sessions = append(sessions, s)
	}
	return sessions, nil
}

func (r *SessionRepository) UpdateStatus(ctx context.Context, id, status string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE sessions SET status = $1, updated_at = NOW() WHERE id = $2`, status, id)
	return err
}

func (r *SessionRepository) UpdateTranscript(ctx context.Context, id, transcript string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE sessions SET transcript = $1, updated_at = NOW() WHERE id = $2`, transcript, id)
	return err
}

func (r *SessionRepository) UpdateSummary(ctx context.Context, id string, summaryData []byte) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE sessions SET summary_data = $1::jsonb, updated_at = NOW() WHERE id = $2`, summaryData, id)
	return err
}

func (r *SessionRepository) UpdateFlashcards(ctx context.Context, id string, data []byte) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE sessions SET flashcard_data = $1::jsonb, updated_at = NOW() WHERE id = $2`, data, id)
	return err
}

func (r *SessionRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM sessions WHERE id = $1`, id)
	return err
}

func nullString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}
