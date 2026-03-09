// Package service contains business logic and app services
package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"employee-management/internal/messaging"
	"employee-management/internal/models"
	"employee-management/internal/repository"
	"employee-management/internal/validator"

	"github.com/sony/gobreaker"
)

// Define service-level erros
var (
	ErrDepartmentRequired            = errors.New("department ID is required")
	ErrNameRequired                  = errors.New("name is required")
	ErrEmailRequired                 = errors.New("email is required")
	ErrDepartmentNotFound            = errors.New("department ID does not exists")
	ErrDepartmentsServiceUnavailable = errors.New("departments service is currently unavailable")
	ErrEmailAlreadyExists            = errors.New("email already exists")
	ErrEmployeeNotFound              = errors.New("employee not found")
	ErrEmailInvalidFormat            = errors.New("invalid email format")
	ErrNameInvalidFormat             = errors.New("invalid name format")
)

// EmployeeService handles business logic for employee operations
// It acts as an intermediary between API handlers and the data repository
type EmployeeService struct {
	repo           repository.EmployeeRepository
	publisher      *messaging.Publisher
	httpClient     *http.Client
	departmentsURL string
	circuitBreaker *gobreaker.CircuitBreaker
}

// NewEmployeeService creates a new instance of EmployeeService
func NewEmployeeService(repo repository.EmployeeRepository, publisher *messaging.Publisher) *EmployeeService {
	// Circuit breaker settings
	cbSettings := gobreaker.Settings{
		Name:        "Departments Service",
		MaxRequests: 3,
		Interval:    10 * time.Second,
		Timeout:     30 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			trip := counts.Requests >= 5 && failureRatio >= 0.6
			if trip {
				log.Printf("Circuit breaker tripped for departments service")
			}
			return trip
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			log.Printf("Circuit breaker '%s' changed from %s to %s", name, from, to)
		},
	}

	return &EmployeeService{
		repo:      repo,
		publisher: publisher,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		departmentsURL: getDepartmentsServiceURL(),
		circuitBreaker: gobreaker.NewCircuitBreaker(cbSettings),
	}
}

// getDepartmentsServiceURL uses a default value
func getDepartmentsServiceURL() string {
	// TODO put this in an environment variable eventually
	url := "http://departments-service:8082"
	return url
}

// Create adds a new employee to the database
func (s *EmployeeService) Create(ctx context.Context, e *models.Employee) error {
	// Validate field formats using "internal/validator"
	validation := validator.ValidateEmployee(e.Email, e.Name)
	if !validation.IsValid {
		for _, verr := range validation.Errors {
			switch verr.Field {
			case "email":
				if verr.Code == validator.CodeRequired {
					return ErrEmailRequired
				}
				return ErrEmailInvalidFormat
			case "name":
				if verr.Code == validator.CodeRequired {
					return ErrNameRequired
				}
				return ErrNameInvalidFormat
			}
		}
		// TODO Add as an service level error
		return errors.New("validation failed")
	}

	// Validate DepartmentID (validator doesn't handle this yet).
	// Should it even handle it?
	if e.DepartmentID == "" {
		return ErrDepartmentRequired
	}

	// Check if the department ID exists via HTTP call to the departments service
	deptExists, err := s.checkDepartmentExistsWithRetry(ctx, e.DepartmentID)
	if err != nil {
		return err
	}
	if !deptExists {
		return ErrDepartmentNotFound
	}

	// Check if the email already exists with a call to the repository
	emailExists, err := s.checkEmailExists(ctx, e.Email)
	if err != nil {
		return err
	}
	if emailExists {
		return ErrEmailAlreadyExists
	}

	// Set defaults
	e.Status = models.StatusActive
	e.HireDate = time.Now()

	// Default values for auth integration
	if e.Role == "" {
		e.Role = "USER"
	}

	// Create in database
	errdb := s.repo.Create(ctx, e)
	if errdb != nil {
		return errdb
	}

	log.Print("Attempting to publish event for employee:", e.ID)

	event := messaging.EmployeeCreatedEvent{
		ID:           e.ID,
		Name:         e.Name,
		Email:        e.Email,
		Role:         e.Role,
		DepartmentID: e.DepartmentID,
		Status:       e.Status,
		HireDate:     e.HireDate,
		CreatedAt:    e.CreatedAt,
	}

	log.Printf("Event payload: %+v", event)

	if err := s.publisher.Publish(
		messaging.EventEmployeeCreated,
		event,
	); err != nil {
		log.Printf("failed to publish employee.created event: %v", err)
	} else {
		log.Printf("Event published successfully to exchange: %s, routing key: %s",
			messaging.EmployeesExchange, messaging.EventEmployeeCreated)
	}

	return nil
}

