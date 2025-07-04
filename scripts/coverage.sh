#!/bin/bash

# Code coverage script for asciicam
# Generates comprehensive code coverage reports in multiple formats

set -e

COVERAGE_DIR="coverage"
COVERAGE_PROFILE="$COVERAGE_DIR/coverage.out"
COVERAGE_HTML="$COVERAGE_DIR/coverage.html"
COVERAGE_XML="$COVERAGE_DIR/coverage.xml"
COVERAGE_JSON="$COVERAGE_DIR/coverage.json"
COVERAGE_SUMMARY="$COVERAGE_DIR/coverage-summary.txt"
MIN_COVERAGE=60  # Minimum coverage percentage required

echo "🧪 Running code coverage analysis..."

# Create coverage directory
mkdir -p "$COVERAGE_DIR"

# Clean previous coverage data
rm -f "$COVERAGE_PROFILE" "$COVERAGE_HTML" "$COVERAGE_XML" "$COVERAGE_JSON" "$COVERAGE_SUMMARY"

# Run tests with coverage
echo "📊 Running tests with coverage..."
go test -coverprofile="$COVERAGE_PROFILE" -covermode=atomic ./...

if [ ! -f "$COVERAGE_PROFILE" ]; then
    echo "❌ Failed to generate coverage profile"
    exit 1
fi

# Generate HTML report
echo "🌐 Generating HTML coverage report..."
go tool cover -html="$COVERAGE_PROFILE" -o "$COVERAGE_HTML"
echo "   HTML report: $COVERAGE_HTML"

# Generate coverage summary
echo "📋 Generating coverage summary..."
coverage_output=$(go tool cover -func="$COVERAGE_PROFILE")
echo "$coverage_output" > "$COVERAGE_SUMMARY"

# Extract total coverage percentage
total_coverage=$(echo "$coverage_output" | grep "total:" | awk '{print $3}' | sed 's/%//')

if [ -z "$total_coverage" ]; then
    echo "❌ Could not extract coverage percentage"
    exit 1
fi

echo "📈 Total Coverage: ${total_coverage}%"

# Check if coverage meets minimum requirement
if (( $(echo "$total_coverage < $MIN_COVERAGE" | bc -l) )); then
    echo "❌ Coverage ${total_coverage}% is below minimum required ${MIN_COVERAGE}%"
    echo "   Please add more tests to improve coverage."
    exit 1
else
    echo "✅ Coverage ${total_coverage}% meets minimum requirement (${MIN_COVERAGE}%)"
fi

# Generate XML report for CI tools (if gocov-xml is available)
if command -v gocov &> /dev/null && command -v gocov-xml &> /dev/null; then
    echo "📄 Generating XML coverage report..."
    gocov convert "$COVERAGE_PROFILE" | gocov-xml > "$COVERAGE_XML"
    echo "   XML report: $COVERAGE_XML"
fi

# Generate JSON report (if gocov is available)
if command -v gocov &> /dev/null; then
    echo "📊 Generating JSON coverage report..."
    gocov convert "$COVERAGE_PROFILE" > "$COVERAGE_JSON"
    echo "   JSON report: $COVERAGE_JSON"
fi

# Display coverage breakdown by package
echo ""
echo "📦 Coverage by package:"
echo "$coverage_output" | grep -v "total:" | sort -k3 -nr

echo ""
echo "📁 Coverage reports generated in: $COVERAGE_DIR/"
echo "   • HTML: $COVERAGE_HTML"
echo "   • Text: $COVERAGE_SUMMARY"
if [ -f "$COVERAGE_XML" ]; then
    echo "   • XML: $COVERAGE_XML"
fi
if [ -f "$COVERAGE_JSON" ]; then
    echo "   • JSON: $COVERAGE_JSON"
fi

echo ""
echo "💡 To view HTML report: open $COVERAGE_HTML"
echo "💡 To install additional coverage tools:"
echo "   go install github.com/axw/gocov/gocov@latest"
echo "   go install github.com/AlekSi/gocov-xml@latest"