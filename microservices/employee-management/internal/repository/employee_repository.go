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
	FindAll(ctx context.Context, limit, offset int, filters map[string]interface{}) ([]models.Employee, error)
	Count(ctx context.Context, filters map[string]interface{}) (int, error)
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
func (r *employeeRepository) FindAll(ctx context.Context, limit, offset int, filters map[string]interface{}) ([]models.Employee, error) {
	baseQuery := `SELECT id, first_name, last_name, email, employee_number, 
                         position, department, status, hire_date, created_at, updated_at
                  FROM employee.employees`
	var conditions []string
	var args []interface{}
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
	if pos, ok := filters["position"]; ok && pos != "" {
		conditions = append(conditions, fmt.Sprintf("position = $%d", argPos))
		args = append(args, pos)
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
		// Check for specific PostgreSQL errors
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "42P01": // undefined_table
				return nil, fmt.Errorf("employees table does not exist: %w", err)
			case "42501": // insufficient_privilege
				return nil, fmt.Errorf("insufficient privileges to access employees: %w", err)
			}
		}
		return nil, fmt.Errorf("failed to query employees: %w", err)
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
			return nil, fmt.Errorf("failed to scan employee row: %w", err)
		}
		employees = append(employees, emp)
	}

	// Check for any iteration errors
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating employee rows: %w", err)
	}

	// Returning empty slice (not nil) for no employees is intentional
	// This makes it easier for callers (no nil check needed)
	return employees, nil
}

func (r *employeeRepository) Count(ctx context.Context, filters map[string]interface{}) (int, error) {
	baseQuery := `SELECT COUNT(*) FROM employee.employees`
	var conditions []string
	var args []interface{}
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
        SET first_name = $2, last_name = $3, email = $4, 
            employee_number = $5, position = $6, department = $7,
            status = $8, updated_at = CURRENT_TIMESTAMP
        WHERE id = $1
        RETURNING updated_at
    `

	result, err := r.db.Exec(ctx, query,
		e.ID,
		e.FirstName,
		e.LastName,
		e.Email,
		e.EmployeeNumber,
		e.Position,
		e.Department,
		e.Status,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch {
			case pgErr.Code == "23505" && pgErr.ConstraintName == "employees_email_key":
				return ErrEmailAlreadyExists
			case pgErr.Code == "23505" && pgErr.ConstraintName == "employees_employee_number_key":
				return ErrEmployeeNumberAlreadyExists
			}
		}
		return fmt.Errorf("failed to update employee: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrEmployeeNotFound
	}

	// Get updated_at if needed
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
