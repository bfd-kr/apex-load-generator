# Apex Load Generator

A Go-based HTTP service for creating computational load through various operations including CPU-intensive calculations, memory allocation, and hex string generation. Built with the Gin web framework, this tool is designed for load testing, performance benchmarking, and system stress testing.

## Features

- **CPU Load Testing**: Prime number generation and Fibonacci calculations
- **Memory Allocation Testing**: Configurable memory allocation with automatic cleanup
- **Network/Bandwidth Testing**: Hex string generation for data transfer testing
- **Request-Level Metrics**: Comprehensive performance tracking for each request
- **Input Validation**: Built-in bounds checking to prevent resource exhaustion
- **Docker Ready**: Containerized deployment with Alpine Linux

## Quick Start

### Local Development

1. **Run the service**:
   ```bash
   go run main.go
   ```

2. **Test an endpoint**:
   ```bash
   curl http://localhost:8080/primes/100
   ```

### Docker Deployment

1. **Build and deploy**:
   ```bash
   ./local-build-docker-push
   ```

## API Endpoints

All endpoints return JSON responses with both `data` (the operation results) and `request_metrics` (performance data).

### Single Operations

#### Prime Number Generation
```bash
GET /primes/{p}
```
Generate the first `p` prime numbers (recommended for CPU load testing).

**Example**:
```bash
curl http://localhost:8080/primes/1000
```

**Response**:
```json
{
  "data": {
    "count": 1000,
    "last_prime": 7919,
    "duration_us": 1234,
    "duration_ms": 1.234
  },
  "request_metrics": {
    "duration_us": 1456,
    "duration_ms": 1.456,
    "cpu_usage_percent": 145.6,
    "memory_used_bytes": 8192,
    "goroutines_before": 8,
    "goroutines_after": 8
  }
}
```

#### Memory Allocation
```bash
GET /memory/{m}
```
Allocate `m` kilobytes of memory temporarily.

**Example**:
```bash
curl http://localhost:8080/memory/1024
```

#### Hex String Generation
```bash
GET /hex/{h}
```
Generate a hex string of `h` kilobytes or a random size within a range (returns full data for bandwidth testing).

**Examples**:
```bash
# Fixed size
curl http://localhost:8080/hex/10

# Random size within range
curl http://localhost:8080/hex/100..500
```

#### Fibonacci Calculation (Deprecated)
```bash
GET /fibonacci/{f}
```
Calculate the `f`th Fibonacci number. **Note**: Deprecated due to unpredictable exponential scaling. Use `/primes/{p}` instead.

### Combined Operations

#### Prime + Hex Generation
```bash
GET /primes/hex/{p}/{h}
```
Generate `p` prime numbers and `h` kilobytes of hex data.

**Example**:
```bash
curl http://localhost:8080/primes/hex/500/50
```

#### Prime + Hex + Memory (Full Load Test)
```bash
GET /primes/hex/memory/{p}/{h}/{m}
```
Comprehensive load test: generate `p` primes, `h` KB hex data, and allocate `m` KB memory.

**Example**:
```bash
curl http://localhost:8080/primes/hex/memory/1000/100/2048
```

## Input Limits

To prevent resource exhaustion, all endpoints enforce the following limits:

| Parameter | Endpoint | Range | Description |
|-----------|----------|-------|-------------|
| `p` | Primes | 0-10,000 | Number of prime numbers |
| `f` | Fibonacci | 0-45 | Fibonacci sequence position |
| `h` | Hex | 0-10,000 KB or range | Hex string size or range (e.g., 100..500) |
| `m` | Memory | 0-1,000,000 KB | Memory allocation size |

## Request Metrics

Every response includes detailed performance metrics:

**Request-Level Metrics:**
- **`duration_us`**: Total request duration in microseconds
- **`duration_ms`**: Total request duration in milliseconds
- **`cpu_usage_percent`**: Estimated CPU usage during request
- **`memory_used_bytes`**: Memory delta (allocated - freed)
- **`goroutines_before/after`**: Goroutine count tracking

**Operation-Level Metrics (in data field):**
- **`duration_us`**: Operation-specific timing in microseconds
- **`duration_ms`**: Operation-specific timing in milliseconds

## Load Testing Examples

### Light CPU Load
```bash
curl http://localhost:8080/primes/100
```

### Heavy CPU Load
```bash
curl http://localhost:8080/primes/5000
```

### Memory Pressure Test
```bash
curl http://localhost:8080/memory/100000
```

### Combined Load Test
```bash
curl http://localhost:8080/primes/hex/memory/2000/500/50000
```

### Bandwidth Test
```bash
curl http://localhost:8080/hex/1000
```

### Variable Size Test
```bash
curl http://localhost:8080/hex/500..2000
```

## Error Handling

The service returns appropriate HTTP status codes:

- **400 Bad Request**: Invalid parameters or out-of-range values
- **500 Internal Server Error**: Memory allocation failures or processing errors

**Example Error**:
```json
{
  "message": "p: number out of range (0-10000)"
}
```

## Performance Notes

- **Prime generation**: Linear complexity, predictable scaling
- **Fibonacci**: Exponential complexity, deprecated for unpredictable performance
- **Hex generation**: Optimized for low CPU usage, good for bandwidth testing
- **Memory allocation**: Uses page-boundary touching for realistic allocation patterns
- **Timing precision**: Microsecond and millisecond accuracy for performance analysis

## Development

### Build
```bash
go build -o apex-load-generator .
```

### Cross-compile for Linux
```bash
GOOS=linux GOARCH=amd64 go build -o out/apex-load-generator .
```

### Dependencies
- Go 1.23.3+
- Gin web framework (`github.com/gin-gonic/gin`)

## Testing

Use the included Postman collection (`apex-load-generator.postman_collection.json`) for comprehensive API testing, or use curl commands as shown in the examples above.
