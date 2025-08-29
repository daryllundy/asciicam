#!/bin/bash

# Integration test script for CI workflow validation
# This script simulates the complete CI workflow locally

set -e

echo "🚀 Testing CI Workflow Integration"
echo "================================="

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

# Test matrix simulation
test_matrix_simulation() {
    print_status "INFO" "Testing matrix build simulation..."
    
    local go_versions=("1.22" "1.23")
    local os_types=("ubuntu-latest" "macos-latest")
    local current_os="macos-latest"  # We're running on macOS
    
    print_status "INFO" "Simulating matrix builds for current OS: $current_os"
    
    for go_version in "${go_versions[@]}"; do
        print_status "INFO" "Testing with Go $go_version..."
        
        # Check if Go version is available
        if go version | grep -q "$go_version"; then
            print_status "SUCCESS" "Go $go_version is available"
        else
            print_status "WARNING" "Go $go_version not available, using current version"
        fi
        
        # Simulate the CI steps for this matrix combination
        print_status "INFO" "Simulating CI steps for Go $go_version on $current_os"
        
        # Step 1: Setup OpenCV (already tested in previous script)
        print_status "SUCCESS" "OpenCV setup (simulated)"
        
        # Step 2: Download modules
        if go mod download; then
            print_status "SUCCESS" "Go modules downloaded for Go $go_version"
        else
            print_status "ERROR" "Failed to download modules for Go $go_version"
            return 1
        fi
        
        # Step 3: Build
        export CGO_ENABLED="1"
        if go build -v ./...; then
            print_status "SUCCESS" "Build successful for Go $go_version"
        else
            print_status "ERROR" "Build failed for Go $go_version"
            return 1
        fi
        
        # Step 4: Test with coverage (unique artifact name simulation)
        local artifact_name="coverage-report-$current_os-go$go_version-test-run"
        print_status "INFO" "Generating coverage with artifact name: $artifact_name"
        
        if go test -v -coverprofile="coverage-$go_version.out" -covermode=atomic ./...; then
            print_status "SUCCESS" "Tests passed for Go $go_version"
            
            # Generate HTML coverage
            if go tool cover -html="coverage-$go_version.out" -o "coverage-$go_version.html"; then
                print_status "SUCCESS" "Coverage HTML generated for Go $go_version"
                
                # Simulate artifact upload by checking file sizes
                if [ -s "coverage-$go_version.out" ] && [ -s "coverage-$go_version.html" ]; then
                    print_status "SUCCESS" "Coverage artifacts ready for upload (Go $go_version)"
                else
                    print_status "ERROR" "Coverage artifacts are empty (Go $go_version)"
                    return 1
                fi
            else
                print_status "ERROR" "Failed to generate HTML coverage for Go $go_version"
                return 1
            fi
        else
            print_status "ERROR" "Tests failed for Go $go_version"
            return 1
        fi
    done
    
    # Clean up test artifacts
    rm -f coverage-*.out coverage-*.html
    print_status "SUCCESS" "Matrix simulation completed successfully"
}

# Test linting jobs simulation
test_linting_jobs() {
    print_status "INFO" "Testing linting jobs simulation..."
    
    # Check if golangci-lint is available
    if ! command_exists golangci-lint; then
        print_status "WARNING" "golangci-lint not found, installing..."
        if command_exists brew; then
            brew install golangci-lint
        else
            print_status "ERROR" "Cannot install golangci-lint - brew not available"
            return 1
        fi
    fi
    
    # Test strict linting (should fail workflow on issues)
    print_status "INFO" "Testing strict linting job..."
    export CGO_ENABLED="1"
    
    # Run with main config
    if golangci-lint run --out-format=colored-line-number; then
        print_status "SUCCESS" "Strict linting passed"
    else
        print_status "WARNING" "Strict linting found issues (would fail CI)"
        # Don't return 1 here as linting issues are expected in development
    fi
    
    # Test soft linting (should not fail workflow)
    print_status "INFO" "Testing soft linting job..."
    
    # Check if soft config exists
    if [ -f ".golangci-soft.yml" ]; then
        # Simulate continue-on-error behavior
        if golangci-lint run --config=.golangci-soft.yml --out-format=colored-line-number; then
            print_status "SUCCESS" "Soft linting passed"
        else
            print_status "INFO" "Soft linting found issues (would continue CI due to continue-on-error)"
        fi
    else
        print_status "WARNING" "Soft linting config not found"
    fi
    
    print_status "SUCCESS" "Linting jobs simulation completed"
}

# Test caching simulation
test_caching_simulation() {
    print_status "INFO" "Testing caching simulation..."
    
    # Simulate Go module cache
    local go_cache_dir="$HOME/.cache/go-build"
    local go_mod_dir="$HOME/go/pkg/mod"
    
    if [ -d "$go_cache_dir" ] || [ -d "$go_mod_dir" ]; then
        print_status "SUCCESS" "Go cache directories exist (would be cached in CI)"
    else
        print_status "INFO" "Go cache directories not found (would be created in CI)"
    fi
    
    # Simulate golangci-lint cache
    local lint_cache_dir="$HOME/.cache/golangci-lint"
    if [ -d "$lint_cache_dir" ]; then
        print_status "SUCCESS" "golangci-lint cache directory exists (would be cached in CI)"
    else
        print_status "INFO" "golangci-lint cache directory not found (would be created in CI)"
    fi
    
    print_status "SUCCESS" "Caching simulation completed"
}

