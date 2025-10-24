#ifndef SHAMIR_H
#define SHAMIR_H

#include <stdint.h>
#include <stddef.h>

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

BenchmarkResult benchmark_shamir_c(BenchmarkConfig config);

#endif