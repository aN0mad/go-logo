# go-logo v2 Makefile
# This is a library project - no main binary to build
# Use this Makefile for testing, formatting, and building examples

VERSION=2.0.0

## help: print this help message
.PHONY: help
help:
	@echo "go-logo v2 - Available commands:"
	@echo ""
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

## tidy: format code and tidy modfile
.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v

## test: run all tests
.PHONY: test
test:
	go test -v -race -count=1 ./...

## test-coverage: run tests with coverage
.PHONY: test-coverage
test-coverage:
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

## bench: run benchmarks
.PHONY: bench
bench:
	go test -bench=. -benchmem ./...

## lint: run linter (requires golangci-lint)
.PHONY: lint
lint:
	golangci-lint run ./...

## examples: build all examples
.PHONY: examples
examples:
	@echo "Building all examples..."
	@cd examples/basic && go build -o /dev/null .
	@cd examples/console_colors && go build -o /dev/null .
	@cd examples/console_pretty && go build -o /dev/null .
	@cd examples/dual_output && go build -o /dev/null .
	@cd examples/file_logging && go build -o /dev/null .
	@cd examples/file_rotation && go build -o /dev/null .
	@cd examples/json_logging && go build -o /dev/null .
	@cd examples/multi_handler && go build -o /dev/null .
	@cd examples/cli_basic && go build -o /dev/null .
	@cd examples/cli_full && go build -o /dev/null .
	@cd examples/cli_profiles && go build -o /dev/null .
	@cd examples/http_echo && go build -o /dev/null .
	@cd examples/http_fiber_v2 && go build -o /dev/null .
	@cd examples/http_fiber_v3 && go build -o /dev/null .
	@cd examples/http_gin && go build -o /dev/null .
	@echo "✓ All examples built successfully"

## examples-tidy: tidy all example modules
.PHONY: examples-tidy
examples-tidy:
	@echo "Tidying all example modules..."
	@for dir in examples/*/; do \
		if [ -f "$$dir/go.mod" ]; then \
			echo "  $$dir"; \
			cd "$$dir" && go mod tidy && cd ../..; \
		fi \
	done
	@echo "✓ All examples tidied"

## clean: clean build artifacts and caches
.PHONY: clean
clean:
	go clean -cache
	go clean -testcache
	rm -f coverage.out coverage.html
	@echo "✓ Cleaned"

## all: run tidy, test, and build examples
.PHONY: all
all: tidy test examples
	@echo "✓ All tasks completed"
