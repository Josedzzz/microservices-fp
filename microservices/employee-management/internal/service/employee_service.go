// Package service contains business logic and app services
package service

import (
	"context"
	"time"

	"employee-management/internal/models"
	"employee-management/internal/repository"
)

// EmployeeService handles business logic for employee operations
// It acts as an intermediary between API handlers and the data repository
type EmployeeService struct {
	repo repository.EmployeeRepository
}

// NewEmployeeService creates a new instance of EmployeeService
func NewEmployeeService(repo repository.EmployeeRepository) *EmployeeService {
	return &EmployeeService{repo: repo}
}

// Create adds a new employee to the database
func (s *EmployeeService) Create(ctx context.Context, e *models.Employee) error {
	e.Status = models.StatusActive
	e.HireDate = time.Now()
	return s.repo.Create(ctx, e)
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
