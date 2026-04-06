# -------- Build stage --------
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install git (needed for go modules in some cases)
RUN apk add --no-cache git

# Copy go mod files first (better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build binary natively for the host architecture
RUN CGO_ENABLED=0 GOMEMLIMIT=512MiB go build -p 1 -trimpath -ldflags="-s -w" -o employee-server ./cmd

# -------- Runtime stage --------
FROM alpine:latest

WORKDIR /app

# Create a non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy binary and docs from builder
COPY --from=builder /app/employee-server .
COPY --from=builder /app/docs ./docs

# Change ownership of files to the non-root user
RUN chown -R appuser:appgroup /app

# Switch to the non-root user
USER appuser

# Expose the application port
EXPOSE 8081

# Run the application
CMD ["./employee-server"]
