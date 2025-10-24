# Build configuration
RUST_DIR := rust
GO_DIR := go
RUST_LIB := $(RUST_DIR)/target/release/libshamir_benchmark.a
GO_BINARY := $(GO_DIR)/benchmark_hybrid

PROJECT_NAME := go_rust_hybrid

# Default target
.PHONY: all
all: build

# Check if required files exist
.PHONY: check
check:
	@echo "=== Checking project structure ==="
	@test -f "$(RUST_DIR)/Cargo.toml" || (echo "Error: $(RUST_DIR)/Cargo.toml not found" && exit 1)
	@test -f "$(RUST_DIR)/src/lib.rs" || (echo "Error: $(RUST_DIR)/src/lib.rs not found" && exit 1)
	@test -f "$(GO_DIR)/main.go" || (echo "Error: $(GO_DIR)/main.go not found" && exit 1)
	@echo "✓ Project structure is valid"

# Initialize Go module if needed
.PHONY: init-go
init-go:
	@echo "Initializing Go module..."
	@if [ ! -f "$(GO_DIR)/go.mod" ]; then \
		cd $(GO_DIR) && go mod init $(PROJECT_NAME); \
		echo "Go module initialized."; \
	else \
		echo "Go module already exists."; \
	fi

# Build both Rust library and Go binary
.PHONY: build
build: check init-go $(RUST_LIB) $(GO_BINARY)

# Build Rust library
$(RUST_LIB):
	@echo "Building Rust library..."
	cd $(RUST_DIR) && cargo build --release

# Build Go binary
$(GO_BINARY): $(RUST_LIB)
	@echo "Building Go binary..."
	cd $(GO_DIR) && go build -o benchmark_hybrid main.go

# Run the benchmark
.PHONY: run
run: build
	@echo "Running hybrid benchmark..."
	cd $(GO_DIR) && ./benchmark_hybrid

# Run with timing
.PHONY: run-timed
run-timed: build
	@echo "Running hybrid benchmark with timing..."
	cd $(GO_DIR) && time ./benchmark_hybrid

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	-cd $(RUST_DIR) && cargo clean
	-rm -f $(GO_DIR)/benchmark_hybrid
	-rm -f $(GO_DIR)/benchmark_hybrid.exe
	-rm -f $(GO_DIR)/go.sum
	@echo "Build artifacts cleaned."

# Deep clean - remove all generated files
.PHONY: deep-clean
deep-clean: clean
	@echo "Deep cleaning..."
	-rm -f $(GO_DIR)/go.mod
	@echo "Deep clean completed."

# Clean and rebuild
.PHONY: rebuild
rebuild: clean build

# Install/update dependencies
.PHONY: deps
deps: init-go
	@echo "Installing/updating Rust dependencies..."
	cd $(RUST_DIR) && cargo update
	@echo "Installing/updating Go dependencies..."
	cd $(GO_DIR) && go mod tidy

# Development mode (debug build for faster iteration)
.PHONY: dev
dev: check init-go
	@echo "Building in development mode..."
	cd $(RUST_DIR) && cargo build
	cd $(GO_DIR) && go build -o benchmark_hybrid main.go

# Run tests
.PHONY: test
test: check
	@echo "Running Rust tests..."
	cd $(RUST_DIR) && cargo test
	@echo "Running Go tests..."
	cd $(GO_DIR) && go test -v ./... 2>/dev/null || echo "No Go tests defined"

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting Rust code..."
	cd $(RUST_DIR) && cargo fmt
	@echo "Formatting Go code..."
	cd $(GO_DIR) && go fmt

# Lint code
.PHONY: lint
lint:
	@echo "Linting Rust code..."
	cd $(RUST_DIR) && cargo clippy
	@echo "Linting Go code..."
	cd $(GO_DIR) && go vet

# Profile build
.PHONY: profile
profile: check init-go
	@echo "Building with profiling..."
	cd $(RUST_DIR) && cargo build --release
	cd $(GO_DIR) && go build -o benchmark_hybrid -ldflags="-s -w" main.go

# Cross-compilation
.PHONY: linux
linux: check init-go
	@echo "Cross-compiling for Linux..."
	cd $(RUST_DIR) && cargo build --release --target=x86_64-unknown-linux-gnu
	cd $(GO_DIR) && GOOS=linux GOARCH=amd64 go build -o benchmark_hybrid_linux main.go

.PHONY: windows
windows: check init-go
	@echo "Cross-compiling for Windows..."
	cd $(RUST_DIR) && cargo build --release --target=x86_64-pc-windows-msvc
	cd $(GO_DIR) && GOOS=windows GOARCH=amd64 go build -o benchmark_hybrid.exe main.go

# Show project status
.PHONY: status
status:
	@echo "=== Project Status ==="
	@echo "Project structure:"
	@tree -I 'target|*.rs.bk' . 2>/dev/null || find . -name 'target' -prune -o -type f -print | sort
	@echo ""
	@echo "Build status:"
	@echo "Rust library: $$(if [ -f "$(RUST_LIB)" ]; then echo "✓ Built"; else echo "✗ Not built"; fi)"
	@echo "Go binary: $$(if [ -f "$(GO_BINARY)" ]; then echo "✓ Built"; else echo "✗ Not built"; fi)"
	@echo "Go module: $$(if [ -f "$(GO_DIR)/go.mod" ]; then echo "✓ Initialized"; else echo "✗ Not initialized"; fi)"

# Show file info
.PHONY: files
files:
	@echo "=== Project Files ==="
	@echo "Rust files:"
	@ls -la $(RUST_DIR)/ $(RUST_DIR)/src/ 2>/dev/null || echo "Rust directory not found"
	@echo ""
	@echo "Go files:"
	@ls -la $(GO_DIR)/ 2>/dev/null || echo "Go directory not found"

# Show help
.PHONY: help
help:
	@echo "Available targets for $(PROJECT_NAME):"
	@echo "  all         - Build everything (default)"
	@echo "  build       - Build both Rust library and Go binary"
	@echo "  run         - Build and run the benchmark"
	@echo "  run-timed   - Build and run with timing"
	@echo "  check       - Validate project structure"
	@echo "  clean       - Remove build artifacts"
	@echo "  deep-clean  - Remove ALL generated files including go.mod"
	@echo "  rebuild     - Clean and rebuild"
	@echo "  deps        - Update dependencies"
	@echo "  dev         - Development build"
	@echo "  test        - Run tests"
	@echo "  fmt         - Format code"
	@echo "  lint        - Lint code"
	@echo "  profile     - Build with profiling"
	@echo "  linux       - Cross-compile for Linux"
	@echo "  windows     - Cross-compile for Windows"
	@echo "  status      - Show project status"
	@echo "  files       - Show file information"
	@echo "  help        - Show this help message"

# Fix and rebuild after libc type issues
.PHONY: fix-libc
fix-libc: clean deps build
	@echo "Fixed libc type issues and rebuilt"
	
# Default target
.DEFAULT_GOAL := help