// Package handlers contains HTTP request handlers for the API endpoints
package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"employee-management/internal/models"
	"employee-management/internal/repository"
	"employee-management/internal/service"

	"github.com/gin-gonic/gin"
)

// HAVE IN MIND:
// 1. If validation grows, move it into a helper function.

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
// Success: Returns 200 OK with the created employee
// Error: Returns 400 Bad Request for invalid input or 500 Internal Server Error
func (h *EmployeeHandler) CreateEmployee(c *gin.Context) {
	var req models.Employee

	// Check JSON shape / types
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Input validation
	if req.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is required"})
		return
	}

	if req.EmployeeNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "employee number is required"})
		return
	}

	// Business logic
	if err := h.service.Create(c.Request.Context(), &req); err != nil {
		switch {
		case errors.Is(err, repository.ErrEmailAlreadyExists),
			errors.Is(err, repository.ErrEmployeeNumberAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
	}

	c.JSON(http.StatusOK, req)
}

// GetEmployeeByID handles GET requests to retrieve an employee by ID
//
// Success: Returns 200 OK with the employee data
// Error: Returns 400 Bad Request for invalid ID, 404 Not Found, or 500 Internal Server Error
func (h *EmployeeHandler) GetEmployeeByID(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee id"})
		return
	}

	emp, err := h.service.FindByID(c.Request.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrEmployeeNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "employee with id " + idParam + " does not exist"})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
	}

	c.JSON(http.StatusOK, emp)
}

// HealthCheck handles GET /health
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "UP",
		"service": "employee-management",
	})
}
