APP_NAME  := myapp
BUILD_DIR := bin
MAIN_PATH := ./cmd/api/main.go
LDFLAGS   := -ldflags="-s -w"

.PHONY: help build run dev test test-verbose test-cover lint tidy clean \
        docker-build docker-run docker-up docker-down

help: ## Show available targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| sort \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the binary into ./bin/
	@mkdir -p $(BUILD_DIR)
	@go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)
	@echo "Built: $(BUILD_DIR)/$(APP_NAME)"

run: ## Run the application directly with go run
	@go run $(MAIN_PATH)

dev: ## Hot-reload with Air (install: go install github.com/air-verse/air@latest)
	@air

test: ## Run all tests
	@go test -race ./...

test-verbose: ## Run tests with verbose output
	@go test -race -v ./...

test-cover: ## Run tests and open HTML coverage report
	@go test -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

lint: ## Run golangci-lint (install: https://golangci-lint.run/usage/install/)
	@golangci-lint run ./...

tidy: ## Tidy and verify go modules
	@go mod tidy
	@go mod verify

clean: ## Remove build artifacts
	@rm -rf $(BUILD_DIR) tmp coverage.out coverage.html

docker-build: ## Build the Docker image
	@docker build -t $(APP_NAME):latest .

docker-run: ## Run the Docker container
	@docker run --env-file .env -p 8080:8080 $(APP_NAME):latest

docker-up: ## Start all services with Docker Compose
	@docker compose up --build

docker-down: ## Stop and remove Docker Compose services
	@docker compose down

.DEFAULT_GOAL := help
