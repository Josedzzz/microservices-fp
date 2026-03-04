# Notification Service

Microservice responsible for processing employee-related events and persisting notifications. Built with Java 17, Spring Boot, and PostgreSQL.

## Features

- **Event-Driven Architecture**: Consumes messages from RabbitMQ via the `employees.events` exchange.
- **Dynamic Event Handling**: Automatically distinguishes between `employee.created` and `employee.deleted` events using:
  - Message routing keys (headers).
  - Smart payload inspection (fallback mechanism) to ensure correct processing even if headers are missing.
- **Notification History**: Stores a permanent record of all notifications sent, categorized by type (`WELCOME`, `TERMINATION`).
- **OpenAPI Documentation**: Fully documented REST API via Swagger UI.

## Tech Stack

- **Framework**: Spring Boot 3.x
- **Database**: PostgreSQL (with Flyway for migrations)
- **Messaging**: RabbitMQ (Spring AMQP)
- **Documentation**: Springdoc-openapi (Swagger)
- **Utilities**: Lombok, MapStruct (if used)

## API Endpoints

- `GET /notifications`: Retrieve a full history of all notifications.
- `GET /notifications/{employeeId}`: Retrieve all notifications for a specific employee.
- `GET /swagger-ui/index.html`: Interactive API documentation and testing tool at `http://localhost:8084/swagger-ui/index.html`.

## Event Logic

The service listens for routing keys:
- `employee.created`: Triggers a **Welcome** notification.
- `employee.deleted`: Triggers a **Termination** (deactivation) notification.

If the routing key is not provided in the message headers, the `NotificationService` inspects the payload fields (like `createdAt`) to determine the event type accurately.
