use std::time::{Duration, Instant};
use rand::RngCore;
use sharks::{Sharks, Share};

#[derive(Debug)]
struct BenchmarkResult {
    total_time: Duration,
    split_time: Duration,
    combine_time: Duration,
    throughput: f64,
    success_count: usize,
}

fn generate_random_data(size: usize) -> Vec<u8> {
    let mut data = vec![0u8; size];
    rand::thread_rng().fill_bytes(&mut data);
    data
}

fn benchmark_shamir(data: &[u8], parts: u8, threshold: u8, rounds: usize) -> BenchmarkResult {
    let mut total_split_time = Duration::new(0, 0);
    let mut total_combine_time = Duration::new(0, 0);
    let mut success_count = 0;

    for i in 0..rounds {
        // Benchmark split operation
        let split_start = Instant::now();
        
        let sharks = Sharks(threshold);
        let dealer = sharks.dealer(data);
        let shares: Vec<Share> = dealer.take(parts as usize).collect();
        
        let split_time = split_start.elapsed();

        // Take only threshold number of shares for reconstruction
        let reconstruction_shares: Vec<Share> = shares.into_iter().take(threshold as usize).collect();

        // Benchmark combine operation
        let combine_start = Instant::now();
        let recovered_secret = match sharks.recover(reconstruction_shares.iter()) {
            Ok(secret) => secret,
            Err(e) => {
                eprintln!("Error in combine operation round {}: {:?}", i, e);
                continue;
            }
        };
        let combine_time = combine_start.elapsed();

        // Verify data integrity
        if data == recovered_secret.as_slice() {
            success_count += 1;
            total_split_time += split_time;
            total_combine_time += combine_time;
        } else {
            eprintln!("Data mismatch in round {}", i);
        }
    }

    let total_time = total_split_time + total_combine_time;
    let throughput = if total_time.as_secs_f64() > 0.0 {
        (success_count * data.len()) as f64 / total_time.as_secs_f64()
    } else {
        0.0
    };

    BenchmarkResult {
        total_time,
        split_time: total_split_time,
        combine_time: total_combine_time,
        throughput,
        success_count,
    }
}

fn main() {
    println!("=== Rust Shamir Secret Sharing Benchmark (sharks crate) ===");
    
    let configs = vec![
        (32, 5, 3, 1000),    // Small secret
        (256, 5, 3, 500),    // Medium secret
        (1024, 5, 3, 200),   // Large secret
        (32, 10, 7, 500),    // More parts, higher threshold
    ];

    for (data_size, parts, threshold, rounds) in configs {
        println!("\n--- Benchmark: {} bytes, {}/{} scheme, {} rounds ---",
                 data_size, threshold, parts, rounds);

        let data = generate_random_data(data_size);
        let result = benchmark_shamir(&data, parts, threshold, rounds);

        let success_rate = if rounds > 0 {
            (result.success_count as f64 / rounds as f64) * 100.0
        } else {
            0.0
        };

        let avg_split = if result.success_count > 0 {
            result.split_time / result.success_count as u32
        } else {
            Duration::new(0, 0)
        };

        let avg_combine = if result.success_count > 0 {
            result.combine_time / result.success_count as u32
        } else {
            Duration::new(0, 0)
        };

        let ops_per_sec = if result.total_time.as_secs_f64() > 0.0 {
            result.success_count as f64 / result.total_time.as_secs_f64()
        } else {
            0.0
        };

        println!("Success rate: {}/{} ({:.2}%)", 
                 result.success_count, rounds, success_rate);
        println!("Total time: {:?}", result.total_time);
        println!("Average split time: {:?}", avg_split);
        println!("Average combine time: {:?}", avg_combine);
        println!("Throughput: {:.2} bytes/sec", result.throughput);
        println!("Operations/sec: {:.2}", ops_per_sec);
    }
}

/*
brew install libsodium

cargo clean
cargo build --release
cargo run --release

 % cargo run --release
    Finished `release` profile [optimized] target(s) in 0.02s
     Running `target/release/shamir_benchmark`
=== Rust Shamir Secret Sharing Benchmark (sharks crate) ===

--- Benchmark: 32 bytes, 3/5 scheme, 1000 rounds ---
Success rate: 1000/1000 (100.00%)
Total time: 11.215925ms
Average split time: 8.217µs
Average combine time: 2.998µs
Throughput: 2853086.13 bytes/sec
Operations/sec: 89158.94

--- Benchmark: 256 bytes, 3/5 scheme, 500 rounds ---
Success rate: 500/500 (100.00%)
Total time: 28.130991ms
Average split time: 42.207µs
Average combine time: 14.054µs
Throughput: 4550141.87 bytes/sec
Operations/sec: 17773.99

--- Benchmark: 1024 bytes, 3/5 scheme, 200 rounds ---
Success rate: 200/200 (100.00%)
Total time: 29.851496ms
Average split time: 113.303µs
Average combine time: 35.953µs
Throughput: 6860627.69 bytes/sec
Operations/sec: 6699.83

--- Benchmark: 32 bytes, 7/10 scheme, 500 rounds ---
Success rate: 500/500 (100.00%)
Total time: 6.925431ms
Average split time: 9.448µs
Average combine time: 4.402µs
Throughput: 2310325.52 bytes/sec
Operations/sec: 72197.67

*/