// Package handlers contains HTTP request handlers for the API endpoints
package handlers

import (
	"net/http"

	"employee-management/internal/models"
	"employee-management/internal/service"
	"github.com/gin-gonic/gin"
)

// EmployeeHandler handles HTTP requests for employee operations
type EmployeeHandler struct {
	service *service.EmployeeService // Bussiness logic dependency
}

// NewEmployeeHandler creates a new EmployeeHandler instance
func NewEmployeeHandler(s *service.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{service: s}
}

// CreateEmployee handles POST requests to create a new employee
// Validates the request body and delegates to the service layer
//
// Success: Returns 201 Created with the created employee
// Error: Returns 400 Bad Request for invalid input or 500 Internal Server Error
func (h *EmployeeHandler) CreateEmployee(c *gin.Context) {
	var req models.Employee

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Create(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, req)
}

// HealthCheck handles GET /health
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "UP",
		"service": "employee-management",
	})
}
