// Package messaging provides functionality for handling messaging events
package messaging

import (
	"time"

	"employee-management/internal/models"
)

const (
	EmployeesExchange = "employees.events"

	EventEmployeeCreated = "employee.created"
	EventEmployeeDeleted = "employee.deleted"
)

// EmployeeCreatedEvent represents the event fired when an employee is created
type EmployeeCreatedEvent struct {
	ID           int64                 `json:"id"`
	Name         string                `json:"name"`
	Email        string                `json:"email"`
	Role         string                `json:"role"`
	DepartmentID string                `json:"departmentId"`
	Status       models.EmployeeStatus `json:"status"`
	HireDate     time.Time             `json:"hireDate"`
	CreatedAt    time.Time             `json:"createdAt"`
}

// EmployeeDeletedEvent represents the event fired when an employee is deleted
type EmployeeDeletedEvent struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
