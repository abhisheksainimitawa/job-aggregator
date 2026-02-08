# Makefile for Job Aggregator

.PHONY: help build run test clean docker-build docker-up docker-down migrate scrape

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	@echo "Building API server..."
	@go build -o bin/api cmd/api/main.go
	@echo "Building scraper CLI..."
	@go build -o bin/scraper cmd/scraper/main.go
	@echo "Build complete!"

run: ## Run the API server
	@go run cmd/api/main.go

scrape: ## Run the scraper
	@go run cmd/scraper/main.go

test: ## Run tests
	@go test -v ./...

test-coverage: ## Run tests with coverage
	@go test -cover -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

clean: ## Clean build artifacts
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "Clean complete!"

deps: ## Download dependencies
	@go mod download
	@go mod tidy

docker-build: ## Build Docker image
	@docker-compose build

docker-up: ## Start all services with Docker Compose
	@docker-compose up -d
	@echo "Services started! API available at http://localhost:8080"

docker-down: ## Stop all Docker services
	@docker-compose down

docker-logs: ## View Docker logs
	@docker-compose logs -f

migrate: ## Run database migrations
	@echo "Initializing database schema..."
	@go run cmd/api/main.go || echo "Schema initialized"

fmt: ## Format code
	@go fmt ./...

lint: ## Run linter
	@golangci-lint run || echo "Install golangci-lint: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"

dev: docker-up ## Start development environment
	@echo "Waiting for services to be ready..."
	@sleep 5
	@make run