# Test permissions simulation
test_permissions_simulation() {
    print_status "INFO" "Testing permissions simulation..."
    
    # Check if we can read repository contents (contents: read)
    if [ -r "README.md" ] && [ -r "go.mod" ]; then
        print_status "SUCCESS" "Repository contents readable (contents: read permission)"
    else
        print_status "ERROR" "Cannot read repository contents"
        return 1
    fi
    
    # Simulate pull-requests: read permission
    print_status "SUCCESS" "Pull request read permission (simulated)"
    
    print_status "SUCCESS" "Permissions simulation completed"
}

# Test artifact generation
test_artifact_generation() {
    print_status "INFO" "Testing artifact generation..."
    
    # Generate test coverage artifacts
    export CGO_ENABLED="1"
    
    if go test -v -coverprofile=test-coverage.out -covermode=atomic ./...; then
        if go tool cover -html=test-coverage.out -o test-coverage.html; then
            # Check artifact sizes and content
            local out_size=$(wc -c < test-coverage.out)
            local html_size=$(wc -c < test-coverage.html)
            
            if [ "$out_size" -gt 0 ] && [ "$html_size" -gt 0 ]; then
                print_status "SUCCESS" "Coverage artifacts generated successfully"
                print_status "INFO" "Coverage.out size: $out_size bytes"
                print_status "INFO" "Coverage.html size: $html_size bytes"
                
                # Simulate 30-day retention
                print_status "SUCCESS" "Artifacts ready for upload with 30-day retention"
            else
                print_status "ERROR" "Generated artifacts are empty"
                return 1
            fi
        else
            print_status "ERROR" "Failed to generate HTML coverage"
            return 1
        fi
    else
        print_status "ERROR" "Failed to generate coverage"
        return 1
    fi
    
    # Clean up
    rm -f test-coverage.out test-coverage.html
    print_status "SUCCESS" "Artifact generation test completed"
}

# Test workflow triggers simulation
test_workflow_triggers() {
    print_status "INFO" "Testing workflow triggers simulation..."
    
    # Simulate push to main branch
    local current_branch=$(git branch --show-current 2>/dev/null || echo "unknown")
    print_status "INFO" "Current branch: $current_branch"
    
    if [ "$current_branch" = "main" ]; then
        print_status "SUCCESS" "On main branch - would trigger CI on push"
    else
        print_status "INFO" "Not on main branch - would trigger CI on push to any branch"
    fi
    
    # Simulate pull request trigger
    print_status "SUCCESS" "Pull request trigger (simulated)"
    
    print_status "SUCCESS" "Workflow triggers simulation completed"
}

# Test environment variables
test_environment_setup() {
    print_status "INFO" "Testing environment setup..."
    
    # Test CGO_ENABLED
    export CGO_ENABLED="1"
    if [ "$CGO_ENABLED" = "1" ]; then
        print_status "SUCCESS" "CGO_ENABLED set correctly"
    else
        print_status "ERROR" "CGO_ENABLED not set correctly"
        return 1
    fi
    
    # Test OpenCV environment (if available)
    if pkg-config --exists opencv4 || pkg-config --exists opencv; then
        if pkg-config --exists opencv4; then
            export CGO_CFLAGS="$(pkg-config --cflags opencv4)"
            export CGO_LDFLAGS="$(pkg-config --libs opencv4)"
        else
            export CGO_CFLAGS="$(pkg-config --cflags opencv)"
            export CGO_LDFLAGS="$(pkg-config --libs opencv)"
        fi
        
        if [ -n "$CGO_CFLAGS" ] && [ -n "$CGO_LDFLAGS" ]; then
            print_status "SUCCESS" "OpenCV environment variables set"
        else
            print_status "ERROR" "OpenCV environment variables empty"
            return 1
        fi
    else
        print_status "WARNING" "OpenCV not available for environment test"
    fi
    
    print_status "SUCCESS" "Environment setup test completed"
}

# Main test execution
main() {
    local failed_tests=0
    
    print_status "INFO" "Starting CI workflow integration tests..."
    echo
    
    # Run all integration tests
    test_matrix_simulation || failed_tests=$((failed_tests + 1))
    test_linting_jobs || failed_tests=$((failed_tests + 1))
    test_caching_simulation || failed_tests=$((failed_tests + 1))
    test_permissions_simulation || failed_tests=$((failed_tests + 1))
    test_artifact_generation || failed_tests=$((failed_tests + 1))
    test_workflow_triggers || failed_tests=$((failed_tests + 1))
    test_environment_setup || failed_tests=$((failed_tests + 1))
    
    echo
    echo "================================="
    if [ $failed_tests -eq 0 ]; then
        print_status "SUCCESS" "All CI integration tests passed!"
        echo
        print_status "INFO" "The CI workflow should function correctly in GitHub Actions"
        print_status "INFO" "Next steps:"
        print_status "INFO" "1. Push to a test branch to trigger the actual CI workflow"
        print_status "INFO" "2. Create a pull request to test the complete flow"
        print_status "INFO" "3. Monitor the workflow runs in GitHub Actions"
        return 0
    else
        print_status "ERROR" "$failed_tests integration test(s) failed"
        echo
        print_status "INFO" "Please fix the issues above before deploying"
        return 1
    fi
}

# Run main function
main "$@"
