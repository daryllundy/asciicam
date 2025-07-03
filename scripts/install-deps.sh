#!/bin/bash
# Dependency installation script for asciicam

set -e

echo "Installing dependencies for asciicam..."

# Detect OS
OS="$(uname -s)"
case "${OS}" in
    Darwin*)
        echo "Detected macOS"
        if command -v brew &> /dev/null; then
            echo "Installing OpenCV via Homebrew..."
            brew install opencv
        else
            echo "Error: Homebrew not found. Please install Homebrew first:"
            echo "  /bin/bash -c \"\$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\""
            exit 1
        fi
        ;;
    Linux*)
        echo "Detected Linux"
        if command -v apt-get &> /dev/null; then
            echo "Installing OpenCV via apt-get..."
            sudo apt-get update
            sudo apt-get install -y libopencv-dev pkg-config
        elif command -v yum &> /dev/null; then
            echo "Installing OpenCV via yum..."
            sudo yum install -y opencv-devel pkgconfig
        else
            echo "Error: Neither apt-get nor yum found. Please install OpenCV manually."
            exit 1
        fi
        ;;
    *)
        echo "Unsupported OS: ${OS}"
        exit 1
        ;;
esac

# Verify OpenCV installation
if pkg-config --exists opencv4; then
    echo "OpenCV installation verified successfully!"
    pkg-config --modversion opencv4
else
    echo "Error: OpenCV installation verification failed"
    exit 1
fi

# Install Go dependencies
echo "Installing Go dependencies..."
go mod download
go mod verify

echo "All dependencies installed successfully!"