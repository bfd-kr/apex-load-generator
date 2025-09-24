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
    - **Returns**: FibonacciResult struct with input n, result value, and timing information
    - **Important**: Provides unpredictable scaling - small input changes cause massive CPU usage differences
    - **Migration**: Replace with `generatePrimes()` for predictable CPU testing
  - `generatePrimes()`: Prime number generation for CPU load (linear complexity - predictable scaling)
    - **Purpose**: CPU load testing with predictable, configurable intensity
    - **Behavior**: Generates first n prime numbers using trial division with optimizations (only test odd candidates, early termination)
    - **Complexity**: O(n^1.5) - much more predictable than exponential Fibonacci
    - **Returns**: PrimeResult struct with timing information, count, and last prime found (no full list for memory efficiency)
    - **Timing**: Uses high-resolution timer (time.Now()) not subject to process suspension
    - **Important**: Preferred over Fibonacci for consistent CPU load testing
  - `createHexString()`: Random hex string generation for CPU/memory load (optimized for low CPU usage)
    - **Purpose**: Generate hex strings of specified size for load testing with minimal CPU overhead
    - **Behavior**: Directly generates hex characters (0-9, a-f) using `math/rand` instead of byte-to-hex conversion
    - **Returns**: HexResult struct with size, length, hex string, and timing information
    - **Optimization**: Avoids expensive `hex.EncodeToString()` and crypto-grade random generation for better performance
    - **Important**: Uses `math/rand.Intn(16)` for efficiency - do not revert to `crypto/rand` or `hex.EncodeToString()`
  - `allocateMemory()`: Memory allocation for memory pressure testing
    - **Purpose**: Temporarily allocate memory to create memory pressure, then allow natural garbage collection
    - **Behavior**: Allocates k kilobytes, touches memory at 4KB page boundaries to ensure real allocation, then lets Go's GC handle cleanup naturally
    - **Returns**: MemoryResult struct with size and timing information, plus error if allocation fails
    - **Error Handling**: Returns error if memory allocation fails (e.g., out of memory conditions)
    - **Important**: Do not force garbage collection with `runtime.GC()` - let it happen naturally for realistic load testing

## API Endpoints

- `GET /fibonacci/:f` - **DEPRECATED** - Calculate nth Fibonacci number (returns timing and performance data)
- `GET /primes/:p` - Generate first p prime numbers (returns timing and performance data)
- `GET /hex/:h` - Generate hex string of h kilobytes (returns timing and performance data)
- `GET /memory/:m` - Allocate m kilobytes of memory (returns timing and performance data)
- `GET /fibonacci/hex/:f/:h` - **DEPRECATED** - Combined Fibonacci and hex generation (use /primes/hex instead)
- `GET /primes/hex/:p/:h` - Combined prime generation and hex string creation (includes timing data)
- `GET /fibonacci/hex/memory/:f/:h/:m` - **DEPRECATED** - Combined all three operations with Fibonacci (use /primes/hex/memory instead)
- `GET /primes/hex/memory/:p/:h/:m` - Combined prime generation, hex string creation, and memory allocation (includes timing data)
  - **Input Limits**: p: 0-10,000, h: 0-1,000 KB, m: 0-1,000,000 KB (prevents resource exhaustion)

## Input Validation

All endpoints now have comprehensive bounds checking to prevent resource exhaustion:

- **`/fibonacci/:f`**: f: 0-45 (prevents exponential explosion)
- **`/primes/:p`**: p: 0-10,000 (prevents excessive CPU usage)
- **`/hex/:h`**: h: 0-1,000 KB (prevents large memory allocations)
- **`/memory/:m`**: m: 0-1,000,000 KB (prevents system memory exhaustion)
- **Combined endpoints**: Apply limits for all respective parameters
  - `/fibonacci/hex/:f/:h` - f: 0-45, h: 0-1,000 KB
  - `/primes/hex/:p/:h` - p: 0-10,000, h: 0-1,000 KB
  - `/fibonacci/hex/memory/:f/:h/:m` - f: 0-45, h: 0-1,000 KB, m: 0-1,000,000 KB
  - `/primes/hex/memory/:p/:h/:m` - p: 0-10,000, h: 0-1,000 KB, m: 0-1,000,000 KB

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

## Dependencies

Primary dependency is `github.com/gin-gonic/gin` for the web framework. Uses standard library packages for encoding, math, and HTTP.