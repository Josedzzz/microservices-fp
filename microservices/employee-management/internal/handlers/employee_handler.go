// Package handlers contains HTTP request handlers for the API endpoints
package handlers

import (
	"errors"
	"net/http"
	"time"

	"employee-management/internal/api"
	"employee-management/internal/models"
	"employee-management/internal/repository"
	"employee-management/internal/service"
	"employee-management/internal/validator"

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

// CreateEmployee godoc
//
//	@Summary		Create a new employee
//	@Description	Creates a new employee in the system
//	@Tags			Employees
//	@Accept			json
//	@Produce		json
//	@Param			employee	body		models.Employee		true	"Employee data"
//	@Success		201			{object}	models.Employee		"Employee created successfully"
//	@Failure		400			{object}	api.ErrorResponse	"Invalid JSON format or validation failed"
//	@Failure		409			{object}	api.ErrorResponse	"Email or employee number already exists"
//	@Failure		500			{object}	api.ErrorResponse	"Internal server error"
//	@Router			/employees [post]
func (h *EmployeeHandler) CreateEmployee(c *gin.Context) {
	var req models.Employee

	// Check JSON shape / types
	if err := c.ShouldBindJSON(&req); err != nil {
		api.BadRequest(c, "Invalid JSON format")
		return
	}

	// Input validation
	validation := validator.ValidateEmployee(req.Email, req.EmployeeNumber, req.FirstName, req.LastName)

	if !validation.IsValid {
		api.ValidationError(c, http.StatusBadRequest, "Validation failed", validation.Errors)
		return
	}

	// Business logic
	if err := h.service.Create(c.Request.Context(), &req); err != nil {
		switch {
		case errors.Is(err, repository.ErrEmailAlreadyExists):
			api.Conflict(c, "Email already exist")
		case errors.Is(err, repository.ErrEmployeeNumberAlreadyExists):
			api.Conflict(c, "Employee number already exists")
		default:
			api.InternalServerError(c, "Failed to create employee")
		}
		return
	}

	c.JSON(http.StatusCreated, req)
}

// GetEmployeeByID godoc
//
//	@Summary		Get employee by ID
//	@Description	Retrieves an employee by its ID
//	@Tags			Employees
//	@Produce		json
//	@Param			id	path		int					true	"Employee ID"
//	@Success		200	{object}	models.Employee		"Employee found"
//	@Failure		400	{object}	api.ErrorResponse	"Invalid ID format"
//	@Failure		404	{object}	api.ErrorResponse	"Employee not found"
//	@Failure		500	{object}	api.ErrorResponse	"Internal server error"
//	@Router			/employees/{id} [get]
func (h *EmployeeHandler) GetEmployeeByID(c *gin.Context) {
	idParam := c.Param("id")

	id, errs := validator.ValidateID(idParam)
	if errs != nil {
		api.ValidationError(c, http.StatusBadRequest, "Invalid ID", errs)
		return
	}

	emp, err := h.service.FindByID(c.Request.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrEmployeeNotFound):
			api.NotFound(c, "Employee not found")
		default:
			api.InternalServerError(c, "Failed to retrieve employee")
		}
		return
	}

	c.JSON(http.StatusOK, emp)
}

// GetAllEmployees godoc
// @Summary Get all employees with pagination and filtering
// @Description Retrieves employees with pagination support. Can filter by department, status, position.
// @Tags Employees
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Number of items per page (default: 10, max: 100)"
// @Param department query string false "Filter by department"
// @Param status query string false "Filter by status (ACTIVE, ON_VACATION, RETIRED)"
// @Param position query string false "Filter by position"
// @Success 200 {object} api.PaginatedResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /employees [get]
func (h *EmployeeHandler) GetAllEmployees(c *gin.Context) {
	var query api.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query parameters"})
		return
	}

	// Set defaults for pagination
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 10
	} else if query.PageSize > 100 {
		query.PageSize = 100
	}

	// Build filters map
	filters := make(map[string]interface{})
	if query.Department != "" {
		filters["department"] = query.Department
	}
	if query.Status != "" {
		filters["status"] = query.Status
	}
	if query.Position != "" {
		filters["position"] = query.Position
	}

	employees, total, err := h.service.FindAll(c.Request.Context(), query.Page, query.PageSize, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	totalPages := (total + query.PageSize - 1) / query.PageSize

	response := api.PaginatedResponse{
		Data: employees,
		Pagination: api.PaginationMeta{
			CurrentPage:  query.Page,
			PageSize:     query.PageSize,
			TotalPages:   totalPages,
			TotalRecords: total,
		},
	}

	c.JSON(http.StatusOK, response)
}

