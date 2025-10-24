1. Build the Rust library:
```
cd go_rust_hybrid/rust
cargo build --release
```

2. Build and run the Go benchmark:
```
cd ../go
go build -o benchmark_hybrid main.go
./benchmark_hybrid
```

3. Just Use Makefile

```
# First time setup (just initialize Go module)
make deps

# Check if everything is ready
make check

# Build and run
make run

# Development workflow
make dev && make run

# Check project status
make status

# Clean up
make clean        # Keep go.mod
make deep-clean   # Remove everything including go.mod
```

```
Go main() 
    ↓ (calls)
benchmarkShamirHybrid() 
    ↓ (via CGO)
C.benchmark_shamir_ffi() 
    ↓ (crosses language boundary)
Rust benchmark_shamir_ffi() 
    ↓ (calls)
Rust benchmark_shamir() 
    ↓ (performs Shamir operations)
Rust BenchmarkResult 
    ↓ (returns)
Go BenchmarkResult
    ↓ (displayed)
Console Output
```

### Key Technical Points
From Go → Rust:
```
C.CBytes() converts Go slice to C pointer

Go structs map to C structs via cgo

extern "C" in Rust makes it callable from C

#[no_mangle] keeps the function name recognizable
```

From Rust → Go:
```
#[repr(C)] ensures C-compatible memory layout

Rust returns struct directly to Go via C ABI

Go receives the struct and converts fields to Go types
```


```
% make run       
=== Checking project structure ===
✓ Project structure is valid
Initializing Go module...
Go module already exists.
Building Rust library...
cd rust && cargo build --release
   Compiling proc-macro2 v1.0.103
   Compiling libc v0.2.177
   Compiling unicode-ident v1.0.20
   Compiling quote v1.0.41
   Compiling zerocopy v0.8.27
   Compiling cfg-if v1.0.4
   Compiling ahash v0.4.8
   Compiling hashbrown v0.9.1
   Compiling getrandom v0.2.16
   Compiling syn v2.0.108
   Compiling rand_core v0.6.4
   Compiling ppv-lite86 v0.2.21
   Compiling rand_chacha v0.3.1
   Compiling rand v0.8.5
   Compiling zeroize_derive v1.4.2
   Compiling zeroize v1.8.2
   Compiling sharks v0.5.0
   Compiling shamir_benchmark v0.1.0 (/Users/chanfamily/ursa/git/shamir_share/go_rust_hybrid/rust)
warning: function `generate_random_data` is never used
  --> src/lib.rs:24:4
   |
24 | fn generate_random_data(size: usize) -> Vec<u8> {
   |    ^^^^^^^^^^^^^^^^^^^^
   |
   = note: `#[warn(dead_code)]` on by default

warning: `shamir_benchmark` (lib) generated 1 warning
    Finished `release` profile [optimized] target(s) in 6.12s
Building Go binary...
cd go && go build -o benchmark_hybrid main.go
Running hybrid benchmark...
cd go && ./benchmark_hybrid
=== Go-Rust Hybrid Shamir Secret Sharing Benchmark ===

--- Benchmark: 32 bytes, 3/5 scheme, 1000 rounds ---
Success rate: 1000/1000 (100.00%)
Total time: 2.691716ms
Average split time: 1.753µs
Average combine time: 937ns
Throughput: 11888327.00 bytes/sec
Operations/sec: 371510.22

--- Benchmark: 256 bytes, 3/5 scheme, 500 rounds ---
Success rate: 500/500 (100.00%)
Total time: 14.164783ms
Average split time: 18.688µs
Average combine time: 9.64µs
Throughput: 9036495.65 bytes/sec
Operations/sec: 35298.81

--- Benchmark: 1024 bytes, 3/5 scheme, 200 rounds ---
Success rate: 200/200 (100.00%)
Total time: 26.286119ms
Average split time: 87.764µs
Average combine time: 43.666µs
Throughput: 7791184.39 bytes/sec
Operations/sec: 7608.58

--- Benchmark: 4096 bytes, 3/5 scheme, 100 rounds ---
Success rate: 100/100 (100.00%)
Total time: 36.482083ms
Average split time: 247.84µs
Average combine time: 116.98µs
Throughput: 11227429.09 bytes/sec
Operations/sec: 2741.07

--- Benchmark: 32 bytes, 7/10 scheme, 500 rounds ---
Success rate: 500/500 (100.00%)
Total time: 4.562823ms
Average split time: 5.268µs
Average combine time: 3.857µs
Throughput: 3506601.07 bytes/sec
Operations/sec: 109581.28
```

Key Features:
```
Go-Rust FFI: Uses CGO to call Rust functions from Go

Memory Safety: Proper memory management with C.CBytes and C.free

Type Conversion: Converts between Go and Rust types including slices and structs

Performance: Rust performs the computationally intensive Shamir operations

Flexibility: Go handles the benchmarking framework and result display
```