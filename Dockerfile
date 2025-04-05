# user-service-qubool-kallyaanam/Dockerfile
FROM golang:1.23.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o user-service ./cmd/server/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/user-service .

EXPOSE 8082

CMD ["./user-service"]