// UpdateEmployee godoc
//
//	@Summary		Update employee
//	@Description	Updates an existing employee
//	@Tags			Employees
//	@Accept			json
//	@Produce		json
//	@Param			id			path		int					true	"Employee ID"
//	@Param			employee	body		models.Employee		true	"Updated employee data"
//	@Success		200			{object}	models.Employee		"Employee updated successfully"
//	@Failure		400			{object}	api.ErrorResponse	"Invalid JSON format or validation failed"
//	@Failure		404			{object}	api.ErrorResponse	"Employee not found"
//	@Failure		409			{object}	api.ErrorResponse	"Email or employee number already exists"
//	@Failure		500			{object}	api.ErrorResponse	"Internal server error"
//	@Router			/employees/{id} [put]
func (h *EmployeeHandler) UpdateEmployee(c *gin.Context) {
	idParam := c.Param("id")

	id, errs := validator.ValidateID(idParam)
	if errs != nil {
		api.ValidationError(c, http.StatusBadRequest, "Invalid ID", errs)
		return
	}

	var req models.Employee
	if err := c.ShouldBindJSON(&req); err != nil {
		api.BadRequest(c, "Invalid JSON format")
		return
	}

	req.ID = id

	validation := validator.ValidateEmployee(
		req.Email,
		req.EmployeeNumber,
		req.FirstName,
		req.LastName,
	)

	if !validation.IsValid {
		api.ValidationError(c, http.StatusBadRequest, "Validation failed", validation.Errors)
		return
	}

	if err := h.service.Update(c.Request.Context(), &req); err != nil {
		switch {
		case errors.Is(err, repository.ErrEmployeeNotFound):
			api.NotFound(c, "Employee not found")
		case errors.Is(err, repository.ErrEmailAlreadyExists):
			api.Conflict(c, "Email already exists")
		case errors.Is(err, repository.ErrEmployeeNumberAlreadyExists):
			api.Conflict(c, "Employee number already exists")
		default:
			api.InternalServerError(c, "Failed to update employee")
		}
		return
	}

	c.JSON(http.StatusOK, req)
}

// DeleteEmployee godoc
//
//	@Summary		Delete employee
//	@Description	Deletes an employee by ID
//	@Tags			Employees
//	@Param			id	path	int	true	"Employee ID"
//	@Success		204	"Employee deleted successfully (no content)"
//	@Failure		400	{object}	api.ErrorResponse	"Invalid ID format"
//	@Failure		404	{object}	api.ErrorResponse	"Employee not found"
//	@Failure		500	{object}	api.ErrorResponse	"Internal server error"
//	@Router			/employees/{id} [delete]
func (h *EmployeeHandler) DeleteEmployee(c *gin.Context) {
	idParam := c.Param("id")

	id, errs := validator.ValidateID(idParam)
	if errs != nil {
		api.ValidationError(c, http.StatusBadRequest, "Invalid ID", errs)
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		switch {
		case errors.Is(err, repository.ErrEmployeeNotFound):
			api.NotFound(c, "Employee not found")
		default:
			api.InternalServerError(c, "Failed to delete employee")
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// HealthCheck handles GET /health
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "UP",
		"service":   "employee-management",
		"timestamp": time.Now().UTC(),
	})
}
