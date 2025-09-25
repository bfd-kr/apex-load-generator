package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	// MaxMemoryKB is the maximum memory allocation limit in kilobytes
	MaxMemoryKB = 1000000
	// MaxFibonacci is the maximum Fibonacci position limit
	MaxFibonacci = 45
	// MaxPrimes is the maximum prime count limit
	MaxPrimes = 10000
	// MaxHexKB is the maximum hex string size limit in kilobytes
	MaxHexKB = 10000
)

// RequestMetrics holds request-level performance metrics
type RequestMetrics struct {
	StartTime        time.Time `json:"-"`
	EndTime          time.Time `json:"-"`
	DurationUs       int64     `json:"duration_us"`
	DurationMs       float64   `json:"duration_ms"`
	CPUUsagePercent  float64   `json:"cpu_usage_percent"`
	MemoryUsedBytes  int64     `json:"memory_used_bytes"`
	GoroutinesBefore int       `json:"goroutines_before"`
	GoroutinesAfter  int       `json:"goroutines_after"`
}

// parseIntOrRange parses a parameter that can be either a single integer or a range.
// Returns the parsed value and whether it was a range.
func parseIntOrRange(param string, maxValue int, paramName string) (int, bool, error) {
	// Parse the parameter (single value or range)
	if strings.Contains(param, "..") {
		parts := strings.Split(param, "..")
		if len(parts) != 2 {
			return 0, false, fmt.Errorf("invalid range format, use min..max")
		}

		min, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil {
			return 0, false, fmt.Errorf("invalid minimum value: %v", err)
		}

		max, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil {
			return 0, false, fmt.Errorf("invalid maximum value: %v", err)
		}

		if min < 0 || max < 0 {
			return 0, false, fmt.Errorf("values must be non-negative")
		}

		if min > max {
			return 0, false, fmt.Errorf("minimum value cannot be greater than maximum")
		}

		if min > maxValue || max > maxValue {
			return 0, false, fmt.Errorf("values must be within range (0-%d)", maxValue)
		}

		actualValue := min + rand.Intn(max-min+1)
		return actualValue, true, nil
	} else {
		// Single value
		value, err := strconv.Atoi(param)
		if err != nil {
			return 0, false, fmt.Errorf("invalid number: %v", err)
		}

		if value < 0 || value > maxValue {
			return 0, false, fmt.Errorf("number out of range (0-%d)", maxValue)
		}

		return value, false, nil
	}
}

// startRequestMetrics initializes request metrics collection
func startRequestMetrics() *RequestMetrics {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	return &RequestMetrics{
		StartTime:        time.Now(),
		GoroutinesBefore: runtime.NumGoroutine(),
		MemoryUsedBytes:  int64(memStats.Alloc),
	}
}

// finishRequestMetrics completes request metrics collection
func (rm *RequestMetrics) finish() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	rm.EndTime = time.Now()
	duration := rm.EndTime.Sub(rm.StartTime)
	rm.DurationUs = duration.Nanoseconds() / 1000
	rm.DurationMs = float64(duration.Nanoseconds()) / 1000000.0
	rm.GoroutinesAfter = runtime.NumGoroutine()

	// Calculate memory used during request
	memoryAfter := int64(memStats.Alloc)
	rm.MemoryUsedBytes = memoryAfter - rm.MemoryUsedBytes

	// Simple CPU usage approximation (not perfect but indicative)
	rm.CPUUsagePercent = float64(rm.DurationUs) / 1000.0 // rough approximation
}

// MemoryResult holds the result of memory allocation including timing
type MemoryResult struct {
	SizeKB         int     `json:"size_kb"`
	RequestedRange string  `json:"requested_range,omitempty"`
	DurationUs     int64   `json:"duration_us"`
	DurationMs     float64 `json:"duration_ms"`
}

// allocateMemory creates a byte slice of size mb and ensures allocation.
// Accepts either a single value (e.g., "1024") or a range (e.g., "500..2000")
func allocateMemory(param string) (MemoryResult, error) {
	start := time.Now()
	var err error

	k, wasRange, err := parseIntOrRange(param, MaxMemoryKB, "memory")
	if err != nil {
		return MemoryResult{}, err
	}

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("memory allocation failed: %v", r)
		}
	}()

	bytes := make([]byte, k*1024)
	// Touch memory to ensure allocation
	for i := 0; i < len(bytes); i += 4096 {
		bytes[i] = 1
	}
	// Memory will be freed naturally by GC

	duration := time.Since(start)

	memoryResult := MemoryResult{
		SizeKB:     k,
		DurationUs: duration.Nanoseconds() / 1000,
		DurationMs: float64(duration.Nanoseconds()) / 1000000.0,
	}

	// Only include requested_range if it was a range
	if wasRange {
		memoryResult.RequestedRange = param
	}

	return memoryResult, err
}

