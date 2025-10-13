SHELL := /bin/bash
GO_BIN ?= $(shell pwd)/.bin
GOCI_LINT_VERSION ?= v2.3.1

PATH := $(GO_BIN):$(PATH)

format-go::
	golangci-lint run --fix ./...

tools::
	mkdir -p $(GO_BIN)
	curl -sSfL "https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh" | sh -s -- -b ${GO_BIN} ${GOCI_LINT_VERSION}

.PHONY: fmt tidy clean build run

# Format all Go files (style, imports, modules)
fmt:
	@echo "🧹 Formatting Go files..."
	@go fmt ./...
	@echo "🧩 Fixing imports..."
	@goimports -w .
	@echo "📦 Tidying go.mod..."
	@go mod tidy
	@echo "✅ Formatting complete!"

# Optional: clean build artifacts
clean:
	@echo "🗑️ Cleaning build artifacts..."
	@go clean
	@echo "✅ Cleaned!"

# Optional: build and run targets for convenience
build:
	@echo "🔨 Building project..."
	@go build -o bin/server ./cmd/server
	@echo "✅ Build complete!"

run:
	@echo "🚀 Running server..."
	@go run ./cmd/server/main.go
