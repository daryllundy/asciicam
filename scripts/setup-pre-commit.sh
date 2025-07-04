#!/bin/bash

# Setup script for pre-commit hooks
# This script installs and configures pre-commit hooks for the asciicam project

set -e

echo "🔧 Setting up pre-commit hooks for asciicam..."

# Check if pre-commit is installed
if ! command -v pre-commit &> /dev/null; then
    echo "📦 Installing pre-commit..."
    
    # Try different installation methods based on available package managers
    if command -v pip3 &> /dev/null; then
        pip3 install pre-commit
    elif command -v pip &> /dev/null; then
        pip install pre-commit
    elif command -v brew &> /dev/null; then
        brew install pre-commit
    elif command -v apt-get &> /dev/null; then
        sudo apt-get update
        sudo apt-get install -y python3-pip
        pip3 install pre-commit
    else
        echo "❌ Could not find a package manager to install pre-commit"
        echo "Please install pre-commit manually: https://pre-commit.com/#installation"
        exit 1
    fi
else
    echo "✅ pre-commit is already installed"
fi

# Check if golangci-lint is installed
if ! command -v golangci-lint &> /dev/null; then
    echo "📦 Installing golangci-lint..."
    
    # Install golangci-lint
    if command -v brew &> /dev/null; then
        brew install golangci-lint
    else
        # Use the official installer
        curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin latest
    fi
else
    echo "✅ golangci-lint is already installed"
fi

# Install Go tools if not present
echo "📦 Installing required Go tools..."
go install golang.org/x/tools/cmd/goimports@latest
go install mvdan.cc/gofumpt@latest

# Install the pre-commit hook
echo "🔗 Installing pre-commit hooks..."
cd "$(dirname "$0")/.."
pre-commit install

# Run the hooks on all files to ensure everything is set up correctly
echo "🧪 Testing pre-commit hooks..."
pre-commit run --all-files || {
    echo "⚠️  Some hooks failed on initial run. This is normal for the first setup."
    echo "   The hooks have been installed and will run on future commits."
}

echo "✅ Pre-commit hooks have been successfully set up!"
echo ""
echo "ℹ️  The following hooks will now run before each commit:"
echo "   • Go code formatting (gofmt, goimports)"
echo "   • Go static analysis (go vet)"
echo "   • Go module tidying"
echo "   • Go linting (golangci-lint)"
echo "   • Go tests"
echo "   • Go build verification"
echo "   • General file checks (trailing whitespace, etc.)"
echo ""
echo "💡 To skip hooks for a specific commit, use: git commit --no-verify"
echo "💡 To run hooks manually: pre-commit run --all-files"
echo "💡 To update hooks: pre-commit autoupdate"