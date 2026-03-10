.PHONY: help build run test test-race test-coverage lint fmt vet check fuzz fuzz-nightly clean deps tidy update-deps generate docs site release snapshot install verify all

# Default target
.DEFAULT_GOAL := help

# Variables
BINARY_NAME=devtui
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE?=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-X github.com/skatkov/devtui/cmd.version=$(VERSION) -X github.com/skatkov/devtui/cmd.commit=$(COMMIT) -X github.com/skatkov/devtui/cmd.date=$(DATE)"
FUZZTIME?=5s
NIGHTLY_FUZZTIME?=45s

help: ## Display this help message
	@echo "DevTUI - Development targets:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'
	@echo ""

build: ## Build the project
	@echo "Building $(BINARY_NAME)..."
	go build -v $(LDFLAGS) -o $(BINARY_NAME) .
	@echo "Build complete: ./$(BINARY_NAME)"

run: ## Run the TUI (main mode)
	@echo "Running DevTUI..."
	go run $(LDFLAGS) main.go

test: ## Run all tests
	@echo "Running tests..."
	go test -v -failfast ./...

test-race: ## Run tests with race detection (what CI uses)
	@echo "Running tests with race detection..."
	go test -v -failfast -race ./...

fuzz: ## Run fuzz tests (set FUZZTIME=30s for deeper runs)
	@echo "Running fuzz tests (FUZZTIME=$(FUZZTIME))..."
	go test ./internal/base64 -run=^$$ -fuzz=FuzzBase64RoundTrip -fuzztime=$(FUZZTIME)
	go test ./internal/base64 -run=^$$ -fuzz=FuzzDecodeMatchesStdlibBehavior -fuzztime=$(FUZZTIME)
	go test ./internal/csv2json -run=^$$ -fuzz=FuzzConvertDoesNotPanic -fuzztime=$(FUZZTIME)
	go test ./internal/converter -run=^$$ -fuzz=FuzzYAMLToJSON -fuzztime=$(FUZZTIME)
	go test ./internal/converter -run=^$$ -fuzz=FuzzJSONToYAML -fuzztime=$(FUZZTIME)
	go test ./internal/converter -run=^$$ -fuzz=FuzzTOMLToJSON -fuzztime=$(FUZZTIME)
	go test ./internal/converter -run=^$$ -fuzz=FuzzJSONToTOML -fuzztime=$(FUZZTIME)
	go test ./internal/converter -run=^$$ -fuzz=FuzzXMLToJSON -fuzztime=$(FUZZTIME)
	go test ./internal/converter -run=^$$ -fuzz=FuzzJSONToXML -fuzztime=$(FUZZTIME)
	go test ./internal/converter -run=^$$ -fuzz=FuzzYAMLToTOML -fuzztime=$(FUZZTIME)
	go test ./internal/converter -run=^$$ -fuzz=FuzzTOMLToYAML -fuzztime=$(FUZZTIME)
	go test ./internal/yamlfmt -run=^$$ -fuzz=FuzzFormatYAML -fuzztime=$(FUZZTIME)
	go test ./internal/htmlfmt -run=^$$ -fuzz=FuzzFormatHTML -fuzztime=$(FUZZTIME)
	go test ./tui/jsonrepair -run=^$$ -fuzz=FuzzRepairJSONProducesValidJSON -fuzztime=$(FUZZTIME)
	go test ./tui/json2toon -run=^$$ -fuzz=FuzzConvertNoPanic -fuzztime=$(FUZZTIME)
	go test ./tui/json2toon -run=^$$ -fuzz=FuzzConvertWithOptionsNoPanic -fuzztime=$(FUZZTIME)

fuzz-nightly: ## Run extended fuzz tests (default NIGHTLY_FUZZTIME=45s)
	@echo "Running nightly fuzz tests (NIGHTLY_FUZZTIME=$(NIGHTLY_FUZZTIME))..."
	$(MAKE) fuzz FUZZTIME=$(NIGHTLY_FUZZTIME)

test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.txt ./...
	go tool cover -html=coverage.txt

lint: ## Run golangci-lint
	@echo "Running linters..."
	golangci-lint run

fmt: ## Format code with gofumpt and goimports
	@echo "Formatting code..."
	golangci-lint run --fix

vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...

check: lint test-race ## Run linters and tests with race detection

docs: ## Generate CLI and TUI documentation for the website
	@echo "Generating documentation..."
	cd docs && go run *.go
	@echo "Documentation generated in site/"

site: docs ## Alias for docs target

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -f $(BINARY_NAME)
	rm -f coverage.txt
	rm -rf dist/
	@echo "Clean complete"

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	go mod download

tidy: ## Run go mod tidy to clean up dependencies
	@echo "Tidying dependencies..."
	go mod tidy
	@echo "Dependencies tidied"

update-deps: ## Update all dependencies to latest versions
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy
	@echo "Dependencies updated"

generate: ## Run go generate
	@echo "Running go generate..."
	go generate ./...

release: ## Build release with goreleaser (creates binaries for all platforms)
	@echo "Building release with goreleaser..."
	goreleaser release --clean

snapshot: ## Build snapshot release with goreleaser (no publish)
	@echo "Building snapshot release..."
	goreleaser release --snapshot --clean
	@echo "Snapshot built in dist/"

install: build ## Build and install to GOPATH/bin
	@echo "Installing $(BINARY_NAME) to $(GOPATH)/bin..."
	go install $(LDFLAGS) .
	@echo "Installed: $(shell which $(BINARY_NAME))"

verify: generate tidy ## Verify go generate and go mod tidy don't create changes
	@if [ -n "$$(git status --porcelain)" ]; then \
		echo "Error: git status is not clean after generate and tidy"; \
		git status --porcelain; \
		exit 1; \
	fi
	@echo "Verification passed"

all: clean generate tidy build test lint ## Run all checks and build
	@echo "All tasks complete"
