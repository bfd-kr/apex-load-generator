# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go-based load generator service that provides HTTP endpoints for creating computational load through Fibonacci calculations, hex string generation, and memory allocation. It's built using the Gin web framework and runs on port 8080.

## Architecture

- **Single-file application**: All code is in `main.go`
- **Web framework**: Uses Gin for HTTP routing and JSON responses
- **Load generation functions**:
  - `fibonacci()`: **DEPRECATED** - Recursive Fibonacci calculation for CPU load (exponential complexity - unpredictable scaling)
    - **Status**: Deprecated in favor of `generatePrimes()`
    - **Purpose**: Legacy CPU load testing with exponential time complexity
    - **Behavior**: Classic recursive implementation with O(2^n) complexity
    - **Returns**: FibonacciResult struct with input n, result value, and timing information (both microseconds and milliseconds)
    - **Important**: Provides unpredictable scaling - small input changes cause massive CPU usage differences
    - **Migration**: Replace with `generatePrimes()` for predictable CPU testing
  - `generatePrimes()`: Prime number generation for CPU load (linear complexity - predictable scaling)
    - **Purpose**: CPU load testing with predictable, configurable intensity
    - **Behavior**: Generates first n prime numbers using trial division with optimizations (only test odd candidates, early termination)
    - **Complexity**: O(n^1.5) - much more predictable than exponential Fibonacci
    - **Returns**: PrimeResult struct with timing information (both microseconds and milliseconds), count, and last prime found (no full list for memory efficiency)
    - **Timing**: Uses high-resolution timer (time.Now()) with microsecond and millisecond precision, not subject to process suspension
    - **Important**: Preferred over Fibonacci for consistent CPU load testing
  - `createHexString()`: Random hex string generation for CPU/memory load (optimized for low CPU usage)
    - **Purpose**: Generate hex strings of specified size or random size within a range for load testing with minimal CPU overhead
    - **Behavior**: Directly generates hex characters (0-9, a-f) using `math/rand` instead of byte-to-hex conversion
    - **Input**: Accepts single values (e.g., "100") or ranges (e.g., "100..500") for variable size testing
    - **Returns**: HexResult struct with actual size, optional requested range, length, hex string, and timing information (both microseconds and milliseconds)
    - **Range Feature**: When range is provided (e.g., 100..500), randomly selects size within range (inclusive) for each request
    - **Data Transfer Testing**: Returns full hex string content for network/bandwidth testing (hex data compresses poorly)
    - **Optimization**: Avoids expensive `hex.EncodeToString()` and crypto-grade random generation for better performance
    - **Important**: Uses `math/rand.Intn(16)` for efficiency - do not revert to `crypto/rand` or `hex.EncodeToString()`
  - `allocateMemory()`: Memory allocation for memory pressure testing
    - **Purpose**: Temporarily allocate memory to create memory pressure, then allow natural garbage collection
    - **Behavior**: Allocates k kilobytes, touches memory at 4KB page boundaries to ensure real allocation, then lets Go's GC handle cleanup naturally
    - **Returns**: MemoryResult struct with size and timing information (both microseconds and milliseconds), plus error if allocation fails
    - **Error Handling**: Returns error if memory allocation fails (e.g., out of memory conditions)
    - **Important**: Do not force garbage collection with `runtime.GC()` - let it happen naturally for realistic load testing

## API Endpoints

- `GET /fibonacci/:f` - **DEPRECATED** - Calculate nth Fibonacci number (returns timing data in both microseconds and milliseconds)
- `GET /primes/:p` - Generate first p prime numbers (returns timing data in both microseconds and milliseconds)
- `GET /hex/:h` - Generate hex string of h kilobytes or random size within range (returns full hex data with timing in both microseconds and milliseconds)
- `GET /memory/:m` - Allocate m kilobytes of memory (returns timing data in both microseconds and milliseconds)
- `GET /fibonacci/hex/:f/:h` - **DEPRECATED** - Combined Fibonacci and hex generation (use /primes/hex instead)
- `GET /primes/hex/:p/:h` - Combined prime generation and hex string creation (includes full hex data with timing in both microseconds and milliseconds)
- `GET /fibonacci/hex/memory/:f/:h/:m` - **DEPRECATED** - Combined all three operations with Fibonacci (use /primes/hex/memory instead)
- `GET /primes/hex/memory/:p/:h/:m` - Combined prime generation, hex string creation, and memory allocation (includes full hex data with timing in both microseconds and milliseconds)
  - **Input Limits**: p: 0-10,000, h: 0-1,000 KB, m: 0-1,000,000 KB (prevents resource exhaustion)

