#!/usr/bin/env python3
"""
Performance comparison analysis between Rust (sharks) and Go (Hashicorp Vault) 
Shamir Secret Sharing implementations
"""

import matplotlib.pyplot as plt
import numpy as np

# Your benchmark results
results = {
    'rust': {
        '32b_5_3': {'split_us': 8.217, 'combine_us': 2.998, 'throughput_mb': 2.85, 'ops_sec': 89158.94},
        '256b_5_3': {'split_us': 42.207, 'combine_us': 14.054, 'throughput_mb': 4.55, 'ops_sec': 17773.99},
        '1024b_5_3': {'split_us': 113.303, 'combine_us': 35.953, 'throughput_mb': 6.86, 'ops_sec': 6699.83},
        '32b_10_7': {'split_us': 9.448, 'combine_us': 4.402, 'throughput_mb': 2.31, 'ops_sec': 72197.67},
    },
    'go': {
        '32b_5_3': {'split_us': 13.783, 'combine_us': 56.653, 'throughput_mb': 0.45, 'ops_sec': 14197.08},
        '256b_5_3': {'split_us': 54.443, 'combine_us': 291.378, 'throughput_mb': 0.74, 'ops_sec': 2891.66},
        '1024b_5_3': {'split_us': 213.964, 'combine_us': 1164.144, 'throughput_mb': 0.74, 'ops_sec': 725.63},
        '4096b_5_3': {'split_us': 853.532, 'combine_us': 4672.214, 'throughput_mb': 0.74, 'ops_sec': 180.97},
        '32b_10_7': {'split_us': 29.315, 'combine_us': 246.447, 'throughput_mb': 0.12, 'ops_sec': 3626.31},
    }
}

