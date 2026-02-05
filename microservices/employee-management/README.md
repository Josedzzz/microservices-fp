# Employee Management Service

This service is responsible for managing employee data.

## Responsibilities

- Create employees
- Update employee information
- Retrieve employees
- Delete employees

## Tech Stack

- Go
- Gin
- PostgreSQL
- pgx

## Database Schema

- Schema: employee
- Table: employees

## Environment Variables

| Variable     | Description                  |
| ------------ | ---------------------------- |
| DATABASE_URL | PostgreSQL connection string |
| SERVER_PORT  | HTTP server port             |

## API Documentation

Swagger available at:
http://localhost:8081/swagger/index.html

## Run locally

go run cmd/main.go
