#!/bin/bash
set -e

# Define variables
DATABASE_NAME="chat"

# Check if the database already exists
if psql -U "$POSTGRES_USER" -lqt | cut -d \| -f 1 | grep -qw "$DATABASE_NAME"; then
    echo "Database '$DATABASE_NAME' already exists. Skipping creation."
else
    echo "Database '$DATABASE_NAME' does not exist. Creating..."
    psql -U "$POSTGRES_USER" -c "CREATE DATABASE $DATABASE_NAME;"
    echo "Database '$DATABASE_NAME' created successfully."
fi