def plot_comprehensive_comparison():
    # Create a 2x3 subplot layout
    fig, axes = plt.subplots(2, 3, figsize=(18, 12))
    fig.suptitle('Rust (sharks) vs Go (Hashicorp Vault) Shamir Secret Sharing Performance Comparison', 
                 fontsize=16, fontweight='bold')
    
    # Data preparation for common scenarios
    common_scenarios = ['32B\n5/3', '256B\n5/3', '1024B\n5/3']
    rust_common = ['32b_5_3', '256b_5_3', '1024b_5_3']
    go_common = ['32b_5_3', '256b_5_3', '1024b_5_3']
    
    # Plot 1: Split Time Comparison (common scenarios)
    x_common = np.arange(len(common_scenarios))
    width = 0.35
    
    rust_split_common = [results['rust'][scenario]['split_us'] for scenario in rust_common]
    go_split_common = [results['go'][scenario]['split_us'] for scenario in go_common]
    
    bars1 = axes[0, 0].bar(x_common - width/2, rust_split_common, width, label='Rust', color='#FF6B6B', alpha=0.8)
    bars2 = axes[0, 0].bar(x_common + width/2, go_split_common, width, label='Go', color='#4ECDC4', alpha=0.8)
    
    axes[0, 0].set_ylabel('Time (microseconds)', fontweight='bold')
    axes[0, 0].set_title('Split Operation Latency', fontweight='bold')
    axes[0, 0].set_xticks(x_common)
    axes[0, 0].set_xticklabels(common_scenarios)
    axes[0, 0].legend()
    axes[0, 0].grid(True, alpha=0.3)
    axes[0, 0].set_axisbelow(True)
    
    # Add value labels on bars
    for bar in bars1:
        height = bar.get_height()
        axes[0, 0].text(bar.get_x() + bar.get_width()/2., height,
                       f'{height:.1f}µs', ha='center', va='bottom', fontsize=8)
    for bar in bars2:
        height = bar.get_height()
        axes[0, 0].text(bar.get_x() + bar.get_width()/2., height,
                       f'{height:.1f}µs', ha='center', va='bottom', fontsize=8)
    
    # Plot 2: Combine Time Comparison (common scenarios)
    rust_combine_common = [results['rust'][scenario]['combine_us'] for scenario in rust_common]
    go_combine_common = [results['go'][scenario]['combine_us'] for scenario in go_common]
    
    bars1 = axes[0, 1].bar(x_common - width/2, rust_combine_common, width, label='Rust', color='#FF6B6B', alpha=0.8)
    bars2 = axes[0, 1].bar(x_common + width/2, go_combine_common, width, label='Go', color='#4ECDC4', alpha=0.8)
    
    axes[0, 1].set_ylabel('Time (microseconds)', fontweight='bold')
    axes[0, 1].set_title('Combine Operation Latency', fontweight='bold')
    axes[0, 1].set_xticks(x_common)
    axes[0, 1].set_xticklabels(common_scenarios)
    axes[0, 1].legend()
    axes[0, 1].grid(True, alpha=0.3)
    axes[0, 1].set_axisbelow(True)
    
    # Add value labels on bars
    for bar in bars1:
        height = bar.get_height()
        axes[0, 1].text(bar.get_x() + bar.get_width()/2., height,
                       f'{height:.1f}µs', ha='center', va='bottom', fontsize=8)
    for bar in bars2:
        height = bar.get_height()
        axes[0, 1].text(bar.get_x() + bar.get_width()/2., height,
                       f'{height:.1f}µs', ha='center', va='bottom', fontsize=8)
    
    # Plot 3: Throughput Comparison (common scenarios)
    rust_throughput_common = [results['rust'][scenario]['throughput_mb'] for scenario in rust_common]
    go_throughput_common = [results['go'][scenario]['throughput_mb'] for scenario in go_common]
    
    bars1 = axes[0, 2].bar(x_common - width/2, rust_throughput_common, width, label='Rust', color='#FF6B6B', alpha=0.8)
    bars2 = axes[0, 2].bar(x_common + width/2, go_throughput_common, width, label='Go', color='#4ECDC4', alpha=0.8)
    
    axes[0, 2].set_ylabel('Throughput (MB/s)', fontweight='bold')
    axes[0, 2].set_title('Throughput Comparison', fontweight='bold')
    axes[0, 2].set_xticks(x_common)
    axes[0, 2].set_xticklabels(common_scenarios)
    axes[0, 2].legend()
    axes[0, 2].grid(True, alpha=0.3)
    axes[0, 2].set_axisbelow(True)
    
    # Add value labels on bars
    for bar in bars1:
        height = bar.get_height()
        axes[0, 2].text(bar.get_x() + bar.get_width()/2., height,
                       f'{height:.2f}MB/s', ha='center', va='bottom', fontsize=8)
    for bar in bars2:
        height = bar.get_height()
        axes[0, 2].text(bar.get_x() + bar.get_width()/2., height,
                       f'{height:.2f}MB/s', ha='center', va='bottom', fontsize=8)
    
    # Plot 4: Operations per Second Comparison
    rust_ops_common = [results['rust'][scenario]['ops_sec'] for scenario in rust_common]
    go_ops_common = [results['go'][scenario]['ops_sec'] for scenario in go_common]
    
    bars1 = axes[1, 0].bar(x_common - width/2, rust_ops_common, width, label='Rust', color='#FF6B6B', alpha=0.8)
    bars2 = axes[1, 0].bar(x_common + width/2, go_ops_common, width, label='Go', color='#4ECDC4', alpha=0.8)
    
    axes[1, 0].set_ylabel('Operations per Second', fontweight='bold')
    axes[1, 0].set_title('Operations per Second', fontweight='bold')
    axes[1, 0].set_xticks(x_common)
    axes[1, 0].set_xticklabels(common_scenarios)
    axes[1, 0].legend()
    axes[1, 0].grid(True, alpha=0.3)
    axes[1, 0].set_axisbelow(True)
    axes[1, 0].set_yscale('log')  # Log scale due to large differences
    
    # Plot 5: Performance Improvement (Speedup)
    speedup_split = [go/rust for go, rust in zip(go_split_common, rust_split_common)]
    speedup_combine = [go/rust for go, rust in zip(go_combine_common, rust_combine_common)]
    speedup_throughput = [rust/go for rust, go in zip(rust_throughput_common, go_throughput_common)]
    
    x_speedup = np.arange(3)
    width_speedup = 0.25
    
    axes[1, 1].bar(x_speedup - width_speedup, speedup_split, width_speedup, label='Split', color='#45B7D1', alpha=0.8)
    axes[1, 1].bar(x_speedup, speedup_combine, width_speedup, label='Combine', color='#F4A259', alpha=0.8)
    axes[1, 1].bar(x_speedup + width_speedup, speedup_throughput, width_speedup, label='Throughput', color='#6A0572', alpha=0.8)
    
    axes[1, 1].set_ylabel('Speedup Factor (Go/Rust)', fontweight='bold')
    axes[1, 1].set_title('Rust Performance Advantage\n(Higher = Better)', fontweight='bold')
    axes[1, 1].set_xticks(x_speedup)
    axes[1, 1].set_xticklabels(common_scenarios)
    axes[1, 1].legend()
    axes[1, 1].grid(True, alpha=0.3)
    axes[1, 1].set_axisbelow(True)
    
    # Add value labels
    for i, (split, combine, throughput) in enumerate(zip(speedup_split, speedup_combine, speedup_throughput)):
        axes[1, 1].text(i - width_speedup, split, f'{split:.1f}x', ha='center', va='bottom', fontsize=8)
        axes[1, 1].text(i, combine, f'{combine:.1f}x', ha='center', va='bottom', fontsize=8)
        axes[1, 1].text(i + width_speedup, throughput, f'{throughput:.1f}x', ha='center', va='bottom', fontsize=8)
    
    # Plot 6: Scheme Complexity Comparison (32B with different schemes)
    scheme_scenarios = ['32B\n5/3', '32B\n10/7']
    rust_schemes = ['32b_5_3', '32b_10_7']
    go_schemes = ['32b_5_3', '32b_10_7']
    
    rust_total_schemes = [results['rust'][scenario]['split_us'] + results['rust'][scenario]['combine_us'] for scenario in rust_schemes]
    go_total_schemes = [results['go'][scenario]['split_us'] + results['go'][scenario]['combine_us'] for scenario in go_schemes]
    
    x_schemes = np.arange(len(scheme_scenarios))
    
    bars1 = axes[1, 2].bar(x_schemes - width/2, rust_total_schemes, width, label='Rust', color='#FF6B6B', alpha=0.8)
    bars2 = axes[1, 2].bar(x_schemes + width/2, go_total_schemes, width, label='Go', color='#4ECDC4', alpha=0.8)
    
    axes[1, 2].set_ylabel('Total Time (microseconds)', fontweight='bold')
    axes[1, 2].set_title('Scheme Complexity Impact\n(32B data)', fontweight='bold')
    axes[1, 2].set_xticks(x_schemes)
    axes[1, 2].set_xticklabels(scheme_scenarios)
    axes[1, 2].legend()
    axes[1, 2].grid(True, alpha=0.3)
    axes[1, 2].set_axisbelow(True)
    
    # Add value labels
    for bar in bars1:
        height = bar.get_height()
        axes[1, 2].text(bar.get_x() + bar.get_width()/2., height,
                       f'{height:.1f}µs', ha='center', va='bottom', fontsize=8)
    for bar in bars2:
        height = bar.get_height()
        axes[1, 2].text(bar.get_x() + bar.get_width()/2., height,
                       f'{height:.1f}µs', ha='center', va='bottom', fontsize=8)
    
    plt.tight_layout()
    plt.savefig('shamir_performance_comparison_detailed.png', dpi=300, bbox_inches='tight')
    plt.show()

