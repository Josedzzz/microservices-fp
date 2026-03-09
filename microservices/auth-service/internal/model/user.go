// Package model contains the data structs and models
package model

import "time"

// User represents a user in the system
type User struct {
	ID           int64
	Email        string
	PasswordHash string
	Role         string
	Status       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
