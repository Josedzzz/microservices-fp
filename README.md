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

_Expected Response: JSON with `access_token`._

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

_Note: If a USER attempts this, the Gateway will return `403 Forbidden`._

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

## Continuous Integration & Deployment

### Conceptual Overview

**Continuous Integration (CI)** is a development practice where code changes are automatically tested and verified as soon as they are committed to the repository. This system integrates **Jenkins** as the CI orchestrator, **SonarQube** for static code analysis and coverage reporting, and **Docker** for containerized builds and multi-service end-to-end testing.

#### Why CI in this Architecture?

This microservices project benefits from CI because:

1. **Automated Verification**: Each commit triggers automated build, unit, and integration tests, catching regressions early.
2. **Code Quality Consistency**: SonarQube enforces minimum coverage thresholds (≥70%) and identifies code smells, security vulnerabilities, and duplications.
3. **Fast Feedback Loops**: Developers receive pipeline results (pass/fail) within minutes, enabling quick iteration.
4. **Reliable Releases**: Multi-language pipelines (Go, Java) ensure language-specific quality gates are met before packaging.
5. **Container Readiness**: Successful builds are automatically packaged as Docker images, ready for deployment.
6. **Multi-Service Validation**: E2E tests verify that all services work together before merging.

---

### Jenkins Access Guide

#### Access URL

```
http://localhost:9090
```

When Jenkins is running via `docker-compose up`, it is accessible on port 9090.

#### Default Admin Credentials

```
Username: admin
Password: (Retrieved from logs during first startup)
```

#### Retrieving the Initial Admin Password

**Option 1: From Container Logs (Recommended)**

```bash
docker logs microservices-fp_jenkins_1 2>&1 | grep -A 5 "initialAdminPassword"
```

Look for output similar to:

```
*************************************************************
*************************************************************
*************************************************************

Jenkins initial setup is required. An admin user has been created and a password generated.
Please use the following password to proceed to installation:

<YOUR_INITIAL_ADMIN_PASSWORD>

This may also be found at: /var/jenkins_home/secrets/initialAdminPassword

*************************************************************
*************************************************************
*************************************************************
```

**Option 2: From Docker Volume**

```bash
docker exec microservices-fp_jenkins_1 cat /var/jenkins_home/secrets/initialAdminPassword
```

**Option 3: From Host Mount**

If you have direct access to the Jenkins volume:

```bash
cat /path/to/jenkins_home/secrets/initialAdminPassword
```

#### Initial Setup Steps

1. Navigate to `http://localhost:9090` in your browser.
2. Enter the initial admin password (from above).
3. Click **"Install suggested plugins"** to install the pre-configured plugin set.
4. Create your first admin user (or skip to use the temporary admin account).
5. Confirm Jenkins URL: `http://localhost:9090` (or adjust for your environment).
6. Click **"Start using Jenkins"**.

After initialization, the Jenkins Configuration as Code (JCasC) configuration will automatically:

- Configure the SonarQube server connection.
- Create pipeline jobs for `notification-service` and `auth-service`.
- Set up credentials for SonarQube token and registry access.

---

### Setup Instructions

#### 1. Start the Entire System (Including Jenkins & SonarQube)

```bash
cd /path/to/microservices-fp

docker-compose up --build
```

This command orchestrates:

- **Jenkins** (port 9090) with Docker-outside-of-Docker (DooD) via socket mount
- **SonarQube** (port 9000) with PostgreSQL backend
- **Docker Registry** (port 5000) for pushing built images
- All microservices and their databases
- **RabbitMQ** for event-driven communication

Wait for all containers to be healthy (typically 1-2 minutes).

#### 2. Access Jenkins and Configure Credentials

1. Open `http://localhost:9090` and complete the initial setup (see **Jenkins Access Guide** above).

