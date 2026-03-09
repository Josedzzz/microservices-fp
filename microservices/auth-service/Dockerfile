# -------- Build stage --------
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install git for go modules
RUN apk add --no-cache git

# Copy go mod files first (better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary natively for the host architecture
RUN CGO_ENABLED=0 GOMEMLIMIT=512MiB go build -p 1 -trimpath -ldflags="-s -w" -o auth-server ./cmd

# -------- Runtime stage --------
FROM alpine:latest

WORKDIR /app

# Create non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy binary from builder
COPY --from=builder /app/auth-server .

# Change ownership
RUN chown -R appuser:appgroup /app

# Use non-root user
USER appuser

# Expose port
EXPOSE 8083

# Run application
CMD ["./auth-server"]
