// Package repository defines data access layer interfaces and implementations
package repository

import (
	"context"
	"fmt"
	"errors"

	"employee-management/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/pgconn"
)

// EmployeeRepository defines the interface for employee data operations
type EmployeeRepository interface {
	Create(ctx context.Context, e *models.Employee) error
	FindByID(ctx context.Context, id int64) (*models.Employee, error)
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

// Declaration of domain errors.
var (
    ErrEmailAlreadyExists          = errors.New("email already exists")
    ErrEmployeeNumberAlreadyExists = errors.New("employee number already exists")
    ErrEmployeeAlreadyExists       = errors.New("employee already exists")
	ErrEmployeeNotFound            = errors.New("employee not found")
)

// Create adds a new employee to the database
func (r *employeeRepository) Create(ctx context.Context, e *models.Employee) error {
    query := `
        INSERT INTO employee.employees
        (first_name, last_name, email, employee_number, position, department, status, hire_date)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id, created_at, updated_at
    `

    err := r.db.QueryRow(ctx, query,
        e.FirstName,
        e.LastName,
        e.Email,
        e.EmployeeNumber,
        e.Position,
        e.Department,
        e.Status,
        e.HireDate,
    ).Scan(&e.ID, &e.CreatedAt, &e.UpdatedAt)

    if err != nil {
        var pgErr *pgconn.PgError
        if errors.As(err, &pgErr) && pgErr.Code == "23505" {
            switch pgErr.ConstraintName {
            case "employees_email_key":
                return ErrEmailAlreadyExists
            case "employees_employee_number_key":
                return ErrEmployeeNumberAlreadyExists
            default:
                return ErrEmployeeAlreadyExists
            }
        }
        return err
    }

    return nil
}

// FindByID retrieves an employee by their id
func (r *employeeRepository) FindByID(ctx context.Context, id int64) (*models.Employee, error) {
    query := `
        SELECT id, first_name, last_name, email, employee_number, 
               position, department, status, hire_date, created_at, updated_at
        FROM employee.employees 
        WHERE id = $1
    `
    
    var emp models.Employee
    err := r.db.QueryRow(ctx, query, id).Scan(
        &emp.ID,
        &emp.FirstName,
        &emp.LastName,
        &emp.Email,
        &emp.EmployeeNumber,
        &emp.Position,
        &emp.Department,
        &emp.Status,
        &emp.HireDate,
        &emp.CreatedAt,
        &emp.UpdatedAt,
    )
    
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, ErrEmployeeNotFound
        }
        return nil, err
    }

    return &emp, nil
}

// FindAll retrives all employees from the db
func (r *employeeRepository) FindAll(ctx context.Context) ([]models.Employee, error) {
    query := `
        SELECT id, first_name, last_name, email, employee_number, 
               position, department, status, hire_date, created_at, updated_at
        FROM employee.employees
        ORDER BY created_at DESC
    `
    
    rows, err := r.db.Query(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var employees []models.Employee
    for rows.Next() {
        var emp models.Employee
        err := rows.Scan(
            &emp.ID,
            &emp.FirstName,
            &emp.LastName,
            &emp.Email,
            &emp.EmployeeNumber,
            &emp.Position,
            &emp.Department,
            &emp.Status,
            &emp.HireDate,
            &emp.CreatedAt,
            &emp.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }
        employees = append(employees, emp)
    }
    
    return employees, nil
}

// Update modifies an existing employee record
func (r *employeeRepository) Update(ctx context.Context, e *models.Employee) error {
    query := `
        UPDATE employee.employees 
        SET first_name = $2, last_name = $3, email = $4, 
            employee_number = $5, position = $6, department = $7,
            status = $8, updated_at = CURRENT_TIMESTAMP
        WHERE id = $1
        RETURNING updated_at
    `
    
    return r.db.QueryRow(ctx, query,
        e.ID,
        e.FirstName,
        e.LastName,
        e.Email,
        e.EmployeeNumber,
        e.Position,
        e.Department,
        e.Status,
    ).Scan(&e.UpdatedAt)
}

// Delete removes an employee from the db by id
func (r *employeeRepository) Delete(ctx context.Context, id string) error {
    query := `DELETE FROM employee.employees WHERE id = $1`
    result, err := r.db.Exec(ctx, query, id)
    if err != nil {
        return err
    }
    
    if result.RowsAffected() == 0 {
        return fmt.Errorf("employee not found")
    }
    
    return nil
}

