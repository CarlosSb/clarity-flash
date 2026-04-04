package repository

import (
	"context"
	"database/sql"
)

type User struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	Email     string       `json:"email"`
	Mode      string       `json:"mode"` // student, professional
	CreatedAt sql.NullTime `json:"created_at"`
}

type UserRepository interface {
	Create(ctx context.Context, u *User) error
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, u *User) error
}