2. Create a **SonarQube Token** credential:

   - Go to **Manage Jenkins** → **Manage Credentials** → **System** → **Global credentials (unrestricted)**.
   - Click **Add Credentials**.
   - **Kind**: Secret text
   - **Secret**: Retrieve from SonarQube:
     ```bash
     # Open http://localhost:9000 (default user: admin / password: admin)
     # User → My Account → Security → Generate token
     ```
   - **ID**: `sonar-token`
   - Click **Create**.

3. Verify **SonarQube Server Connection**:
   - Go to **Manage Jenkins** → **System Configuration** → Find section "SonarQube servers".
   - Ensure URL is `http://sonarqube:9000` and token is set.
   - Click **Test Connection** (should show success).

#### 3. Verify Pipelines Are Auto-Provisioned

Navigate to Jenkins dashboard. You should see two pipeline jobs:

- **notification-service-pipeline** (Java/Maven)
- **auth-service-pipeline** (Go)

These are auto-created by JCasC (`jenkins/casc.yaml`). If they don't appear:

```bash
# Restart Jenkins to reload JCasC config
docker restart microservices-fp_jenkins_1
```

#### 4. Trigger a Pipeline Manually

**Option 1: Via Jenkins UI**

1. Click on the pipeline job (e.g., `notification-service-pipeline`).
2. Click **Build Now** (top-left).
3. Jenkins clones the repo, checks out the specified branch, and runs the Jenkinsfile.

**Option 2: Via Git Webhook (Optional)**

Configure GitHub to send push events to Jenkins:

```
Jenkins → Manage Jenkins → System Configuration → GitHub Server
Repo URL: https://github.com/Josedzzz/microservices-fp
Webhook: http://<your-jenkins-ip>:9090/github-webhook/
```

Then any push to the repo triggers the pipeline automatically.

**Option 3: Via curl**

```bash
curl -X POST http://localhost:9090/job/notification-service-pipeline/build \
  -u admin:<password>
```

---

### Pipeline Stage Breakdown

Both pipelines (notification-service and auth-service) follow a **Declarative Pipeline** pattern with stages that ensure build quality, coverage, and readiness for deployment.

#### Stage 1: **Checkout**

- **Purpose**: Fetch the latest code from the repository.
- **What it verifies**:
  - Git repository is accessible.
  - Specified branch/commit is available.
  - Jenkinsfile is present in the service directory.
- **Failure reason**: Repository not reachable, branch doesn't exist, or network issue.

#### Stage 2: **Build**

**For notification-service (Java/Maven):**

- Runs `mvn clean compile` inside a Maven 3.9 + OpenJDK 17 Docker container.
- Compiles Java sources and resolves Maven dependencies.
- Output: Compiled `.class` files and assembled JAR dependencies.

**For auth-service (Go):**

- Runs `go build` inside a Go 1.25 Docker container.
- Downloads Go modules and compiles the binary.
- Output: `auth-server` executable in the service directory.
- **Caching**: Go module cache and build cache are stored locally in the workspace to speed up subsequent builds.

#### Stage 3: **Test**

**For notification-service:**

- Runs `mvn verify` (includes `mvn test`).
- Executes all JUnit tests and generates coverage via **JaCoCo** plugin.
- Output: JUnit XML reports in `target/surefire-reports/` and coverage in `target/site/jacoco/`.

**For auth-service:**

- Runs `go test ./... -coverprofile=coverage.out`.
- Executes all `*_test.go` files and generates Go coverage profile.
- Output: Coverage report in `.scannerwork/coverage.out`.
- **Failure reason**: Test failures, assertion errors, or uncaught panics.

#### Stage 4: **SonarQube Analysis**

- Runs `sonar-scanner` (or Maven sonar plugin) to analyze code quality.
- Submits coverage reports, duplication metrics, and security findings to SonarQube.
- Communicates with SonarQube server (`http://sonarqube:9000`) to create a new analysis.
- Output: SonarQube dashboard updated; `report-task.txt` generated with task URL.
- **Failure reason**: Sonar server unreachable, invalid token, or configuration error.

#### Stage 5: **Quality Gate**

- **Purpose**: Enforce quality standards before allowing the build to proceed.
- **What it checks**:
  - Code coverage is above minimum threshold (configured per project).
  - No critical or blocker-level issues detected.
  - Security vulnerabilities meet acceptance criteria.
- **How it works**:
  - Polls the Sonar CE task endpoint to check analysis completion.
  - Retrieves quality gate status from SonarQube API.
  - **Pass**: Status is `OK` → proceed to next stage.
  - **Fail**: Status is `ERROR` → pipeline stops; no deployment occurs.
- **Output**: Console log shows "Quality Gate status: OK" or failure details.
- **Failure reason**: Coverage too low, critical issues present, or gate configuration not met.

#### Stage 6: **Package**

- **Purpose**: Build a Docker image and push to the local registry.
- **Steps**:
  1. Reads the service's `Dockerfile` (multi-stage: build → runtime).
  2. Runs `docker build -t localhost:5000/<service>:<BUILD_NUMBER> .`.
  3. Tags the image with `latest` tag as well.
  4. Pushes both tags to the local Docker registry (`http://localhost:5000`).
- **Output**: Image available in registry; can be deployed to other environments.
- **Failure reason**: Docker daemon unreachable, build layer fails, or registry unavailable.

#### Stage 7: **E2E Tests** _(notification-service only, can be extended)_

- **Purpose**: Verify the microservice works correctly in a full-stack environment.
- **Steps**:
  1. Cleans up any orphaned service containers from previous runs.
  2. Runs `docker-compose up` with the newly built image to start:
     - The microservice (using the built image from registry)
     - Its PostgreSQL database
     - Dependent services (RabbitMQ, other services)
  3. Executes end-to-end tests (e.g., API health checks, event consumption).
  4. Collects logs and outputs for debugging.
- **Output**: Service logs and test results; confirms service is production-ready.
- **Failure reason**: Service failed to start, database migration error, E2E test assertion failed, or inter-service communication broken.

#### Stage 8: **Post Actions**

- **Purpose**: Archive build artifacts and logs for future reference.
- **Actions**:
  - Archive coverage reports, test results, and SonarQube task report.
  - Collect logs from service containers for debugging failures.
- **Output**: Artifacts accessible in the Jenkins job's "Archive" tab.

---

### Interpreting Results

#### Jenkins Dashboard View

After a pipeline run completes, Jenkins displays the stage view:

![Pipeline Success](./screenshots/pipeline-success.png)

**Stage Status Indicators:**

| Status         | Color     | Meaning                                                        |
| -------------- | --------- | -------------------------------------------------------------- |
| ✅ SUCCESS     | **Green** | Stage completed without errors.                                |
| ❌ FAILURE     | **Red**   | Stage failed; pipeline stopped. Subsequent stages are skipped. |
| ⊘ SKIPPED      | **Gray**  | Stage was not executed (usually due to earlier failure).       |
| ⏳ IN PROGRESS | **Blue**  | Stage is currently running.                                    |

#### Interpreting Common Failures

**Red stage: "Build"**

- **Cause**: Compilation error (syntax, missing imports, or incompatible dependency).
- **Fix**: Check console output for compiler errors; correct code and push fix.

**Red stage: "Test"**

- **Cause**: Unit test failed (assertion, NPE, timeout).
- **Fix**: Review test logs; debug failing test and push fix.

**Red stage: "Quality Gate"**

- **Cause**: Coverage below threshold or critical issues detected.
- **Fix**: Increase test coverage or resolve Sonar issues; see SonarQube dashboard for details.

**Red stage: "Package"**

- **Cause**: Docker build failed (missing file, layer failure).
- **Fix**: Verify Dockerfile syntax and dependencies; check Docker daemon logs.

**Red stage: "E2E Tests"**

- **Cause**: Service failed to start, test assertion failed, or inter-service communication broken.
- **Fix**: Review service startup logs; check database migrations and network connectivity.

#### Viewing Console Output

