// Package db provides database connection management
package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPostgresPool creates and return a new Postgresql connection pool
// It validates the connection by pinging the hb and will terminate the
// app if connection or ping fails
func NewPostgresPool(dbURL string) *pgxpool.Pool {
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("failed to create db pool: %v", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := ensureSchemaAndTable(context.Background(), pool); err != nil {
		log.Fatalf("database initialization failed: %v", err)
	}

	return pool
}

// ensureSchemaAndTable validates if the schema and table exists
// If not, creates the schema an table
func ensureSchemaAndTable(ctx context.Context, db *pgxpool.Pool) error {
	schemaQuery := `
	CREATE SCHEMA IF NOT EXISTS employee;
	`

	tableQuery := `
	CREATE TABLE IF NOT EXISTS employee.employees (
		id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
		first_name VARCHAR(255) NOT NULL,
		last_name VARCHAR(255) NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		employee_number VARCHAR(50) UNIQUE NOT NULL,
		position VARCHAR(255) NOT NULL,
		department VARCHAR(255) NOT NULL,
		status VARCHAR(20) NOT NULL,
		hire_date TIMESTAMP NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	if _, err := db.Exec(ctx, schemaQuery); err != nil {
		return err
	}

	if _, err := db.Exec(ctx, tableQuery); err != nil {
		return err
	}

	return nil
}
