#!/bin/bash
# Setup script for AulaFlash backend

set -e

echo "=== AulaFlash Backend Setup ==="

# Check Go
if ! command -v go &> /dev/null; then
  echo "ERROR: Go is not installed"
  exit 1
fi

echo "Go version: $(go version)"

# Create directories
mkdir -p /tmp/aulaflash-uploads

# Install deps
echo "Installing Go dependencies..."
cd "$(dirname "$0")/.."
go mod tidy

# Database
echo "Creating database..."
createdb aulaflash 2>/dev/null || echo "Database may already exist"

# Run migrations
echo "Running migrations..."
go run cmd/migrate/main.go

echo ""
echo "=== Setup complete ==="
echo "Run: make run"
