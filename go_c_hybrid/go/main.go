package main

/*
#cgo CFLAGS: -I../c
#cgo LDFLAGS: -L../c -lshamir -lm
#include "shamir.h"
#include <stdlib.h>  // Add this for free()
*/
import "C"
import (
	"crypto/rand"
	"fmt"
	"log"
	"time"
	"unsafe"
)

type BenchmarkResult struct {
	TotalTime    time.Duration
	SplitTime    time.Duration
	CombineTime  time.Duration
	Throughput   float64
	SuccessCount int
}

func generateRandomData(size int) []byte {
	data := make([]byte, size)
	_, err := rand.Read(data)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

func benchmarkShamirC(data []byte, parts, threshold int, rounds int) BenchmarkResult {
	config := C.BenchmarkConfig{
		data:      (*C.uint8_t)(C.CBytes(data)),
		data_len:  C.uint64_t(len(data)),
		parts:     C.uint8_t(parts),
		threshold: C.uint8_t(threshold),
		rounds:    C.uint64_t(rounds),
	}
	defer C.free(unsafe.Pointer(config.data))

	result := C.benchmark_shamir_c(config)

	return BenchmarkResult{
		TotalTime:    time.Duration(result.total_time_ns) * time.Nanosecond,
		SplitTime:    time.Duration(result.split_time_ns) * time.Nanosecond,
		CombineTime:  time.Duration(result.combine_time_ns) * time.Nanosecond,
		Throughput:   float64(result.throughput),
		SuccessCount: int(result.success_count),
	}
}

func main() {
	fmt.Println("=== Go-C Hybrid Shamir Secret Sharing Benchmark ===")
	
	configs := []struct {
		dataSize   int
		parts      int
		threshold  int
		rounds     int
	}{
		{32, 5, 3, 1000},    // Standard 32-byte secret
		{48, 5, 3, 500},     // Medium secret
		{64, 5, 3, 200},     // Larger secret
		{32, 10, 7, 500},    // More parts, higher threshold
	}

	for _, config := range configs {
		fmt.Printf("\n--- Benchmark: %d bytes, %d/%d scheme, %d rounds ---\n",
			config.dataSize, config.threshold, config.parts, config.rounds)

		data := generateRandomData(config.dataSize)
		result := benchmarkShamirC(data, config.parts, config.threshold, config.rounds)

		fmt.Printf("Success rate: %d/%d (%.2f%%)\n", 
			result.SuccessCount, config.rounds, 
			float64(result.SuccessCount)/float64(config.rounds)*100)
		fmt.Printf("Total time: %v\n", result.TotalTime)
		if result.SuccessCount > 0 {
			fmt.Printf("Average split time: %v\n", result.SplitTime/time.Duration(result.SuccessCount))
			fmt.Printf("Average combine time: %v\n", result.CombineTime/time.Duration(result.SuccessCount))
		} else {
			fmt.Printf("Average split time: N/A (no successful runs)\n")
			fmt.Printf("Average combine time: N/A (no successful runs)\n")
		}
		fmt.Printf("Throughput: %.2f bytes/sec\n", result.Throughput)
		if result.TotalTime > 0 {
			fmt.Printf("Operations/sec: %.2f\n", float64(result.SuccessCount)/result.TotalTime.Seconds())
		} else {
			fmt.Printf("Operations/sec: N/A\n")
		}
	}
}

/*
% make run                        
Building C library...
cd c && gcc -c -O3 -I. shamir.c -o shamir.o
cd c && ar rcs libshamir.a shamir.o
Building Go binary...
cd go && go build -o benchmark_c_hybrid main.go
Running Go-C hybrid benchmark...
cd go && ./benchmark_c_hybrid
=== Go-C Hybrid Shamir Secret Sharing Benchmark ===

--- Benchmark: 32 bytes, 3/5 scheme, 1000 rounds ---
Success rate: 1000/1000 (100.00%)
Total time: 8.846907s
Average split time: 8.84642ms
Average combine time: 487ns
Throughput: 3617.08 bytes/sec
Operations/sec: 113.03

--- Benchmark: 48 bytes, 3/5 scheme, 500 rounds ---
Success rate: 500/500 (100.00%)
Total time: 6.619224s
Average split time: 13.237736ms
Average combine time: 712ns
Throughput: 3625.80 bytes/sec
Operations/sec: 75.54

--- Benchmark: 64 bytes, 3/5 scheme, 200 rounds ---
Success rate: 200/200 (100.00%)
Total time: 3.527311s
Average split time: 17.635605ms
Average combine time: 950ns
Throughput: 3628.83 bytes/sec
Operations/sec: 56.70

--- Benchmark: 32 bytes, 7/10 scheme, 500 rounds ---
Success rate: 500/500 (100.00%)
Total time: 13.274243s
Average split time: 26.54598ms
Average combine time: 2.506µs
Throughput: 1205.34 bytes/sec
Operations/sec: 37.67

*/