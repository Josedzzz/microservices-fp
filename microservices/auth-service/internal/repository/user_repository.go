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
	SELECT id,email,password_hash,role,status,recovery_token
	FROM users
	WHERE email=$1
	`

	row := r.db.QueryRowContext(ctx, query, email)

	user := &model.User{}

	var recoveryToken sql.NullString
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.Status,
		&recoveryToken,
	)
	if err != nil {
		return nil, err
	}

	if recoveryToken.Valid {
		user.RecoveryToken = recoveryToken.String
	}

	return user, nil
}

// UpdateRecoveryToken stores a token for a user by their email
func (r *UserRepository) UpdateRecoveryToken(ctx context.Context, email, token string) error {
	query := `
	UPDATE users
	SET recovery_token=$1
	WHERE email=$2
	`

	_, err := r.db.ExecContext(ctx, query, token, email)

	return err
}

// UpdatePasswordWithToken resets the password hash and clears the recovery token
func (r *UserRepository) UpdatePasswordWithToken(ctx context.Context, email, token, hash string) error {
	query := `
	UPDATE users
	SET password_hash=$1,status='ACTIVE',recovery_token=NULL
	WHERE email=$2 AND recovery_token=$3
	`

	res, err := r.db.ExecContext(ctx, query, hash, email, token)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
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

// DisableUser sets the status of a user identified by their email to DISABLED and clears recovery tokens
func (r *UserRepository) DisableUser(ctx context.Context, email string) error {
	query := `
	UPDATE users
	SET status='DISABLED', recovery_token=NULL
	WHERE email=$1
	`

	_, err := r.db.ExecContext(ctx, query, email)

	return err
}
