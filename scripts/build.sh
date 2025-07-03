#!/bin/bash
# Build script for asciicam

set -e

echo "Building asciicam..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed"
    exit 1
fi

# Check if OpenCV is available
if ! pkg-config --exists opencv4; then
    echo "Error: OpenCV4 not found. Please install OpenCV:"
    echo "  macOS: brew install opencv"
    echo "  Linux: sudo apt-get install libopencv-dev"
    exit 1
fi

# Build the application
echo "Building for current platform..."
go build -o asciicam ./cmd/asciicam

echo "Build complete! Binary created: asciicam"