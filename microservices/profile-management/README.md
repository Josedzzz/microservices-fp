# Profile Management Microservice

The **Profile Management** service is a TypeScript-based microservice designed to combine asynchronous and synchronous communication. It automatically creates user profiles by consuming events and provides a REST API to query and update them.

## Features

- **Asynchronous Event Consumption**: Listens for `employee.created` events via RabbitMQ to automatically scaffold new profiles.
- **RESTful API**: Exposes endpoints to retrieve paginated profiles, find profiles by employee ID, and update profile information.
- **Data Persistence**: Uses PostgreSQL with TypeORM for structured data management.
- **Validation**: Implements strict input validation using Zod.
- **Containerized**: Fully Dockerized for seamless deployment.

## Tech Stack

- **Language**: TypeScript
- **Framework**: Express.js
- **ORM**: TypeORM
- **Database**: PostgreSQL
- **Messaging**: RabbitMQ (amqplib)
- **Validation**: Zod

## API Endpoints

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/profiles` | Returns a paginated list of all profiles. |
| `GET` | `/profiles/{employeeId}` | Retrieves a specific profile by its associated Employee ID. |
| `PUT` | `/profiles/{employeeId}` | Fully or partially updates a profile's information (Name, Email, Phone, Address, City, Bio). |
| `GET` | `/health` | Service health check. |

### Expected Responses

- **200 OK**: Successful retrieval or update.
- **400 Bad Request**: Invalid input data (e.g., malformed email).
- **404 Not Found**: The requested profile does not exist.
- **500 Internal Server Error**: Unexpected server-side issues.

## Event Consumption

The service consumes the following event from the `employees.events` exchange:

- **Routing Key**: `employee.created`
- **Action**: Creates a new profile record with the `employeeId`, `name`, and `email` provided in the event payload. Default empty values are assigned to the bio, phone, and address fields.

## Configuration

Configuration is managed via environment variables:

| Variable | Default | Description |
| :--- | :--- | :--- |
| `PORT` | `8085` | The port the service runs on. |
| `DB_HOST` | `localhost` | PostgreSQL host. |
| `DB_PORT` | `5432` | PostgreSQL port. |
| `DB_USER` | `postgres` | PostgreSQL username. |
| `DB_PASSWORD` | `postgres` | PostgreSQL password. |
| `DB_NAME` | `profiles_db` | PostgreSQL database name. |
| `RABBITMQ_HOST` | `localhost` | RabbitMQ broker host. |
| `RABBITMQ_PORT` | `5672` | RabbitMQ broker port. |
| `RABBITMQ_USER` | `guest` | RabbitMQ username. |
| `RABBITMQ_PASSWORD`| `guest` | RabbitMQ password. |

## Development and Deployment

### Running with Docker

This service is integrated into the root `docker-compose.yml`. To start the entire ecosystem, including the profile database:

```bash
docker-compose up --build profiles-service
```

### Local Development

1. Install dependencies:
   ```bash
   npm install
   ```
2. Configure `.env` file with local credentials.
3. Run in development mode:
   ```bash
   npm run dev
   ```
4. Build for production:
   ```bash
   npm run build
   ```
