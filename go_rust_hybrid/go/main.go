package main

/*
#cgo LDFLAGS: -L../rust/target/release -lshamir_benchmark -ldl -lm -lc
#include <stdint.h>
#include <stdlib.h>

typedef struct {
    uint64_t total_time_ns;
    uint64_t split_time_ns;
    uint64_t combine_time_ns;
    double throughput;
    uint64_t success_count;
} BenchmarkResult;

typedef struct {
    uint8_t* data;
    uint64_t data_len;
    uint8_t parts;
    uint8_t threshold;
    uint64_t rounds;
} BenchmarkConfig;

extern BenchmarkResult benchmark_shamir_ffi(BenchmarkConfig config);
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

func benchmarkShamirHybrid(data []byte, parts, threshold int, rounds int) BenchmarkResult {
	config := C.BenchmarkConfig{
		data:      (*C.uint8_t)(C.CBytes(data)),
		data_len:  C.uint64_t(len(data)),
		parts:     C.uint8_t(parts),
		threshold: C.uint8_t(threshold),
		rounds:    C.uint64_t(rounds),
	}
	defer C.free(unsafe.Pointer(config.data))

	result := C.benchmark_shamir_ffi(config)

	return BenchmarkResult{
		TotalTime:    time.Duration(result.total_time_ns) * time.Nanosecond,
		SplitTime:    time.Duration(result.split_time_ns) * time.Nanosecond,
		CombineTime:  time.Duration(result.combine_time_ns) * time.Nanosecond,
		Throughput:   float64(result.throughput),
		SuccessCount: int(result.success_count),
	}
}

func main() {
	fmt.Println("=== Go-Rust Hybrid Shamir Secret Sharing Benchmark ===")
	
	configs := []struct {
		dataSize   int
		parts      int
		threshold  int
		rounds     int
	}{
		{32, 5, 3, 1000},    // Small secret
		{256, 5, 3, 500},    // Medium secret
		{1024, 5, 3, 200},   // Large secret
		{4096, 5, 3, 100},   // Very large secret
		{32, 10, 7, 500},    // More parts, higher threshold
	}

	for _, config := range configs {
		fmt.Printf("\n--- Benchmark: %d bytes, %d/%d scheme, %d rounds ---\n",
			config.dataSize, config.threshold, config.parts, config.rounds)

		data := generateRandomData(config.dataSize)
		result := benchmarkShamirHybrid(data, config.parts, config.threshold, config.rounds)

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