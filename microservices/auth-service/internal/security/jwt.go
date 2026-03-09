// Package security provides utilities for managing JWT tokens
package security

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTManager holds the secret key used for signing tokens
type JWTManager struct {
	secret string
}

// NewJWTManager creates a new instance of JWTManager
func NewJWTManager(secret string) *JWTManager {
	return &JWTManager{secret: secret}
}

// GenerateToken generates a new JWT tokuen with the given email and role
func (j *JWTManager) GenerateToken(email, role string) (string, error) {
	claims := jwt.MapClaims{
		"sub":  email,
		"role": role,
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(j.secret))
}
