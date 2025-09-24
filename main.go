package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestMetrics holds request-level performance metrics
type RequestMetrics struct {
	StartTime        time.Time `json:"-"`
	EndTime          time.Time `json:"-"`
	DurationUs       int64     `json:"duration_us"`
	CPUUsagePercent  float64   `json:"cpu_usage_percent"`
	MemoryUsedBytes  int64     `json:"memory_used_bytes"`
	GoroutinesBefore int       `json:"goroutines_before"`
	GoroutinesAfter  int       `json:"goroutines_after"`
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
	rm.DurationUs = rm.EndTime.Sub(rm.StartTime).Nanoseconds() / 1000
	rm.GoroutinesAfter = runtime.NumGoroutine()

	// Calculate memory used during request
	memoryAfter := int64(memStats.Alloc)
	rm.MemoryUsedBytes = memoryAfter - rm.MemoryUsedBytes

	// Simple CPU usage approximation (not perfect but indicative)
	rm.CPUUsagePercent = float64(rm.DurationUs) / 1000.0 // rough approximation
}

// MemoryResult holds the result of memory allocation including timing
type MemoryResult struct {
	SizeKB     int   `json:"size_kb"`
	DurationUs int64 `json:"duration_us"`
}

// allocateMemory creates a byte slice of size mb and ensures allocation.
func allocateMemory(k int) (MemoryResult, error) {
	start := time.Now()
	var err error

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

	return MemoryResult{
		SizeKB:     k,
		DurationUs: time.Since(start).Nanoseconds() / 1000,
	}, err
}

// getMemory handles GET requests to allocate memory.
func getMemory(c *gin.Context) {
	metrics := startRequestMetrics()

	m := c.Param("m")
	num, err := strconv.Atoi(m)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid number"})
		return
	}
	if num < 0 || num > 1000000 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "m: number out of range (0-1000000)"})
		return
	}
	result, err := allocateMemory(num)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "memory allocation failed"})
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
	N          int   `json:"n"`
	Result     int   `json:"result"`
	DurationUs int64 `json:"duration_us"`
}

// fibonacci calculates the nth Fibonacci number.
// Deprecated: Use generatePrimes() for more predictable CPU load testing.
func fibonacci(n int) FibonacciResult {
	start := time.Now()

	var result int
	if n <= 1 {
		result = n
	} else {
		result = fibonacciRecursive(n)
	}

	return FibonacciResult{
		N:          n,
		Result:     result,
		DurationUs: time.Since(start).Nanoseconds() / 1000,
	}
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
	Count      int   `json:"count"`
	LastPrime  int   `json:"last_prime"`
	DurationUs int64 `json:"duration_us"`
}

// generatePrimes generates the first n prime numbers and returns timing information.
func generatePrimes(n int) PrimeResult {
	start := time.Now()

	if n <= 0 {
		return PrimeResult{
			Count:      0,
			LastPrime:  0,
			DurationUs: time.Since(start).Nanoseconds() / 1000,
		}
	}

	if n == 1 {
		return PrimeResult{
			Count:      1,
			LastPrime:  2,
			DurationUs: time.Since(start).Nanoseconds() / 1000,
		}
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

	return PrimeResult{
		Count:      count,
		LastPrime:  lastPrime,
		DurationUs: time.Since(start).Nanoseconds() / 1000,
	}
}

// getFibonacci handles GET requests to calculate the nth Fibonacci number.
// Deprecated: Use getPrimes() for more predictable CPU load testing.
func getFibonacci(c *gin.Context) {
	metrics := startRequestMetrics()

	f := c.Param("f")
	num, err := strconv.Atoi(f)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid number"})
		return
	}
	if num < 0 || num > 45 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "f: number out of range (0-45)"})
		return
	}
	result := fibonacci(num)
	metrics.finish()
	c.IndentedJSON(http.StatusOK, gin.H{
		"data":            result,
		"request_metrics": metrics,
	})
}

// getPrimes handles GET requests to generate the first n prime numbers.
func getPrimes(c *gin.Context) {
	metrics := startRequestMetrics()

	p := c.Param("p")
	num, err := strconv.Atoi(p)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid number"})
		return
	}
	if num < 0 || num > 10000 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "p: number out of range (0-10000)"})
		return
	}
	result := generatePrimes(num)
	metrics.finish()
	c.IndentedJSON(http.StatusOK, gin.H{
		"data":            result,
		"request_metrics": metrics,
	})
}