1. Click on the failed pipeline job in Jenkins.
2. Click **Console Output** (or stage-specific logs).
3. Scroll to the error; identify the root cause.
4. Example error output:
   ```
   FAILURE: Post-Conditions Check Failed
   Quality Gate has failed. See SonarQube for details.
   ```

#### Checking SonarQube Dashboard

For code quality issues:

1. Open `http://localhost:9000` (Admin credentials: `admin` / `admin`).
2. Navigate to project (e.g., `notification-service`).
3. Review:
   - **Code Smells**: Maintainability issues.
   - **Security Hotspots**: Potential vulnerabilities.
   - **Coverage**: Current test coverage percentage.
   - **Duplications**: Repeated code blocks.

---

### Architecture: Docker-outside-of-Docker (DooD)

Jenkins runs inside a container but needs to build and run other containers (for tests, E2E, packaging). This is achieved via **Docker-outside-of-Docker (DooD)**:

- Jenkins container mounts the host's Docker socket: `/var/run/docker.sock:/var/run/docker.sock`.
- Jenkins runs Docker commands directly on the host daemon (not a nested Docker instance).
- Agents use `docker` CLI and `docker compose` to orchestrate services.

**Benefit**: No performance overhead; fully isolated build environments per job.

---

### Environment Variables & Configuration

#### Jenkins Environment Variables

| Variable              | Value                                    | Purpose                         |
| --------------------- | ---------------------------------------- | ------------------------------- |
| `JAVA_OPTS`           | `-Djenkins.install.runSetupWizard=false` | Skip setup wizard on first run. |
| `CASC_JENKINS_CONFIG` | `/var/jenkins_home/casc_configs`         | JCasC config directory.         |

#### SonarQube Configuration

- **Server URL**: `http://sonarqube:9000`
- **Default Credentials**: `admin` / `admin` (changeable via UI)
- **Token**: Generated per user; required for CI integration.

#### Docker Registry

- **Local Registry URL**: `localhost:5000`
- **Images pushed**: `localhost:5000/notification-service:<BUILD_NUMBER>`
- **Cleanup**: Manual `docker image rm` or registry garbage collection.

---

## Technical Challenge Documentation (Reto 6)

### Punto 1: Jenkins Configuration - Technical Q&A

#### Q: What additional plugins are pre-installed in the Dockerfile based on tools used in later points?

**A:** The Jenkins Dockerfile pre-installs seven key plugins via `jenkins-plugin-cli`:

1. **workflow-aggregator**: Enables declarative and scripted pipeline support
2. **git**: Git integration for repository checkout
3. **docker-pipeline**: Docker integration within Jenkinsfile (for `docker.image()` commands)
4. **pipeline-stage-view**: Visual stage representation in Jenkins UI
5. **configuration-as-code (JCasC)**: Auto-configuration from `jenkins/casc.yaml`
6. **job-dsl**: Programmatic job creation (required by JCasC for pipeline jobs)
7. **sonar**: SonarQube scanner integration for quality analysis

These plugins enable the complete CI/CD workflow: containerized builds, static analysis, and automated job provisioning.

#### Q: Does it mount the Docker socket or use DinD? How does it resolve Docker socket permissions?

**A:** The implementation uses **Docker-outside-of-Docker (DooD)** with socket mounting:

