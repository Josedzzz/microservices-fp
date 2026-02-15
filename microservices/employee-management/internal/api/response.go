// Package api handle the response of the handlers
package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ErrorDetail represents a specific validation error
type ErrorDetail struct {
	Field         string `json:"field"`
	Message       string `json:"message"`
	RejectedValue string `json:"rejectedValue,omitempty"`
}

// ErrorResponse is the standart struct for error response
//
//	@Description	Standard error response structure
type ErrorResponse struct {
	Status    int           `json:"status"`
	Error     string        `json:"error"`
	Message   string        `json:"message"`
	Timestamp time.Time     `json:"timestamp"`
	Path      string        `json:"path"`
	Errors    []ErrorDetail `json:"errors,omitempty"`
}

// ValidationError creates a validation error response
func ValidationError(c *gin.Context, status int, message string, errors []ErrorDetail) {
	response := ErrorResponse{
		Status:    status,
		Error:     http.StatusText(status),
		Message:   message,
		Timestamp: time.Now().UTC(),
		Path:      c.Request.URL.Path,
		Errors:    errors,
	}
	c.JSON(status, response)
}

// Error creates a simple error response
func Error(c *gin.Context, status int, message string) {
	response := ErrorResponse{
		Status:    status,
		Error:     http.StatusText(status),
		Message:   message,
		Timestamp: time.Now().UTC(),
		Path:      c.Request.URL.Path,
	}
	c.JSON(status, response)
}

// InternalServerError for 500 errors
func InternalServerError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, message)
}

// BadRequest for 400 errors
func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, message)
}

// NotFound for 404 errors
func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, message)
}

// Conflict for 409 errors
func Conflict(c *gin.Context, message string) {
	Error(c, http.StatusConflict, message)
}
