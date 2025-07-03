#!/bin/bash
# Test script for asciicam

set -e

echo "Running tests for asciicam..."

# Run unit tests
echo "Running unit tests..."
go test -v ./...

# Run tests with coverage
echo "Running tests with coverage..."
go test -cover ./...

# Run race condition tests
echo "Running race condition tests..."
go test -race ./...

echo "All tests completed successfully!"