// getMemory handles GET requests to allocate memory of m kilobytes or a random size within a range.
func getMemory(c *gin.Context) {
	metrics := startRequestMetrics()

	m := c.Param("m")
	result, err := allocateMemory(m)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("m: %v", err)})
		return
	}
	metrics.finish()
	c.IndentedJSON(http.StatusOK, gin.H{
		"data":            result,
		"request_metrics": metrics,
	})
}

// FibonacciResult holds the result of Fibonacci calculation including timing
type FibonacciResult struct {
	N              int     `json:"n"`
	RequestedRange string  `json:"requested_range,omitempty"`
	Result         int     `json:"result"`
	DurationUs     int64   `json:"duration_us"`
	DurationMs     float64 `json:"duration_ms"`
}

// fibonacci calculates the nth Fibonacci number.
// Accepts either a single value (e.g., "30") or a range (e.g., "25..35")
// Deprecated: Use generatePrimes() for more predictable CPU load testing.
func fibonacci(param string) (FibonacciResult, error) {
	start := time.Now()

	n, wasRange, err := parseIntOrRange(param, MaxFibonacci, "fibonacci")
	if err != nil {
		return FibonacciResult{}, err
	}

	var result int
	if n <= 1 {
		result = n
	} else {
		result = fibonacciRecursive(n)
	}

	duration := time.Since(start)

	fibResult := FibonacciResult{
		N:          n,
		Result:     result,
		DurationUs: duration.Nanoseconds() / 1000,
		DurationMs: float64(duration.Nanoseconds()) / 1000000.0,
	}

	// Only include requested_range if it was a range
	if wasRange {
		fibResult.RequestedRange = param
	}

	return fibResult, nil
}

// fibonacciRecursive is the actual recursive implementation
func fibonacciRecursive(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacciRecursive(n-1) + fibonacciRecursive(n-2)
}

// PrimeResult holds the result of prime generation including timing
type PrimeResult struct {
	Count          int     `json:"count"`
	RequestedRange string  `json:"requested_range,omitempty"`
	LastPrime      int     `json:"last_prime"`
	DurationUs     int64   `json:"duration_us"`
	DurationMs     float64 `json:"duration_ms"`
}

// generatePrimes generates the first n prime numbers and returns timing information.
// Accepts either a single value (e.g., "100") or a range (e.g., "100..1000")
func generatePrimes(param string) (PrimeResult, error) {
	start := time.Now()

	n, wasRange, err := parseIntOrRange(param, MaxPrimes, "primes")
	if err != nil {
		return PrimeResult{}, err
	}

	if n <= 0 {
		duration := time.Since(start)
		result := PrimeResult{
			Count:      0,
			LastPrime:  0,
			DurationUs: duration.Nanoseconds() / 1000,
			DurationMs: float64(duration.Nanoseconds()) / 1000000.0,
		}
		if wasRange {
			result.RequestedRange = param
		}
		return result, nil
	}

	if n == 1 {
		duration := time.Since(start)
		result := PrimeResult{
			Count:      1,
			LastPrime:  2,
			DurationUs: duration.Nanoseconds() / 1000,
			DurationMs: float64(duration.Nanoseconds()) / 1000000.0,
		}
		if wasRange {
			result.RequestedRange = param
		}
		return result, nil
	}

	// Keep track of primes found so far for trial division, but only store what we need
	primes := []int{2}
	lastPrime := 2
	count := 1

	for candidate := 3; count < n; candidate += 2 {
		isPrime := true
		for _, prime := range primes {
			if prime*prime > candidate {
				break
			}
			if candidate%prime == 0 {
				isPrime = false
				break
			}
		}
		if isPrime {
			primes = append(primes, candidate)
			lastPrime = candidate
			count++
		}
	}

	duration := time.Since(start)
	result := PrimeResult{
		Count:      count,
		LastPrime:  lastPrime,
		DurationUs: duration.Nanoseconds() / 1000,
		DurationMs: float64(duration.Nanoseconds()) / 1000000.0,
	}
	if wasRange {
		result.RequestedRange = param
	}
	return result, nil
}

