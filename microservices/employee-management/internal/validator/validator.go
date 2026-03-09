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

type ValidationErrorCode string

const (
	CodeRequired      ValidationErrorCode = "REQUIRED"
	CodeInvalidFormat ValidationErrorCode = "INVALID_FORMAT"
	CodeInvalidValue  ValidationErrorCode = "INVALID_VALUE"
)

type ValidationError struct {
	Field   string
	Code    ValidationErrorCode
	Message string
	Value   any
}

// ToAPIErrorDetail converts ValidationError to api.ErrorDetail
func (ve ValidationError) ToAPIErrorDetail() api.ErrorDetail {
	return api.ErrorDetail{
		Field:         ve.Field,
		Message:       ve.Message,
		RejectedValue: ve.Value,
	}
}

type ValidationResult struct {
	IsValid bool
	Errors  []ValidationError
}

// ValidateEmployee validates employee data
func ValidateEmployee(email, name string) ValidationResult {
	result := ValidationResult{IsValid: true, Errors: []ValidationError{}}

	// Validate email
	if email == "" {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "email",
			Code:    CodeRequired,
			Message: "email is required",
		})
		result.IsValid = false
	} else if !IsValidEmail(email) {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "email",
			Code:    CodeInvalidFormat,
			Message: "email format is invalid",
			Value:   email,
		})
		result.IsValid = false
	}

	// Validate name
	if strings.TrimSpace(name) == "" {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "name",
			Code:    CodeRequired,
			Message: "name is required",
		})
		result.IsValid = false
	}

	return result
}

// ToAPIErrorDetails converts validation errors to API error details (for handlers)
func (vr ValidationResult) ToAPIErrorDetails() []api.ErrorDetail {
	details := make([]api.ErrorDetail, len(vr.Errors))
	for i, err := range vr.Errors {
		details[i] = err.ToAPIErrorDetail()
	}
	return details
}

// IsValidEmail validates the format of a email
func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil && emailRegex.MatchString(email)
}

// ValidateID validates an id
func ValidateID(idStr string) (int64, *ValidationResult) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, &ValidationResult{
			IsValid: false,
			Errors: []ValidationError{
				{
					Field:   "id",
					Code:    CodeInvalidFormat,
					Message: "ID must be a valid integer",
					Value:   idStr,
				},
			},
		}
	}

	if id <= 0 {
		return 0, &ValidationResult{
			IsValid: false,
			Errors: []ValidationError{
				{
					Field:   "id",
					Code:    CodeInvalidValue,
					Message: "ID must be a positive number",
				},
			},
		}
	}

	return id, &ValidationResult{IsValid: true}
}

// ValidateField validates a single field
func ValidateField(field, value string) ValidationResult {
	result := ValidationResult{IsValid: true, Errors: []ValidationError{}}

	switch field {
	case "email":
		if value == "" {
			result.Errors = append(result.Errors, ValidationError{
				Field:   "email",
				Code:    CodeRequired,
				Message: "email is required",
			})
			result.IsValid = false
		} else if !IsValidEmail(value) {
			result.Errors = append(result.Errors, ValidationError{
				Field:   "email",
				Code:    CodeInvalidFormat,
				Message: "email format is invalid",
				Value:   value,
			})
			result.IsValid = false
		}

	case "name":
		if strings.TrimSpace(value) == "" {
			result.Errors = append(result.Errors, ValidationError{
				Field:   "name",
				Code:    CodeRequired,
				Message: "name is required",
			})
			result.IsValid = false
		}
	}

	return result
}
