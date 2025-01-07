SHELL=bash

# Default platform parameters if not specified
OS ?= $(GOOS)
ARCH ?= $(GOARCH)
OUT ?= get-release

# Install dependencies
.PHONY: install-deps
install-deps:
	go mod tidy

# Run tests
.PHONY: tests
tests:
	@set -e
	go test ./... -v || exit 1;

.PHONY: cover
cover:
	@set -e
	go test -short -count=1 -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm ./coverage.out

# Build single binary
.PHONY: build
build:
	@set -e
	@echo "Building $(OUT) for OS=$(OS) ARCH=$(ARCH)"
	CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) go build -o ./tmp/$(OUT) ./src/main.go || exit 1;


# Build single binary
.PHONY: run
run:
	@set -e
	@echo "Run $(OUT) for OS=$(OS) ARCH=$(ARCH)"
	CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) go run ./src/main.go || exit 1;

# Help (optional)
help:
	@echo "Available targets:"
	@echo "  install-deps    Install dependencies"
	@echo "  tests           Run tests"
	@echo "  cover           Run coverage"
	@echo "  build           Build the application (single platform)"
	@echo "  run			 Run application"
	@echo "  clean           Clean build artifacts"
