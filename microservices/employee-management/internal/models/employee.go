// Package models define the core data structures for the employee management
package models

import "time"

// EmployeeStatus represents the current employment status
type EmployeeStatus string

const (
	StatusActive     EmployeeStatus = "ACTIVE"
	StatusOnVacation EmployeeStatus = "ON_VACATION"
	StatusRetired    EmployeeStatus = "RETIRED"
)

// Employee represents an employee record in the system
// All fields are tagged for JSON serialization
type Employee struct {
	ID           int64          `json:"id"`
	Name         string         `json:"name"`
	Email        string         `json:"email"`
	Role         string         `json:"role,omitempty"` // ADMIN or USER
	DepartmentID string         `json:"departmentID"`
	Status       EmployeeStatus `json:"status"`
	HireDate     time.Time      `json:"hireDate"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
}
