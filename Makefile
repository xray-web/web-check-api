$(shell cp -n .env.example .env)
include .env
export

# Run the application
run:
	@go run main.go
.PHONY: run

# Test the application
test:
	@go test ./...
.PHONY: test

# Build the application
build:
	@go build -o bin/app main.go
.PHONY: build

# Clean the build artifacts
clean:
	@rm -rf bin/
.PHONY: clean

# Lint the codebase
lint:
	@golangci-lint run
.PHONY: lint

# Format the codebase
format:
	@go fmt ./...
.PHONY: format

# Install dependencies
deps:
	@go mod tidy
	@go mod vendor
.PHONY: deps

# Ensure .env file is sourced
env:
	@source .env
.PHONY: env
