#!/bin/bash

# Test script for GitHub Actions workflow validation
# This script validates workflow syntax, dependencies, and configuration

set -e

echo "🧪 Testing GitHub Actions Workflow Functionality"
echo "================================================"

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

# Test 1: Validate workflow syntax
test_workflow_syntax() {
    print_status "INFO" "Testing workflow syntax validation..."
    
    # Check if actionlint is available
    if command_exists actionlint; then
        print_status "INFO" "Running actionlint on workflow files..."
        if actionlint .github/workflows/*.yml; then
            print_status "SUCCESS" "All workflow files have valid syntax"
        else
            print_status "ERROR" "Workflow syntax validation failed"
            return 1
        fi
    else
        print_status "WARNING" "actionlint not found, skipping syntax validation"
        print_status "INFO" "Install actionlint with: go install github.com/rhymond/actionlint/cmd/actionlint@latest"
    fi
}

# Test 2: Validate matrix configuration
test_matrix_configuration() {
    print_status "INFO" "Testing matrix build configuration..."
    
    # Check if CI workflow has proper matrix setup
    if grep -q "go-version.*1\.22.*1\.23" .github/workflows/ci.yml && \
       grep -q "os.*ubuntu-latest.*macos-latest" .github/workflows/ci.yml; then
        print_status "SUCCESS" "Matrix configuration includes required Go versions and platforms"
    else
        print_status "ERROR" "Matrix configuration missing required Go versions (1.22, 1.23) or platforms"
        return 1
    fi
    
    # Check for proper matrix artifact naming
    if grep -q "coverage-report-.*matrix\.os.*matrix\.go-version" .github/workflows/ci.yml; then
        print_status "SUCCESS" "Matrix builds use unique artifact names"
    else
        print_status "ERROR" "Matrix builds may have conflicting artifact names"
        return 1
    fi
}

# Test 3: Validate OpenCV setup action
test_opencv_setup() {
    print_status "INFO" "Testing OpenCV setup action..."
    
    # Check if OpenCV action exists
    if [ -f ".github/actions/setup-opencv/action.yml" ]; then
        print_status "SUCCESS" "OpenCV setup action exists"
        
        # Check for OS-specific steps
        if grep -q "runner\.os.*Linux" .github/actions/setup-opencv/action.yml && \
           grep -q "runner\.os.*macOS" .github/actions/setup-opencv/action.yml; then
            print_status "SUCCESS" "OpenCV action handles both Linux and macOS"
        else
            print_status "ERROR" "OpenCV action missing OS-specific handling"
            return 1
        fi
        
        # Check for pkg-config setup
        if grep -q "pkg-config" .github/actions/setup-opencv/action.yml; then
            print_status "SUCCESS" "OpenCV action includes pkg-config setup"
        else
            print_status "ERROR" "OpenCV action missing pkg-config configuration"
            return 1
        fi
    else
        print_status "ERROR" "OpenCV setup action not found"
        return 1
    fi
}

# Test 4: Validate caching configuration
test_caching_setup() {
    print_status "INFO" "Testing caching configuration..."
    
    # Check for Go module caching
    if grep -q "cache.*go-build\|cache.*pkg/mod" .github/workflows/ci.yml; then
        print_status "SUCCESS" "Go build and module caching configured"
    else
        print_status "ERROR" "Go caching not properly configured"
        return 1
    fi
    
    # Check for golangci-lint caching
    if grep -q "cache.*golangci-lint" .github/workflows/ci.yml; then
        print_status "SUCCESS" "golangci-lint caching configured"
    else
        print_status "ERROR" "golangci-lint caching not configured"
        return 1
    fi
}

# Test 5: Validate coverage reporting
test_coverage_reporting() {
    print_status "INFO" "Testing coverage reporting configuration..."
    
    # Check for coverage generation
    if grep -q "coverprofile.*coverage\.out" .github/workflows/ci.yml && \
       grep -q "cover.*html.*coverage\.html" .github/workflows/ci.yml; then
        print_status "SUCCESS" "Coverage generation configured for both text and HTML"
    else
        print_status "ERROR" "Coverage generation not properly configured"
        return 1
    fi
    
    # Check for artifact upload
    if grep -q "upload-artifact" .github/workflows/ci.yml && \
       grep -q "coverage\.out\|coverage\.html" .github/workflows/ci.yml; then
        print_status "SUCCESS" "Coverage artifacts upload configured"
    else
        print_status "ERROR" "Coverage artifact upload not configured"
        return 1
    fi
}

# Test 6: Validate linting configuration
test_linting_setup() {
    print_status "INFO" "Testing linting configuration..."
    
    # Check for strict linting job
    if grep -q "lint-strict" .github/workflows/ci.yml; then
        print_status "SUCCESS" "Strict linting job configured"
    else
        print_status "ERROR" "Strict linting job not found"
        return 1
    fi
    
    # Check for soft linting job
    if grep -q "lint-soft" .github/workflows/ci.yml && \
       grep -q "continue-on-error.*true" .github/workflows/ci.yml; then
        print_status "SUCCESS" "Soft linting job configured with continue-on-error"
    else
        print_status "ERROR" "Soft linting job not properly configured"
        return 1
    fi
    
    # Check for PR-only soft linting
    if grep -q "if.*github\.event_name.*pull_request" .github/workflows/ci.yml; then
        print_status "SUCCESS" "Soft linting runs only on pull requests"
    else
        print_status "ERROR" "Soft linting not restricted to pull requests"
        return 1
    fi
}

# Test 7: Validate GitLab mirror workflow
test_gitlab_mirror() {
    print_status "INFO" "Testing GitLab mirror workflow..."
    
    # Check if mirror workflow exists
    if [ -f ".github/workflows/mirror-to-gitlab.yml" ]; then
        print_status "SUCCESS" "GitLab mirror workflow exists"
        
        # Check for secret validation
        if grep -q "GITLAB_USERNAME\|GITLAB_TOKEN" .github/workflows/mirror-to-gitlab.yml; then
            print_status "SUCCESS" "GitLab mirror includes secret validation"
        else
            print_status "ERROR" "GitLab mirror missing secret validation"
            return 1
        fi
        
        # Check for updated checkout action
        if grep -q "actions/checkout@v4" .github/workflows/mirror-to-gitlab.yml; then
            print_status "SUCCESS" "GitLab mirror uses updated checkout action"
        else
            print_status "WARNING" "GitLab mirror may be using outdated checkout action"
        fi
    else
        print_status "ERROR" "GitLab mirror workflow not found"
        return 1
    fi
}

# Test 8: Validate badge update workflow
test_badge_update() {
    print_status "INFO" "Testing badge update workflow..."
    
    # Check if badge workflow exists
    if [ -f ".github/workflows/update-readme-badge.yml" ]; then
        print_status "SUCCESS" "Badge update workflow exists"
        
        # Check for workflow dependency
        if grep -q "workflow_run" .github/workflows/update-readme-badge.yml; then
            print_status "SUCCESS" "Badge update triggered by mirror workflow"
        else
            print_status "ERROR" "Badge update not properly triggered by mirror workflow"
            return 1
        fi
        
        # Check for updated checkout action
        if grep -q "actions/checkout@v4" .github/workflows/update-readme-badge.yml; then
            print_status "SUCCESS" "Badge update uses updated checkout action"
        else
            print_status "WARNING" "Badge update may be using outdated checkout action"
        fi
    else
        print_status "ERROR" "Badge update workflow not found"
        return 1
    fi
}

# Test 9: Validate permissions configuration
test_permissions() {
    print_status "INFO" "Testing workflow permissions..."
    
    # Check CI workflow permissions
    if grep -q "permissions:" .github/workflows/ci.yml; then
        print_status "SUCCESS" "CI workflow has permissions configured"
        
        # Check for minimal permissions
        if grep -A5 "permissions:" .github/workflows/ci.yml | grep -q "contents.*read"; then
            print_status "SUCCESS" "CI workflow uses minimal read permissions"
        else
            print_status "WARNING" "CI workflow permissions may be too broad"
        fi
    else
        print_status "ERROR" "CI workflow missing permissions configuration"
        return 1
    fi
}

# Test 10: Validate action versions
test_action_versions() {
    print_status "INFO" "Testing action version pinning..."
    
    # Check for pinned action versions
    local unpinned_actions=0
    
    # Common actions that should be pinned
    for workflow in .github/workflows/*.yml; do
        if grep -q "uses:.*@main\|uses:.*@master" "$workflow"; then
            print_status "WARNING" "Found unpinned actions in $workflow"
            unpinned_actions=$((unpinned_actions + 1))
        fi
    done
    
    if [ $unpinned_actions -eq 0 ]; then
        print_status "SUCCESS" "All actions appear to be version pinned"
    else
        print_status "WARNING" "Some actions may not be properly version pinned"
    fi
}

# Main test execution
main() {
    local failed_tests=0
    
    echo "Starting workflow validation tests..."
    echo
    
    # Run all tests
    test_workflow_syntax || failed_tests=$((failed_tests + 1))
    test_matrix_configuration || failed_tests=$((failed_tests + 1))
    test_opencv_setup || failed_tests=$((failed_tests + 1))
    test_caching_setup || failed_tests=$((failed_tests + 1))
    test_coverage_reporting || failed_tests=$((failed_tests + 1))
    test_linting_setup || failed_tests=$((failed_tests + 1))
    test_gitlab_mirror || failed_tests=$((failed_tests + 1))
    test_badge_update || failed_tests=$((failed_tests + 1))
    test_permissions || failed_tests=$((failed_tests + 1))
    test_action_versions || failed_tests=$((failed_tests + 1))
    
    echo
    echo "================================================"
    if [ $failed_tests -eq 0 ]; then
        print_status "SUCCESS" "All workflow validation tests passed!"
        echo
        print_status "INFO" "Workflows are ready for testing. Consider:"
        print_status "INFO" "1. Creating a test branch and pushing to trigger CI"
        print_status "INFO" "2. Creating a pull request to test soft linting"
        print_status "INFO" "3. Configuring GitLab secrets for mirror testing"
        return 0
    else
        print_status "ERROR" "$failed_tests test(s) failed"
        echo
        print_status "INFO" "Please fix the issues above before deploying workflows"
        return 1
    fi
}

# Run main function
main "$@"
