// Package service contains the business logic of the application
package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"auth-service/internal/messaging"
	"auth-service/internal/repository"
	"auth-service/internal/security"
)

// AuthService provides auth-related services
type AuthService struct {
	repo     *repository.UserRepository
	jwt      *security.JWTManager
	producer *messaging.Producer
}

// NewAuthService creates a new instance of AuthService with the ginev user, JWT and producer
func NewAuthService(repo *repository.UserRepository, jwt *security.JWTManager, producer *messaging.Producer) *AuthService {
	return &AuthService{repo: repo, jwt: jwt, producer: producer}
}

// CreateUser creates a new user with its email, role, status and password hash
func (s *AuthService) CreateUser(ctx context.Context, email, role, status string) error {
	err := s.repo.CreateUser(ctx, email, role, status)
	if err != nil {
		return err
	}

	token := uuid.New().String()
	return s.producer.PublishRecoverPasswordEvent(ctx, email, token)
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

// RecoverPassword generates a token and publishes an event for password recovery
func (s *AuthService) RecoverPassword(ctx context.Context, email string) error {
	// Check if user exists
	_, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return errors.New("user not found")
	}

	token := uuid.New().String()

	return s.producer.PublishRecoverPasswordEvent(ctx, email, token)
}

// ResetPassword updates the password of a user by their email
func (s *AuthService) ResetPassword(ctx context.Context, email, password string) error {
	hash, err := security.HashPassword(password)
	if err != nil {
		return err
	}

	return s.repo.UpdatePassword(ctx, email, hash)
}
