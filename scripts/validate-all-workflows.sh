#!/bin/bash

# Master validation script for all GitHub Actions workflows
# This script runs all validation tests and generates a comprehensive report

set -e

echo "🔍 Comprehensive GitHub Actions Workflow Validation"
echo "=================================================="

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

# Function to run a test script and capture results
run_test() {
    local test_name=$1
    local script_path=$2
    local log_file="validation-${test_name}.log"
    
    print_status "INFO" "Running $test_name validation..."
    
    if [ -f "$script_path" ]; then
        if bash "$script_path" > "$log_file" 2>&1; then
            print_status "SUCCESS" "$test_name validation passed"
            return 0
        else
            print_status "ERROR" "$test_name validation failed (see $log_file)"
            return 1
        fi
    else
        print_status "ERROR" "$test_name script not found: $script_path"
        return 1
    fi
}

# Main validation function
main() {
    local start_time=$(date +%s)
    local failed_tests=0
    local total_tests=0
    
    print_status "INFO" "Starting comprehensive workflow validation..."
    print_status "INFO" "Test environment: $(uname -s) $(uname -m)"
    print_status "INFO" "Go version: $(go version)"
    echo
    
    # Test 1: Basic workflow validation
    total_tests=$((total_tests + 1))
    if ! run_test "workflow-syntax" "./scripts/test-workflows.sh"; then
        failed_tests=$((failed_tests + 1))
    fi
    
    # Test 2: OpenCV setup validation
    total_tests=$((total_tests + 1))
    if ! run_test "opencv-setup" "./scripts/test-opencv-setup.sh"; then
        failed_tests=$((failed_tests + 1))
    fi
    
    # Test 3: CI integration validation
    total_tests=$((total_tests + 1))
    if ! run_test "ci-integration" "./scripts/test-ci-integration.sh"; then
        failed_tests=$((failed_tests + 1))
    fi
    
    # Calculate test duration
    local end_time=$(date +%s)
    local duration=$((end_time - start_time))
    
    echo
    echo "=================================================="
    print_status "INFO" "Validation Summary"
    echo "=================================================="
    print_status "INFO" "Total tests run: $total_tests"
    print_status "INFO" "Tests passed: $((total_tests - failed_tests))"
    print_status "INFO" "Tests failed: $failed_tests"
    print_status "INFO" "Duration: ${duration}s"
    echo
    
    if [ $failed_tests -eq 0 ]; then
        print_status "SUCCESS" "All workflow validations passed!"
        echo
        print_status "INFO" "🚀 Workflows are ready for deployment"
        print_status "INFO" "📋 See scripts/workflow-validation-report.md for detailed results"
        echo
        print_status "INFO" "Next steps:"
        print_status "INFO" "1. Push to a test branch to trigger actual CI workflows"
        print_status "INFO" "2. Create a pull request to test the complete flow"
        print_status "INFO" "3. Configure GitLab secrets (GITLAB_USERNAME, GITLAB_TOKEN)"
        print_status "INFO" "4. Monitor first workflow runs in GitHub Actions"
        
        # Generate success report
        generate_success_report "$duration"
        
        return 0
    else
        print_status "ERROR" "Workflow validation failed!"
        echo
        print_status "INFO" "❌ $failed_tests test(s) failed"
        print_status "INFO" "📋 Check individual log files for details:"
        for log in validation-*.log; do
            if [ -f "$log" ]; then
                print_status "INFO" "   - $log"
            fi
        done
        
        return 1
    fi
}

# Generate success report
generate_success_report() {
    local duration=$1
    local report_file="validation-success-$(date +%Y%m%d-%H%M%S).md"
    
    cat > "$report_file" << EOF
# Workflow Validation Success Report

**Date**: $(date)
**Duration**: ${duration}s
**Environment**: $(uname -s) $(uname -m)
**Go Version**: $(go version)

## Test Results

✅ **Workflow Syntax Validation**: PASSED
✅ **OpenCV Setup Validation**: PASSED  
✅ **CI Integration Validation**: PASSED

## Validated Components

### CI Workflow (.github/workflows/ci.yml)
- Matrix builds for Go 1.22, 1.23 on ubuntu-latest, macos-latest
- OpenCV setup action integration
- Build, test, and coverage generation
- Strict and soft linting jobs
- Proper caching configuration
- Artifact upload with retention

### OpenCV Setup Action (.github/actions/setup-opencv/action.yml)
- Cross-platform OpenCV installation (Ubuntu, macOS)
- pkg-config detection and configuration
- CGO environment variable setup
- Error handling and fallback mechanisms

### GitLab Integration Workflows
- Mirror workflow with secret validation
- Badge update workflow with dependency triggers
- Updated to latest action versions
- Comprehensive error handling

### Security and Best Practices
- Minimal required permissions
- Pinned action versions
- Proper secret handling
- Caching optimization

## Deployment Readiness

🚀 **READY FOR DEPLOYMENT**

All workflows have been validated and are ready for production use.

## Next Steps

1. Push to test branch to trigger CI
2. Create pull request to test complete flow
3. Configure GitLab secrets if needed
4. Monitor first workflow runs

---
Generated by: scripts/validate-all-workflows.sh
EOF

    print_status "SUCCESS" "Success report generated: $report_file"
}

# Cleanup function
cleanup() {
    print_status "INFO" "Cleaning up temporary files..."
    rm -f validation-*.log 2>/dev/null || true
}

# Set up cleanup on exit
trap cleanup EXIT

# Run main function
main "$@"
