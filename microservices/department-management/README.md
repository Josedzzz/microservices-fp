# Departments Service

This service is responsible for managing departments. It is built with FastAPI and PostgreSQL.

## Responsibilities

- Create departments
- Retrieve department by ID
- List all departments
- Provide department validation for Employee Service

## Tech Stack

- Python 3.11
- FastAPI
- SQLAlchemy
- Pydantic
- PostgreSQL 15
- Uvicorn
- Swagger/OpenAPI for documentation

## Database Schema

Schema: public
Table: departments

Columns:

- id (varchar, primary key, UUID)
- name (varchar, not null)
- description (text, nullable)

## Environment Variables

| Variable     | Description       | Default        | Required |
| ------------ | ----------------- | -------------- | -------- |
| DB_USER      | Database username | postgres       | Yes      |
| DB_PASSWORD  | Database password | postgres       | Yes      |
| DB_HOST      | Database host     | localhost      | Yes      |
| DB_PORT      | Database port     | 5432           | Yes      |
| DB_NAME      | Database name     | departments_db | Yes      |
| SERVICE_PORT | HTTP server port  | 8082           | No       |

## Local Development

### Prerequisites

- Python 3.11 or higher
- PostgreSQL 15 or higher
- pip (Python package manager)
- Git

### Setup Steps

1. Create PostgreSQL database

```bash
CREATE DATABASE departments_db;
```

2. Create and active virtual environment

```bash
python -m venv venv
venv\Scripts\activate
```

3. Install dependencies

```bash
pip install -r requirements.txt
```

4. Create .env file

```bash
# Database configuration for local development
DB_USER=postgres
DB_PASSWORD=postgres
DB_HOST=localhost
DB_PORT=5432
DB_NAME=departments_db

# Service configuration
SERVICE_PORT=8082
```

5. Run the service

```bash
uvicorn app.main:app --reload --port 8082
```

or use:

```bash
python -m uvicorn app.main:app --reload --port 8082
```
