#!/bin/sh

echo "Running migrations..."
migrate -source "file://${DB_MIGRATIONS_PATH}" -database "${DB_URI}/${DB_NAME}" up

echo "Starting the app..."
exec ./main
