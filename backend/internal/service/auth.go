package service

import (
	"context"
	"errors"

	"github.com/aulaflash/backend/internal/domain/repository"
)

// AuthServiceImpl implements the AuthService interface for auth handler
type AuthServiceImpl struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) *AuthServiceImpl {
	return &AuthServiceImpl{
		userRepo: userRepo,
	}
}

func (s *AuthServiceImpl) Register(ctx context.Context, name, email, password, mode string) (*repository.User, error) {
	// Check if user already exists
	existing, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil && existing != nil {
		return nil, errors.New("user already exists")
	}

	user := &repository.User{
		Name:  name,
		Email: email,
		Mode:  mode,
	}

	if err := user.SetPassword(password); err != nil {
		return nil, err
	}

	// Generate user ID
	id, err := generateID()
	if err != nil {
		return nil, err
	}
	user.ID = id

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthServiceImpl) Login(ctx context.Context, email, password string) (*repository.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil || user == nil {
		return nil, errors.New("invalid credentials")
	}

	if !user.CheckPassword(password) {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}
