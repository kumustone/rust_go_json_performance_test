# Rust vs Go JSON Serialization Performance Comparison

A comprehensive benchmark comparing JSON serialization/deserialization performance between Rust and Go implementations.

## 📊 Benchmark Scope

This benchmark compares:
- **Go easyjson**: Code-generated JSON marshaler/unmarshaler
- **Go stdlib**: Standard library `encoding/json` (v1)
- **Rust serde_json**: Popular serde-based JSON library (optimized with FxHashMap + LTO)

> **Note**: Go's experimental `encoding/json/v2` is not included in the default benchmark due to its unstable status. See [JSONV2_NOTES.md](JSONV2_NOTES.md) for details and manual testing instructions.

## 🏗️ Test Structure

The benchmark uses an HTTP request-like structure with varying payload sizes:

```go
type HttpRequest struct {
    Mark          string
    Method        string
    Scheme        string
    Url           string
    Proto         string
    Host          string
    RemoteAddr    string
    ContentLength uint64
    Header        map[string][]string
    Body          []byte
}
```

### Test Cases

- **S (Small)**: 512B body, 6 headers
- **M (Medium)**: 8KB body, 12 headers  
- **L (Large)**: 64KB body, 20 headers

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- Rust 1.70+
- Python 3.8+ with matplotlib (for visualization)

### Setup

```bash
make setup
```

### Run Benchmarks

```bash
# Run all benchmarks and generate plots
make all

# Or run individually
make bench-go      # Go benchmarks only
make bench-rust    # Rust benchmarks only
make plot          # Generate comparison charts
```

## 📈 Results

After running benchmarks, you'll get:
- `benchmark_comparison.png`: Absolute performance comparison
- `benchmark_ratio.png`: Performance ratio vs Go stdlib

### Sample Results (Apple M1 Pro)

| Case | Operation | Go easyjson | Go stdlib (v1) | Go json/v2 | Rust serde_json |
|------|-----------|-------------|----------------|------------|-----------------|
| S | encode | 1.20 µs | 1.49 µs | 2.64 µs | **1.02 µs** |
| S | decode | **1.57 µs** | 7.27 µs | 2.78 µs | 1.62 µs |
| M | encode | 13.89 µs | **7.97 µs** | 23.29 µs | 7.90 µs |
| M | decode | **7.39 µs** | 58.17 µs | 15.84 µs | 7.01 µs |
| L | encode | 97.56 µs | **47.23 µs** | 159.05 µs | 59.06 µs |
| L | decode | 43.14 µs | 412.09 µs | 99.30 µs | **41.50 µs** |

**Key Findings:**
- **Rust serde_json**: Fastest for small encode and large decode
- **Go easyjson**: Excels at small/medium decode (code generation wins)
- **Go stdlib (v1)**: Fast for large encode but very slow for decode (4-10x slower)
- **Go json/v2**: Decode is 2.6-4.1x faster than v1, but still slower than easyjson; encode needs optimization

## 🔧 Optimization Details

### Rust Optimizations Applied

1. **FxHashMap**: Faster non-cryptographic hash for HashMap
2. **LTO (Link Time Optimization)**: `lto = "thin"`
3. **Single codegen unit**: Better optimization opportunities
4. **Release profile**: `opt-level = 3`

### Go Optimizations

- easyjson uses code generation for zero-reflection overhead
- Standard library uses reflection (slower but more flexible)

## 📁 Project Structure

```
.
├── go-bench/                          # Go benchmarks
│   ├── http_request.go
│   ├── http_request_benchmark_test.go
│   ├── http_request_jsonv2_test.go   # Optional JSON v2 benchmark
│   └── go.mod
├── rust-bench/                        # Rust benchmarks
│   ├── src/
│   │   └── lib.rs
│   ├── benches/
│   │   └── waf_http_request_json.rs
│   └── Cargo.toml
├── plot_benchmark.py                  # Visualization script
├── PERFORMANCE_ANALYSIS.md            # Detailed performance analysis
├── JSONV2_NOTES.md                    # Go JSON v2 information
├── Makefile
└── README.md
```

## 🧪 Running Custom Tests

### Go

```bash
cd go-bench
go test -bench . -benchmem -benchtime=1s
```

### Rust

```bash
cd rust-bench
cargo bench --bench waf_http_request_json
```

## 📝 License

MIT

## 📚 Additional Documentation

- **[PERFORMANCE_ANALYSIS.md](PERFORMANCE_ANALYSIS.md)**: Deep dive into performance differences, reflection overhead, compiler optimizations, and memory allocation patterns
- **[JSONV2_NOTES.md](JSONV2_NOTES.md)**: Information about Go's experimental JSON v2 library and how to test it

## 🤝 Contributing

Contributions welcome! Feel free to:
- Add more serialization libraries (e.g., simd-json, sonic, JSON v2)
- Test different data structures
- Improve benchmark methodology
- Submit JSON v2 benchmark results
