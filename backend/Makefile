# Makefile for educ-retro backend

.PHONY: help test test-unit test-all bench clean deps build run swagger swagger-serve swagger-clean

# Default target
help:
	@echo "Available targets:"
	@echo "  test          - Run all tests (unit + schema)"
	@echo "  test-unit     - Run unit tests only"
	@echo "  test-schema   - Run schema validation tests"
	@echo "  bench         - Run benchmarks"
	@echo "  deps          - Install dependencies"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application"
	@echo "  clean         - Clean build artifacts"
	@echo "  swagger       - Generate Swagger documentation"
	@echo "  swagger-serve - Generate docs and show serve instructions"
	@echo "  swagger-clean - Clean and regenerate Swagger docs"
	@echo "  swagger-validate - Generate and validate Swagger docs"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy
	@echo "Installing test dependencies..."
	go get github.com/stretchr/testify
	go get github.com/mattn/go-sqlite3
	@echo "Installing Swagger dependencies..."
	go get github.com/swaggo/swag/cmd/swag
	go get github.com/swaggo/gin-swagger
	go get github.com/swaggo/files

# Run unit tests
test-unit:
	@echo "Running unit tests..."
	go test ./... -v -short

# Run schema tests
test-schema:
	@echo "Running schema tests..."
	go test ./tests/schema -v

# Run all tests
test: test-unit test-schema

# Build the application
build:
	@echo "Building application..."
	go build -o bin/server cmd/server/main.go

# Run the application
run:
	@echo "Running application..."
	go run cmd/server/main.go

# Stop the application
stop:
	@echo "Stopping application..."
	@pkill -f "go run cmd/server/main.go" || echo "No server process found"
	@pkill -f "bin/server" || echo "No server binary process found"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	go clean

# Generate Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	@./scripts/generate-swagger.sh

# Generate Swagger docs and show serve instructions
swagger-serve:
	@echo "Generating Swagger documentation..."
	@./scripts/generate-swagger.sh -s

# Clean and regenerate Swagger docs
swagger-clean:
	@echo "Cleaning and regenerating Swagger documentation..."
	@./scripts/generate-swagger.sh -c

# Generate and validate Swagger docs
swagger-validate:
	@echo "Generating and validating Swagger documentation..."
	@./scripts/generate-swagger.sh -v

# Test coverage
coverage:
	@echo "Running tests with coverage..."
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"


# Lint code
lint:
	@echo "Running linter..."
	golangci-lint run

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Vet code
vet:
	@echo "Vetting code..."
	go vet ./...

# Security scan
security:
	@echo "Running security scan..."
	gosec ./...

# Full test suite with coverage and linting
test-full: deps fmt vet lint test coverage
	@echo "Full test suite completed"

# CI/CD pipeline
ci: deps fmt vet lint test
	@echo "CI pipeline completed"

# Development setup
dev-setup: deps swagger
	@echo "Development environment setup completed"
	@echo "To start the server: make run"
	@echo "To view API docs: http://localhost:8080/swagger/index.html"

# Production build
prod-build: deps test build
	@echo "Production build completed"

# Quick start (for new developers)
quick-start: deps swagger build
	@echo "Quick start completed!"
	@echo ""
	@echo "Next steps:"
	@echo "  1. Start the server: make run"
	@echo "  2. View API docs: http://localhost:8080/swagger/index.html"
	@echo "  3. Run tests: make test"
