// Package validator define structure validations
package validator

import (
	"net/mail"
	"regexp"
	"strconv"
	"strings"

	"employee-management/internal/api"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// ValidationResult contains the result of a validation
type ValidationResult struct {
	IsValid bool
	Errors  []api.ErrorDetail
}

// ValidateEmployee validates employee data
func ValidateEmployee(email, employeeNumber, firstName, lastName string) ValidationResult {
	result := ValidationResult{IsValid: true, Errors: []api.ErrorDetail{}}

	// Validate email
	if email == "" {
		result.Errors = append(result.Errors, api.ErrorDetail{
			Field:   "email",
			Message: "Email is required",
		})
		result.IsValid = false
	} else if !IsValidEmail(email) {
		result.Errors = append(result.Errors, api.ErrorDetail{
			Field:         "email",
			Message:       "Email format is invalid",
			RejectedValue: email,
		})
		result.IsValid = false
	}

	// Validate employee number
	if employeeNumber == "" {
		result.Errors = append(result.Errors, api.ErrorDetail{
			Field:   "employeeNumber",
			Message: "Employee number is required",
		})
		result.IsValid = false
	}

	// Validate name
	if strings.TrimSpace(firstName) == "" {
		result.Errors = append(result.Errors, api.ErrorDetail{
			Field:   "firstName",
			Message: "First name is required",
		})
		result.IsValid = false
	}

	// Validate last name
	if strings.TrimSpace(lastName) == "" {
		result.Errors = append(result.Errors, api.ErrorDetail{
			Field:   "lastName",
			Message: "Last name is required",
		})
		result.IsValid = false
	}

	return result
}

// IsValidEmail validates the format of a email
func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil && emailRegex.MatchString(email)
}

// ValidateID validates an id
func ValidateID(idStr string) (int64, []api.ErrorDetail) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, []api.ErrorDetail{{
			Field:         "id",
			Message:       "ID must be a valid integer",
			RejectedValue: idStr,
		}}
	}

	if id <= 0 {
		return 0, []api.ErrorDetail{{
			Field:   "id",
			Message: "ID must be a positive number",
		}}
	}

	return id, nil
}