def print_performance_summary():
    print("=" * 100)
    print("SHAMIR SECRET SHARING PERFORMANCE SUMMARY")
    print("=" * 100)
    
    print("\nKEY OBSERVATIONS:")
    print("-" * 50)
    
    # Calculate average improvements
    common_scenarios = ['32b_5_3', '256b_5_3', '1024b_5_3']
    
    split_improvements = []
    combine_improvements = [] 
    throughput_improvements = []
    ops_improvements = []
    
    for scenario in common_scenarios:
        split_improvements.append(results['go'][scenario]['split_us'] / results['rust'][scenario]['split_us'])
        combine_improvements.append(results['go'][scenario]['combine_us'] / results['rust'][scenario]['combine_us'])
        throughput_improvements.append(results['rust'][scenario]['throughput_mb'] / results['go'][scenario]['throughput_mb'])
        ops_improvements.append(results['rust'][scenario]['ops_sec'] / results['go'][scenario]['ops_sec'])
    
    avg_split_improvement = sum(split_improvements) / len(split_improvements)
    avg_combine_improvement = sum(combine_improvements) / len(combine_improvements)
    avg_throughput_improvement = sum(throughput_improvements) / len(throughput_improvements)
    avg_ops_improvement = sum(ops_improvements) / len(ops_improvements)
    
    print(f"1. Rust is {avg_split_improvement:.1f}x faster for split operations")
    print(f"2. Rust is {avg_combine_improvement:.1f}x faster for combine operations") 
    print(f"3. Rust achieves {avg_throughput_improvement:.1f}x higher throughput")
    print(f"4. Rust processes {avg_ops_improvement:.1f}x more operations per second")
    
    print("\nDETAILED COMPARISON TABLE:")
    print("-" * 90)
    print(f"{'Scenario':<12} {'Metric':<12} {'Rust':<12} {'Go':<12} {'Improvement':<12} {'Notes':<20}")
    print("-" * 90)
    
    scenarios = [
        ('32B 5/3', '32b_5_3'),
        ('256B 5/3', '256b_5_3'), 
        ('1KB 5/3', '1024b_5_3'),
        ('32B 10/7', '32b_10_7')
    ]
    
    for name, key in scenarios:
        if key in results['rust'] and key in results['go']:
            rust_split = results['rust'][key]['split_us']
            go_split = results['go'][key]['split_us']
            improvement = go_split / rust_split
            print(f"{name:<12} {'Split':<12} {rust_split:<8.1f}µs {go_split:<8.1f}µs {improvement:<8.1f}x {'Rust faster' if improvement > 1 else 'Go faster':<20}")
            
            rust_combine = results['rust'][key]['combine_us']
            go_combine = results['go'][key]['combine_us']
            improvement = go_combine / rust_combine
            print(f"{name:<12} {'Combine':<12} {rust_combine:<8.1f}µs {go_combine:<8.1f}µs {improvement:<8.1f}x {'Rust faster' if improvement > 1 else 'Go faster':<20}")
            
            rust_ops = results['rust'][key]['ops_sec']
            go_ops = results['go'][key]['ops_sec']
            improvement = rust_ops / go_ops
            print(f"{name:<12} {'Ops/sec':<12} {rust_ops:<8.0f} {go_ops:<8.0f} {improvement:<8.1f}x {'Rust faster' if improvement > 1 else 'Go faster':<20}")
            print("-" * 90)

