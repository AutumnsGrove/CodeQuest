# CodeQuest Makefile
# Build automation for the CodeQuest RPG productivity tool

# Variables
BINARY_NAME=codequest
MAIN_PATH=cmd/codequest/main.go
BUILD_DIR=build
DIST_DIR=dist
COVERAGE_FILE=coverage.out
COVERAGE_HTML=coverage.html

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet

# Build variables
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
COMMIT_HASH=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.CommitHash=$(COMMIT_HASH)"

# Colors for output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
NC=\033[0m # No Color

.PHONY: all build clean test coverage coverage-html run install dev fmt vet lint \
        deps dev-deps help docker-build docker-run release-all bench prof \
        check pre-commit init install-global uninstall-global install-deps

# Default target
all: clean fmt vet test build

## help: Show this help message
help:
	@echo "CodeQuest Build System"
	@echo "====================="
	@echo ""
	@echo "Available targets:"
	@echo ""
	@grep -E '^##' Makefile | sed 's/## /  /' | column -t -s ':'
	@echo ""
	@echo "Examples:"
	@echo "  make build    - Build the binary"
	@echo "  make test     - Run all tests"
	@echo "  make dev      - Run in development mode"
	@echo "  make install  - Install globally"

## init: Initialize the Go module and download dependencies
init:
	@echo "$(GREEN)Initializing Go module...$(NC)"
	$(GOMOD) init github.com/AutumnsGrove/codequest || true
	$(GOMOD) tidy
	@echo "$(GREEN)Module initialized!$(NC)"

## deps: Download module dependencies
deps:
	@echo "$(GREEN)Downloading dependencies...$(NC)"
	$(GOMOD) download
	$(GOMOD) tidy

## dev-deps: Install development tools
dev-deps:
	@echo "$(GREEN)Installing development tools...$(NC)"
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/rakyll/hey@latest
	go install github.com/securego/gosec/v2/cmd/gosec@latest
	@echo "$(GREEN)Development tools installed!$(NC)"

## build: Build the binary
build:
	@echo "$(GREEN)Building $(BINARY_NAME)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "$(GREEN)Build complete: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"

## build-all: Build for multiple platforms
build-all:
	@echo "$(GREEN)Building for multiple platforms...$(NC)"
	@mkdir -p $(DIST_DIR)
	# macOS AMD64
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	# macOS ARM64 (M1/M2)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	# Linux AMD64
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	# Linux ARM64
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_PATH)
	# Windows AMD64
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	@echo "$(GREEN)Multi-platform build complete!$(NC)"

## run: Run the application
run: build
	@echo "$(GREEN)Running $(BINARY_NAME)...$(NC)"
	./$(BUILD_DIR)/$(BINARY_NAME)

## dev: Run in development mode with hot reload (requires air)
dev:
	@command -v air > /dev/null || (echo "$(RED)Please install air: go install github.com/cosmtrek/air@latest$(NC)" && exit 1)
	@echo "$(GREEN)Running in development mode...$(NC)"
	air

## install: Install the binary globally
install: build
	@echo "$(GREEN)Installing $(BINARY_NAME) globally...$(NC)"
	go install $(LDFLAGS) ./$(MAIN_PATH)
	@echo "$(GREEN)$(BINARY_NAME) installed to $(GOPATH)/bin$(NC)"

## uninstall: Remove the globally installed binary
uninstall:
	@echo "$(YELLOW)Uninstalling $(BINARY_NAME)...$(NC)"
	rm -f $(GOPATH)/bin/$(BINARY_NAME)
	@echo "$(GREEN)$(BINARY_NAME) uninstalled$(NC)"

## install-global: Install codequest to /usr/local/bin for global access
install-global: build
	@echo "$(GREEN)Installing $(BINARY_NAME) to /usr/local/bin...$(NC)"
	@cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)
	@chmod +x /usr/local/bin/$(BINARY_NAME)
	@echo "$(GREEN)✅ $(BINARY_NAME) installed! Run 'codequest' from anywhere.$(NC)"

## uninstall-global: Remove codequest from /usr/local/bin
uninstall-global:
	@echo "$(YELLOW)Removing $(BINARY_NAME) from /usr/local/bin...$(NC)"
	@rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "$(GREEN)✅ $(BINARY_NAME) uninstalled.$(NC)"

## install-deps: Check for required dependencies (Skate, optional Mods)
install-deps:
	@echo "$(GREEN)Checking dependencies...$(NC)"
	@which skate > /dev/null || echo "$(RED)⚠️  Skate not found. Install: brew install charmbracelet/tap/skate$(NC)"
	@which mods > /dev/null || echo "$(YELLOW)ℹ️  Mods not found (optional). Install: brew install charmbracelet/tap/mods$(NC)"
	@echo "$(GREEN)✅ Dependency check complete$(NC)"

