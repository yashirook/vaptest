# Makefile

BINARY_NAME=vaptest
BUILD_DIR=bin
VERSION=0.1.0
GIT_COMMIT=$(shell git rev-parse --short HEAD)
BUILD_DATE=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS=-ldflags "-X 'github.com/yashirook/vaptest/cmd.version=$(VERSION)' -X 'github.com/yashirook/vaptest/cmd.gitCommit=$(GIT_COMMIT)' -X 'github.com/yashirook/vaptest/cmd.buildDate=$(BUILD_DATE)'"

.PHONY: all build run test clean

all: build

build:
	@echo "Building $(BINARY_NAME)..."
	mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)
	@echo "Built $(BINARY_NAME) version $(VERSION) at $(BUILD_DATE)"
	@echo "Binary is located at $(BUILD_DIR)/$(BINARY_NAME)"

run: build
	@echo "Running $(BINARY_NAME)..."
	$(BUILD_DIR)/$(BINARY_NAME) $(filter-out $@,$(MAKECMDGOALS))

test:
	@echo "Running tests..."
	go test ./...

clean:
	@echo "Cleaning up..."
	rm -rf $(BUILD_DIR)

# この行は、引数をターゲットと解釈しないようにするためのものです
%:
	@:
