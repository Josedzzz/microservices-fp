# Employee Onboarding & Offboarding System

A microservices-based system for managing employee onboarding and offboarding processes. This project contains three independent microservices that communicate via REST APIs and asynchronous messaging.

## Services Overview

### Employees Service (Go) - Port 8081

A Go-based service for managing employee records. Each employee is associated with a department.

### Departments Service (Python/FastAPI) - Port 8082

A FastAPI-based service for managing departments. Used by the employees service to validate department existence.

### Notifications Service (Java/Spring Boot) - Port 8084

A Spring Boot-based service that listens for employee events (creation/deletion) via RabbitMQ and persists a history of notifications sent to employees.

## Architecture

The system consists of the following components:

- Employees Service (Go) - Manages employee data on port 8081
- Departments Service (Python/FastAPI) - Manages department data on port 8082
- Notifications Service (Java/Spring Boot) - Manages notification history on port 8084
- RabbitMQ - Message broker for event-driven communication
- PostgreSQL databases (one for each service)
- Docker Compose for orchestration

## Prerequisites

- Docker and Docker Compose
- Go 1.24+ (for local development of employees service)
- Python 3.11+ (for local development of departments service)
- Java 17+ (for local development of notifications service)
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

This command will build the microservices, create the PostgreSQL databases, set up a network for service communication, and start all containers.

Verify the services running:

```bash
docker ps
```

You should see seven containers running:

- employees-service
- departments-service
- notifications-service
- rabbitmq
- database-employees
- database-departments
- database-notifications

## Services API documentation

- Employees Service Swagger UI: http://localhost:8081/swagger/index.html
- Departments Service Swagger UI: http://localhost:8082/docs
- Notifications Service Swagger UI: http://localhost:8084/swagger-ui/index.html

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

Notifications Database:

```text
Host: localhost
Port: 5434
User: postgres
Password: postgres
Database: notifications_db
Schema: public
Table: notifications
```

## Message Broker - RabbitMQ

The system uses RabbitMQ as the message broker to enable asynchronous communication between microservices.

### RabbitMQ Management UI

- URL: http://localhost:15672
- Username: guest
- Password: guest

### Why RabbitMQ?

RabbitMQ was selected as the message broker for this project based on the scope, requirements, and educational goals of the system.

Key Reasons:

- Simple to set up with Docker
- Excellent support for event-based communication
- Mature ecosystem and documentation
- Easy integration with Go, Python, TypeScript, and Java
- Includes a web-based management UI, which is ideal for learning and debugging

### Why not kafka?

Kafka is better suited for high-throughput data pipelines and event streaming.
For this project, Kafka would be overkill and introduce unnecessary operational complexity.

### Why not Redis Streams?

Redis Streams are powerful but are not primarily designed as a message broker.
RabbitMQ provides clearer semantics for messaging and event routing.

### Why not NATS?

NATS excels in low-latency, high-performance systems, but RabbitMQ is more suitable for learning event-driven architectures with visibility and control.

## Tasks

### Branch rabbit

Implement rabbitmq with the employee service

### Branch notifications

Email sending.

### Branch profiles

Implement the base of the profile service (Ts)
