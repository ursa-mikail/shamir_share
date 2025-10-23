package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/vault/shamir"
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

func benchmarkShamir(data []byte, parts, threshold int, rounds int) BenchmarkResult {
	var totalSplitTime, totalCombineTime time.Duration
	successCount := 0

	for i := 0; i < rounds; i++ {
		// Benchmark split operation
		startSplit := time.Now()
		shares, err := shamir.Split(data, parts, threshold)
		splitTime := time.Since(startSplit)

		if err != nil {
			log.Printf("Error in split operation: %v", err)
			continue
		}

		// Benchmark combine operation
		startCombine := time.Now()
		recovered, err := shamir.Combine(shares[:threshold])
		combineTime := time.Since(startCombine)

		if err != nil {
			log.Printf("Error in combine operation: %v", err)
			continue
		}

		// Verify data integrity
		if hex.EncodeToString(data) == hex.EncodeToString(recovered) {
			successCount++
			totalSplitTime += splitTime
			totalCombineTime += combineTime
		} else {
			log.Printf("Data mismatch in round %d", i)
		}
	}

	totalTime := totalSplitTime + totalCombineTime
	throughput := float64(successCount*len(data)) / totalTime.Seconds()

	return BenchmarkResult{
		TotalTime:    totalTime,
		SplitTime:    totalSplitTime,
		CombineTime:  totalCombineTime,
		Throughput:   throughput,
		SuccessCount: successCount,
	}
}

func main() {
	fmt.Println("=== Go Shamir Secret Sharing Benchmark (Hashicorp Vault) ===")
	
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
		result := benchmarkShamir(data, config.parts, config.threshold, config.rounds)

		fmt.Printf("Success rate: %d/%d (%.2f%%)\n", 
			result.SuccessCount, config.rounds, 
			float64(result.SuccessCount)/float64(config.rounds)*100)
		fmt.Printf("Total time: %v\n", result.TotalTime)
		fmt.Printf("Average split time: %v\n", result.SplitTime/time.Duration(result.SuccessCount))
		fmt.Printf("Average combine time: %v\n", result.CombineTime/time.Duration(result.SuccessCount))
		fmt.Printf("Throughput: %.2f bytes/sec\n", result.Throughput)
		fmt.Printf("Operations/sec: %.2f\n", float64(result.SuccessCount)/result.TotalTime.Seconds())
	}
}

/*
go mod init benchmark_standalone
go mod tidy
go run main.go

% go run main.go
go: downloading go1.25.1 (darwin/arm64)
=== Go Shamir Secret Sharing Benchmark (Hashicorp Vault) ===

--- Benchmark: 32 bytes, 3/5 scheme, 1000 rounds ---
Success rate: 1000/1000 (100.00%)
Total time: 70.437024ms
Average split time: 13.783µs
Average combine time: 56.653µs
Throughput: 454306.53 bytes/sec
Operations/sec: 14197.08

--- Benchmark: 256 bytes, 3/5 scheme, 500 rounds ---
Success rate: 500/500 (100.00%)
Total time: 172.911001ms
Average split time: 54.443µs
Average combine time: 291.378µs
Throughput: 740265.22 bytes/sec
Operations/sec: 2891.66

--- Benchmark: 1024 bytes, 3/5 scheme, 200 rounds ---
Success rate: 200/200 (100.00%)
Total time: 275.621706ms
Average split time: 213.964µs
Average combine time: 1.164144ms
Throughput: 743047.43 bytes/sec
Operations/sec: 725.63

--- Benchmark: 4096 bytes, 3/5 scheme, 100 rounds ---
Success rate: 100/100 (100.00%)
Total time: 552.574714ms
Average split time: 853.532µs
Average combine time: 4.672214ms
Throughput: 741257.23 bytes/sec
Operations/sec: 180.97

--- Benchmark: 32 bytes, 7/10 scheme, 500 rounds ---
Success rate: 500/500 (100.00%)
Total time: 137.881296ms
Average split time: 29.315µs
Average combine time: 246.447µs
Throughput: 116041.85 bytes/sec
Operations/sec: 3626.31

*/