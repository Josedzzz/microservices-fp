// Package repository defines data access layer interfaces and implementations
package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"employee-management/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// EmployeeRepository defines the interface for employee data operations
type EmployeeRepository interface {
	Create(ctx context.Context, e *models.Employee) error
	FindByID(ctx context.Context, id int64) (*models.Employee, error)
	FindAll(ctx context.Context, limit, offset int, filters map[string]any) ([]models.Employee, error)
	Count(ctx context.Context, filters map[string]any) (int, error)
	Update(ctx context.Context, e *models.Employee) error
	Delete(ctx context.Context, id int64) error
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
        (name, email, department, status, hire_date)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, created_at, updated_at
    `

	err := r.db.QueryRow(ctx, query,
		e.Name,
		e.Email,
		e.DepartmentID,
		e.Status,
		e.HireDate,
	).Scan(&e.ID, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			switch pgErr.ConstraintName {
			case "employees_email_key":
				return ErrEmailAlreadyExists
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
		SELECT id, name, email, department, status, hire_date, created_at, updated_at
		FROM employee.employees 
		WHERE id = $1
	`

	var emp models.Employee
	err := r.db.QueryRow(ctx, query, id).Scan(
		&emp.ID,
		&emp.Name,
		&emp.Email,
		&emp.DepartmentID,
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

// FindAll retrieves all employees from the database
func (r *employeeRepository) FindAll(ctx context.Context, limit, offset int, filters map[string]any) ([]models.Employee, error) {
	baseQuery := `SELECT id, name, email, department, status, hire_date, created_at, updated_at
				  FROM employee.employees`
	var conditions []string
	var args []any
	argPos := 1

	if dept, ok := filters["department"]; ok && dept != "" {
		conditions = append(conditions, fmt.Sprintf("department = $%d", argPos))
		args = append(args, dept)
		argPos++
	}
	if status, ok := filters["status"]; ok && status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argPos))
		args = append(args, status)
		argPos++
	}

	if len(conditions) > 0 {
		baseQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	baseQuery += " ORDER BY created_at DESC"
	baseQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argPos, argPos+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query employees: %w", err)
	}
	defer rows.Close()

	var employees []models.Employee
	for rows.Next() {
		var emp models.Employee
		err := rows.Scan(
			&emp.ID,
			&emp.Name,
			&emp.Email,
			&emp.DepartmentID,
			&emp.Status,
			&emp.HireDate,
			&emp.CreatedAt,
			&emp.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan employee row: %w", err)
		}
		employees = append(employees, emp)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating employee rows: %w", err)
	}

	return employees, nil
}

func (r *employeeRepository) Count(ctx context.Context, filters map[string]any) (int, error) {
	baseQuery := `SELECT COUNT(*) FROM employee.employees`
	var conditions []string
	var args []any
	argPos := 1

	// same filter logic
	if dept, ok := filters["department"]; ok && dept != "" {
		conditions = append(conditions, fmt.Sprintf("department = $%d", argPos))
		args = append(args, dept)
		argPos++
	}
	if status, ok := filters["status"]; ok && status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argPos))
		args = append(args, status)
		argPos++
	}
	if pos, ok := filters["position"]; ok && pos != "" {
		conditions = append(conditions, fmt.Sprintf("position = $%d", argPos))
		args = append(args, pos)
		argPos++
	}

	if len(conditions) > 0 {
		baseQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	var count int
	err := r.db.QueryRow(ctx, baseQuery, args...).Scan(&count)
	return count, err
}

// Update modifies an existing employee record
func (r *employeeRepository) Update(ctx context.Context, e *models.Employee) error {
	query := `
		UPDATE employee.employees 
		SET name = $2, email = $3, department = $4, 
			status = $5, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING updated_at
	`

	result, err := r.db.Exec(ctx, query,
		e.ID,
		e.Name,
		e.Email,
		e.DepartmentID,
		e.Status,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch {
			case pgErr.Code == "23505" && pgErr.ConstraintName == "employees_email_key":
				return ErrEmailAlreadyExists
			}
		}
		return fmt.Errorf("failed to update employee: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrEmployeeNotFound
	}

	err = r.db.QueryRow(ctx, "SELECT updated_at FROM employee.employees WHERE id = $1", e.ID).Scan(&e.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to get updated timestamp: %w", err)
	}

	return nil
}

// Delete removes an employee from the db by id
func (r *employeeRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM employee.employees WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23503" { // foreign_key_violation
				return fmt.Errorf("employee has related records and cannot be deleted: %w", err)
			}
		}
		return fmt.Errorf("failed to delete employee: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrEmployeeNotFound
	}

	return nil
}
