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

// CreateEmployee godoc
// @Summary Create a new employee
// @Description Creates a new employee in the system
// @Tags Employees
// @Accept json
// @Produce json
// @Param employee body models.Employee true "Employee data"
// @Success 200 {object} models.Employee
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /employees [post]
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

// GetEmployeeByID godoc
// @Summary Get employee by ID
// @Description Retrieves an employee by its ID
// @Tags Employees
// @Produce json
// @Param id path int true "Employee ID"
// @Success 200 {object} models.Employee
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /employees/{id} [get]
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

// GetAllEmployees godoc
// @Summary Get all employees
// @Description Retrieves all employees
// @Tags Employees
// @Produce json
// @Success 200 {array} models.Employee
// @Failure 500 {object} map[string]string
// @Router /employees [get]
func (h *EmployeeHandler) GetAllEmployees(c *gin.Context) {
	employees, err := h.service.FindAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, employees)
}

// UpdateEmployee godoc
// @Summary Update employee
// @Description Updates an existing employee
// @Tags Employees
// @Accept json
// @Produce json
// @Param id path int true "Employee ID"
// @Param employee body models.Employee true "Updated employee data"
// @Success 204 "Employee updated successfully"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /employees/{id} [put]
func (h *EmployeeHandler) UpdateEmployee(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee id"})
		return
	}

	var req models.Employee
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.ID = id

	if err := h.service.Update(c.Request.Context(), &req); err != nil {
		switch {
		case errors.Is(err, repository.ErrEmployeeNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "employee with id " + idParam + " does not exist"})
			return
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

// DeleteEmployee godoc
// @Summary Delete employee
// @Description Deletes an employee by ID
// @Tags Employees
// @Param id path int true "Employee ID"
// @Success 204 "Employee deleted successfully"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /employees/{id} [delete]
func (h *EmployeeHandler) DeleteEmployee(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee id"})
		return
	}

	if err := h.service.Delete(
		c.Request.Context(),
		strconv.FormatInt(id, 10),
	); err != nil {
		switch {
		case errors.Is(err, repository.ErrEmployeeNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "employee with id " + idParam + " does not exist"})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
	}

	c.Status(http.StatusNoContent)
}

// HealthCheck handles GET /health
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "UP",
		"service": "employee-management",
	})
}
