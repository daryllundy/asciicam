# Makefile for asciicam

.PHONY: build clean test lint fmt deps install help run dev cross-build check coverage

# Variables
BINARY_NAME = asciicam
CMD_PATH = ./cmd/asciicam
BUILD_DIR = build
MAIN_PACKAGE = $(CMD_PATH)
VERSION ?= $(shell git describe --tags --always --dirty)
LDFLAGS = -ldflags "-X main.version=$(VERSION)"

# Default target
all: build

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	go build $(LDFLAGS) -o $(BINARY_NAME) $(MAIN_PACKAGE)

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f $(BINARY_NAME)
	rm -rf $(BUILD_DIR)
	go clean

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
coverage:
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run linting
lint:
	@echo "Running linter..."
	golangci-lint run

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	goimports -w .

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod verify

# Install build tools
install-tools:
	@echo "Installing build tools..."
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Install the binary
install: build
	@echo "Installing $(BINARY_NAME)..."
	cp $(BINARY_NAME) /usr/local/bin/

# Run the application
run: build
	./$(BINARY_NAME)

# Development mode - run with live reload (requires 'entr')
dev:
	@echo "Starting development mode..."
	find . -name "*.go" | entr -r make run

# Cross-compilation for multiple platforms
cross-build:
	@echo "Cross-compiling for multiple platforms..."
	mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PACKAGE)
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_PACKAGE)
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PACKAGE)
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PACKAGE)
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PACKAGE)
	@echo "Cross-compilation complete. Binaries in $(BUILD_DIR)/"

# Check for OpenCV dependencies
check-deps:
	@echo "Checking dependencies..."
	@which go > /dev/null || (echo "Go is not installed" && exit 1)
	@pkg-config --exists opencv4 || (echo "OpenCV4 not found. Run 'make install-opencv'" && exit 1)
	@echo "All dependencies are satisfied!"

# Install OpenCV (macOS only)
install-opencv-macos:
	@echo "Installing OpenCV via Homebrew..."
	brew install opencv

# Install OpenCV (Linux - Ubuntu/Debian)
install-opencv-linux:
	@echo "Installing OpenCV via apt..."
	sudo apt-get update
	sudo apt-get install -y libopencv-dev pkg-config

# Quality checks
check: fmt lint test
	@echo "All quality checks passed!"

# CI/CD targets
ci: deps check-deps check coverage
	@echo "CI pipeline completed successfully!"

# Release preparation
release-prep: clean fmt lint test cross-build
	@echo "Release preparation complete!"

# Run benchmarks
bench:
	@echo "Running benchmarks..."
	go test -bench=. ./...

# Generate documentation
docs:
	@echo "Generating documentation..."
	godoc -http=:6060

# Update go.mod
mod-update:
	@echo "Updating go.mod..."
	go get -u ./...
	go mod tidy

# Vendor dependencies
vendor:
	@echo "Vendoring dependencies..."
	go mod vendor

# Help target
help:
	@echo "Available targets:"
	@echo "  build          - Build the application"
	@echo "  clean          - Clean build artifacts"
	@echo "  test           - Run tests"
	@echo "  coverage       - Run tests with coverage report"
	@echo "  lint           - Run linter"
	@echo "  fmt            - Format code"
	@echo "  deps           - Install dependencies"
	@echo "  install-tools  - Install build tools"
	@echo "  install        - Install binary to /usr/local/bin"
	@echo "  run            - Build and run the application"
	@echo "  dev            - Development mode with live reload"
	@echo "  cross-build    - Cross-compile for multiple platforms"
	@echo "  check-deps     - Check for required dependencies"
	@echo "  install-opencv-macos - Install OpenCV on macOS"
	@echo "  install-opencv-linux - Install OpenCV on Linux"
	@echo "  check          - Run all quality checks"
	@echo "  ci             - Run CI pipeline"
	@echo "  release-prep   - Prepare for release"
	@echo "  bench          - Run benchmarks"
	@echo "  docs           - Generate documentation"
	@echo "  mod-update     - Update go.mod"
	@echo "  vendor         - Vendor dependencies"
	@echo "  help           - Show this help message"