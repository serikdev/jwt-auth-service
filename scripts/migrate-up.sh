#!/bin/bash

set -e

echo "Running DB migrations...."

source .env
goose -dir ./migrations postgres "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}" up

echo "Migrations completed successfully"