#!/bin/bash
# Cross-compilation build script for asciicam

set -e

echo "Cross-compiling asciicam for multiple platforms..."

# Create build output directory
mkdir -p build

# Build for different platforms
platforms=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
)

for platform in "${platforms[@]}"; do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    
    output_name="asciicam-${GOOS}-${GOARCH}"
    if [ $GOOS = "windows" ]; then
        output_name+=".exe"
    fi
    
    echo "Building for $GOOS/$GOARCH..."
    env GOOS=$GOOS GOARCH=$GOARCH go build -o build/$output_name ./cmd/asciicam
    
    if [ $? -ne 0 ]; then
        echo "Failed to build for $GOOS/$GOARCH"
        exit 1
    fi
done

echo "Cross-compilation complete! Binaries created in build/ directory"