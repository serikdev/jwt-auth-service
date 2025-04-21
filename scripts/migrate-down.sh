#!/bin/bash

set -e

echo "Reverting DB migrations...."

source .env
goose -dir ./migrations postgres "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}" down

echo "Migrations reverted successfully"
