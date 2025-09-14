# Variables
APP_NAME := bagr-backend
DOCKER_IMAGE := $(APP_NAME):latest
DOCKER_CONTAINER := $(APP_NAME)-container
GO_VERSION := 1.21

# Default target
.DEFAULT_GOAL := help

# Help target
.PHONY: help
help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Development commands
.PHONY: run
run: ## Run the application locally
	@echo "Starting $(APP_NAME)..."
	go run cmd/main.go

.PHONY: run-with-config
run-with-config: ## Run the application with config file
	@echo "Starting $(APP_NAME) with config..."
	go run cmd/main.go -config=config.yaml

.PHONY: build
build: ## Build the application binary
	@echo "Building $(APP_NAME)..."
	go build -o bin/$(APP_NAME) cmd/main.go

.PHONY: clean
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	go clean

# Testing commands
.PHONY: test
test: ## Run tests
	@echo "Running tests..."
	go test -v ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Code quality commands
.PHONY: fmt
fmt: ## Format Go code
	@echo "Formatting code..."
	go fmt ./...

.PHONY: vet
vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...

.PHONY: lint
lint: ## Run golangci-lint (requires golangci-lint to be installed)
	@echo "Running golangci-lint..."
	golangci-lint run

.PHONY: tidy
tidy: ## Tidy go modules
	@echo "Tidying go modules..."
	go mod tidy

# Docker commands
.PHONY: docker-build
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE) .

.PHONY: docker-run
docker-run: ## Run Docker container
	@echo "Running Docker container..."
	docker run --rm --name $(DOCKER_CONTAINER) -p 8080:8080 $(DOCKER_IMAGE)

.PHONY: docker-run-detached
docker-run-detached: ## Run Docker container in detached mode
	@echo "Running Docker container in detached mode..."
	docker run -d --name $(DOCKER_CONTAINER) -p 8080:8080 $(DOCKER_IMAGE)

.PHONY: docker-stop
docker-stop: ## Stop Docker container
	@echo "Stopping Docker container..."
	docker stop $(DOCKER_CONTAINER) || true

.PHONY: docker-clean
docker-clean: docker-stop ## Clean Docker container and image
	@echo "Cleaning Docker container and image..."
	docker rm $(DOCKER_CONTAINER) || true
	docker rmi $(DOCKER_IMAGE) || true

# Database commands (for future use)
.PHONY: db-up
db-up: ## Start database services (requires docker-compose.yml)
	@echo "Starting database services..."
	docker-compose up -d postgres redis

.PHONY: db-down
db-down: ## Stop database services
	@echo "Stopping database services..."
	docker-compose down

.PHONY: db-migrate
db-migrate: ## Run database migrations (placeholder)
	@echo "Running database migrations..."
	@echo "TODO: Implement database migrations"

# Development workflow
.PHONY: dev-setup
dev-setup: tidy ## Setup development environment
	@echo "Setting up development environment..."
	@echo "Installing development dependencies..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: dev-check
dev-check: fmt vet lint test ## Run all development checks
	@echo "All development checks passed!"

.PHONY: dev-run
dev-run: dev-check run ## Run development checks and start the application

# Production commands
.PHONY: prod-build
prod-build: clean build ## Build for production
	@echo "Production build completed!"

.PHONY: prod-docker
prod-docker: docker-clean docker-build ## Build production Docker image
	@echo "Production Docker image built!"

# Utility commands
.PHONY: deps
deps: ## Download dependencies
	@echo "Downloading dependencies..."
	go mod download

.PHONY: deps-update
deps-update: ## Update dependencies
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy

.PHONY: version
version: ## Show Go version
	@go version

.PHONY: env
env: ## Show environment info
	@echo "Go version: $(shell go version)"
	@echo "Go environment:"
	@go env
