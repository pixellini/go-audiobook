.PHONY: build run test clean help

# Build the application
build:
	go build -o bin/go-audiobook main.go

# Run the application (requires a text file argument)
run: build
	./bin/go-audiobook

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Install dependencies
deps:
	go mod tidy
	go mod download

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	golangci-lint run

# Show help
help:
	@echo "Available commands:"
	@echo "  build  - Build the application"
	@echo "  run    - Run the application (requires text file argument)"
	@echo "  test   - Run tests"
	@echo "  clean  - Clean build artifacts"
	@echo "  deps   - Install dependencies"
	@echo "  fmt    - Format code"
	@echo "  lint   - Run linter"
	@echo "  help   - Show this help message" 