// Package service contains business logic and app services
package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"employee-management/internal/models"
	"employee-management/internal/repository"
	"employee-management/internal/validator"
)

// Define service-level erros
var (
	ErrDepartmentRequired = errors.New("department ID is required")
	ErrNameRequired       = errors.New("name is required")
	ErrEmailRequired      = errors.New("email is required")
	ErrDepartmentNotFound = errors.New("department ID does not exists")
)

// EmployeeService handles business logic for employee operations
// It acts as an intermediary between API handlers and the data repository
type EmployeeService struct {
	repo       repository.EmployeeRepository
	httpClient *http.Client
}

// NewEmployeeService creates a new instance of EmployeeService
func NewEmployeeService(repo repository.EmployeeRepository) *EmployeeService {
	return &EmployeeService{
		repo: repo,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// Create adds a new employee to the database
func (s *EmployeeService) Create(ctx context.Context, e *models.Employee) error {
	// Validate field formats using "internal/validator"
	validation := validator.ValidateEmployee(e.Email, e.Name)
	if !validation.IsValid {
		// Convert validator errors to service errors
		for _, errDetail := range validation.Errors {
			switch errDetail.Field {
			case "email":
				if errDetail.Message == "Email is required" {
					return ErrEmailRequired
				}
				return fmt.Errorf("invalid email: %s", errDetail.Message)
			case "name":
				if errDetail.Message == "Name is required" {
					return ErrNameRequired
				}
				return fmt.Errorf("invalid name: %s", errDetail.Message)
			}
		}
		return errors.New("validation failed")
	}

	// Validate DepartmentID (validator doesn't handle this yet)
	if e.DepartmentID == "" {
		return ErrDepartmentRequired
	}

	// Check if the department ID exists via HTTP call to the departments service
	exists, err := s.checkDepartmentExists(ctx, e.DepartmentID)
	if err != nil {
		return fmt.Errorf("department validation failed: %w", err)
	}
	if !exists {
		return ErrDepartmentNotFound // FIXED: using exported error
	}

	// Set defaults
	e.Status = models.StatusActive
	e.HireDate = time.Now()

	// Create in database
	return s.repo.Create(ctx, e)
}

// checkDepartmentExists calls the departments microservice
func (s *EmployeeService) checkDepartmentExists(ctx context.Context, departmentID string) (bool, error) {
	// Should this be an ENV variable? Revealing this in production could be dangerous!!!
	departmentsServiceURL := "http://departments-service:8082"
	url := fmt.Sprintf("%s/departments-service/api/departments/%s", departmentsServiceURL, departmentID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return false, err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return true, nil
	case http.StatusNotFound:
		return false, nil
	default:
		return false, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}
}

// FindByID retrieves an employee by id
func (s *EmployeeService) FindByID(ctx context.Context, id int64) (*models.Employee, error) {
	return s.repo.FindByID(ctx, id)
}

// FindAll retrieves all employees
func (s *EmployeeService) FindAll(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]models.Employee, int, error) {
	// Defensive programming protocols!!!
	// Validate and set defaults
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10 // default page size
	}
	// Enforcing max page size might result useful?
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	employees, err := s.repo.FindAll(ctx, pageSize, offset, filters)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.Count(ctx, filters)
	if err != nil {
		return nil, 0, err
	}

	return employees, total, nil
}

// Update updates an employee
func (s *EmployeeService) Update(ctx context.Context, e *models.Employee) error {
	return s.repo.Update(ctx, e)
}

// Delete removes an employee
func (s *EmployeeService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