// getFibonacci handles GET requests to calculate the nth Fibonacci number or a random position within a range.
// Deprecated: Use getPrimes() for more predictable CPU load testing.
func getFibonacci(c *gin.Context) {
	metrics := startRequestMetrics()

	f := c.Param("f")
	result, err := fibonacci(f)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("f: %v", err)})
		return
	}
	metrics.finish()
	c.IndentedJSON(http.StatusOK, gin.H{
		"data":            result,
		"request_metrics": metrics,
	})
}

// getPrimes handles GET requests to generate the first n prime numbers or a random count within a range.
func getPrimes(c *gin.Context) {
	metrics := startRequestMetrics()

	p := c.Param("p")
	result, err := generatePrimes(p)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("p: %v", err)})
		return
	}
	metrics.finish()
	c.IndentedJSON(http.StatusOK, gin.H{
		"data":            result,
		"request_metrics": metrics,
	})
}

// HexResult holds the result of hex string generation including timing
type HexResult struct {
	SizeKB         int     `json:"size_kb"`
	RequestedRange string  `json:"requested_range,omitempty"`
	Length         int     `json:"length"`
	HexString      string  `json:"hex_string"`
	DurationUs     int64   `json:"duration_us"`
	DurationMs     float64 `json:"duration_ms"`
}

// createHexString generates a hex string of specified size in kilobytes.
// Accepts either a single value (e.g., "100") or a range (e.g., "100..500")
func createHexString(param string) (HexResult, error) {
	start := time.Now()

	n, wasRange, err := parseIntOrRange(param, MaxHexKB, "hex")
	if err != nil {
		return HexResult{}, err
	}

	hexChars := "0123456789abcdef"
	result := make([]byte, n*1024)
	for i := range result {
		result[i] = hexChars[rand.Intn(16)]
	}

	hexString := string(result)
	duration := time.Since(start)

	hexResult := HexResult{
		SizeKB:     n,
		Length:     len(hexString),
		HexString:  hexString,
		DurationUs: duration.Nanoseconds() / 1000,
		DurationMs: float64(duration.Nanoseconds()) / 1000000.0,
	}

	// Only include requested_range if it was a range
	if wasRange {
		hexResult.RequestedRange = param
	}

	return hexResult, nil
}


// getHexString handles GET requests to generate a hex string of n kilobytes or a random size within a range.
func getHexString(c *gin.Context) {
	metrics := startRequestMetrics()

	h := c.Param("h")
	result, err := createHexString(h)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("h: %v", err)})
		return
	}
	metrics.finish()
	c.IndentedJSON(http.StatusOK, gin.H{
		"data":            result,
		"request_metrics": metrics,
	})
}

func getFibonacciHex(c *gin.Context) {
	metrics := startRequestMetrics()

	f := c.Param("f")
	h := c.Param("h")

	fResult, err := fibonacci(f)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("f: %v", err)})
		return
	}

	hResult, err := createHexString(h)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("h: %v", err)})
		return
	}

	metrics.finish()
	c.IndentedJSON(http.StatusOK, gin.H{
		"data":            map[string]interface{}{"fibonacci_result": fResult, "hex_result": hResult},
		"request_metrics": metrics,
	})
}

// getPrimesHex handles GET requests to generate primes and hex string.
func getPrimesHex(c *gin.Context) {
	metrics := startRequestMetrics()

	p := c.Param("p")
	h := c.Param("h")

	pResult, err := generatePrimes(p)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("p: %v", err)})
		return
	}

	hResult, err := createHexString(h)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("h: %v", err)})
		return
	}

	metrics.finish()
	c.IndentedJSON(http.StatusOK, gin.H{
		"data":            map[string]interface{}{"prime_result": pResult, "hex_result": hResult},
		"request_metrics": metrics,
	})
}

// create function fibonacci, hex, memory
func fibonacciHexMemory(c *gin.Context) {
	metrics := startRequestMetrics()

	f := c.Param("f")
	h := c.Param("h")
	m := c.Param("m")

	fResult, err := fibonacci(f)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("f: %v", err)})
		return
	}

	hResult, err := createHexString(h)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("h: %v", err)})
		return
	}

	mResult, err := allocateMemory(m)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("m: %v", err)})
		return
	}

	metrics.finish()
	c.IndentedJSON(http.StatusOK, gin.H{
		"data":            map[string]interface{}{"fibonacci_result": fResult, "hex_result": hResult, "memory_result": mResult},
		"request_metrics": metrics,
	})
}

