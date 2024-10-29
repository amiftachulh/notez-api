#!/bin/bash

# Create a new migration file
# Usage: ./create_migration.sh <migration_name>
# Example: ./create_migration.sh create_users_table

# Check if migration name is provided
if [ -z "$1" ]; then
  echo "Please provide a migration name"
  echo "Usage: ./create_migration.sh <name>"
  exit 1
fi

# Create a new migration file
# The migration file will be created in db/migrations
migrate create -ext sql -dir db/migrations $1

# Success message with the migration name
echo "Migration \`$1\` created successfully"
