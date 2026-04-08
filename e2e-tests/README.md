# E2E BDD Testing Suite - Microservices Ecosystem

This project contains the automated functional testing suite for the microservices ecosystem, developed using the **BDD (Behavior-Driven Development)** methodology as required for **Challenge 5**.

## BDD Methodology
We have selected BDD to align technical requirements with expected business behavior. Scenarios are written in **Gherkin** (natural language) in English, allowing the documentation to be "living" and executable.

## Tech Stack
- **Node.js**: Runtime environment.
- **Cucumber.js**: BDD framework for executing Gherkin scenarios.
- **Axios**: HTTP client for performing real requests against the API Gateway.
- **Chai**: Assertion library for validating results.

## Project Structure
- `features/`: Contains `.feature` files with test scenarios.
- `features/step_definitions/`: Code implementation for each Gherkin step.
- `features/step_definitions/world.js`: Shared context object to handle tokens and responses.
- `features/step_definitions/polling.js`: Retry mechanism to handle asynchronous eventual consistency from RabbitMQ.

## Covered Scenarios
1. **Smoke Tests**: Basic verification of Gateway connectivity.
2. **Security & RBAC**: Strict validation of access rules (401 Unauthorized and 403 Forbidden for the USER role).
3. **Onboarding Flow**: Employee registration, asynchronous event waiting, and initial access.
4. **Offboarding Flow**: Employee removal and automatic access revocation.

## How to Run the Tests

### Prerequisites
- Have the complete system running via Docker Compose (`docker-compose up -d`).
- Node.js v18+ installed locally.

### Installation
```bash
cd e2e-tests
npm install
```

### Execution
```bash
# Run all tests
npm test

# Run a specific feature
npx cucumber-js features/onboarding.feature
```

## Polling Configuration
For asynchronous events (RabbitMQ), we have configured a polling system with a maximum of **15 attempts** every **2 seconds** (Total 30s). This time is sufficient for the Auth and Notifications services to process messages in a development environment.
