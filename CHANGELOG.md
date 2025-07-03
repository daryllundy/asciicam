# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Project restructure with proper Go project layout
- Internal packages for better code organization:
  - `internal/ascii/` - ASCII conversion logic
  - `internal/camera/` - Camera handling
  - `internal/config/` - Configuration management
  - `internal/greenscreen/` - Greenscreen functionality
- Build automation scripts in `scripts/` directory
- Usage examples in `examples/` directory
- Comprehensive documentation structure
- Contributing guidelines (CONTRIBUTING.md)
- Code of Conduct (CODE_OF_CONDUCT.md)
- Enhanced .gitignore for Go projects
- Build automation with Makefile

### Changed
- Main application moved to `cmd/asciicam/main.go`
- Screenshots moved to `docs/` directory
- Improved code organization and modularity

### Fixed
- [List any bug fixes here]

### Removed
- [List any removed features here]

## [v0.1.0] - 2024-XX-XX

### Added
- Initial release with basic ASCII webcam functionality
- ASCII art conversion from webcam feed
- ANSI color block output mode
- Virtual greenscreen functionality
- Camera device selection
- Zoom control (25%, 50%, 75%, 100%)
- Monochrome color output
- FPS counter
- Background sample generation for greenscreen
- Terminal size auto-detection
- Configurable output dimensions
- Real-time webcam to ASCII conversion

### Features
- **ASCII Mode**: Convert webcam feed to ASCII art
- **ANSI Mode**: High-resolution color blocks using ANSI escape codes
- **Greenscreen**: Virtual background removal with sample-based detection
- **Color Options**: Monochrome output with custom colors
- **Performance**: Real-time processing with FPS display
- **Flexibility**: Configurable camera, output dimensions, and zoom levels

### Technical Details
- Built with Go 1.22+
- OpenCV integration via gocv
- Cross-platform support (macOS, Linux)
- Terminal-based output with escape sequence handling
- Efficient image processing pipeline

## [v0.0.1] - 2024-XX-XX

### Added
- Initial proof of concept
- Basic webcam capture
- Simple ASCII conversion

---

## Release Notes

### Version Numbering
- **Major**: Breaking changes or significant feature additions
- **Minor**: New features, backward compatible
- **Patch**: Bug fixes and small improvements

### Supported Platforms
- macOS (Intel and Apple Silicon)
- Linux (x86_64, ARM64)
- Windows (experimental)

### Dependencies
- Go 1.22 or newer
- OpenCV 4.x
- pkg-config

### Installation
See README.md for detailed installation instructions.

### Migration Guide
For breaking changes, migration guides will be provided in the release notes.

---

**Links:**
- [Keep a Changelog](https://keepachangelog.com/en/1.0.0/)
- [Semantic Versioning](https://semver.org/spec/v2.0.0.html)
- [GitHub Releases](https://github.com/your-username/asciicam/releases)