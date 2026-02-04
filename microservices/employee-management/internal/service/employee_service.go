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
func (s *EmployeeService) FindAll(ctx context.Context) ([]models.Employee, error) {
	return s.repo.FindAll(ctx)
}

// Update updates an employee
func (s *EmployeeService) Update(ctx context.Context, e *models.Employee) error {
	return s.repo.Update(ctx, e)
}

// Delete removes an employee
func (s *EmployeeService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

