# Go project common targets

# Binary name
BINARY_NAME=modelgate

# Go build variables
GO=go
GOCGO=CGO_ENABLED=1
GOOS=$(shell $(GO) env GOOS)
GOARCH=$(shell $(GO) env GOARCH)
LDFLAGS=-ldflags="-s -w"

# Docker
DOCKER_IMAGE=modelgate/modelgate
DOCKER_TAG?=latest

# Default target
.PHONY: all
all: help

# Build the binary
.PHONY: build
build:
	$(GOCGO) $(GO) build $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd/main.go

# Build for different platforms
.PHONY: build-linux
build-linux:
	GOOS=linux GOARCH=amd64 $(GOCGO) $(GO) build $(LDFLAGS) -o bin/$(BINARY_NAME)-linux-amd64 ./cmd/main.go

.PHONY: build-darwin
build-darwin:
	GOOS=darwin GOARCH=amd64 $(GOCGO) $(GO) build $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin-amd64 ./cmd/main.go

# Run the application
.PHONY: run
run:
	$(GO) run ./cmd/main.go all

# Clean build artifacts
.PHONY: clean
clean:
	rm -rf bin/
	$(GO) clean

# Download dependencies
.PHONY: deps
deps:
	$(GO) mod download

# Tidy dependencies
.PHONY: tidy
tidy:
	$(GO) mod tidy

# Format code
.PHONY: fmt
fmt:
	$(GO) fmt ./...

# Run go vet
.PHONY: vet
vet:
	$(GO) vet ./...

# Run tests
.PHONY: test
test:
	$(GO) test -v ./...

# Run tests with coverage
.PHONY: test-cover
test-cover:
	$(GO) test -v -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html

# Run tests in race mode
.PHONY: test-race
test-race:
	$(GO) test -race -v ./...

# Run unit tests only
.PHONY: test-unit
test-unit:
	$(GO) test -v -short ./...

# Generate code (protobuf, mocks, etc.)
.PHONY: generate
generate:
	$(GO) generate ./...

# Install the binary
.PHONY: install
install:
	$(GOCGO) $(GO) install $(LDFLAGS) ./cmd/main.go

# Run linter (golangci-lint)
.PHONY: lint
lint:
	golangci-lint run ./...

# Run linter with auto-fix
.PHONY: lint-fix
lint-fix:
	golangci-lint run ./... --fix

# Docker: build image
.PHONY: docker-build
docker-build:
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

# Docker: build image with latest tag
.PHONY: docker-build-latest
docker-build-latest:
	docker build -t $(DOCKER_IMAGE):latest -t $(DOCKER_TAG):$(shell git describe --tags --always --dirty) .

# Docker: run container
.PHONY: docker-run
docker-run:
	docker run -p 8080:8080 -p 8088:8088 $(DOCKER_IMAGE):$(DOCKER_TAG)

# Docker: push image
.PHONY: docker-push
docker-push:
	docker push $(DOCKER_IMAGE):$(DOCKER_TAG)

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build          - Build the binary"
	@echo "  build-linux    - Build for Linux AMD64"
	@echo "  build-darwin   - Build for Darwin AMD64"
	@echo "  run            - Run the application"
	@echo "  run-config     - Run with custom config file"
	@echo "  clean          - Clean build artifacts"
	@echo "  deps           - Download dependencies"
	@echo "  tidy           - Tidy dependencies"
	@echo "  fmt            - Format code"
	@echo "  vet            - Run go vet"
	@echo "  test           - Run tests"
	@echo "  test-cover     - Run tests with coverage"
	@echo "  test-race      - Run tests in race mode"
	@echo "  test-unit      - Run unit tests only"
	@echo "  generate       - Generate code"
	@echo "  install        - Install the binary"
	@echo "  lint           - Run linter"
	@echo "  lint-fix       - Run linter with auto-fix"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run Docker container"
	@echo "  docker-push    - Push Docker image"
