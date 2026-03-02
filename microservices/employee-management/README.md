# Employee Management Service

This service is responsible for managing employee data. It is built with Go and PostgreSQL.

## Responsibilities

- Create employees
- Update employee information
- Retrieve employees by ID
- List all employees
- Delete employees
- Validate department existence with Departments Service

## Tech Stack

- Go 1.24
- Gin Web Framework
- PostgreSQL 15
- pgx database driver
- Swagger for API documentation
- RabbitMQ (event broker)

## Database Schema

- Schema: employee
- Table: employees

Columns:

- id (integer, auto-generated)
- first_name (varchar)
- last_name (varchar)
- email (varchar, unique)
- employee_number (varchar, unique)
- position (varchar)
- department (varchar)
- status (varchar)
- hire_date (timestamp)
- created_at (timestamp)
- updated_at (timestamp)

## Environment Variables

| Variable                | Description                 | Default               | Required |
| ----------------------- | --------------------------- | --------------------- | -------- |
| SERVER_PORT             | HTTP server port            | 8081                  | No       |
| DB_HOST                 | Database host               | localhost             | Yes      |
| DB_PORT                 | Database port               | 5432                  | Yes      |
| DB_NAME                 | Database name               | employees_db          | Yes      |
| DB_USER                 | Database username           | postgres              | Yes      |
| DB_PASSWORD             | Database password           | postgres              | Yes      |
| DB_SSLMODE              | SSL mode for database       | disable               | No       |
| DEPARTMENTS_SERVICE_URL | URL for departments service | http://localhost:8082 | Yes      |
| RABBITMQ_HOST           | Rabbit host                 | rabbitmq              | Yes      |
| RABBITMQ_PORT           | Rabbit port                 | 5672                  | Yes      |
| RABBITMQ_USER           | Rabbit user                 | guest                 | Yes      |
| RABBITMQ_PASSWORD       | Rabbit password             | guest                 | Yes      |
| RABBITMQ_VHOST          | Rabbit vhost                | /                     | Yes      |

## Local Development

### Prerequisites

- Go 1.24 or higher
- PostgreSQL 15 or higher
- Git

### Setup Steps

1. Create PostgreSQL database

```bash
CREATE DATABASE employees_db;
```

2. Create .env file

```bash
# Server configuration
SERVER_PORT=8081

# Database configuration for local development
DB_HOST=localhost
DB_PORT=5432
DB_NAME=employees_db
DB_USER=postgres
DB_PASSWORD=postgres
DB_SSLMODE=disable

# Departments service URL
DEPARTMENTS_SERVICE_URL=http://localhost:8082

# RabbitMQ
RABBITMQ_HOST=rabbitmq
RABBITMQ_PORT=5672
RABBITMQ_USER=guest
RABBITMQ_PASSWORD=guest
RABBITMQ_VHOST=/
```

3. Install dependencies

```bash
go mod download
```

4. Run the service

```bash
go run cmd/main.go
```
