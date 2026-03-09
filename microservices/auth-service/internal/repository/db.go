package repository

import (
	"context"
	"database/sql"
	"log"
)

// EnsureSchema validates if the table exists
// If not, creates the table
func EnsureSchema(db *sql.DB) error {
	tableQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
		password_hash VARCHAR(255),
		role VARCHAR(50) NOT NULL,
		status VARCHAR(20) NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);
	`

	if _, err := db.ExecContext(context.Background(), tableQuery); err != nil {
		return err
	}

	log.Println("Schema ensured for auth-service")
	return nil
}
