// Package service contains the business logic of the application
package service

import (
	"context"
	"errors"

	"auth-service/internal/repository"
	"auth-service/internal/security"
)

// AuthService provides auth-related services
type AuthService struct {
	repo *repository.UserRepository
	jwt  *security.JWTManager
}

// NewAuthService creates a new instance of AuthService with the ginev user and JWT
func NewAuthService(repo *repository.UserRepository, jwt *security.JWTManager) *AuthService {
	return &AuthService{repo: repo, jwt: jwt}
}

// CreateUser creates a new user with its email, role, status and password hash
func (s *AuthService) CreateUser(ctx context.Context, email, role, status string) error {
	return s.repo.CreateUser(ctx, email, role, status)
}

// Login authenticates a user by their email and password
func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	if user.Status != "ACTIVE" {
		return "", errors.New("user disabled")
	}

	err = security.CheckPassword(password, user.PasswordHash)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	return s.jwt.GenerateToken(user.Email, user.Role)
}

// ResetPassword updates the password of a user by their email
func (s *AuthService) ResetPassword(ctx context.Context, email, password string) error {
	hash, err := security.HashPassword(password)
	if err != nil {
		return err
	}

	return s.repo.UpdatePassword(ctx, email, hash)
}
