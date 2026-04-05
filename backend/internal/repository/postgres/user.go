package postgres

import (
	"context"
	"database/sql"

	"github.com/aulaflash/backend/internal/domain/repository"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, u *repository.User) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO users (id, name, email, password_hash, mode) VALUES ($1, $2, $3, $4, $5)`,
		u.ID, u.Name, u.Email, u.PasswordHash, u.Mode)
	return err
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*repository.User, error) {
	u := &repository.User{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, email, password_hash, mode, created_at FROM users WHERE id = $1`, id).Scan(
		&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Mode, &u.CreatedAt)
	return u, err
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*repository.User, error) {
	u := &repository.User{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, email, password_hash, mode, created_at FROM users WHERE email = $1`, email).Scan(
		&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Mode, &u.CreatedAt)
	return u, err
}

func (r *UserRepository) Update(ctx context.Context, u *repository.User) error {
	q := `UPDATE users SET name = $1, email = $2, mode = $3, updated_at = NOW() WHERE id = $4`
	args := []any{u.Name, u.Email, u.Mode, u.ID}

	// Optionally update password_hash too
	if u.PasswordHash != "" {
		q = `UPDATE users SET name = $1, email = $2, password_hash = $3, mode = $4, updated_at = NOW() WHERE id = $5`
		args = []any{u.Name, u.Email, u.PasswordHash, u.Mode, u.ID}
	}

	_, err := r.db.ExecContext(ctx, q, args...)
	return err
}