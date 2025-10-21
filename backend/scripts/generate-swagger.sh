#!/bin/bash

# Swagger Documentation Generator Script
# This script generates Swagger documentation for the educ-retro API

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if Go is installed
check_go() {
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed or not in PATH"
        exit 1
    fi
    
    GO_VERSION=$(go version | awk '{print $3}')
    print_success "Go version: $GO_VERSION"
}

# Function to check if swag is available
check_swag() {
    if ! go run github.com/swaggo/swag/cmd/swag@latest version &> /dev/null; then
        print_warning "Swag not found, installing..."
        go install github.com/swaggo/swag/cmd/swag@latest
    fi
    print_success "Swag is available"
}

# Function to generate Swagger docs
generate_docs() {
    print_status "Generating Swagger documentation..."
    
    # Remove old docs if they exist
    if [ -d "docs" ]; then
        print_status "Removing old documentation..."
        rm -rf docs
    fi
    
    # Generate new docs
    if go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/server/main.go -o docs; then
        print_success "Swagger documentation generated successfully"
    else
        print_error "Failed to generate Swagger documentation"
        exit 1
    fi
}

# Function to validate generated docs
validate_docs() {
    print_status "Validating generated documentation..."
    
    if [ ! -f "docs/swagger.json" ]; then
        print_error "swagger.json not found"
        exit 1
    fi
    
    if [ ! -f "docs/swagger.yaml" ]; then
        print_error "swagger.yaml not found"
        exit 1
    fi
    
    if [ ! -f "docs/docs.go" ]; then
        print_error "docs.go not found"
        exit 1
    fi
    
    print_success "All documentation files generated successfully"
}

# Function to show generated files
show_files() {
    print_status "Generated files:"
    echo "  📄 docs/swagger.json - JSON format"
    echo "  📄 docs/swagger.yaml - YAML format"
    echo "  📄 docs/docs.go - Go package"
    echo ""
    print_status "Documentation size:"
    if [ -f "docs/swagger.json" ]; then
        echo "  📊 swagger.json: $(wc -c < docs/swagger.json) bytes"
    fi
    if [ -f "docs/swagger.yaml" ]; then
        echo "  📊 swagger.yaml: $(wc -c < docs/swagger.yaml) bytes"
    fi
}

# Function to show help
show_help() {
    echo "Swagger Documentation Generator for educ-retro API"
    echo ""
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -h, --help              Show this help message"
    echo "  -c, --clean             Clean generated docs before generating"
    echo "  -v, --validate          Validate generated documentation"
    echo "  -s, --serve             Start server to view docs (requires server to be running)"
    echo ""
    echo "Examples:"
    echo "  $0                      # Generate documentation"
    echo "  $0 -c                   # Clean and generate documentation"
    echo "  $0 -v                   # Generate and validate documentation"
    echo ""
    echo "After generation, start the server and visit:"
    echo "  http://localhost:8080/swagger/index.html"
}

# Main function
main() {
    local clean=false
    local validate=false
    local serve=false
    
    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                exit 0
                ;;
            -c|--clean)
                clean=true
                shift
                ;;
            -v|--validate)
                validate=true
                shift
                ;;
            -s|--serve)
                serve=true
                shift
                ;;
            *)
                print_error "Unknown option: $1"
                show_help
                exit 1
                ;;
        esac
    done
    
    # Start
    print_status "Starting Swagger documentation generation..."
    
    # Check if we're in the right directory
    if [ ! -f "go.mod" ]; then
        print_error "go.mod not found. Please run this script from the backend directory"
        exit 1
    fi
    
    # Check Go installation
    check_go
    
    # Check swag availability
    check_swag
    
    # Clean if requested
    if [ "$clean" = true ]; then
        print_status "Cleaning old documentation..."
        rm -rf docs
    fi
    
    # Generate documentation
    generate_docs
    
    # Validate if requested
    if [ "$validate" = true ]; then
        validate_docs
    fi
    
    # Show generated files
    show_files
    
    # Show serve instructions
    if [ "$serve" = true ]; then
        print_status "To view the documentation:"
        echo "  1. Start the server: go run cmd/server/main.go"
        echo "  2. Open browser: http://localhost:8080/swagger/index.html"
    else
        print_success "Documentation generated successfully!"
        print_status "To view the documentation:"
        echo "  1. Start the server: go run cmd/server/main.go"
        echo "  2. Open browser: http://localhost:8080/swagger/index.html"
    fi
}

# Run main function with all arguments
main "$@"
