# Build configuration
C_DIR := c
GO_DIR := go
C_LIB := $(C_DIR)/libshamir.a
GO_BINARY := $(GO_DIR)/benchmark_c_hybrid

# Default target
.PHONY: all
all: build

# Build C library
$(C_LIB):
	@echo "Building C library..."
	cd $(C_DIR) && gcc -c -O3 -I. shamir.c -o shamir.o
	cd $(C_DIR) && ar rcs libshamir.a shamir.o

# Build Go binary
$(GO_BINARY): $(C_LIB)
	@echo "Building Go binary..."
	cd $(GO_DIR) && go build -o benchmark_c_hybrid main.go

# Build everything
.PHONY: build
build: $(GO_BINARY)

# Run benchmark
.PHONY: run
run: build
	@echo "Running Go-C hybrid benchmark..."
	cd $(GO_DIR) && ./benchmark_c_hybrid

# Clean
.PHONY: clean
clean:
	-rm -f $(C_DIR)/*.o $(C_DIR)/*.a
	-rm -f $(GO_DIR)/benchmark_c_hybrid

.PHONY: rebuild
rebuild: clean build

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all     - Build everything"
	@echo "  build   - Build C library and Go binary"
	@echo "  run     - Build and run benchmark"
	@echo "  clean   - Remove build artifacts"
	@echo "  rebuild - Clean and rebuild"

.DEFAULT_GOAL := help