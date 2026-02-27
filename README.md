# Employee Onboarding & Offboarding System

A microservices-based system for managing employee onboarding and offboarding processes. This project contains two independent microservices that communicate via REST APIs.

## Services Overview

### Employees Service (Go) - Port 8081

A Go-based service for managing employee records. Each employee is associated with a department.

### Departments Service (Python/FastAPI) - Port 8082

A FastAPI-based service for managing departments. Used by the employees service to validate department existence.

## Architecture

The system consists of the following components:

- Employees Service (Go) - Manages employee data on port 8081
- Departments Service (Python/FastAPI) - Manages department data on port 8082
- PostgreSQL databases (one for each service)
- Docker Compose for orchestration

## Prerequisites

- Docker and Docker Compose
- Go 1.24+ (for local development of employees service)
- Python 3.11+ (for local development of departments service)
- PostgreSQL 15+ (for local development)

## Quick Start with Docker Compose

Clone the repository:

```bash
git clone https://github.com/Josedzzz/microservices-fp.git
cd microservices-fp
```

## Start all services

```bash
docker-compose up --build
```

This command will build both microservices, create two PostgreSQL databases, set up a network for service communication, and start all containers.

Verify the services running:

```bash
docker ps
```

You should see four containers running:

- employees-service
- departments-service
- database-employees
- database-departments

## Services API documentation

- Employees Service Swagger UI: http://localhost:8081/swagger/index.html
- Departments Service Swagger UI: http://localhost:8082/docs

## Stop the services

```bash
docker-compose down
```

To also remove the volumes (database data):

```bash
docker-compose down -v
```

## Database Connections

You can connect to the databases using any PostgreSQL client (like Beekeeper Studio).

Employees Database:

```text
Host: localhost
Port: 5432
User: postgres
Password: postgres
Database: employees_db
Schema: employee
Table: employees
```

Departments Database:

```text
Host: localhost
Port: 5433
User: postgres
Password: postgres
Database: departments_db
Schema: public
Table: departments
```

## Tasks

### Branch rabbit

Implement rabbitmq with the employee service

### Branch notifications

Implement the base of the notification service

## Branch profiles

Implement the base of the profile service
