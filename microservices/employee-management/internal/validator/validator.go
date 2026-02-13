// Package validator define structure validations
package validator

import (
	"net/mail"
	"regexp"
	"strconv"

	"employee-management/internal/api"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

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
