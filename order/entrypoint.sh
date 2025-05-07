#!/bin/sh

echo "Running migrations..."
migrate -source "${DB_MIGRATIONS_PATH}" -database "${DB_URI}/${DB_NAME}" up

echo "Starting the app..."
exec ./main