if __name__ == "__main__":
    plot_comprehensive_comparison()
    print_performance_summary()

"""
====================================================================================================
SHAMIR SECRET SHARING PERFORMANCE SUMMARY
====================================================================================================

KEY OBSERVATIONS:
--------------------------------------------------
1. Rust is 1.6x faster for split operations
2. Rust is 24.0x faster for combine operations
3. Rust achieves 7.3x higher throughput
4. Rust processes 7.2x more operations per second

DETAILED COMPARISON TABLE:
------------------------------------------------------------------------------------------
Scenario     Metric       Rust         Go           Improvement  Notes               
------------------------------------------------------------------------------------------
32B 5/3      Split        8.2     µs 13.8    µs 1.7     x Rust faster         
32B 5/3      Combine      3.0     µs 56.7    µs 18.9    x Rust faster         
32B 5/3      Ops/sec      89159    14197    6.3     x Rust faster         
------------------------------------------------------------------------------------------
256B 5/3     Split        42.2    µs 54.4    µs 1.3     x Rust faster         
256B 5/3     Combine      14.1    µs 291.4   µs 20.7    x Rust faster         
256B 5/3     Ops/sec      17774    2892     6.1     x Rust faster         
------------------------------------------------------------------------------------------
1KB 5/3      Split        113.3   µs 214.0   µs 1.9     x Rust faster         
1KB 5/3      Combine      36.0    µs 1164.1  µs 32.4    x Rust faster         
1KB 5/3      Ops/sec      6700     726      9.2     x Rust faster         
------------------------------------------------------------------------------------------
32B 10/7     Split        9.4     µs 29.3    µs 3.1     x Rust faster         
32B 10/7     Combine      4.4     µs 246.4   µs 56.0    x Rust faster         
32B 10/7     Ops/sec      72198    3626     19.9    x Rust faster         
------------------------------------------------------------------------------------------
"""    