## clean: Clean build artifacts
clean:
	@echo "$(YELLOW)Cleaning build artifacts...$(NC)"
	$(GOCLEAN)
	rm -rf $(BUILD_DIR) $(DIST_DIR)
	rm -f $(COVERAGE_FILE) $(COVERAGE_HTML)
	rm -f *.prof
	@echo "$(GREEN)Clean complete!$(NC)"

## test: Run all tests
test:
	@echo "$(GREEN)Running tests...$(NC)"
	$(GOTEST) -v -race -timeout 30s ./...

## test-short: Run short tests only
test-short:
	@echo "$(GREEN)Running short tests...$(NC)"
	$(GOTEST) -v -short ./...

## coverage: Run tests with coverage
coverage:
	@echo "$(GREEN)Running tests with coverage...$(NC)"
	$(GOTEST) -v -race -coverprofile=$(COVERAGE_FILE) -covermode=atomic ./...
	@echo "$(GREEN)Coverage report:$(NC)"
	@go tool cover -func=$(COVERAGE_FILE)

## coverage-html: Generate HTML coverage report
coverage-html: coverage
	@echo "$(GREEN)Generating HTML coverage report...$(NC)"
	go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "$(GREEN)Coverage report generated: $(COVERAGE_HTML)$(NC)"
	@command -v open > /dev/null && open $(COVERAGE_HTML) || true

## bench: Run benchmarks
bench:
	@echo "$(GREEN)Running benchmarks...$(NC)"
	$(GOTEST) -bench=. -benchmem ./...

## prof: Run CPU profiling
prof:
	@echo "$(GREEN)Running CPU profiling...$(NC)"
	$(GOTEST) -cpuprofile=cpu.prof -bench=. ./internal/game
	go tool pprof cpu.prof

## fmt: Format the code
fmt:
	@echo "$(GREEN)Formatting code...$(NC)"
	$(GOFMT) ./...
	@command -v goimports > /dev/null && goimports -w . || true

## vet: Run go vet
vet:
	@echo "$(GREEN)Running go vet...$(NC)"
	$(GOVET) ./...

## lint: Run golangci-lint
lint:
	@echo "$(GREEN)Running linter...$(NC)"
	@command -v golangci-lint > /dev/null || (echo "$(RED)Please install golangci-lint$(NC)" && exit 1)
	golangci-lint run

## security: Run security checks
security:
	@echo "$(GREEN)Running security checks...$(NC)"
	@command -v gosec > /dev/null || (echo "$(RED)Please install gosec$(NC)" && exit 1)
	gosec -fmt=text ./...

## check: Run all checks (fmt, vet, lint, test)
check: fmt vet lint test
	@echo "$(GREEN)All checks passed!$(NC)"

## pre-commit: Run pre-commit checks
pre-commit: fmt vet lint test-short
	@echo "$(GREEN)Pre-commit checks passed!$(NC)"

## docker-build: Build Docker image
docker-build:
	@echo "$(GREEN)Building Docker image...$(NC)"
	docker build -t $(BINARY_NAME):$(VERSION) .
	docker tag $(BINARY_NAME):$(VERSION) $(BINARY_NAME):latest

## docker-run: Run in Docker container
docker-run:
	@echo "$(GREEN)Running in Docker...$(NC)"
	docker run -it --rm $(BINARY_NAME):latest

## release: Create a release build
release: clean test build-all
	@echo "$(GREEN)Creating release $(VERSION)...$(NC)"
	@mkdir -p $(DIST_DIR)/release-$(VERSION)
	@cp $(DIST_DIR)/* $(DIST_DIR)/release-$(VERSION)/
	@cp README.md LICENSE $(DIST_DIR)/release-$(VERSION)/
	@cd $(DIST_DIR) && tar -czf codequest-$(VERSION).tar.gz release-$(VERSION)
	@echo "$(GREEN)Release created: $(DIST_DIR)/codequest-$(VERSION).tar.gz$(NC)"

## docs: Generate documentation
docs:
	@echo "$(GREEN)Generating documentation...$(NC)"
	@command -v godoc > /dev/null || (echo "$(YELLOW)Installing godoc...$(NC)" && go install golang.org/x/tools/cmd/godoc@latest)
	@echo "$(GREEN)Documentation available at http://localhost:6060$(NC)"
	godoc -http=:6060

## setup: Complete development environment setup
setup: init deps dev-deps
	@echo "$(GREEN)Development environment ready!$(NC)"
	@echo ""
	@echo "Next steps:"
	@echo "  1. Create the project structure: make init-structure"
	@echo "  2. Run tests: make test"
	@echo "  3. Build the project: make build"
	@echo "  4. Run the application: make run"

# Print variable values (for debugging)
print-%:
	@echo $* = $($*)

.DEFAULT_GOAL := help