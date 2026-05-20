// Package handler contains HTTP handlers for processing API requests
package handler

import (
	"net/http"
	"time"

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
//
//	@Summary		User login
//	@Description	Receives credentials and returns a JWT
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		loginRequest	true	"Login credentials"
//	@Success		200		{object}	map[string]string	"Successful login"
//	@Failure		401		{object}	map[string]string	"Invalid credentials"
//	@Router			/login [post]
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
//
//	@Summary		Password recovery
//	@Description	Initiates the process of account recovery by receiving an email
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		recoverPasswordRequest	true	"Recovery email"
//	@Success		200		{object}	map[string]string		"Recovery process initiated"
//	@Failure		400		{object}	map[string]string		"Invalid request"
//	@Router			/recover-password [post]
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
//
//	@Summary		Reset password
//	@Description	Receives the recovery token and new password to update it in the database
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		resetPasswordRequest	true	"New password data"
//	@Success		200		{object}	map[string]string		"Password updated successfully"
//	@Failure		400		{object}	map[string]string		"Invalid token or request"
//	@Router			/reset-password [post]
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

// HealthCheck handles GET /health
//
//	@Summary		Health check
//	@Description	Returns the service health status
//	@Tags			System
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}	"Service is healthy"
//	@Router			/health [get]
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "UP",
		"service":   "auth-service",
		"timestamp": time.Now().UTC(),
	})
}
