# Technology Stack

## Language & Runtime
- **Go 1.23+** - Primary language (go.mod specifies 1.23.0)
- **CGO enabled** - Required for OpenCV bindings

## Key Dependencies
- **gocv.io/x/gocv** - OpenCV Go bindings for camera capture and image processing
- **github.com/muesli/termenv** - Terminal environment detection and ANSI color support
- **github.com/lucasb-eyer/go-colorful** - Color manipulation and hex parsing
- **github.com/nfnt/resize** - Image resizing utilities
- **golang.org/x/term** - Terminal size detection

## System Dependencies
- **OpenCV 4.x** - Computer vision library for camera operations
- **pkg-config** - For OpenCV library detection during build

## Build System

### Make Targets
```bash
# Development
make build          # Build the application
make run            # Build and run
make dev            # Development mode with live reload (requires entr)

# Testing & Quality
make test           # Run tests
make coverage       # Run comprehensive coverage analysis
make lint           # Run golangci-lint
make fmt            # Format code with go fmt and goimports
make check          # Run fmt, lint, and test

# Dependencies
make deps           # Install Go dependencies
make check-deps     # Verify system dependencies (OpenCV, pkg-config)
make install-tools  # Install development tools (goimports, golangci-lint)

# Installation
make install        # Install binary to /usr/local/bin
make clean          # Clean build artifacts

# Cross-compilation
make cross-build    # Build for multiple platforms (Linux, macOS, Windows)

# CI/CD
make ci             # Full CI pipeline (deps, check-deps, check, coverage)
```

### Build Configuration
- **Binary name**: `asciicam`
- **Main package**: `./cmd/asciicam`
- **Build flags**: Includes version from git tags
- **Cross-platform**: Supports Linux (amd64, arm64), macOS (amd64, arm64), Windows (amd64)

## Code Quality Tools
- **golangci-lint** - Comprehensive linting with custom configuration
- **goimports** - Import formatting and organization
- **go fmt** - Standard Go formatting
- **Pre-commit hooks** - Automated quality checks

## Testing
- Standard Go testing framework
- Benchmark tests for performance-critical code
- Coverage reporting via custom script
