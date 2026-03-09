// Package handler contains HTTP handlers for processing API requests
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"auth-service/internal/service"
)

// AuthHandler interacts with the AuthService to perform logic operations
type AuthHandler struct {
	service *service.AuthService
}

// NewAuthHandler creates a new instance of AuthHandler
func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// loginRequest represents the structure of the login request payload
type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login handles the HTTP POST request for user login
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	token, err := h.service.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": token,
	})
}
