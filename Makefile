.PHONY: help proto build run test clean install-tools

# Variables
PROTO_DIR=api/proto
PROTO_GEN_DIR=api/proto/gen/go
PROTO_FILES=$(wildcard $(PROTO_DIR)/*.proto)
BINARY_NAME=user-service
BINARY_PATH=./bin/$(BINARY_NAME)
CLIENT_BINARY_NAME=user-client
CLIENT_BINARY_PATH=./bin/$(CLIENT_BINARY_NAME)

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install-tools: ## Install required tools for proto generation
	@echo "Installing protoc plugins..."
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "Tools installed. Make sure protoc is installed: https://grpc.io/docs/protoc-installation/"

proto: ## Generate Go code from proto files
	@echo "Generating proto files..."
	@mkdir -p $(PROTO_GEN_DIR)/user
	@protoc --go_out=$(PROTO_GEN_DIR)/user \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_GEN_DIR)/user \
		--go-grpc_opt=paths=source_relative \
		--proto_path=$(PROTO_DIR) \
		$(PROTO_FILES)
	@echo "Proto files generated in $(PROTO_GEN_DIR)/user"

build: proto ## Build the application
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p bin
	@go build -o $(BINARY_PATH) ./cmd/server
	@echo "Binary built: $(BINARY_PATH)"

build-client: proto ## Build the client application
	@echo "Building $(CLIENT_BINARY_NAME)..."
	@mkdir -p bin
	@go build -o $(CLIENT_BINARY_PATH) ./cmd/client
	@echo "Client binary built: $(CLIENT_BINARY_PATH)"

run: build ## Build and run the application
	@echo "Running $(BINARY_NAME)..."
	@$(BINARY_PATH)

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

clean: ## Clean generated files and binaries
	@echo "Cleaning..."
	@rm -rf $(PROTO_GEN_DIR)
	@rm -rf bin
	@echo "Clean complete"

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

