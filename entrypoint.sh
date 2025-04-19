#!/bin/sh

# Exit immediately if a command exits with a non-zero status.
set -e

echo "Running database migrations..."
# The DB_SOURCE environment variable should be set in docker-compose.yml
# Example format: postgresql://db_user:db_password@db_host:db_port/db_name?sslmode=disable
migrate -path /app/db/migration -database "$DB_SOURCE" -verbose up

echo "Migrations applied successfully."

# Execute the main command (passed as arguments to this script)
# This will be ["/app/main"] by default (from the Dockerfile CMD)
exec "$@"