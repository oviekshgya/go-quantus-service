#!/bin/sh

echo "🔄 Menunggu PostgreSQL..."
until nc -z postgres 5432; do
  sleep 1
done

echo "✅ PostgreSQL siap."

echo "🔄 Menunggu RabbitMQ..."
until nc -z rabbitmq 5672; do
  sleep 1
done

echo "✅ RabbitMQ siap."

# Jalankan app Go kamu
./quantus-app
