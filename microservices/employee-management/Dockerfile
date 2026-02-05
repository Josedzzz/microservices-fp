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

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o employee-server ./cmd

# -------- Runtime stage --------
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/employee-server .

# Expose the application port
EXPOSE 8081

# Run the application
CMD ["./employee-server"]
