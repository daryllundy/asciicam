# Product Overview

asciicam is a real-time webcam-to-ASCII art converter written in Go. It transforms live video feeds into ASCII characters or ANSI color blocks directly in the terminal.

## Core Features

- **Real-time ASCII conversion** - Live webcam feed to ASCII art
- **ANSI color mode** - High-resolution color blocks using ANSI escape codes  
- **Virtual greenscreen** - Background removal with sample-based detection
- **Zoom control** - 25%, 50%, 75%, or 100% zoom levels
- **Color options** - Monochrome output with custom hex colors
- **Auto-sizing** - Automatically detects terminal dimensions
- **Performance monitoring** - Built-in FPS counter

## Target Platforms

- Primary: macOS (Intel and Apple Silicon), Linux (x86_64, ARM64)
- Experimental: Windows support

## Dependencies

- Go 1.22+
- OpenCV 4.x for webcam capture
- pkg-config for OpenCV detection
- Modern terminal with ANSI support (recommended)
