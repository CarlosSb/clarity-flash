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
		`INSERT INTO users (id, name, email, mode) VALUES ($1, $2, $3, $4)`,
		u.ID, u.Name, u.Email, u.Mode)
	return err
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*repository.User, error) {
	u := &repository.User{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, email, mode, created_at FROM users WHERE id = $1`, id).Scan(
		&u.ID, &u.Name, &u.Email, &u.Mode, &u.CreatedAt)
	return u, err
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*repository.User, error) {
	u := &repository.User{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, email, mode, created_at FROM users WHERE email = $1`, email).Scan(
		&u.ID, &u.Name, &u.Email, &u.Mode, &u.CreatedAt)
	return u, err
}

func (r *UserRepository) Update(ctx context.Context, u *repository.User) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET name = $1, email = $2, mode = $3 WHERE id = $4`,
		u.Name, u.Email, u.Mode, u.ID)
	return err
}
