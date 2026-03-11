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

// recoverPasswordRequest represents the structure of the password recovery request payload
type recoverPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// RecoverPassword handles the HTTP POST request for initiating password recovery
func (h *AuthHandler) RecoverPassword(c *gin.Context) {
	var req recoverPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err := h.service.RecoverPassword(c.Request.Context(), req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "recovery email sent if user exists",
	})
}

// resetPasswordRequest represents the structure of the password reset request payload
type resetPasswordRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// ResetPassword handles the HTTP POST request for resetting the password using a token
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req resetPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err := h.service.ResetPasswordWithToken(c.Request.Context(), req.Email, req.Token, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "password updated successfully",
	})
}
