FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy all source first so tidy can see main.go
COPY . .

# Fix go.mod and tidy up
RUN go mod edit -go=1.25
RUN go mod tidy
RUN go mod download

# Use memory limits and p=1 to prevent OOM during build
RUN CGO_ENABLED=0 GOMEMLIMIT=512MiB go build -p 1 -trimpath -ldflags="-s -w" -o gateway main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/gateway .
COPY --from=builder /app/templates ./templates

EXPOSE 8080

CMD ["./gateway"]