// primesHexMemory handles GET requests to generate primes, hex string, and allocate memory.
func primesHexMemory(c *gin.Context) {
	metrics := startRequestMetrics()

	p := c.Param("p")
	h := c.Param("h")
	m := c.Param("m")

	pResult, err := generatePrimes(p)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("p: %v", err)})
		return
	}

	hResult, err := createHexString(h)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("h: %v", err)})
		return
	}

	mResult, err := allocateMemory(m)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("m: %v", err)})
		return
	}

	metrics.finish()
	c.IndentedJSON(http.StatusOK, gin.H{
		"data":            map[string]interface{}{"prime_result": pResult, "hex_result": hResult, "memory_result": mResult},
		"request_metrics": metrics,
	})
}

// getIndex serves the API documentation homepage
func getIndex(c *gin.Context) {
	html := `<!DOCTYPE html>
<html>
<head>
    <title>Apex Load Generator API</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; background-color: #f5f5f5; }
        .container { background: white; padding: 30px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        h1 { color: #333; border-bottom: 3px solid #4CAF50; padding-bottom: 10px; }
        h2 { color: #555; margin-top: 30px; }
        .endpoint { background: #f9f9f9; padding: 15px; margin: 10px 0; border-radius: 5px; border-left: 4px solid #4CAF50; }
        .method { background: #4CAF50; color: white; padding: 3px 8px; border-radius: 3px; font-weight: bold; }
        .deprecated { border-left-color: #ff9800; }
        .deprecated .method { background: #ff9800; }
        .example { background: #e8f5e8; padding: 10px; margin: 8px 0; border-radius: 3px; font-family: monospace; }
        .limits { background: #fff3cd; padding: 10px; margin: 8px 0; border-radius: 3px; border: 1px solid #ffeaa7; }
        a { color: #4CAF50; text-decoration: none; font-weight: bold; }
        a:hover { text-decoration: underline; }
        .note { background: #d1ecf1; padding: 10px; margin: 8px 0; border-radius: 3px; border: 1px solid #bee5eb; }
    </style>
</head>
<body>
    <div class="container">
        <h1>🚀 Apex Load Generator API</h1>
        <p>A Go-based HTTP service for creating computational load through CPU calculations, memory allocation, and hex string generation. All endpoints return JSON with both operation results and request-level performance metrics.</p>

        <h2>🔢 Single Operations</h2>

        <div class="endpoint">
            <span class="method">GET</span> <strong>/primes/{p}</strong> - Generate Prime Numbers (Recommended)
            <div class="example">
                Example: <a href="/primes/100">/primes/100</a> - Generate first 100 prime numbers<br>
                Range: <a href="/primes/50..200">/primes/50..200</a> - Generate random count between 50-200 primes
            </div>
            <div class="limits">Limits: p = 0-10,000 or range (e.g., 50..200) | Best for predictable CPU load testing</div>
        </div>

        <div class="endpoint">
            <span class="method">GET</span> <strong>/memory/{m}</strong> - Allocate Memory
            <div class="example">
                Example: <a href="/memory/1024">/memory/1024</a> - Allocate 1MB (1024 KB) of memory<br>
                Range: <a href="/memory/500..2000">/memory/500..2000</a> - Allocate random size between 500KB-2MB
            </div>
            <div class="limits">Limits: m = 0-1,000,000 KB or range (e.g., 500..2000) | Memory is allocated then freed naturally</div>
        </div>

        <div class="endpoint">
            <span class="method">GET</span> <strong>/hex/{h}</strong> - Generate Hex String
            <div class="example">
                Example: <a href="/hex/10">/hex/10</a> - Generate 10KB of hex data<br>
                Range: <a href="/hex/100..500">/hex/100..500</a> - Generate random size between 100-500KB
            </div>
            <div class="limits">Limits: h = 0-10,000 KB or range (e.g., 100..500) | Returns full hex data for bandwidth testing</div>
        </div>

        <div class="endpoint deprecated">
            <span class="method">GET</span> <strong>/fibonacci/{f}</strong> - Fibonacci Calculation (Deprecated)
            <div class="example">
                Example: <a href="/fibonacci/30">/fibonacci/30</a> - Calculate 30th Fibonacci number<br>
                Range: <a href="/fibonacci/25..35">/fibonacci/25..35</a> - Calculate random position between 25-35
            </div>
            <div class="limits">Limits: f = 0-45 or range (e.g., 25..35) | ⚠️ Deprecated: Use /primes for predictable CPU testing</div>
        </div>

        <h2>🔄 Combined Operations</h2>

        <div class="endpoint">
            <span class="method">GET</span> <strong>/primes/hex/{p}/{h}</strong> - Prime + Hex Generation
            <div class="example">
                Example: <a href="/primes/hex/500/50">/primes/hex/500/50</a> - 500 primes + 50KB hex
            </div>
            <div class="limits">Limits: p = 0-10,000, h = 0-10,000 KB</div>
        </div>

        <div class="endpoint">
            <span class="method">GET</span> <strong>/primes/hex/memory/{p}/{h}/{m}</strong> - Full Load Test
            <div class="example">
                Example: <a href="/primes/hex/memory/1000/100/2048">/primes/hex/memory/1000/100/2048</a> - Complete load test
            </div>
            <div class="limits">Limits: p = 0-10,000, h = 0-10,000 KB, m = 0-1,000,000 KB</div>
        </div>

        <div class="endpoint deprecated">
            <span class="method">GET</span> <strong>/fibonacci/hex/{f}/{h}</strong> - Fibonacci + Hex (Deprecated)
            <div class="example">
                Example: <a href="/fibonacci/hex/25/50">/fibonacci/hex/25/50</a> - Fibonacci + 50KB hex
            </div>
            <div class="limits">Limits: f = 0-45, h = 0-10,000 KB | ⚠️ Use /primes/hex instead</div>
        </div>

        <div class="endpoint deprecated">
            <span class="method">GET</span> <strong>/fibonacci/hex/memory/{f}/{h}/{m}</strong> - Fibonacci + Hex + Memory (Deprecated)
            <div class="example">
                Example: <a href="/fibonacci/hex/memory/20/50/1024">/fibonacci/hex/memory/20/50/1024</a> - All operations
            </div>
            <div class="limits">Limits: f = 0-45, h = 0-10,000 KB, m = 0-1,000,000 KB | ⚠️ Use /primes/hex/memory instead</div>
        </div>

        <h2>📊 Response Format</h2>
        <div class="note">
            All endpoints return JSON with:
            <ul>
                <li><strong>data</strong>: Operation results (timing in both microseconds and milliseconds, counts, generated content)</li>
                <li><strong>request_metrics</strong>: Performance data (duration_us, duration_ms, cpu_usage_percent, memory_used_bytes, goroutine counts)</li>
            </ul>
        </div>

        <h2>🎯 Load Testing Examples</h2>
        <div class="example">
            Light CPU: <a href="/primes/100">/primes/100</a><br>
            Heavy CPU: <a href="/primes/5000">/primes/5000</a><br>
            Variable CPU: <a href="/primes/100..1000">/primes/100..1000</a><br>
            Memory Test: <a href="/memory/100000">/memory/100000</a><br>
            Variable Memory: <a href="/memory/50000..200000">/memory/50000..200000</a><br>
            Bandwidth Test: <a href="/hex/1000">/hex/1000</a><br>
            Variable Size Test: <a href="/hex/500..2000">/hex/500..2000</a><br>
            Combined Load: <a href="/primes/hex/memory/2000/500/50000">/primes/hex/memory/2000/500/50000</a>
        </div>

        <div class="note">
            <strong>💡 Performance Tips:</strong>
            <ul>
                <li>Prime generation: Linear complexity, predictable scaling</li>
                <li>Fibonacci: Exponential complexity, deprecated for unpredictable performance</li>
                <li>Hex generation: Optimized for low CPU usage, good for bandwidth testing</li>
                <li>Memory allocation: Uses page-boundary touching for realistic patterns</li>
                <li>All timing: Microsecond and millisecond precision for accurate performance analysis</li>
            </ul>
        </div>
    </div>
</body>
</html>`

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(200, html)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	router := gin.Default()
	router.GET("/", getIndex)
	router.GET("/fibonacci/:f", getFibonacci)
	router.GET("/primes/:p", getPrimes)
	router.GET("/hex/:h", getHexString)
	router.GET("/memory/:m", getMemory)
	router.GET("/fibonacci/hex/:f/:h", getFibonacciHex)
	router.GET("/primes/hex/:p/:h", getPrimesHex)
	router.GET("/fibonacci/hex/memory/:f/:h/:m", fibonacciHexMemory)
	router.GET("/primes/hex/memory/:p/:h/:m", primesHexMemory)

	router.Run(":8080")
}