## Input Validation

All endpoints now have comprehensive bounds checking to prevent resource exhaustion:

- **`/fibonacci/:f`**: f: 0-45 (prevents exponential explosion)
- **`/primes/:p`**: p: 0-10,000 (prevents excessive CPU usage)
- **`/hex/:h`**: h: 0-10,000 KB or range (e.g., 100..500) (prevents excessive memory allocations)
- **`/memory/:m`**: m: 0-1,000,000 KB (prevents system memory exhaustion)
- **Combined endpoints**: Apply limits for all respective parameters
  - `/fibonacci/hex/:f/:h` - f: 0-45, h: 0-10,000 KB
  - `/primes/hex/:p/:h` - p: 0-10,000, h: 0-10,000 KB
  - `/fibonacci/hex/memory/:f/:h/:m` - f: 0-45, h: 0-10,000 KB, m: 0-1,000,000 KB
  - `/primes/hex/memory/:p/:h/:m` - p: 0-10,000, h: 0-10,000 KB, m: 0-1,000,000 KB

## Error Handling

- **Memory allocation failures**: All endpoints that use `allocateMemory()` now handle allocation failures gracefully
  - Returns HTTP 500 with "memory allocation failed" message
  - Uses panic recovery to catch out-of-memory conditions
  - Affected endpoints: `/memory/:m`, `/fibonacci/hex/memory/:f/:h/:m`, `/primes/hex/memory/:p/:h/:m`

## Development Commands

### Build
```bash
go build -o apex-load-generator .
```

### Cross-compile for Linux (for Docker)
```bash
GOOS=linux GOARCH=amd64 go build -o out/apex-load-generator .
```

### Run locally
```bash
go run main.go
```

### Docker build and deployment
```bash
./local-build-docker-push
```

## Files Structure

- `main.go` - Main application with all handlers and logic
- `go.mod/go.sum` - Go module dependencies
- `Dockerfile` - Alpine-based container definition
- `local-build-docker-push` - Build script for Docker deployment
- `apex-load-generator.postman_collection.json` - Postman collection for API testing
- `out/` - Build output directory

## Request-Level Metrics

**New Feature**: All API endpoints now include request-level performance metrics in their responses. This enables comprehensive analysis of request performance including CPU usage, memory consumption, and timing data.

### RequestMetrics Structure

Each request now includes a `request_metrics` field in the JSON response containing:

- **`duration_us`**: Request duration in microseconds (from start to completion)
- **`duration_ms`**: Request duration in milliseconds (from start to completion)
- **`cpu_usage_percent`**: Approximated CPU usage percentage during request processing
- **`memory_used_bytes`**: Memory delta in bytes during request execution (memory consumed - memory freed)
- **`goroutines_before`**: Number of goroutines before request processing
- **`goroutines_after`**: Number of goroutines after request processing

### Response Format

All endpoints now return JSON with two top-level fields:
```json
{
  "data": {
    // Original response data (FibonacciResult, PrimeResult, HexResult, MemoryResult, or combinations)
  },
  "request_metrics": {
    "duration_us": 1234,
    "duration_ms": 1.234,
    "cpu_usage_percent": 25.5,
    "memory_used_bytes": 1048576,
    "goroutines_before": 8,
    "goroutines_after": 8
  }
}
```

### Implementation Details

- **Timing**: Uses high-resolution `time.Now()` with both microsecond and millisecond precision
- **Memory Tracking**: Captures `runtime.MemStats.Alloc` before and after request processing
- **CPU Estimation**: Simple approximation based on request duration (not perfect but indicative)
- **Goroutine Monitoring**: Tracks goroutine count changes during request processing
- **Zero Overhead**: Metrics collection adds minimal performance overhead

### Use Cases

- **Performance Analysis**: Track request performance across different load levels
- **Resource Monitoring**: Monitor memory usage patterns and goroutine behavior
- **Load Testing**: Compare performance metrics under various load conditions
- **Optimization**: Identify performance bottlenecks in load generation functions

## Dependencies

Primary dependency is `github.com/gin-gonic/gin` for the web framework. Uses standard library packages for encoding, math, and HTTP.