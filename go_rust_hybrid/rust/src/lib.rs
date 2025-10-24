use std::time::{Duration, Instant};
use rand::RngCore;
use sharks::{Sharks, Share};
use libc::c_double;

#[repr(C)]
pub struct BenchmarkResult {
    total_time_ns: u64,
    split_time_ns: u64,
    combine_time_ns: u64,
    throughput: c_double,
    success_count: u64,
}

#[repr(C)]
pub struct BenchmarkConfig {
    data: *const u8,
    data_len: u64,
    parts: u8,
    threshold: u8,
    rounds: u64,
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
        total_time_ns: total_time.as_nanos() as u64,
        split_time_ns: total_split_time.as_nanos() as u64,
        combine_time_ns: total_combine_time.as_nanos() as u64,
        throughput,
        success_count: success_count as u64,
    }
}

#[no_mangle]
pub extern "C" fn benchmark_shamir_ffi(config: BenchmarkConfig) -> BenchmarkResult {
    let data_slice = unsafe {
        std::slice::from_raw_parts(config.data, config.data_len as usize)
    };
    
    benchmark_shamir(
        data_slice,
        config.parts,
        config.threshold,
        config.rounds as usize,
    )
}