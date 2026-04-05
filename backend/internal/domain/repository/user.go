package repository

import (
	"context"
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID            string       `json:"id"`
	Name          string       `json:"name"`
	Email         string       `json:"email"`
	PasswordHash  string       `json:"-"`
	Mode          string       `json:"mode"` // student, professional
	CreatedAt     sql.NullTime `json:"created_at"`
}

// SetPassword hashes and stores the password
func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	return nil
}

// CheckPassword verifies password against stored hash
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

type UserRepository interface {
	Create(ctx context.Context, u *User) error
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, u *User) error
}
