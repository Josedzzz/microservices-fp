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

You should see **twelve** containers running:

- **api-gateway**: Central entry point (Port 8000)
- **auth-service**: Identity provider
- **employees-service**: Core employee logic
- **departments-service**: Department management
- **notifications-service**: Event-driven notifications
- **profiles-service**: User profiles
- **rabbitmq**: Message broker
- **database-auth**: Auth data
- **database-employees**: Employee data
- **database-departments**: Department data
- **database-notifications**: Notification data
- **database-profiles**: Profile data

## Services API documentation

All documentation is accessible both directly (for development) and through the API Gateway:

### Directly (Localhost)
- **Auth Service**: http://localhost:8083/swagger/index.html
- **Employees Service**: http://localhost:8081/swagger/index.html
- **Departments Service**: http://localhost:8082/docs
- **Notifications Service**: http://localhost:8084/swagger-ui/index.html
- **Profiles Service**: http://localhost:8085/swagger

### Through API Gateway (Port 8000)
- **Employees Docs**: http://localhost:8000/employees-service/swagger/index.html
- **Departments Docs**: http://localhost:8000/departments-service/docs
- **Notifications Docs**: http://localhost:8000/notifications-service/swagger-ui/index.html
- **Profiles Docs**: http://localhost:8000/profiles-service/swagger

---

## Token Instructions & Usage Examples

### 1. Login (Public)
The Gateway proxies the `/auth-service` directly.
```bash
curl -X POST http://localhost:8000/auth-service/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@onboarding.com",
    "password": "admin123"
  }'
```
*Expected Response: JSON with `access_token`.*

### 2. Protected Endpoint (Read-only for USER, Total for ADMIN)
```bash
# Get all employees
curl -X GET http://localhost:8000/employees-service/api/employees \
  -H "Authorization: Bearer <YOUR_JWT_TOKEN>"
```

### 3. Admin-Only Endpoint (Write operations)
```bash
# Create a new department (Requires ADMIN role in JWT)
curl -X POST http://localhost:8000/departments-service/api/departments \
  -H "Authorization: Bearer <ADMIN_JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Engineering",
    "description": "Technical team"
  }'
```
*Note: If a USER attempts this, the Gateway will return `403 Forbidden`.*

### 4. Password Recovery (Public)
```bash
curl -X POST http://localhost:8000/auth-service/api/recover-password \
  -H "Content-Type: application/json" \
  -d '{"email": "juan@empresa.com"}'
```

---

#### Security Flow:
1. **Bootstrap Admin:** A seed admin is created automatically (`admin@onboarding.com` / `admin123`).
2. **Gateway Entrance:** All traffic enters through `http://localhost:8000`. The Gateway validates the JWT signature and checks the expiration.
3. **Onboarding:** Create an employee via `POST http://localhost:8000/employees-service/api/employees`.
4. **Activation:** Use the token from the logs to set a password via `POST http://localhost:8000/auth-service/api/reset-password`.
5. **RBAC:** The Gateway enforces: `ADMIN` role required for POST/PUT/DELETE. `USER` role has read-only (GET) access.
