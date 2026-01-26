Key Differences from Go-Rust:
```
Simpler FFI: CGO is designed for C interoperability, so no complex type conversions needed

Direct C Calls: Go calls C functions directly without name mangling concerns

Memory Management: C memory management is more manual but straightforward

Performance: C can be faster for cryptographic operations due to minimal abstraction
```

```

# Build and run
make clean    # Start fresh
make run
```

```
Go main()
    ↓ (calls)
benchmarkShamirC()
    ↓ (via CGO)
C.benchmark_shamir_c()  ← Direct C function call!
    ↓ (executes in C)
C performs Shamir operations
    ↓ (returns)
Go receives C struct
    ↓ (displays)
Results
```
