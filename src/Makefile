# Makefile for monorepo project

.PHONY: help test build run clean docker-build docker-run

# Default target
help:
	@echo "Available commands:"
	@echo "  make test pkg=<package>     - Run tests for the specified API"
	@echo "  make build pkg=<package>    - Build the binary for specified API"
	@echo "  make run pkg=<package>      - Run API and DB with Docker"
	@echo "  make clean pkg=<package>    - Clean build artifacts for specified API"
	@echo "  make docker-build pkg=<package> - Build Docker image for specified API"
	@echo "  make docker-run pkg=<package>   - Run with Docker Compose for specified API"
	@echo ""
	@echo "Example: make run pkg=facture"

# Validate pkg parameter
check-pkg:
	@if [ -z "$(pkg)" ]; then \
		echo "Error: pkg parameter is required. Usage: make <command> pkg=<package>"; \
		exit 1; \
	fi
	@if [ ! -d "$(pkg)" ]; then \
		echo "Error: Package $(pkg) does not exist"; \
		exit 1; \
	fi

# Run tests
test: check-pkg
	@echo "Running tests for $(pkg) API..."
	@if [ -d "$(pkg)" ]; then \
		cd $(pkg) && go test -v ./...; \
	else \
		echo "No tests found for $(pkg) API"; \
	fi

# Build the binary
build: check-pkg
	@echo "Building $(pkg) binary..."
	@mkdir -p bin/$(pkg)
	@go build -o bin/$(pkg)/$(pkg) ./$(pkg)/bin/$(pkg)/main.go
	@echo "Binary built: bin/$(pkg)/$(pkg)"

# Run API and DB with Docker
run: check-pkg
	@echo "Starting $(pkg) API and receipts DB..."
	@cd $(pkg) && docker compose up --build

# Clean build artifacts
clean: check-pkg
	@echo "Cleaning build artifacts for $(pkg)..."
	@rm -rf bin/$(pkg)/
	@cd $(pkg) && docker compose down -v
	@echo "Clean completed for $(pkg)"

# Build Docker image
docker-build: check-pkg
	@echo "Building Docker image for $(pkg)..."
	@docker build -f $(pkg)/Dockerfile -t $(pkg)-api .

# Run with Docker Compose
docker-run: check-pkg
	@echo "Running services with Docker Compose for $(pkg)..."
	@cd $(pkg) && docker compose up --build
