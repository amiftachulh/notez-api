#!/bin/bash

# Migrate
# Usage: ./migrate.sh <up|down>

# Check if the direction is valid
if [ "$1" != "up" ] && [ "$1" != "down" ]; then
  echo "Please provide a valid direction (up or down)"
  echo "Usage: ./migrate.sh <up|down>"
  exit 1
fi

# Get .env file and check if it exists
# Check if .env file exists and inside of it there is a DATABASE_URL variable
if [ -f .env ]; then
  source .env
  if [ -z "$DATABASE_URL" ]; then
    echo "Please provide a DATABASE_URL in the .env file"
    exit 1
  fi
else
  echo "Please provide a .env file with a DATABASE_URL variable"
  exit 1
fi

migrate -verbose  -database $DATABASE_URL -path db/migrations $1 
