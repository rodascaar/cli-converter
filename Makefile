.PHONY: build test clean install uninstall fmt vet mod-tidy

# Build the binary
build:
	go build -o bin/audioconv main.go

# Build for development (with debug info)
dev:
	go build -gcflags="all=-N -l" -o bin/audioconv main.go

# Run tests
test:
	go test ./...

# Run tests with coverage
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run tests with race detection
test-race:
	go test -race ./...

# Format code
fmt:
	go fmt ./...

# Vet code
vet:
	go vet ./...

# Tidy modules
mod-tidy:
	go mod tidy

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html
	rm -f audioconv

# Install to system (requires sudo)
install: build
	sudo cp bin/audioconv /usr/local/bin/

# Uninstall from system (requires sudo)
uninstall:
	sudo rm -f /usr/local/bin/audioconv

# Run linter (requires golangci-lint)
lint:
	golangci-lint run

# Cross-compile for multiple platforms
cross-compile:
	GOOS=linux GOARCH=amd64 go build -o bin/audioconv-linux-amd64 main.go
	GOOS=linux GOARCH=arm64 go build -o bin/audioconv-linux-arm64 main.go
	GOOS=darwin GOARCH=amd64 go build -o bin/audioconv-darwin-amd64 main.go
	GOOS=darwin GOARCH=arm64 go build -o bin/audioconv-darwin-arm64 main.go
	GOOS=windows GOARCH=amd64 go build -o bin/audioconv-windows-amd64.exe main.go

# Create release archives
release: cross-compile
	mkdir -p release
	tar -czf release/audioconv-linux-amd64.tar.gz -C bin audioconv-linux-amd64
	tar -czf release/audioconv-linux-arm64.tar.gz -C bin audioconv-linux-arm64
	tar -czf release/audioconv-darwin-amd64.tar.gz -C bin audioconv-darwin-amd64
	tar -czf release/audioconv-darwin-arm64.tar.gz -C bin audioconv-darwin-arm64
	zip release/audioconv-windows-amd64.zip bin/audioconv-windows-amd64.exe

# Development workflow
dev-setup: mod-tidy fmt vet

# CI workflow
ci: mod-tidy fmt vet test

# All-in-one development command
all: clean dev-setup test build

# Help
help:
	@echo "Available targets:"
	@echo "  build          - Build the binary"
	@echo "  dev            - Build for development with debug info"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  test-race      - Run tests with race detection"
	@echo "  fmt            - Format code"
	@echo "  vet            - Vet code"
	@echo "  mod-tidy       - Tidy modules"
	@echo "  clean          - Clean build artifacts"
	@echo "  install        - Install to system (requires sudo)"
	@echo "  uninstall      - Uninstall from system (requires sudo)"
	@echo "  lint           - Run linter (requires golangci-lint)"
	@echo "  cross-compile  - Cross-compile for multiple platforms"
	@echo "  release        - Create release archives"
	@echo "  dev-setup      - Setup for development"
	@echo "  ci             - CI workflow"
	@echo "  all            - All-in-one development command"
	@echo "  help           - Show this help"