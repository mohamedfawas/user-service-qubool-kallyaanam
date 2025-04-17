# user-service-qubool-kallyaanam/Dockerfile
FROM golang:1.23.4-alpine AS builder

# Install required packages
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Copy go.mod and go.sum files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . .

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o user-service ./cmd/server/main.go

# Use a minimal alpine image for the final stage
FROM alpine:latest

RUN apk add --no-cache ca-certificates postgresql-client

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/user-service .

# Copy the migrations directory
COPY --from=builder /app/migrations ./migrations

# Expose the application port
EXPOSE 8082

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget -q --spider http://localhost:8082/health || exit 1

# Run the application
CMD ["./user-service"]