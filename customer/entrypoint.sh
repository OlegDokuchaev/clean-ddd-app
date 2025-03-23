#!/bin/sh

echo "Running migrations..."
DATABASE_URL="postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable"
migrate -path ${DB_MIGRATIONS_PATH} -database "$DATABASE_URL" up

echo "Starting the app..."
exec ./main
