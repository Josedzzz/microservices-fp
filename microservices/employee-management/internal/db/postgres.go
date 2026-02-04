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
		log.Fatalf("unable to connect to database: %v", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("database ping failed: %v", err)
	}

	log.Println("PostgreSQL connected")
	return pool
}
