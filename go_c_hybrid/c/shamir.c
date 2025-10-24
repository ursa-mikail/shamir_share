#include "shamir.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>

// Simple GF(256) arithmetic for Shamir's Secret Sharing
#define GF_SIZE 256
#define PRIMITIVE_POLY 0x11D  // x^8 + x^4 + x^3 + x^2 + 1

static unsigned char gf_exp[512];
static unsigned char gf_log[256];
static int gf_tables_initialized = 0;

void init_gf_tables() {
    if (gf_tables_initialized) return;
    
    unsigned char x = 1;
    for (int i = 0; i < 255; i++) {
        gf_exp[i] = x;
        gf_log[x] = i;
        x = (x << 1) ^ ((x & 0x80) ? PRIMITIVE_POLY : 0);
    }
    
    for (int i = 255; i < 512; i++) {
        gf_exp[i] = gf_exp[i - 255];
    }
    
    gf_tables_initialized = 1;
}

unsigned char gf_add(unsigned char a, unsigned char b) {
    return a ^ b;
}

unsigned char gf_mul(unsigned char a, unsigned char b) {
    if (a == 0 || b == 0) return 0;
    return gf_exp[gf_log[a] + gf_log[b]];
}

unsigned char gf_pow(unsigned char a, unsigned char n) {
    if (a == 0) return 0;
    if (n == 0) return 1;
    return gf_exp[(gf_log[a] * n) % 255];
}

// Simple random bytes implementation
static void randombytes(uint8_t *buf, size_t len) {
    FILE *f = fopen("/dev/urandom", "rb");
    if (f != NULL) {
        fread(buf, 1, len, f);
        fclose(f);
    } else {
        // Fallback
        for (size_t i = 0; i < len; i++) {
            buf[i] = rand() % 256;
        }
    }
}

static uint64_t get_time_ns() {
    struct timespec ts;
    clock_gettime(CLOCK_MONOTONIC, &ts);
    return (uint64_t)ts.tv_sec * 1000000000ULL + (uint64_t)ts.tv_nsec;
}

// Basic Shamir's Secret Sharing implementation
void shamir_split(uint8_t **shares, const uint8_t *secret, int secret_len, 
                  int total_shares, int threshold) {
    init_gf_tables();
    
    for (int i = 0; i < total_shares; i++) {
        // Share index (x-value)
        shares[i][0] = i + 1;
    }
    
    // For each byte position in the secret
    for (int pos = 0; pos < secret_len; pos++) {
        // Generate random coefficients for the polynomial
        uint8_t coefficients[threshold];
        coefficients[0] = secret[pos];  // Constant term is the secret
        
        for (int j = 1; j < threshold; j++) {
            randombytes(&coefficients[j], 1);
        }
        
        // Evaluate polynomial for each share
        for (int i = 0; i < total_shares; i++) {
            uint8_t x = shares[i][0];
            uint8_t y = 0;
            
            for (int j = 0; j < threshold; j++) {
                y = gf_add(y, gf_mul(coefficients[j], gf_pow(x, j)));
            }
            
            shares[i][pos + 1] = y;
        }
    }
}

int shamir_combine(uint8_t *secret, uint8_t **shares, int secret_len, 
                   int num_shares) {
    init_gf_tables();
    
    for (int pos = 0; pos < secret_len; pos++) {
        // Lagrange interpolation
        uint8_t result = 0;
        
        for (int i = 0; i < num_shares; i++) {
            uint8_t xi = shares[i][0];
            uint8_t yi = shares[i][pos + 1];
            
            uint8_t numerator = 1;
            uint8_t denominator = 1;
            
            for (int j = 0; j < num_shares; j++) {
                if (i == j) continue;
                
                uint8_t xj = shares[j][0];
                numerator = gf_mul(numerator, xj);
                denominator = gf_mul(denominator, gf_add(xi, xj));
            }
            
            uint8_t li = gf_mul(numerator, gf_pow(denominator, 254)); // 254 is -1 in GF(256)
            result = gf_add(result, gf_mul(yi, li));
        }
        
        secret[pos] = result;
    }
    
    return 0; // Success
}

BenchmarkResult benchmark_shamir_c(BenchmarkConfig config) {
    // Validate input
    if (config.parts > 255 || config.threshold > config.parts || 
        config.threshold < 1 || config.data_len > 255) {
        return (BenchmarkResult){0, 0, 0, 0.0, 0};
    }

    // Initialize random seed
    srand((unsigned int)time(NULL));
    
    // Allocate memory for shares (each share: 1 byte index + data_len bytes)
    int share_size = config.data_len + 1;
    uint8_t **shares = malloc(config.parts * sizeof(uint8_t*));
    for (int i = 0; i < config.parts; i++) {
        shares[i] = malloc(share_size);
    }
    
    uint8_t *recovered = malloc(config.data_len);
    uint8_t **used_shares = malloc(config.threshold * sizeof(uint8_t*));
    for (int i = 0; i < config.threshold; i++) {
        used_shares[i] = malloc(share_size);
    }
    
    uint64_t total_split_time = 0;
    uint64_t total_combine_time = 0;
    uint64_t success_count = 0;

    for (uint64_t round = 0; round < config.rounds; round++) {
        // Benchmark split operation
        uint64_t split_start = get_time_ns();
        
        shamir_split(shares, config.data, config.data_len, 
                    config.parts, config.threshold);
        
        uint64_t split_time = get_time_ns() - split_start;

        // Select threshold number of shares for reconstruction
        for (int i = 0; i < config.threshold; i++) {
            memcpy(used_shares[i], shares[i], share_size);
        }

        // Benchmark combine operation
        uint64_t combine_start = get_time_ns();
        
        int result = shamir_combine(recovered, used_shares, 
                                   config.data_len, config.threshold);
        
        uint64_t combine_time = get_time_ns() - combine_start;

        // Verify data integrity
        if (result == 0 && memcmp(config.data, recovered, config.data_len) == 0) {
            success_count++;
            total_split_time += split_time;
            total_combine_time += combine_time;
        }
    }

    // Cleanup
    for (int i = 0; i < config.parts; i++) {
        free(shares[i]);
    }
    free(shares);
    
    for (int i = 0; i < config.threshold; i++) {
        free(used_shares[i]);
    }
    free(used_shares);
    
    free(recovered);

    uint64_t total_time = total_split_time + total_combine_time;
    double throughput = (total_time > 0) ? 
        (double)(success_count * config.data_len) / ((double)total_time / 1e9) : 0.0;

    return (BenchmarkResult){
        .total_time_ns = total_time,
        .split_time_ns = total_split_time,
        .combine_time_ns = total_combine_time,
        .throughput = throughput,
        .success_count = success_count
    };
}