.DEFAULT_GOAL := help

BINARY_NAME=authkeeper
VERSION?=dev
BUILD_DIR=.

help: ## Show this help message
	@awk 'BEGIN {FS = ":.*## "; printf "\nUsage:\n  make <target>\n\nTargets:\n"} \
		/^([a-zA-Z_-]+):.*## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the binary
	go build -ldflags "-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/authkeeper

install: ## Install the binary to $GOPATH/bin
	go install -ldflags "-X main.version=$(VERSION)" ./cmd/authkeeper

run: ## Run the application
	go run ./cmd/authkeeper

test: ## Run tests
	go test -v -race ./...

test-coverage: ## Run tests with coverage
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	go tool cover -func=coverage.out | grep total

generate-mocks: ## Generate mocks using mockery
	go run github.com/vektra/mockery/v2@latest

clean: ## Clean build artifacts
	rm -f $(BINARY_NAME)
	rm -rf dist/
	rm -f coverage.out coverage.html

tidy: ## Run go mod tidy
	go mod tidy

fmt: ## Format code
	go fmt ./...

lint: ## Run golangci-lint
	golangci-lint run

.PHONY: help build install run test test-coverage generate-mocks clean tidy fmt lint
