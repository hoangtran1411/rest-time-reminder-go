# RestTimeReminder Makefile
# =========================

# Build variables
BINARY_NAME=rest-time-reminder
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
BUILD_DATE?=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.buildDate=$(BUILD_DATE)"

# Go commands
GO=go
GOFMT=gofmt
GOLINT=golangci-lint

# Directories
CMD_DIR=./cmd/reminder
BUILD_DIR=./build

.PHONY: all build clean test lint fix fmt help install run dev

# Default target
all: clean lint test build

## Build Commands

build: ## Build the application for current platform
	@echo "Building $(BINARY_NAME)..."
	$(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)

build-all: build-windows build-linux build-darwin ## Build for all platforms

build-windows: ## Build for Windows
	@echo "Building for Windows..."
	GOOS=windows GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(CMD_DIR)

build-linux: ## Build for Linux
	@echo "Building for Linux..."
	GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(CMD_DIR)

build-darwin: ## Build for macOS
	@echo "Building for macOS..."
	GOOS=darwin GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(CMD_DIR)
	GOOS=darwin GOARCH=arm64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(CMD_DIR)

## Development Commands

run: ## Run the application in development mode
	$(GO) run $(CMD_DIR)

dev: ## Run with live reload (requires air)
	air

install: ## Install the application locally
	$(GO) install $(LDFLAGS) $(CMD_DIR)

## Quality Commands

# Test flags
# Use -race only if CGO is enabled
CGO_ENABLED_VAL=$(shell go env CGO_ENABLED)
ifeq ($(CGO_ENABLED_VAL),1)
	TEST_FLAGS=-race
else
	TEST_FLAGS=
endif

test: ## Run tests
	$(GO) test -v $(TEST_FLAGS) -cover ./...

test-coverage: ## Run tests with coverage report
	$(GO) test -v $(TEST_FLAGS) -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html

lint: ## Run linters
	$(GOLINT) run ./...

fix: ## Run linters and fix auto-fixable issues
	$(GOLINT) run --fix ./...

fmt: ## Format code
	$(GOFMT) -s -w .

vet: ## Run go vet
	$(GO) vet ./...

check: fmt vet lint test ## Run all checks

## Utility Commands

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

deps: ## Download dependencies
	$(GO) mod download

tidy: ## Tidy go.mod
	$(GO) mod tidy

update: ## Update dependencies
	$(GO) get -u ./...
	$(GO) mod tidy

## Release Commands

release: clean build-all ## Create release builds
	@echo "Creating release archives..."
	cd $(BUILD_DIR) && zip $(BINARY_NAME)-windows-amd64.zip $(BINARY_NAME)-windows-amd64.exe
	cd $(BUILD_DIR) && tar -czvf $(BINARY_NAME)-linux-amd64.tar.gz $(BINARY_NAME)-linux-amd64
	cd $(BUILD_DIR) && tar -czvf $(BINARY_NAME)-darwin-amd64.tar.gz $(BINARY_NAME)-darwin-amd64
	cd $(BUILD_DIR) && tar -czvf $(BINARY_NAME)-darwin-arm64.tar.gz $(BINARY_NAME)-darwin-arm64

## Help

help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)
