// Package repository provides the data access layer for interacting with the database
package repository

import (
	"context"
	"database/sql"

	"auth-service/internal/model"
)

// UserRepository is responsible for interacting with the users table
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new instance of UserRepository with the db connection
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// FindByEmail gets a user from the db by their email address
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
	SELECT id,email,password_hash,role,status
	FROM users
	WHERE email=$1
	`

	row := r.db.QueryRowContext(ctx, query, email)

	user := &model.User{}

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.Status,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// CreateUser inserts a new user into the db with the given email, hash, role, and status
func (r *UserRepository) CreateUser(ctx context.Context, email, role, status string) error {
	query := `
	INSERT INTO users(email,role,status)
	VALUES($1,$2,$3)
	`

	_, err := r.db.ExecContext(ctx, query, email, role, status)

	return err
}

// UpdatePassword updates the password hash of a user by their email
func (r *UserRepository) UpdatePassword(ctx context.Context, email, hash string) error {
	query := `
	UPDATE users
	SET password_hash=$1,status='ACTIVE'
	WHERE email=$2
	`

	_, err := r.db.ExecContext(ctx, query, hash, email)

	return err
}

// DisableUser sets the status of a user identified by their email to DISABLED
func (r *UserRepository) DisableUser(ctx context.Context, email string) error {
	query := `
	UPDATE users
	SET status='DISABLED'
	WHERE email=$1
	`

	_, err := r.db.ExecContext(ctx, query, email)

	return err
}
