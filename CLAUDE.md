# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go-based load generator service that provides HTTP endpoints for creating computational load through Fibonacci calculations, hex string generation, and memory allocation. It's built using the Gin web framework and runs on port 8080.

## Architecture

- **Single-file application**: All code is in `main.go`
- **Web framework**: Uses Gin for HTTP routing and JSON responses
- **Load generation functions**:
  - `fibonacci()`: Recursive Fibonacci calculation for CPU load
  - `createHexString()`: Random hex string generation for CPU/memory load
  - `allocateMemory()`: Memory allocation for memory pressure testing

## API Endpoints

- `GET /fibonacci/:f` - Calculate nth Fibonacci number
- `GET /hex/:h` - Generate hex string of h kilobytes
- `GET /memory/:m` - Allocate m kilobytes of memory
- `GET /fibonacci/hex/:f/:h` - Combined Fibonacci and hex generation
- `GET /fibonacci/hex/memory/:f/:h/:m` - Combined all three operations

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