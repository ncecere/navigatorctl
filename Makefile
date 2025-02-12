.PHONY: build install clean test lint

# Build the binary
build:
	go build -o navigatorctl

# Install the binary
install:
	go install

# Clean build artifacts
clean:
	rm -f navigatorctl
	go clean

# Run tests
test:
	go test ./...

# Run linter
lint:
	go vet ./...
	test -z "$$(gofmt -l .)"

# Build and run locally
run: build
	./navigatorctl

# Generate default config
config:
	mkdir -p $(HOME)/.config/navigatorctl
	cp config/default.yaml $(HOME)/.navigatorctl.yaml
	@echo "Default config copied to $(HOME)/.navigatorctl.yaml"
	@echo "Please edit the file and add your API credentials"

# Help target
help:
	@echo "Available targets:"
	@echo "  build    - Build the navigatorctl binary"
	@echo "  install  - Install the binary to GOPATH"
	@echo "  clean    - Remove build artifacts"
	@echo "  test     - Run tests"
	@echo "  lint     - Run linter"
	@echo "  run      - Build and run locally"
	@echo "  config   - Generate default config file"
	@echo "  help     - Show this help message"