func (s *EmployeeService) checkEmailExists(ctx context.Context, email string) (bool, error) {
	employee, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return false, err
	}
	return employee != nil, nil
}

func (s *EmployeeService) checkDepartmentExistsWithRetry(ctx context.Context, departmentID string) (bool, error) {
	const maxRetries = 3
	backoff := 100 * time.Millisecond

	for attempt := range [maxRetries]int{} {
		result, err := s.circuitBreaker.Execute(func() (any, error) {
			return s.checkDepartmentExists(ctx, departmentID)
		})

		if err == nil {
			return result.(bool), nil
		}

		// Circuit breaker open
		if err == gobreaker.ErrOpenState {
			log.Printf("Circuit breaker open for departments service")
			return false, ErrDepartmentsServiceUnavailable
		}

		// For other errors, retry
		if attempt < maxRetries-1 {
			log.Printf("Departments service call failed (attempt %d): %v. Retrying...", attempt+1, err)
			time.Sleep(backoff)
			backoff *= 2
		}
	}
	log.Printf("All retries exhausted for departments service")
	return false, ErrDepartmentsServiceUnavailable
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
	employee, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve employee: %w", err)
	}
	if employee == nil {
		return nil, ErrEmployeeNotFound
	}
	return employee, nil
}

// FindAll retrieves all employees
func (s *EmployeeService) FindAll(ctx context.Context, page, pageSize int, filters map[string]any) ([]models.Employee, int, error) {
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

// Update updates an employee and returns the updated employee
func (s *EmployeeService) Update(ctx context.Context, e *models.Employee) (*models.Employee, error) {
	log.Printf("UPDATE - Request received: ID=%d, Name='%s', Email='%s', Dept='%s', Status='%s'",
		e.ID, e.Name, e.Email, e.DepartmentID, e.Status)

	// Check if employee exists
	existing, err := s.repo.FindByID(ctx, e.ID)
	if err != nil {
		log.Printf("UPDATE - FindByID error: %v", err)
		if strings.Contains(err.Error(), "not found") {
			return nil, ErrEmployeeNotFound
		}
		return nil, fmt.Errorf("failed to verify employee existence: %w", err)
	}
	log.Printf("UPDATE - Existing employee: Name='%s', Email='%s', Dept='%s', Status='%s'",
		existing.Name, existing.Email, existing.DepartmentID, existing.Status)

	// Track what needs validation
	emailChanged := false
	deptChanged := false

	// Name update (if provided)
	if e.Name != "" {
		log.Printf("UPDATE - Name provided: '%s'", e.Name)
		if e.Name != existing.Name {
			log.Printf("UPDATE - Name changing from '%s' to '%s'", existing.Name, e.Name)
			// Validate name
			validation := validator.ValidateField("name", e.Name)
			if !validation.IsValid {
				log.Printf("UPDATE - Name validation failed: %v", validation.Errors)
				return nil, fmt.Errorf("invalid name: %s", validation.Errors[0].Message)
			}
			existing.Name = e.Name
		} else {
			log.Printf("UPDATE - Name unchanged")
		}
	} else {
		log.Printf("UPDATE - Name not provided (empty string), keeping existing: '%s'", existing.Name)
	}

	// Email update (if provided)
	if e.Email != "" {
		log.Printf("UPDATE - Email provided: '%s'", e.Email)
		if e.Email != existing.Email {
			log.Printf("UPDATE - Email changing from '%s' to '%s'", existing.Email, e.Email)
			// Validate email format
			validation := validator.ValidateField("email", e.Email)
			if !validation.IsValid {
				log.Printf("UPDATE - Email validation failed: %v", validation.Errors)
				return nil, fmt.Errorf("invalid email: %s", validation.Errors[0].Message)
			}
			emailChanged = true
		} else {
			log.Printf("UPDATE - Email unchanged")
		}
	} else {
		log.Printf("UPDATE - Email not provided (empty string), keeping existing: '%s'", existing.Email)
	}

	// Department update (if provided)
	if e.DepartmentID != "" {
		log.Printf("UPDATE - Department provided: '%s'", e.DepartmentID)
		if e.DepartmentID != existing.DepartmentID {
			log.Printf("UPDATE - Department changing from '%s' to '%s'", existing.DepartmentID, e.DepartmentID)
			deptChanged = true
		} else {
			log.Printf("UPDATE - Department unchanged")
		}
	} else {
		log.Printf("UPDATE - Department not provided (empty string), keeping existing: '%s'", existing.DepartmentID)
	}

	// Status update (if provided)
	if e.Status != "" {
		log.Printf("UPDATE - Status provided: '%s'", e.Status)
		if e.Status != existing.Status {
			log.Printf("UPDATE - Status changing from '%s' to '%s'", existing.Status, e.Status)
			existing.Status = e.Status
		} else {
			log.Printf("UPDATE - Status unchanged")
		}
	} else {
		log.Printf("UPDATE - Status not provided (empty string), keeping existing: '%s'", existing.Status)
	}

	// Perform validations that need external calls
	if emailChanged {
		log.Printf("UPDATE - Checking if email '%s' already exists", e.Email)
		emailExists, err := s.checkEmailExists(ctx, e.Email)
		if err != nil {
			log.Printf("UPDATE - Email check error: %v", err)
			return nil, err
		}
		if emailExists {
			log.Printf("UPDATE - Email '%s' already exists", e.Email)
			return nil, ErrEmailAlreadyExists
		}
		log.Printf("UPDATE - Email '%s' is available", e.Email)
		existing.Email = e.Email
	}

	if deptChanged {
		log.Printf("UPDATE - Checking if department '%s' exists", e.DepartmentID)
		deptExists, err := s.checkDepartmentExistsWithRetry(ctx, e.DepartmentID)
		if err != nil {
			log.Printf("UPDATE - Department check error: %v", err)
			return nil, err
		}
		if !deptExists {
			log.Printf("UPDATE - Department '%s' not found", e.DepartmentID)
			return nil, ErrDepartmentNotFound
		}
		log.Printf("UPDATE - Department '%s' exists", e.DepartmentID)
		existing.DepartmentID = e.DepartmentID
	}

	log.Printf("UPDATE - Final employee to save: Name='%s', Email='%s', Dept='%s', Status='%s'",
		existing.Name, existing.Email, existing.DepartmentID, existing.Status)

	// Perform update with merged data
	err = s.repo.Update(ctx, existing)
	if err != nil {
		log.Printf("UPDATE - Repository update error: %v", err)
		return nil, err
	}

	// Fetch the updated employee to return
	updated, err := s.repo.FindByID(ctx, e.ID)
	if err != nil {
		log.Printf("UPDATE - Failed to fetch updated employee: %v", err)
		return nil, fmt.Errorf("failed to fetch updated employee: %w", err)
	}

	log.Printf("UPDATE - Successfully updated employee ID=%d", e.ID)
	return updated, nil
}

// Delete retires an employee and publishes employee.eliminated event
func (s *EmployeeService) Delete(ctx context.Context, id int64) error {
	// Try to fetch employee first (for event data)
	employee, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to retrieve employee: %w", err)
	}

	// If employee does not exist → idempotent success
	if employee == nil {
		return nil
	}

	// If already retired → idempotent success
	if employee.Status == models.StatusRetired {
		return nil
	}

	// Soft delete (set status = RETIRED)
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to retire employee: %w", err)
	}

	// Build event
	event := messaging.EmployeeDeletedEvent{
		ID:    employee.ID,
		Name:  employee.Name,
		Email: employee.Email,
	}

	// Publish event (non-blocking / non-transactional)
	if err := s.publisher.Publish(
		messaging.EventEmployeeDeleted,
		event,
	); err != nil {
		log.Printf("failed to publish employee.deleted event: %v", err)
	}

	return nil
}
