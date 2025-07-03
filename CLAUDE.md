# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

asciicam is a Go application that displays your webcam feed as ASCII art in the terminal. It uses OpenCV (via gocv) for webcam capture and supports both ASCII and ANSI color output modes.

## Build Commands

- **Build**: `go build`
- **Run**: `./asciicam` (after building)
- **Test**: `go test ./...`
- **Clean**: `go clean`

## Dependencies

- Go 1.23 or newer
- OpenCV 4.x core modules (install with `brew install opencv` on macOS)
- pkg-config must be able to find OpenCV: `pkg-config --cflags --libs opencv4`
- **Note**: Uses gocv v0.35.0 to avoid ArUco compilation issues

## Architecture

The application follows a modular Go project structure:

### cmd/asciicam/main.go
- **Application entry point** with graceful shutdown handling
- **Orchestrates** all internal packages
- **Terminal management** and main rendering loop

### internal/ascii/
- **ASCII/ANSI conversion logic** with Converter struct
- Supports both character-based ASCII and color ANSI block output
- Configurable character sets and color handling

### internal/camera/
- **Camera capture functionality** with Capture struct
- OpenCV integration via gocv for webcam access
- Image format conversion and resizing

### internal/config/
- **Configuration management** with Config struct
- Command-line flag parsing and validation
- Terminal size detection and zoom calculation

### internal/greenscreen/
- **Virtual greenscreen processing** with Processor struct
- Background sample generation and loading
- Color-based background removal using Lab color space

## Key Features

- **Webcam capture** with configurable resolution and device selection
- **ASCII/ANSI output** with automatic terminal size detection
- **Color modes**: Monochrome, truecolor, or ANSI color blocks
- **Zoom control** (25%, 50%, 75%, 100%)
- **Virtual greenscreen** with background sample generation
- **FPS counter** for performance monitoring

## Command Line Options

Key flags for development and testing:
- `-dev=N`: Camera device ID (default: 0)
- `-ansi=true`: Use ANSI color blocks instead of ASCII
- `-gen=true -sample bgdata/`: Generate background samples
- `-greenscreen=true`: Enable virtual greenscreen
- `-fps=true`: Show FPS counter
- `-width=N -height=N`: Override terminal dimensions

## Development Notes

- The application uses terminal escape codes for cursor positioning and screen clearing
- Frame processing includes zoom scaling and optional greenscreen effects
- Background sample generation creates PNG files for greenscreen functionality
- Terminal resize detection updates display dimensions dynamically