// HexResult holds the result of hex string generation including timing
type HexResult struct {
	SizeKB     int    `json:"size_kb"`
	Length     int    `json:"length"`
	HexString  string `json:"hex_string"`
	DurationUs int64  `json:"duration_us"`
}

// createHexString generates a hex string of n kilobytes.
func createHexString(n int) (HexResult, error) {
	start := time.Now()

	hexChars := "0123456789abcdef"
	result := make([]byte, n*1024)
	for i := range result {
		result[i] = hexChars[rand.Intn(16)]
	}

	hexString := string(result)
	return HexResult{
		SizeKB:     n,
		Length:     len(hexString),
		HexString:  hexString,
		DurationUs: time.Since(start).Nanoseconds() / 1000,
	}, nil
}

// getHexString handles GET requests to generate a hex string of n kilobytes.
func getHexString(c *gin.Context) {
	metrics := startRequestMetrics()

	h := c.Param("h")
	num, err := strconv.Atoi(h)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid number"})
		return
	}
	if num < 0 || num > 1000 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "h: number out of range (0-1000)"})
		return
	}
	result, err := createHexString(num)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "could not generate hex string"})
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
	fNum, err := strconv.Atoi(f)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid number"})
		return
	}
	if fNum < 0 || fNum > 45 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "f: number out of range (0-45)"})
		return
	}
	hNum, err := strconv.Atoi(h)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid number"})
		return
	}
	if hNum < 0 || hNum > 1000 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "h: number out of range (0-1000)"})
		return
	}
	fResult := fibonacci(fNum)
	hResult, err := createHexString(hNum)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "could not generate hex string"})
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
	pNum, err := strconv.Atoi(p)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid number"})
		return
	}
	if pNum < 0 || pNum > 10000 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "p: number out of range (0-10000)"})
		return
	}
	hNum, err := strconv.Atoi(h)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid number"})
		return
	}
	if hNum < 0 || hNum > 1000 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "h: number out of range (0-1000)"})
		return
	}
	pResult := generatePrimes(pNum)
	hResult, err := createHexString(hNum)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "could not generate hex string"})
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
	fNum, err := strconv.Atoi(f)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "f: invalid number"})
		return
	}

	if fNum < 0 || fNum > 45 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "f: number out of range (0-45)"})
		return
	}

	hNum, err := strconv.Atoi(h)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "h: invalid number"})
		return
	}

	if hNum < 0 || hNum > 1000 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "h: number out of range (0-1000)"})
		return
	}

	mNum, err := strconv.Atoi(m)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "m: invalid number"})
		return
	}

	if mNum < 0 || mNum > 1000000 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "m: number out of range (0-1000000)"})
		return
	}

	fResult := fibonacci(fNum)
	hResult, err := createHexString(hNum)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "could not generate hex string"})
		return
	}

	mResult, err := allocateMemory(mNum)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "memory allocation failed"})
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
	pNum, err := strconv.Atoi(p)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "p: invalid number"})
		return
	}

	if pNum < 0 || pNum > 10000 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "p: number out of range (0-10000)"})
		return
	}

	hNum, err := strconv.Atoi(h)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "h: invalid number"})
		return
	}

	if hNum < 0 || hNum > 1000 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "h: number out of range (0-1000)"})
		return
	}

	mNum, err := strconv.Atoi(m)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "m: invalid number"})
		return
	}

	if mNum < 0 || mNum > 1000000 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "m: number out of range (0-1000000)"})
		return
	}

	pResult := generatePrimes(pNum)
	hResult, err := createHexString(hNum)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "could not generate hex string"})
		return
	}

	mResult, err := allocateMemory(mNum)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "memory allocation failed"})
		return
	}

	metrics.finish()
	c.IndentedJSON(http.StatusOK, gin.H{
		"data":            map[string]interface{}{"prime_result": pResult, "hex_result": hResult, "memory_result": mResult},
		"request_metrics": metrics,
	})
}

func main() {
	rand.Seed(time.Now().UnixNano())
	router := gin.Default()
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
