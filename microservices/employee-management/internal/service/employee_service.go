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

	"github.com/sony/gobreaker"
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
	repo           repository.EmployeeRepository
	httpClient     *http.Client
	departmentsURL string
	circuitBreaker *gobreaker.CircuitBreaker
}

// NewEmployeeService creates a new instance of EmployeeService
func NewEmployeeService(repo repository.EmployeeRepository) *EmployeeService {
	// Circuit breaker settings
	cbSettings := gobreaker.Settings{
		Name:        "Departments Service",
		MaxRequests: 3,                // Max requests allowed when half-open
		Interval:    10 * time.Second, // Interval for clearing counts
		Timeout:     30 * time.Second, // Time to wait before moving from open to half-open
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 5 && failureRatio >= 0.6
		},
	}

	return &EmployeeService{
		repo: repo,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		departmentsURL: getDepartmentsServiceURL(),
		circuitBreaker: gobreaker.NewCircuitBreaker(cbSettings),
	}
}

// getDepartmentsServiceURL reads from environment or uses default
// Is this stupid? TODO put this in an environment variable
func getDepartmentsServiceURL() string {
	url := "http://departments-service:8082"
	return url
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
	exists, err := s.checkDepartmentExistsWithRetry(ctx, e.DepartmentID)
	if err != nil {
		return fmt.Errorf("department validation failed after retries: %w", err)
	}
	if !exists {
		return ErrDepartmentNotFound
	}

	// Set defaults
	e.Status = models.StatusActive
	e.HireDate = time.Now()

	// Create in database
	return s.repo.Create(ctx, e)
}

// checkDepartmentExistsWithRetry wraps the call with retries and circuit breaker
func (s *EmployeeService) checkDepartmentExistsWithRetry(ctx context.Context, departmentID string) (bool, error) {
	const maxRetries = 3
	backoff := 100 * time.Millisecond

	for attempt := 0; attempt < maxRetries; attempt++ {
		// Execute through circuit breaker
		result, err := s.circuitBreaker.Execute(func() (interface{}, error) {
			return s.checkDepartmentExists(ctx, departmentID)
		})

		if err == nil {
			return result.(bool), nil
		}

		// If circuit breaker is open, fail fast
		if err == gobreaker.ErrOpenState {
			return false, fmt.Errorf("circuit breaker open: departments service unavailable")
		}

		// For other errors, retry after backoff
		if attempt < maxRetries-1 {
			time.Sleep(backoff)
			backoff *= 2 // exponential backoff
		}
	}
	return false, fmt.Errorf("max retries exceeded")
}

// checkDepartmentExists calls the departments microservice (actual HTTP call)
func (s *EmployeeService) checkDepartmentExists(ctx context.Context, departmentID string) (bool, error) {
	url := fmt.Sprintf("%s/departments-service/api/departments/%s", s.departmentsURL, departmentID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return false, err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		// Network errors, timeouts, etc.
		return false, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return true, nil
	case http.StatusNotFound:
		return false, nil
	default:
		// Treat 5xx as errors that should trigger retries/circuit breaker
		if resp.StatusCode >= 500 {
			return false, fmt.Errorf("departments service returned 5xx: %d", resp.StatusCode)
		}
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
