// Package handlers contains HTTP request handlers for the API endpoints
package handlers

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"employee-management/internal/api"
	"employee-management/internal/models"
	"employee-management/internal/service"
	"employee-management/internal/validator"

	"github.com/gin-gonic/gin"
)

// TODO Check if the handler in any LOC does not just VALIDATE THE REQUEST, VALIDATE THE ID FORMAT, CALL THE SERVICE.

// EmployeeHandler handles HTTP requests for employee operations
type EmployeeHandler struct {
	service *service.EmployeeService
}

// NewEmployeeHandler creates a new EmployeeHandler instance
func NewEmployeeHandler(s *service.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{service: s}
}

// CreateEmployee godoc
//
//	@Summary		Create a new employee
//	@Description	Creates a new employee in the system. Only name, email and departmentID are required.
//	@Tags			Employees
//	@Accept			json
//	@Produce		json
//	@Param			employee	body		object{name=string,email=string,departmentID=string}	true	"Employee data (only name, email and departmentID needed)"
//	@Success		201			{object}	models.Employee		"Employee created successfully"
//
// @Failure		400			{object}	api.ErrorResponse	"Invalid JSON format, missing fields, invalid email format, or department not found"
//
//	@Failure		409			{object}	api.ErrorResponse	"Email already exists"
//	@Failure		503			{object}	api.ErrorResponse	"Departments service unavailable"
//	@Failure		500			{object}	api.ErrorResponse	"Internal server error"
//	@Security		BearerAuth
//	@Router			/employees [post]
func (h *EmployeeHandler) CreateEmployee(c *gin.Context) {
	var req models.Employee
	// Check JSON shape / types - service handles validation
	if err := c.ShouldBindJSON(&req); err != nil {
		api.BadRequest(c, "invalid JSON format")
		return
	}

	// Extract token from Authorization header
	authHeader := c.GetHeader("Authorization")
	tokenStr := ""
	if strings.HasPrefix(authHeader, "Bearer ") {
		tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
	}

	// Call service to handle business logic
	if err := h.service.Create(c.Request.Context(), &req, tokenStr); err != nil {
		switch {
		case errors.Is(err, service.ErrDepartmentRequired),
			errors.Is(err, service.ErrNameRequired),
			errors.Is(err, service.ErrEmailRequired),
			errors.Is(err, service.ErrDepartmentNotFound):
			api.BadRequest(c, err.Error())

		case errors.Is(err, service.ErrEmailAlreadyExists):
			api.Conflict(c, err.Error())

		case errors.Is(err, service.ErrDepartmentsServiceUnavailable):
			api.ServiceUnavailable(c, err.Error())
		case errors.Is(err, service.ErrEmailInvalidFormat),
			errors.Is(err, service.ErrNameInvalidFormat):
			api.BadRequest(c, err.Error())
		default:
			log.Printf("unexpected error creating employee: %v", err) // Log in docker-compose or when compiled
			api.InternalServerError(c, "failed to create employee")
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
//	@Param			id	path		int	true	"Employee ID"
//	@Success		200	{object}	models.Employee
//	@Failure		400	{object}	api.ErrorResponse	"Invalid ID format"
//	@Failure		404	{object}	api.ErrorResponse	"Employee not found"
//	@Failure		500	{object}	api.ErrorResponse	"Internal server error"
//	@Security		BearerAuth
//	@Router			/employees/{id} [get]
func (h *EmployeeHandler) GetEmployeeByID(c *gin.Context) {
	idParam := c.Param("id")

	id, validationResult := validator.ValidateID(idParam)
	if !validationResult.IsValid {
		api.ValidationError(c, http.StatusBadRequest,
			"invalid ID",
			validationResult.ToAPIErrorDetails())
		return
	}

	emp, err := h.service.FindByID(c.Request.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrEmployeeNotFound):
			api.NotFound(c, err.Error())
		default:
			api.InternalServerError(c, "failed to retrieve employee")
		}
		return
	}

	c.JSON(http.StatusOK, emp)
}

// GetAllEmployees godoc
//
//	@Summary		Get all employees with pagination and filtering
//	@Description	Retrieves employees with pagination support. Can filter by department, status.
//	@Tags			Employees
//	@Produce		json
//	@Param			page		query		int		false	"Page number (default: 1)"
//	@Param			page_size	query		int		false	"Number of items per page (default: 10, max: 100)"
//	@Param			department	query		string	false	"Filter by department ID"
//	@Param			status		query		string	false	"Filter by status (ACTIVE, ON_VACATION, RETIRED)"
//	@Success		200			{object}	api.PaginatedResponse
//	@Failure		400			{object}	api.ErrorResponse	"Invalid query parameters"
//	@Failure		500			{object}	api.ErrorResponse	"Internal server error"
//	@Security		BearerAuth
//	@Router			/employees [get]
func (h *EmployeeHandler) GetAllEmployees(c *gin.Context) {
	var query api.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		api.BadRequest(c, "invalid query parameters")
		return
	}

	// Set defaults for pagination (service also does this, but early validation is fine)
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 10
	} else if query.PageSize > 100 {
		query.PageSize = 100
	}

	// Build filters map
	filters := make(map[string]any)
	if query.Department != "" {
		filters["department"] = query.Department
	}
	if query.Status != "" {
		filters["status"] = query.Status
	}

	employees, total, err := h.service.FindAll(c.Request.Context(), query.Page, query.PageSize, filters)
	if err != nil {
		api.InternalServerError(c, "failed to retrieve employees")
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
//	@Param			employee	body		models.Employee		true	"Updated employee data (name, email, departmentID)"
//	@Success		200			{object}	models.Employee
//	@Failure		400			{object}	api.ErrorResponse	"Invalid ID, JSON format, or validation failed"
//	@Failure		404			{object}	api.ErrorResponse	"Employee not found"
//	@Failure		409			{object}	api.ErrorResponse	"Email already exists"
//	@Failure		500			{object}	api.ErrorResponse	"Internal server error"
//	@Security		BearerAuth
//	@Router			/employees/{id} [put]
func (h *EmployeeHandler) UpdateEmployee(c *gin.Context) {
	idParam := c.Param("id")

	id, validationResult := validator.ValidateID(idParam)
	if !validationResult.IsValid {
		api.ValidationError(c, http.StatusBadRequest,
			"invalid ID",
			validationResult.ToAPIErrorDetails())
		return
	}

	var req models.Employee
	if err := c.ShouldBindJSON(&req); err != nil {
		api.BadRequest(c, "invalid JSON format")
		return
	}

	// Extract token from Authorization header
	authHeader := c.GetHeader("Authorization")
	tokenStr := ""
	if strings.HasPrefix(authHeader, "Bearer ") {
		tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
	}

	req.ID = id

	updatedEmployee, err := h.service.Update(c.Request.Context(), &req, tokenStr)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrEmployeeNotFound):
			api.NotFound(c, err.Error())
		case errors.Is(err, service.ErrEmailAlreadyExists):
			api.Conflict(c, err.Error())
		case errors.Is(err, service.ErrDepartmentNotFound):
			api.BadRequest(c, err.Error())
		default:
			api.InternalServerError(c, "failed to update employee")
		}
		return
	}

	c.JSON(http.StatusOK, updatedEmployee)
}

// DeleteEmployee godoc
//
//	@Summary		Delete employee
//	@Description	Deletes an employee by ID
//	@Tags			Employees
//	@Param			id	path		int	true	"Employee ID"
//	@Success		204	"No Content"
//	@Failure		400	{object}	api.ErrorResponse	"Invalid ID format"
//	@Failure		404	{object}	api.ErrorResponse	"Employee not found"
//	@Failure		500	{object}	api.ErrorResponse	"Internal server error"
//	@Security		BearerAuth
//	@Router			/employees/{id} [delete]
func (h *EmployeeHandler) DeleteEmployee(c *gin.Context) {
	idParam := c.Param("id")

	id, validationResult := validator.ValidateID(idParam)
	if !validationResult.IsValid {
		api.ValidationError(c, http.StatusBadRequest,
			"invalid ID",
			validationResult.ToAPIErrorDetails())
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		switch {
		case errors.Is(err, service.ErrEmployeeNotFound):
			api.NotFound(c, err.Error())
		default:
			api.InternalServerError(c, "failed to delete employee")
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
