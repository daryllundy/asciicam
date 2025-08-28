# Project Structure

## Directory Layout

```
asciicam/
├── cmd/asciicam/          # Main application entry point
│   └── main.go            # Application bootstrap and main loop
├── internal/              # Private application packages
│   ├── ascii/             # ASCII/ANSI conversion logic
│   ├── camera/            # Camera capture and image handling
│   ├── config/            # Configuration management and CLI parsing
│   ├── errors/            # Custom error types and handling
│   └── greenscreen/       # Virtual greenscreen functionality
├── docs/                  # Documentation and assets
│   └── screenshots/       # Demo images
├── examples/              # Usage example scripts
├── scripts/               # Build and development scripts
├── .kiro/                 # Kiro AI assistant configuration
└── .github/               # GitHub workflows and templates
```

## Package Organization

### `cmd/asciicam/`
- **Purpose**: Application entry point and main execution loop
- **Key responsibilities**: Signal handling, context management, component orchestration
- **Pattern**: Minimal main.go that delegates to internal packages

### `internal/` Packages
- **Visibility**: Private to this module (Go internal package convention)
- **Structure**: Domain-driven organization by functionality

#### `internal/config/`
- Configuration struct and CLI flag parsing
- Terminal size detection and validation
- Default value management
- Color parsing and validation

#### `internal/ascii/`
- Image-to-ASCII character conversion
- ANSI color block rendering
- Pixel intensity mapping
- Character set management

#### `internal/camera/`
- OpenCV camera capture interface
- Image resizing and preprocessing
- Frame reading with context support
- Camera device management

#### `internal/greenscreen/`
- Background sample generation and storage
- Virtual greenscreen processing
- Threshold-based background removal
- Sample-based background detection

#### `internal/errors/`
- Custom error types
- Error wrapping and context
- Application-specific error handling

## Code Conventions

### File Naming
- `package_name.go` - Main package implementation
- `package_name_test.go` - Unit tests
- Use descriptive names that match functionality

### Package Structure
- Each package should have a clear, single responsibility
- Expose minimal public API surface
- Use constructor functions (e.g., `NewConverter()`, `NewCapture()`)
- Include comprehensive test coverage with benchmarks

### Testing
- Unit tests in `*_test.go` files alongside source
- Benchmark tests for performance-critical code
- Table-driven tests for multiple scenarios
- Test both success and error cases

### Documentation
- Package-level documentation in doc comments
- Exported functions and types must have doc comments
- Examples in `examples/` directory for common usage patterns

## Build Artifacts
- **Binary output**: `asciicam` (root directory)
- **Cross-compiled binaries**: `build/` directory
- **Temporary files**: Cleaned by `make clean`

## Configuration Files
- `.golangci.yml` - Linting configuration
- `go.mod` / `go.sum` - Go module dependencies
- `Makefile` - Build automation
- `.gitignore` - Version control exclusions
