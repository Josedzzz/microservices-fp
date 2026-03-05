# Employee Onboarding & Offboarding System

A microservices-based system for managing employee onboarding and offboarding processes. This project contains four independent microservices that communicate via REST APIs and asynchronous messaging.

## Services Overview

### Employees Service (Go) - Port 8081
A Go-based service for managing employee records. Each employee is associated with a department.

### Departments Service (Python/FastAPI) - Port 8082
A FastAPI-based service for managing departments. Used by the employees service to validate department existence.

### Notifications Service (Java/Spring Boot) - Port 8084
A Spring Boot-based service that listens for employee events (creation/deletion) via RabbitMQ and persists a history of notifications sent to employees.

### Profile Management Service (TypeScript/Express) - Port 8085
A TypeScript-based service that combines asynchronous and synchronous communication. It consumes `employee.created` events to automatically create profiles and exposes REST endpoints for querying and updating them.

## Architecture

The system consists of the following components:

- **Employees Service (Go)** - Manages employee data on port 8081
- **Departments Service (Python/FastAPI)** - Manages department data on port 8082
- **Notifications Service (Java/Spring Boot)** - Manages notification history on port 8084
- **Profile Management Service (TypeScript)** - Manages user profiles on port 8085
- **RabbitMQ** - Message broker for event-driven communication
- **PostgreSQL databases** (one for each service)
- **Docker Compose** for orchestration

## Prerequisites

- Docker and Docker Compose
- Go 1.24+ (for local development of employees service)
- Python 3.11+ (for local development of departments service)
- Java 17+ (for local development of notifications service)
- Node.js 20+ & TypeScript (for local development of profiles service)
- PostgreSQL 15+ (for local development)

## Quick Start with Docker Compose

Clone the repository:

```bash
git clone https://github.com/Josedzzz/microservices-fp.git
cd microservices-fp
```

### Start all services

```bash
docker-compose up --build
```

This command will build the microservices, create the PostgreSQL databases, set up a network for service communication, and start all containers.

Verify the services running:

```bash
docker ps
```

You should see **nine** containers running:

- employees-service
- departments-service
- notifications-service
- profiles-service
- rabbitmq
- database-employees
- database-departments
- database-notifications
- database-profiles

## Services API documentation

- **Employees Service Swagger UI**: http://localhost:8081/swagger/index.html
- **Departments Service Swagger UI**: http://localhost:8082/docs
- **Notifications Service Swagger UI**: http://localhost:8084/swagger-ui/index.html
- **Profiles Service API**: http://localhost:8085/profiles (GET /profiles, GET /profiles/{employeeId}, PUT /profiles/{employeeId})

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

**Employees Database:**
- Port: 5432 | DB: `employees_db`

**Departments Database:**
- Port: 5433 | DB: `departments_db`

**Notifications Database:**
- Port: 5434 | DB: `notifications_db`

**Profiles Database:**
- Port: 5435 | DB: `profiles_db`
- User: `postgres` | Password: `postgres`

## Message Broker - RabbitMQ

The system uses RabbitMQ as the message broker to enable asynchronous communication between microservices.

### RabbitMQ Management UI
- URL: http://localhost:15672
- Username: `guest`
- Password: `guest`

### Why RabbitMQ?
RabbitMQ was selected for its mature ecosystem, excellent support for event-based communication, and ease of integration with Go, Python, TypeScript, and Java.

## Tasks

### Branch rabbit
Implement rabbitmq with the employee service.

### Branch notifications
Email sending and paginated response of all notifications.

### Branch profiles
Implement the Profile Management service (TypeScript/Express) with hybrid sync/async communication.
