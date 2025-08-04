# Dockerfile
FROM golang:1.24 as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
#RUN go mod download
RUN go mod tidy
COPY . .

RUN go build -o main ./main.go

# Run Stage
FROM debian:bullseye-slim

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 7003

CMD ["./main"]
