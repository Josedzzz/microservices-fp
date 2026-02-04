// Package repository defines data acces layer interfaces and implementations
package repository

import (
	"context"

	"employee-management/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

// EmployeeRepository defines the interface for employee data operations
type EmployeeRepository interface {
	Create(ctx context.Context, e *models.Employee) error
	FindByID(ctx context.Context, id string) (*models.Employee, error)
	FindAll(ctx context.Context) ([]models.Employee, error)
	Update(ctx context.Context, e *models.Employee) error
	Delete(ctx context.Context, id string) error
}

// employeeRepository is the postgresql implementation of EmployeeRepository
type employeeRepository struct {
	db *pgxpool.Pool // db connection pool
}

// NewEmployeeRepository creates a new instance of EmployeeRepository
func NewEmployeeRepository(db *pgxpool.Pool) EmployeeRepository {
	return &employeeRepository{db: db}
}

// TODO
// Create adds a new employee to the database
func (r *employeeRepository) Create(ctx context.Context, e *models.Employee) error {
	return nil
}

// TODO
// FindByID retrieves an employee by their id
func (r *employeeRepository) FindByID(ctx context.Context, id string) (*models.Employee, error) {
	return nil, nil
}

// TODO
// FindAll retrives all employees from the db
func (r *employeeRepository) FindAll(ctx context.Context) ([]models.Employee, error) {
	return nil, nil
}

// TODO
// Update modifies an existing employee record
func (r *employeeRepository) Update(ctx context.Context, e *models.Employee) error {
	return nil
}

// TODO
// Delete removes an employee from the db by id
func (r *employeeRepository) Delete(ctx context.Context, id string) error {
	return nil
}
