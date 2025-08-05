# Dockerfile
FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod tidy

# Copy semua file, termasuk .env
COPY . .
COPY .env .

RUN go build -o main ./main.go

# Run Stage
FROM debian:bookworm

WORKDIR /app

# Install postgres-client agar pg_isready tersedia
RUN apt-get update && apt-get install -y postgresql-client && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/main .
COPY --from=builder /app/.env .env

EXPOSE 7003

# Gunakan script wait-for-postgres sebagai entrypoint
CMD ["./main"]

