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
		recovery_token VARCHAR(255),
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);
	`

	if _, err := db.ExecContext(context.Background(), tableQuery); err != nil {
		return err
	}

	// Migración manual: Asegurar que la columna recovery_token existe (por si la tabla ya existía)
	addColumnQuery := `ALTER TABLE users ADD COLUMN IF NOT EXISTS recovery_token VARCHAR(255);`
	db.ExecContext(context.Background(), addColumnQuery)

	log.Println("Schema ensured for auth-service")

	// Seed admin user if not exists
	adminEmail := "admin@onboarding.com"
	seedQuery := `
	INSERT INTO users (email, password_hash, role, status)
	SELECT $1, $2, 'ADMIN', 'ACTIVE'
	WHERE NOT EXISTS (SELECT 1 FROM users WHERE email = $1);
	`
	hash := "$2a$10$ByI76z5yV9PTY.TMmSpmS.XInpw.SdqZxqCSAnW.ZfqEqWbiSyS.y"

	res, err := db.ExecContext(context.Background(), seedQuery, adminEmail, hash)
	if err != nil {
		log.Printf("CRITICAL: could not seed admin user: %v", err)
	} else {
		rows, _ := res.RowsAffected()
		if rows > 0 {
			log.Printf("SUCCESS: Admin seed user created: %s", adminEmail)
		} else {
			log.Printf("INFO: Admin seed user already exists: %s", adminEmail)
		}
	}

	return nil
}
