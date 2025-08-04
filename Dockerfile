# Dockerfile
FROM golang:1.24 as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod tidy

# Copy semua file, termasuk .env
COPY . .

RUN go build -o main ./main.go

# Run Stage
FROM debian:bookworm

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/.env .

EXPOSE 7003

CMD ["./main"]
