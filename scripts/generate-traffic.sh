#!/bin/bash
# ============================================================
# Traffic Generation Script for Observability Demo
# ============================================================
# This script generates synthetic traffic to populate:
#   - Prometheus metrics (HTTP requests, latency, errors)
#   - Grafana dashboards (service panels, latency graphs)
#   - Zipkin distributed traces
#   - Loki logs
#
# Usage: ./scripts/generate-traffic.sh
# ============================================================

set -eo pipefail

GATEWAY_URL="${GATEWAY_URL:-http://localhost:8080}"
AUTH_EMAIL="${AUTH_EMAIL:-admin@onboarding.com}"
AUTH_PASSWORD="${AUTH_PASSWORD:-admin123}"

echo "============================================"
echo "  Observability Traffic Generator"
echo "============================================"
echo "Gateway: $GATEWAY_URL"
echo ""
sleep 1

# --------------------------------------------------
# Step 1: Authenticate and obtain JWT
# --------------------------------------------------
echo "[1/4] Authenticating as $AUTH_EMAIL ..."

LOGIN_RESPONSE=$(curl -s -L -X POST "$GATEWAY_URL/auth-service/api/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\": \"$AUTH_EMAIL\", \"password\": \"$AUTH_PASSWORD\"}")

TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"access_token"[[:space:]]*:[[:space:]]*"[^"]*"' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
  echo "ERROR: Login failed. Response: $LOGIN_RESPONSE"
  exit 1
fi

echo "       Token obtained successfully"
echo ""

# --------------------------------------------------
# Step 2: Create a department
# --------------------------------------------------
echo "[2/4] Creating department (ENG) ..."

curl -s -L -X POST "$GATEWAY_URL/departments-service/api/departments" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "id": "ENG",
    "name": "Engineering",
    "description": "Engineering department for traffic generation"
  }' > /dev/null

echo "       Department ENG created (or already exists)"
echo ""

# --------------------------------------------------
# Step 3: Create employees
# --------------------------------------------------
echo "[3/4] Creating sample employees ..."

CREATED=0
for entry in \
  "Alice Johnson,alice.johnson@example.com" \
  "Bob Smith,bob.smith@example.com" \
  "Carol Williams,carol.williams@example.com" \
  "David Brown,david.brown@example.com" \
  "Eva Martinez,eva.martinez@example.com" \
  "Frank Garcia,frank.garcia@example.com" \
  "Grace Lee,grace.lee@example.com" \
  "Henry Wilson,henry.wilson@example.com" \
  "Ivy Anderson,ivy.anderson@example.com" \
  "Jack Taylor,jack.taylor@example.com"; do

  NAME="${entry%,*}"
  EMAIL="${entry#*,}"

  HTTP_CODE=$(curl -s --max-time 10 -L -o /dev/null -w "%{http_code}" \
    -X POST "$GATEWAY_URL/employees-service/api/employees" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d "{\"name\": \"$NAME\", \"email\": \"$EMAIL\", \"departmentID\": \"ENG\"}")

  if [ "$HTTP_CODE" = "201" ]; then
    echo "       + Created: $NAME <$EMAIL>"
    CREATED=$((CREATED + 1))
  else
    echo "       - Skipped: $NAME (HTTP $HTTP_CODE)"
  fi
done
echo "       Created $CREATED new employees"
echo ""

# --------------------------------------------------
# Step 4: Generate read traffic
# --------------------------------------------------
echo "[4/4] Generating read traffic (30 iterations) ..."
echo ""

for i in $(seq 1 30); do
  # List all employees
  HTTP_CODE=$(curl -s -L -o /dev/null -w "%{http_code}" -X GET "$GATEWAY_URL/employees-service/api/employees" \
    -H "Authorization: Bearer $TOKEN")
  printf "       [%02d/30] List employees -> HTTP %s\n" "$i" "$HTTP_CODE"

  # Get a specific employee by ID (cycle through first 5)
  EMP_ID=$(( (i % 5) + 1 ))
  HTTP_CODE=$(curl -s -L -o /dev/null -w "%{http_code}" -X GET "$GATEWAY_URL/employees-service/api/employees/$EMP_ID" \
    -H "Authorization: Bearer $TOKEN")
  printf "       [%02d/30] Get employee   -> HTTP %s (id=%d)\n" "$i" "$HTTP_CODE" "$EMP_ID"

  sleep 0.5
done

echo ""
echo "============================================"
echo "  Traffic generation complete!"
echo "============================================"
echo ""
echo "What to check next:"
echo "  - Grafana:  http://localhost:3000"
echo "  - Prometheus: http://localhost:9090"
echo "  - Zipkin:   http://localhost:9411"
echo "  - Loki:    http://localhost:3100"
echo ""
