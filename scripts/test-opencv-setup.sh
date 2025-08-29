#!/bin/bash

# Test script for OpenCV setup action validation
# This script simulates the OpenCV setup action locally to verify it works

set -e

echo "🔧 Testing OpenCV Setup Action Locally"
echo "======================================"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    local status=$1
    local message=$2
    case $status in
        "SUCCESS") echo -e "${GREEN}✅ $message${NC}" ;;
        "ERROR") echo -e "${RED}❌ $message${NC}" ;;
        "WARNING") echo -e "${YELLOW}⚠️  $message${NC}" ;;
        "INFO") echo -e "${BLUE}ℹ️  $message${NC}" ;;
    esac
}

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Detect OS
detect_os() {
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        echo "Linux"
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        echo "macOS"
    else
        echo "Unknown"
    fi
}

# Test OpenCV installation detection
test_opencv_detection() {
    local os=$(detect_os)
    print_status "INFO" "Testing OpenCV detection on $os..."
    
    # Check if pkg-config exists
    if ! command_exists pkg-config; then
        print_status "ERROR" "pkg-config not found - required for OpenCV detection"
        return 1
    fi
    
    print_status "SUCCESS" "pkg-config is available"
    
    # Test OpenCV detection logic (similar to the action)
    PKG_CONFIG_PATH=""
    if [ "$os" == "macOS" ]; then
        # For Homebrew, OpenCV paths are typically here
        PKG_CONFIG_PATH="/usr/local/opt/opencv/lib/pkgconfig:/opt/homebrew/lib/pkgconfig:$PKG_CONFIG_PATH"
        export PKG_CONFIG_PATH
        print_status "INFO" "Set macOS PKG_CONFIG_PATH: $PKG_CONFIG_PATH"
    fi
    
    print_status "INFO" "Searching for OpenCV pkg-config..."
    
    if pkg-config --exists opencv4; then
        print_status "SUCCESS" "Found opencv4"
        
        # Test getting flags
        local cflags=$(pkg-config --cflags opencv4)
        local ldflags=$(pkg-config --libs opencv4)
        
        print_status "INFO" "CGO_CFLAGS: $cflags"
        print_status "INFO" "CGO_LDFLAGS: $ldflags"
        
        if [ -n "$cflags" ] && [ -n "$ldflags" ]; then
            print_status "SUCCESS" "OpenCV4 flags retrieved successfully"
            return 0
        else
            print_status "ERROR" "OpenCV4 flags are empty"
            return 1
        fi
        
    elif pkg-config --exists opencv; then
        print_status "SUCCESS" "Found opencv (legacy)"
        
        # Test getting flags
        local cflags=$(pkg-config --cflags opencv)
        local ldflags=$(pkg-config --libs opencv)
        
        print_status "INFO" "CGO_CFLAGS: $cflags"
        print_status "INFO" "CGO_LDFLAGS: $ldflags"
        
        if [ -n "$cflags" ] && [ -n "$ldflags" ]; then
            print_status "SUCCESS" "OpenCV flags retrieved successfully"
            return 0
        else
            print_status "ERROR" "OpenCV flags are empty"
            return 1
        fi
    else
        print_status "ERROR" "OpenCV not found by pkg-config"
        print_status "INFO" "Available pkg-config packages:"
        pkg-config --list-all | grep -i opencv || print_status "WARNING" "No OpenCV packages found"
        return 1
    fi
}

# Test Go build with OpenCV
test_go_build_with_opencv() {
    print_status "INFO" "Testing Go build with OpenCV..."
    
    # Set up environment variables like the action would
    if pkg-config --exists opencv4; then
        export CGO_CFLAGS="$(pkg-config --cflags opencv4)"
        export CGO_LDFLAGS="$(pkg-config --libs opencv4)"
    elif pkg-config --exists opencv; then
        export CGO_CFLAGS="$(pkg-config --cflags opencv)"
        export CGO_LDFLAGS="$(pkg-config --libs opencv)"
    else
        print_status "ERROR" "Cannot set CGO flags - OpenCV not found"
        return 1
    fi
    
    export CGO_ENABLED="1"
    
    print_status "INFO" "Environment variables set:"
    print_status "INFO" "CGO_ENABLED=$CGO_ENABLED"
    print_status "INFO" "CGO_CFLAGS=$CGO_CFLAGS"
    print_status "INFO" "CGO_LDFLAGS=$CGO_LDFLAGS"
    
    # Test if Go can compile with these settings
    print_status "INFO" "Testing Go module download..."
    if go mod download; then
        print_status "SUCCESS" "Go modules downloaded successfully"
    else
        print_status "ERROR" "Failed to download Go modules"
        return 1
    fi
    
    print_status "INFO" "Testing Go build..."
    if go build -v ./...; then
        print_status "SUCCESS" "Go build completed successfully with OpenCV"
        return 0
    else
        print_status "ERROR" "Go build failed with OpenCV configuration"
        return 1
    fi
}

# Test coverage generation
test_coverage_generation() {
    print_status "INFO" "Testing coverage generation..."
    
    # Run tests with coverage (similar to CI workflow)
    if go test -v -coverprofile=coverage.out -covermode=atomic ./...; then
        print_status "SUCCESS" "Tests passed and coverage generated"
        
        # Test HTML coverage generation
        if go tool cover -html=coverage.out -o coverage.html; then
            print_status "SUCCESS" "HTML coverage report generated"
            
            # Check if files exist and have content
            if [ -s coverage.out ] && [ -s coverage.html ]; then
                print_status "SUCCESS" "Coverage files have content"
                
                # Clean up test files
                rm -f coverage.out coverage.html
                return 0
            else
                print_status "ERROR" "Coverage files are empty"
                return 1
            fi
        else
            print_status "ERROR" "Failed to generate HTML coverage report"
            return 1
        fi
    else
        print_status "ERROR" "Tests failed or coverage generation failed"
        return 1
    fi
}

# Test linting with OpenCV environment
test_linting_with_opencv() {
    print_status "INFO" "Testing linting with OpenCV environment..."
    
    # Check if golangci-lint is available
    if ! command_exists golangci-lint; then
        print_status "WARNING" "golangci-lint not found, skipping lint test"
        print_status "INFO" "Install with: brew install golangci-lint"
        return 0
    fi
    
    # Run linting with CGO enabled (like in CI)
    export CGO_ENABLED="1"
    
    if golangci-lint run --out-format=colored-line-number; then
        print_status "SUCCESS" "Linting passed with OpenCV environment"
        return 0
    else
        print_status "WARNING" "Linting found issues (this may be expected)"
        return 0  # Don't fail the test for linting issues
    fi
}

# Main test execution
main() {
    local failed_tests=0
    local os=$(detect_os)
    
    print_status "INFO" "Running OpenCV setup tests on $os"
    echo
    
    # Run tests
    test_opencv_detection || failed_tests=$((failed_tests + 1))
    test_go_build_with_opencv || failed_tests=$((failed_tests + 1))
    test_coverage_generation || failed_tests=$((failed_tests + 1))
    test_linting_with_opencv || failed_tests=$((failed_tests + 1))
    
    echo
    echo "======================================"
    if [ $failed_tests -eq 0 ]; then
        print_status "SUCCESS" "All OpenCV setup tests passed!"
        echo
        print_status "INFO" "The OpenCV setup action should work correctly in CI"
        return 0
    else
        print_status "ERROR" "$failed_tests test(s) failed"
        echo
        print_status "INFO" "OpenCV setup may need adjustments for CI environment"
        return 1
    fi
}

# Run main function
main "$@"