- **Socket Mounting**: `/var/run/docker.sock:/var/run/docker.sock` is mounted in docker-compose.yml
- **Permission Resolution**:
  1. Jenkins Dockerfile installs `docker.io` as root user
  2. Adds the `jenkins` user to the Docker group: `usermod -aG docker jenkins`
  3. The group ID is mapped in docker-compose: `group_add: ["991"]` (Docker daemon's group ID on host)
  4. This allows jenkins user to access the socket without needing `sudo`

**Advantage over DinD**: Simpler setup, better performance (no Docker-in-Docker overhead), and direct host Docker daemon access for image building and container management.

#### Q: Does it disable the Setup Wizard or keep it for manual admin user creation?

**A:** The Setup Wizard is **disabled** via the environment variable:

```yaml
environment:
  - JAVA_OPTS=-Djenkins.install.runSetupWizard=false
```

This means:

- Jenkins skips the interactive setup wizard on first run
- Admin user credentials must be retrieved from `/var/jenkins_home/secrets/initialAdminPassword`
- Jenkins Configuration as Code (JCasC) automatically configures everything (no manual UI setup needed)
- This is ideal for CI/CD environments where automation and reproducibility are critical

---

### Punto 2: Build & Test Jenkinsfiles - Technical Q&A

#### Q: How does dependency management avoid downloading all dependencies from scratch on each run?

**A:** The pipeline uses **Docker volume mounting** for dependency caches:

**For Maven (notification-service)**:

```groovy
docker.image('maven:3.9.6-eclipse-temurin-17')
  .inside("-v /var/jenkins_home/.m2:/root/.m2 --network ${DOCKER_NET}") {
    // Maven caches to /root/.m2 which maps to Jenkins volume
  }
```

In docker-compose.yml:

```yaml
volumes:
  - maven_cache:/var/jenkins_home/.m2
```

**For Go (auth-service)**:

```bash
export GOMODCACHE="$WORKSPACE/.go/pkg/mod"
export GOCACHE="$WORKSPACE/.go/cache"
mkdir -p "$GOMODCACHE" "$GOCACHE"
go mod download  # Downloads to the cached location
```

**Benefits**:

- Maven dependencies (~400MB) are downloaded once and reused across builds
- Go modules and build cache avoid repeated downloads
- Builds complete 2-3x faster after the initial run
- Reduces bandwidth and network latency

---

### Punto 3: SonarQube Integration - Technical Q&A

#### Q: How does the pipeline manage the SonarQube token?

**A:** The SonarQube token is managed using **Jenkins Credentials**:

1. **Storage**: Token is stored in Jenkins Credentials as `sonar-token` (type: `Secret text`)

   - Retrieved in Jenkinsfile: `withCredentials([string(credentialsId: 'sonar-token', variable: 'SONAR_TOKEN')])`
   - Not stored in docker-compose.yml or git repositories (no hardcoding)

2. **Token Generation**:

   - Created in SonarQube dashboard: `http://localhost:9000 > Account > Security > Tokens`
   - Must be configured manually after SonarQube first startup

3. **Usage in Pipeline**:
   - Auth-service: `sonar-scanner -Dsonar.login="$SONAR_TOKEN"`
   - Notification-service: `mvn sonar:sonar -Dsonar.login="$SONAR_TOKEN"`

This approach keeps credentials secure and separate from infrastructure configuration.

#### Q: Does it use the default Quality Gate ("Sonar way") or a custom one? What are the thresholds?

**A:** The implementation uses **SonarQube's default "Sonar way" Quality Gate** with these thresholds:

- **Code Coverage**: ≥ 70% (enforced via `sonar-project.properties`)
- **Code Duplication**: < 3%
- **Code Smells**: Blockers and Critical issues must be 0
- **Security Hotspots**: Must be reviewed
- **Reliability**: No critical or blocker issues

To verify the gate:

1. Navigate to `http://localhost:9000 > Quality Gates`
2. The "Sonar way" gate appears as default
3. Each project's "Quality Gate" setting is configured in `sonar-project.properties`

**Customization Option**: To create a custom gate:

- SonarQube UI > Quality Gates > Create
- Set thresholds per language (Go, Java, etc.)

#### Q: How does it configure the SonarQube webhook for quality gate checks? For non-Java languages, how is sonar-scanner installed?

**A:**

**Quality Gate Polling (instead of webhooks)**:

- The pipeline does NOT use webhooks; instead it **polls the CE task status** with exponential backoff:
  ```bash
  for i in $(seq 1 60); do
    CE_STATUS=$(curl -u "$SONAR_TOKEN:" "$CE_TASK_URL")
    if [ "$CE_STATUS" = "SUCCESS" ]; then break; fi
    sleep 5
  done
  ```
- This avoids webhook configuration complexity and works in isolated environments

**For non-Java languages (Go in auth-service)**:

The pipeline runs `sonar-scanner` inside a Docker container:

```groovy
docker.image('sonarsource/sonar-scanner-cli:latest')
  .inside("--network ${DOCKER_NET}") {
    sh 'sonar-scanner -Dsonar.projectKey=auth-service ...'
  }
```

- **Scanner Installation**: Pre-installed in `sonarsource/sonar-scanner-cli:latest` image
- **No setup required**: Docker handles the scanner availability
- **Go Coverage**: Configured in `sonar-project.properties`:
  ```properties
  sonar.go.coverage.reportPaths=coverage.out
  ```

For Java (notification-service), Maven plugin handles scanning:

```bash
mvn sonar:sonar -Dsonar.projectKey=notification-service
```

---

### Punto 4: Docker Packaging & E2E Testing - Technical Q&A

#### Q: What naming convention is used for Docker images? Does it include project prefix, semantic versioning, or commit hash?

**A:** The naming convention follows: `localhost:5000/<SERVICE_NAME>:<BUILD_NUMBER>` with additional tags:

```bash
docker build -t "localhost:5000/${SERVICE_NAME}:${BUILD_NUMBER}" \
             -t "localhost:5000/${SERVICE_NAME}:latest" \
             "${SERVICE_DIR}"
docker push "localhost:5000/${SERVICE_NAME}:${BUILD_NUMBER}"
docker push "localhost:5000/${SERVICE_NAME}:latest"
```

**Convention breakdown**:

- **Registry**: `localhost:5000` (local Docker registry)
- **Service Name**: `notification-service`, `auth-service` (environment variable: `${SERVICE_NAME}`)
- **Build Number Tag**: Jenkins build ID (ensures unique builds)
- **Latest Tag**: Always points to most recent build

**Example image names**:

- `localhost:5000/notification-service:42` (build #42)
- `localhost:5000/notification-service:latest` (current)

**Note**: This project does NOT use:

- Semantic versioning (v1.2.3) - managed by external release process
- Commit hashes - BUILD_NUMBER provides sufficient traceability
- Project prefix - single registry serves all services

#### Q: Do the Dockerfiles use multi-stage builds? How does this integrate with the pipeline?

**A:** **YES**, all microservices use multi-stage builds to optimize image sizes:

**Example 1: Java (notification-service)**

```dockerfile
FROM maven:3.9-eclipse-temurin-17 AS build
# Build stage: compile, test, package (~800MB intermediate)
COPY pom.xml ./
RUN mvn -B -DskipTests clean package

FROM eclipse-temurin:17-jre
# Runtime stage: ~300MB final image
COPY --from=build /app/target/notification-service-*.jar /app/app.jar
```

**Example 2: Go (auth-service, employee-management)**

```dockerfile
FROM golang:1.24-alpine AS builder
# Build stage: compile (~700MB intermediate)
RUN go build -o app-server ./cmd

FROM alpine:latest
# Runtime stage: ~50MB final image
COPY --from=builder /app/app-server .
```

**Example 3: Node.js (profile-management)**

```dockerfile
FROM node:20-alpine AS build
# Build stage: npm install, build (~500MB intermediate)
RUN npm run build

FROM node:20-alpine
# Runtime stage: ~200MB final image (prod dependencies only)
RUN npm install --omit=dev
```

**Pipeline Integration**:

- Docker build happens in "Package" stage of Jenkinsfile
- Multi-stage builds are automatic (no special Jenkinsfile configuration)
- Benefits: Smaller images (faster push/pull), reduced attack surface

#### Q: Is it a local or remote registry? How is authentication configured?

**A:** The implementation uses a **local Docker registry** running in docker-compose:

```yaml
registry:
  image: registry:2
  container_name: registry
  ports:
    - "5000:5000"
  volumes:
    - registry_data:/var/lib/registry
  networks:
    - microservices-network
```

**Registry Configuration**:

- **Location**: `localhost:5000` (accessible within docker-compose network as `registry:5000`)
- **Authentication**: NONE (unauthenticated registry for development)
- **Storage**: Persists to Docker volume `registry_data` (survives `docker-compose down`)

**Pipeline Integration**:

```bash
docker push "localhost:5000/${SERVICE_NAME}:${BUILD_NUMBER}"
```

**For Production**: Replace with Docker Hub or private registry:

- Docker Hub: Use `docker login` with credentials stored in Jenkins
- Private Registry (Nexus, Artifactory): Configure authentication in Jenkins credentials

#### Q: How does the pipeline verify all services are ready before running E2E tests?

**A:** The E2E stage performs **multi-level health checks**:

**1. Service Startup** (docker-compose up):

```bash
run_compose --project-name "${CI_PROJECT}" --profile test up -d --build \
  rabbitmq \
  database-employees employees-service \
  database-departments departments-service \
  database-notifications notifications-service \
  database-profiles profiles-service \
  database-auth auth-service \
  api-gateway
```

**2. HTTP Polling** (explicit checks before BDD tests):

```bash
wait_for_http "API Gateway" "http://api-gateway:8080/" 60
wait_for_http "Auth Service" "http://auth-service:8083/api/swagger.yaml" 60
wait_for_http "Employees Service" "http://employees-service:8081/api/health" 60
wait_for_http "Departments Service" "http://departments-service:8082/api/" 60
wait_for_http "Notifications Service" "http://notifications-service:8084/v3/api-docs" 60
wait_for_http "Profiles Service" "http://profiles-service:8085/health" 60
```

Where `wait_for_http` is a polling function:

```bash
for i in $(seq 1 "${max_tries}"); do
  if docker run --rm --network "${NET}" curlimages/curl:8.10.1 \
     -fsS "${url}" >/dev/null 2>&1; then
    echo "${name} is ready"
    return 0
  fi
  sleep 5
done
```

**3. Why this approach**:

- `depends_on` with `condition: service_started` only checks container startup, not readiness
- HTTP polling ensures services are actually accepting connections
- 60 retries × 5 seconds = 5 minutes max wait time
- Prevents flaky tests by waiting for services to be truly ready

#### Q: How does the pipeline guarantee E2E tests don't interfere with data from previous runs?

**A:** Complete isolation is achieved through **project-specific containers and data cleanup**:

**1. Unique Project Naming** (per build):

```bash
CI_PROJECT="ci-${BUILD_NUMBER}"
run_compose --project-name "${CI_PROJECT}" ...
```

Results in:

- Network: `ci-42_onboarding-network` (for build #42)
- Containers: `ci-42_database-notifications_1`, `ci-42_auth-service_1`, etc.
- Volumes: All test-specific, separate from production/previous runs

**2. Container Cleanup Before Tests**:

```bash
for container in rabbitmq database-employees database-departments \
                 database-notifications database-profiles database-auth \
                 api-gateway notifications-service employees-service \
                 departments-service profiles-service auth-service; do
  docker rm -f "$container" 2>/dev/null || true
done
docker container prune -f --filter "label!=keep=true" >/dev/null 2>&1 || true
```

**3. Volume Cleanup After Tests** (post section):

```bash
run_compose --project-name "${CI_PROJECT}" down -v --remove-orphans
```

The `down -v` flag removes:

- All containers created by the compose file
- Named volumes specific to that project
- Network specific to that project

**Result**: Each E2E test run is completely isolated; no data bleeds between builds.

---

### Punto 5: General CI/CD Architecture & Benefits

#### Q: Why should each microservice have its own Jenkinsfile instead of a monolithic Jenkins job?

**A:** Individual Jenkinsfiles per service provide:

1. **Decoupled Pipelines**: Each service's build is independent; failures in one don't block others
2. **Language-Specific Stages**: Go services use different build tools than Java or Node.js
3. **Scalability**: New services can be added without modifying central job definitions
4. **Debugging**: Service-specific logs make it easy to identify which service failed
5. **Parallel Execution**: Jenkins can trigger multiple service pipelines simultaneously

#### Q: What is the purpose of using Docker agents in the pipeline instead of bare Jenkins executors?

**A:** Docker agents provide:

1. **Clean Environment**: Each build starts fresh; no leftover artifacts from previous builds
2. **Language Isolation**: Go build environment doesn't interfere with Java or Node.js builds
3. **Reproducibility**: Builds are identical regardless of Jenkins host setup
4. **Security**: Services don't have direct host access; only what's mounted explicitly
5. **Scalability**: Can spin up arbitrary build images without installing them on Jenkins

#### Q: How does the pipeline ensure code quality before allowing deployments?

**A:** Multi-layered quality gates prevent poor code from reaching production:

1. **Unit Tests**: Each service's test stage runs unit tests and generates coverage reports
2. **SonarQube Analysis**: Static code analysis identifies bugs, security vulnerabilities, and code smells
3. **Quality Gate**: Enforces minimum coverage (≥70%), blocks if issues exceed thresholds
4. **E2E Tests**: Full system integration tests verify all services work together
5. **Build Artifact Archiving**: Test results are preserved for audit and troubleshooting

Failure at any stage stops the pipeline and prevents image push to registry.

#### Q: What is the CI/CD workflow from commit to production-ready Docker image?

**A:** The complete workflow is:

```
1. Developer commits code to GitHub
   ↓
2. GitHub webhook triggers Jenkins pipeline (if configured)
   ↓
3. Checkout Stage: Git clone the repository
   ↓
4. Build Stage: Compile code in Docker agent, cache dependencies
   ↓
5. Test Stage: Run unit tests, generate coverage reports
   ↓
6. SonarQube Analysis: Upload coverage and run static analysis
   ↓
7. Quality Gate: Poll SonarQube until status is SUCCESS or FAILED
   ↓ (if any stage fails, stop here and notify developer)
   ↓
8. Package Stage: Build Docker image, push to localhost:5000 registry
   ↓
9. E2E Tests (notification-service only): Spin up entire microservices stack
   ↓
10. Run BDD tests against live services
    ↓ (if tests fail, cleanup and stop)
    ↓
11. Cleanup: Remove test containers and volumes
    ↓
12. Success: Production-ready Docker image tagged and available in registry
```

**Time to Production**: ~5-10 minutes per service (first run slower due to dependency downloads).

---

### Troubleshooting

#### Jenkins Won't Start

```bash
# Check logs
docker logs microservices-fp_jenkins_1

# Verify Docker socket is mounted
docker inspect microservices-fp_jenkins_1 | grep -A 5 "Volumes"

# Restart Jenkins
docker restart microservices-fp_jenkins_1
```

#### Pipelines Not Appearing

- JCasC may not have been applied. Restart Jenkins.
- Check JCasC config file: `jenkins/casc.yaml` is mounted correctly.

#### Pipeline Fails with "Docker daemon not accessible"

```bash
# Verify Docker socket mount in docker-compose.yml
docker inspect microservices-fp_jenkins_1 --format='{{json .Mounts}}'

# Verify Jenkins user can access Docker
docker exec microservices-fp_jenkins_1 docker ps
```

#### SonarQube Analysis Timeout

- SonarQube may still be initializing (first run takes 1-2 minutes).
- Check SonarQube logs: `docker logs microservices-fp_sonarqube_1`

#### Quality Gate Always Fails

1. Check Sonar token: `http://localhost:9000/account/security`
2. Verify coverage threshold in SonarQube project settings.
3. Review SonarQube dashboard for critical issues.

---

### Next Steps

1. **Monitor Pipelines**: After the first successful run, pipelines are ready for production-like workflows.
2. **Set Up Webhooks**: Configure GitHub to trigger builds on push (optional).
3. **Scale to More Services**: Add additional microservices to the CI pipeline using the same Jenkinsfile pattern.
4. **Customize Quality Gates**: Adjust SonarQube thresholds based on team standards.
5. **Deployment Integration**: Extend pipelines with deployment stages (Dev, Staging, Prod) as needed.
