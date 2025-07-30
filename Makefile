# Change these variables as necessary.
MAIN_PACKAGE_PATH=cmd/main.go
BINARY_NAME=go-logo
BIN_DIR=./bin
VERSION=1.0.0

## help: print this help message
.PHONY: help
help:
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

## tidy: format code and tidy modfile
.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v

## build: build the unix version
.PHONY: build
build:
    # Include additional build steps, like TypeScript, SCSS or Tailwind compilation here...
	go build -ldflags="-s -w -X main.VERSION=${VERSION}" -o=${BIN_DIR}/${BINARY_NAME}.elf ${MAIN_PACKAGE_PATH}

## buildwin: build the windows version
.PHONY: build
buildwin:
    # Include additional build steps, like TypeScript, SCSS or Tailwind compilation here...
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -X main.VERSION=${VERSION}" -o=${BIN_DIR}/${BINARY_NAME}.exe ${MAIN_PACKAGE_PATH}

## all: build all applications for unix and windows
.PHONY: all
all:
	make build
	make buildwin

## clean: clean the repository
.PHONY: clean
clean:
	rm -rf ${BIN_DIR}/*
	go clean -cache
	go mod tidy -v
	mv ./logs /tmp/${BINARY_NAME}-logs-$(shell date +%Y%m%d-%H%M%S) || true