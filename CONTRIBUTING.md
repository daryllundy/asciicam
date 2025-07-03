# Contributing to asciicam

Thank you for your interest in contributing to asciicam! This document provides guidelines for contributing to the project.

## Getting Started

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/your-username/asciicam.git
   cd asciicam
   ```
3. **Install dependencies**:
   ```bash
   ./scripts/install-deps.sh
   ```
4. **Build the project**:
   ```bash
   ./scripts/build.sh
   ```

## Development Workflow

### Setting up Development Environment

1. **Prerequisites**:
   - Go 1.22 or newer
   - OpenCV 4.x
   - pkg-config
   - A webcam for testing

2. **Install development tools**:
   ```bash
   go install golang.org/x/tools/cmd/goimports@latest
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   ```

### Making Changes

1. **Create a feature branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes** following the coding standards below

3. **Test your changes**:
   ```bash
   ./scripts/test.sh
   ```

4. **Build and test manually**:
   ```bash
   ./scripts/build.sh
   ./asciicam  # Test basic functionality
   ```

### Coding Standards

- **Go Code Style**: Follow standard Go conventions
- **Formatting**: Use `gofmt` and `goimports`
- **Linting**: Code must pass `golangci-lint` checks
- **Testing**: Add tests for new functionality
- **Documentation**: Update documentation for new features

### Code Structure

- `cmd/asciicam/` - Main application entry point
- `internal/ascii/` - ASCII conversion logic
- `internal/camera/` - Camera handling
- `internal/config/` - Configuration management
- `internal/greenscreen/` - Greenscreen functionality
- `docs/` - Documentation and screenshots
- `examples/` - Usage examples
- `scripts/` - Build and deployment scripts

## Types of Contributions

### Bug Reports

Before creating a bug report:
1. Check existing issues to avoid duplicates
2. Test with the latest version
3. Provide minimal reproduction steps

Include in your bug report:
- Operating system and version
- Go version
- OpenCV version
- Steps to reproduce
- Expected vs actual behavior
- Terminal output or screenshots

### Feature Requests

Before submitting a feature request:
1. Check if it already exists in issues
2. Consider if it fits the project's scope
3. Think about implementation complexity

Include in your feature request:
- Clear description of the feature
- Use case and benefits
- Possible implementation approach
- Examples or mockups if applicable

### Code Contributions

**Good first issues**:
- Documentation improvements
- Test coverage improvements
- Bug fixes
- Performance optimizations

**Areas needing help**:
- Cross-platform compatibility
- Test coverage
- Performance optimization
- Documentation
- New ASCII art modes

## Pull Request Process

1. **Before submitting**:
   - Run tests: `./scripts/test.sh`
   - Run linting: `golangci-lint run`
   - Update documentation if needed
   - Add tests for new functionality

2. **Pull Request Description**:
   - Describe what changes you made
   - Link to relevant issues
   - Include screenshots for UI changes
   - Mention any breaking changes

3. **Review Process**:
   - All PRs require review
   - Address review feedback promptly
   - Keep PRs focused and atomic
   - Squash commits before merging

## Testing

### Running Tests

```bash
# Run all tests
./scripts/test.sh

# Run specific package tests
go test ./internal/ascii/

# Run tests with coverage
go test -cover ./...

# Run benchmarks
go test -bench=. ./...
```

### Writing Tests

- Add tests for new functionality
- Use table-driven tests when appropriate
- Mock external dependencies
- Test edge cases and error conditions
- Follow the naming convention: `TestFunctionName`

### Integration Testing

For features involving camera input:
- Test with different camera resolutions
- Test with different operating systems
- Verify memory usage and performance

## Documentation

### Updating Documentation

- Update README.md for user-facing changes
- Update CLAUDE.md for development changes
- Add examples for new features
- Update inline code comments

### Documentation Style

- Use clear, concise language
- Include code examples
- Add screenshots for visual features
- Keep documentation up to date

## Release Process

1. Update CHANGELOG.md
2. Update version in relevant files
3. Create release PR
4. Tag release after merge
5. Update documentation

## Community

### Code of Conduct

This project adheres to a Code of Conduct. By participating, you are expected to uphold this code.

### Getting Help

- **Documentation**: Check README.md and docs/
- **Issues**: Search existing issues first
- **Questions**: Open a GitHub issue with the "question" label

### Communication

- Be respectful and constructive
- Provide context in issues and PRs
- Follow up on your contributions
- Help others in the community

## Recognition

Contributors will be acknowledged in:
- CHANGELOG.md for releases
- README.md contributor section
- GitHub contributor statistics

Thank you for contributing to asciicam! 🎥✨