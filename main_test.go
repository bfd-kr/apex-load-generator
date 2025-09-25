package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

// TestParseIntOrRange tests the abstracted range parsing function
func TestParseIntOrRange(t *testing.T) {
	tests := []struct {
		name        string
		param       string
		maxValue    int
		paramName   string
		expectError bool
		minExpected int
		maxExpected int
		expectRange bool
	}{
		{
			name:        "Valid single value",
			param:       "100",
			maxValue:    1000,
			paramName:   "test",
			expectError: false,
			minExpected: 100,
			maxExpected: 100,
			expectRange: false,
		},
		{
			name:        "Valid range",
			param:       "50..150",
			maxValue:    1000,
			paramName:   "test",
			expectError: false,
			minExpected: 50,
			maxExpected: 150,
			expectRange: true,
		},
		{
			name:        "Invalid number",
			param:       "invalid",
			maxValue:    1000,
			paramName:   "test",
			expectError: true,
		},
		{
			name:        "Negative value",
			param:       "-10",
			maxValue:    1000,
			paramName:   "test",
			expectError: true,
		},
		{
			name:        "Value exceeds max",
			param:       "2000",
			maxValue:    1000,
			paramName:   "test",
			expectError: true,
		},
		{
			name:        "Invalid range format",
			param:       "50-150",
			maxValue:    1000,
			paramName:   "test",
			expectError: true,
		},
		{
			name:        "Range with too many parts",
			param:       "50..100..150",
			maxValue:    1000,
			paramName:   "test",
			expectError: true,
		},
		{
			name:        "Range with invalid min",
			param:       "invalid..150",
			maxValue:    1000,
			paramName:   "test",
			expectError: true,
		},
		{
			name:        "Range with invalid max",
			param:       "50..invalid",
			maxValue:    1000,
			paramName:   "test",
			expectError: true,
		},
		{
			name:        "Range with negative values",
			param:       "-10..50",
			maxValue:    1000,
			paramName:   "test",
			expectError: true,
		},
		{
			name:        "Range with min > max",
			param:       "150..50",
			maxValue:    1000,
			paramName:   "test",
			expectError: true,
		},
		{
			name:        "Range exceeds max",
			param:       "500..2000",
			maxValue:    1000,
			paramName:   "test",
			expectError: true,
		},
		{
			name:        "Edge case - zero value",
			param:       "0",
			maxValue:    1000,
			paramName:   "test",
			expectError: false,
			minExpected: 0,
			maxExpected: 0,
			expectRange: false,
		},
		{
			name:        "Edge case - max value",
			param:       "1000",
			maxValue:    1000,
			paramName:   "test",
			expectError: false,
			minExpected: 1000,
			maxExpected: 1000,
			expectRange: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, isRange, err := parseIntOrRange(tt.param, tt.maxValue, tt.paramName)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if isRange != tt.expectRange {
				t.Errorf("Expected isRange=%v, got %v", tt.expectRange, isRange)
			}

			if tt.expectRange {
				if val < tt.minExpected || val > tt.maxExpected {
					t.Errorf("Expected value between %d-%d, got %d", tt.minExpected, tt.maxExpected, val)
				}
			} else {
				if val != tt.minExpected {
					t.Errorf("Expected value %d, got %d", tt.minExpected, val)
				}
			}
		})
	}
}

// TestAllocateMemory tests memory allocation function
func TestAllocateMemory(t *testing.T) {
	tests := []struct {
		name        string
		param       string
		expectError bool
		expectSize  int
	}{
		{
			name:        "Valid single value",
			param:       "10",
			expectError: false,
			expectSize:  10,
		},
		{
			name:        "Zero allocation",
			param:       "0",
			expectError: false,
			expectSize:  0,
		},
		{
			name:        "Invalid parameter",
			param:       "invalid",
			expectError: true,
		},
		{
			name:        "Exceeds max memory",
			param:       "2000000",
			expectError: true,
		},
		{
			name:        "Valid range",
			param:       "5..15",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := allocateMemory(tt.param)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if !strings.Contains(tt.param, "..") {
				// Single value test
				if result.SizeKB != tt.expectSize {
					t.Errorf("Expected size %d, got %d", tt.expectSize, result.SizeKB)
				}
			} else {
				// Range test - just verify it's within bounds
				if result.SizeKB < 5 || result.SizeKB > 15 {
					t.Errorf("Expected size between 5-15, got %d", result.SizeKB)
				}
				if result.RequestedRange != tt.param {
					t.Errorf("Expected RequestedRange=%s, got %s", tt.param, result.RequestedRange)
				}
			}

			// Verify timing fields are set (may be 0 for very fast operations)
			if result.DurationUs < 0 {
				t.Errorf("Expected non-negative DurationUs, got %d", result.DurationUs)
			}
			if result.DurationMs < 0 {
				t.Errorf("Expected non-negative DurationMs, got %f", result.DurationMs)
			}
		})
	}
}

// TestFibonacci tests Fibonacci calculation function
func TestFibonacci(t *testing.T) {
	tests := []struct {
		name           string
		param          string
		expectError    bool
		expectedResult int
	}{
		{
			name:           "Fibonacci 0",
			param:          "0",
			expectError:    false,
			expectedResult: 0,
		},
		{
			name:           "Fibonacci 1",
			param:          "1",
			expectError:    false,
			expectedResult: 1,
		},
		{
			name:           "Fibonacci 5",
			param:          "5",
			expectError:    false,
			expectedResult: 5,
		},
		{
			name:           "Fibonacci 10",
			param:          "10",
			expectError:    false,
			expectedResult: 55,
		},
		{
			name:        "Invalid parameter",
			param:       "invalid",
			expectError: true,
		},
		{
			name:        "Exceeds max fibonacci",
			param:       "50",
			expectError: true,
		},
		{
			name:        "Valid range",
			param:       "5..10",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fibonacci(tt.param)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if !strings.Contains(tt.param, "..") {
				// Single value test - verify the position (N) and result
				expectedN, _ := strconv.Atoi(tt.param)
				if result.N != expectedN {
					t.Errorf("Expected N=%d, got %d", expectedN, result.N)
				}
				if result.Result != tt.expectedResult {
					t.Errorf("Expected Result=%d, got %d", tt.expectedResult, result.Result)
				}
			} else {
				// Range test - just verify it's within bounds
				if result.N < 5 || result.N > 10 {
					t.Errorf("Expected N between 5-10, got %d", result.N)
				}
				if result.RequestedRange != tt.param {
					t.Errorf("Expected RequestedRange=%s, got %s", tt.param, result.RequestedRange)
				}
			}

			// Verify timing fields are set (may be 0 for very fast operations)
			if result.DurationUs < 0 {
				t.Errorf("Expected non-negative DurationUs, got %d", result.DurationUs)
			}
			if result.DurationMs < 0 {
				t.Errorf("Expected non-negative DurationMs, got %f", result.DurationMs)
			}
		})
	}
}

// TestGeneratePrimes tests prime number generation function
func TestGeneratePrimes(t *testing.T) {
	tests := []struct {
		name          string
		param         string
		expectError   bool
		expectedCount int
		expectedPrime int
	}{
		{
			name:          "Generate 0 primes",
			param:         "0",
			expectError:   false,
			expectedCount: 0,
			expectedPrime: 0,
		},
		{
			name:          "Generate 1 prime",
			param:         "1",
			expectError:   false,
			expectedCount: 1,
			expectedPrime: 2,
		},
		{
			name:          "Generate 5 primes",
			param:         "5",
			expectError:   false,
			expectedCount: 5,
			expectedPrime: 11,
		},
		{
			name:        "Invalid parameter",
			param:       "invalid",
			expectError: true,
		},
		{
			name:        "Exceeds max primes",
			param:       "20000",
			expectError: true,
		},
		{
			name:        "Valid range",
			param:       "5..10",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := generatePrimes(tt.param)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if !strings.Contains(tt.param, "..") {
				// Single value test
				if result.Count != tt.expectedCount {
					t.Errorf("Expected Count=%d, got %d", tt.expectedCount, result.Count)
				}
				if tt.expectedCount > 0 && result.LastPrime != tt.expectedPrime {
					t.Errorf("Expected LastPrime=%d, got %d", tt.expectedPrime, result.LastPrime)
				}
			} else {
				// Range test - just verify it's within bounds
				if result.Count < 5 || result.Count > 10 {
					t.Errorf("Expected Count between 5-10, got %d", result.Count)
				}
				if result.RequestedRange != tt.param {
					t.Errorf("Expected RequestedRange=%s, got %s", tt.param, result.RequestedRange)
				}
			}

			// Verify timing fields are set (may be 0 for very fast operations)
			if result.DurationUs < 0 {
				t.Errorf("Expected non-negative DurationUs, got %d", result.DurationUs)
			}
			if result.DurationMs < 0 {
				t.Errorf("Expected non-negative DurationMs, got %f", result.DurationMs)
			}
		})
	}
}

// TestCreateHexString tests hex string generation function
func TestCreateHexString(t *testing.T) {
	tests := []struct {
		name         string
		param        string
		expectError  bool
		expectedSize int
	}{
		{
			name:         "Generate 1KB hex",
			param:        "1",
			expectError:  false,
			expectedSize: 1,
		},
		{
			name:         "Generate 0KB hex",
			param:        "0",
			expectError:  false,
			expectedSize: 0,
		},
		{
			name:        "Invalid parameter",
			param:       "invalid",
			expectError: true,
		},
		{
			name:        "Exceeds max hex",
			param:       "20000",
			expectError: true,
		},
		{
			name:        "Valid range",
			param:       "1..5",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := createHexString(tt.param)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if !strings.Contains(tt.param, "..") {
				// Single value test
				if result.SizeKB != tt.expectedSize {
					t.Errorf("Expected SizeKB=%d, got %d", tt.expectedSize, result.SizeKB)
				}
				expectedLength := tt.expectedSize * 1024
				if result.Length != expectedLength {
					t.Errorf("Expected Length=%d, got %d", expectedLength, result.Length)
				}
				if len(result.HexString) != expectedLength {
					t.Errorf("Expected HexString length=%d, got %d", expectedLength, len(result.HexString))
				}
			} else {
				// Range test - just verify it's within bounds
				if result.SizeKB < 1 || result.SizeKB > 5 {
					t.Errorf("Expected SizeKB between 1-5, got %d", result.SizeKB)
				}
				if result.RequestedRange != tt.param {
					t.Errorf("Expected RequestedRange=%s, got %s", tt.param, result.RequestedRange)
				}
			}

			// Verify hex string contains only valid hex characters
			if result.SizeKB > 0 {
				for _, char := range result.HexString {
					if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f')) {
						t.Errorf("Invalid hex character found: %c", char)
						break
					}
				}
			}

			// Verify timing fields are set (may be 0 for very fast operations)
			if result.DurationUs < 0 {
				t.Errorf("Expected non-negative DurationUs, got %d", result.DurationUs)
			}
			if result.DurationMs < 0 {
				t.Errorf("Expected non-negative DurationMs, got %f", result.DurationMs)
			}
		})
	}
}

// TestFibonacciRecursive tests the recursive Fibonacci implementation
func TestFibonacciRecursive(t *testing.T) {
	tests := []struct {
		input    int
		expected int
	}{
		{0, 0},
		{1, 1},
		{2, 1},
		{3, 2},
		{4, 3},
		{5, 5},
		{6, 8},
		{7, 13},
		{8, 21},
		{9, 34},
		{10, 55},
	}

	for _, tt := range tests {
		t.Run(string(rune(tt.input)), func(t *testing.T) {
			result := fibonacciRecursive(tt.input)
			if result != tt.expected {
				t.Errorf("fibonacciRecursive(%d) = %d, expected %d", tt.input, result, tt.expected)
			}
		})
	}
}

// TestRequestMetrics tests request metrics functionality
func TestStartRequestMetrics(t *testing.T) {
	metrics := startRequestMetrics()

	if metrics == nil {
		t.Fatal("Expected non-nil metrics")
	}

	if metrics.StartTime.IsZero() {
		t.Error("Expected StartTime to be set")
	}

	if metrics.GoroutinesBefore < 0 {
		t.Error("Expected non-negative GoroutinesBefore")
	}

	if metrics.MemoryUsedBytes < 0 {
		t.Error("Expected non-negative MemoryUsedBytes")
	}

	// Test finish method
	time.Sleep(1 * time.Millisecond) // Ensure some time passes
	metrics.finish()

	if metrics.EndTime.IsZero() {
		t.Error("Expected EndTime to be set after finish()")
	}

	if metrics.DurationUs <= 0 {
		t.Error("Expected positive DurationUs after finish()")
	}

	if metrics.DurationMs <= 0 {
		t.Error("Expected positive DurationMs after finish()")
	}

	if metrics.GoroutinesAfter < 0 {
		t.Error("Expected non-negative GoroutinesAfter")
	}
}

// BenchmarkParseIntOrRange benchmarks the abstracted parsing function
func BenchmarkParseIntOrRange(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parseIntOrRange("100..500", 1000, "test")
	}
}

// BenchmarkAllocateMemory benchmarks memory allocation
func BenchmarkAllocateMemory(b *testing.B) {
	for i := 0; i < b.N; i++ {
		allocateMemory("1")
	}
}

// BenchmarkFibonacci benchmarks Fibonacci calculation
func BenchmarkFibonacci(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fibonacci("10")
	}
}

// BenchmarkGeneratePrimes benchmarks prime generation
func BenchmarkGeneratePrimes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		generatePrimes("10")
	}
}

// BenchmarkCreateHexString benchmarks hex string generation
func BenchmarkCreateHexString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		createHexString("1")
	}
}

// setupRouter creates a test router with all routes
func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/", getIndex)
	router.GET("/fibonacci/:f", getFibonacci)
	router.GET("/primes/:p", getPrimes)
	router.GET("/hex/:h", getHexString)
	router.GET("/memory/:m", getMemory)
	router.GET("/fibonacci/hex/:f/:h", getFibonacciHex)
	router.GET("/primes/hex/:p/:h", getPrimesHex)
	router.GET("/fibonacci/hex/memory/:f/:h/:m", fibonacciHexMemory)
	router.GET("/primes/hex/memory/:p/:h/:m", primesHexMemory)
	return router
}

// TestGetIndex tests the HTML homepage endpoint
func TestGetIndex(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if !strings.Contains(w.Body.String(), "Apex Load Generator API") {
		t.Error("Expected homepage to contain 'Apex Load Generator API'")
	}

	if !strings.Contains(w.Header().Get("Content-Type"), "text/html") {
		t.Error("Expected Content-Type to be text/html")
	}
}

// TestGetMemory tests the memory allocation endpoint
func TestGetMemory(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		name           string
		param          string
		expectedStatus int
		expectError    bool
	}{
		{
			name:           "Valid memory allocation",
			param:          "10",
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "Valid memory range",
			param:          "5..15",
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "Invalid parameter",
			param:          "invalid",
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:           "Exceeds maximum",
			param:          "2000000",
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/memory/"+tt.param, nil)
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if !tt.expectError {
				// Parse successful response
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("Failed to parse JSON response: %v", err)
				}

				// Verify response structure
				if _, ok := response["data"]; !ok {
					t.Error("Expected 'data' field in response")
				}
				if _, ok := response["request_metrics"]; !ok {
					t.Error("Expected 'request_metrics' field in response")
				}
			}
		})
	}
}

// TestGetFibonacci tests the Fibonacci calculation endpoint
func TestGetFibonacci(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		name           string
		param          string
		expectedStatus int
		expectError    bool
	}{
		{
			name:           "Valid Fibonacci calculation",
			param:          "5",
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "Valid Fibonacci range",
			param:          "3..7",
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "Invalid parameter",
			param:          "invalid",
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:           "Exceeds maximum",
			param:          "50",
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/fibonacci/"+tt.param, nil)
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if !tt.expectError {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("Failed to parse JSON response: %v", err)
				}

				if _, ok := response["data"]; !ok {
					t.Error("Expected 'data' field in response")
				}
				if _, ok := response["request_metrics"]; !ok {
					t.Error("Expected 'request_metrics' field in response")
				}
			}
		})
	}
}

// TestGetPrimes tests the prime generation endpoint
func TestGetPrimes(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		name           string
		param          string
		expectedStatus int
		expectError    bool
	}{
		{
			name:           "Valid prime generation",
			param:          "5",
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "Valid prime range",
			param:          "3..7",
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "Invalid parameter",
			param:          "invalid",
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:           "Exceeds maximum",
			param:          "50000",
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/primes/"+tt.param, nil)
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if !tt.expectError {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("Failed to parse JSON response: %v", err)
				}

				if _, ok := response["data"]; !ok {
					t.Error("Expected 'data' field in response")
				}
				if _, ok := response["request_metrics"]; !ok {
					t.Error("Expected 'request_metrics' field in response")
				}
			}
		})
	}
}

// TestGetHexString tests the hex string generation endpoint
func TestGetHexString(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		name           string
		param          string
		expectedStatus int
		expectError    bool
	}{
		{
			name:           "Valid hex generation",
			param:          "1",
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "Valid hex range",
			param:          "1..3",
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "Invalid parameter",
			param:          "invalid",
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:           "Exceeds maximum",
			param:          "50000",
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/hex/"+tt.param, nil)
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if !tt.expectError {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("Failed to parse JSON response: %v", err)
				}

				if _, ok := response["data"]; !ok {
					t.Error("Expected 'data' field in response")
				}
				if _, ok := response["request_metrics"]; !ok {
					t.Error("Expected 'request_metrics' field in response")
				}
			}
		})
	}
}

// TestGetFibonacciHex tests the combined Fibonacci + hex endpoint
func TestGetFibonacciHex(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		name           string
		fibParam       string
		hexParam       string
		expectedStatus int
		expectError    bool
	}{
		{
			name:           "Valid combined operation",
			fibParam:       "5",
			hexParam:       "1",
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "Valid with ranges",
			fibParam:       "3..7",
			hexParam:       "1..3",
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "Invalid fibonacci parameter",
			fibParam:       "invalid",
			hexParam:       "1",
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:           "Invalid hex parameter",
			fibParam:       "5",
			hexParam:       "invalid",
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/fibonacci/hex/"+tt.fibParam+"/"+tt.hexParam, nil)
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if !tt.expectError {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("Failed to parse JSON response: %v", err)
				}

				data, ok := response["data"].(map[string]interface{})
				if !ok {
					t.Error("Expected 'data' field to be an object")
				}

				if _, ok := data["fibonacci_result"]; !ok {
					t.Error("Expected 'fibonacci_result' in data")
				}
				if _, ok := data["hex_result"]; !ok {
					t.Error("Expected 'hex_result' in data")
				}
			}
		})
	}
}

// TestGetPrimesHex tests the combined primes + hex endpoint
func TestGetPrimesHex(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		name           string
		primeParam     string
		hexParam       string
		expectedStatus int
		expectError    bool
	}{
		{
			name:           "Valid combined operation",
			primeParam:     "5",
			hexParam:       "1",
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "Valid with ranges",
			primeParam:     "3..7",
			hexParam:       "1..3",
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "Invalid prime parameter",
			primeParam:     "invalid",
			hexParam:       "1",
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:           "Invalid hex parameter",
			primeParam:     "5",
			hexParam:       "invalid",
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/primes/hex/"+tt.primeParam+"/"+tt.hexParam, nil)
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if !tt.expectError {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("Failed to parse JSON response: %v", err)
				}

				data, ok := response["data"].(map[string]interface{})
				if !ok {
					t.Error("Expected 'data' field to be an object")
				}

				if _, ok := data["prime_result"]; !ok {
					t.Error("Expected 'prime_result' in data")
				}
				if _, ok := data["hex_result"]; !ok {
					t.Error("Expected 'hex_result' in data")
				}
			}
		})
	}
}

// TestFibonacciHexMemory tests the triple combined endpoint
func TestFibonacciHexMemory(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		name           string
		fibParam       string
		hexParam       string
		memParam       string
		expectedStatus int
		expectError    bool
	}{
		{
			name:           "Valid triple operation",
			fibParam:       "5",
			hexParam:       "1",
			memParam:       "10",
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "Valid with ranges",
			fibParam:       "3..7",
			hexParam:       "1..3",
			memParam:       "5..15",
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "Invalid fibonacci parameter",
			fibParam:       "invalid",
			hexParam:       "1",
			memParam:       "10",
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:           "Invalid memory parameter",
			fibParam:       "5",
			hexParam:       "1",
			memParam:       "invalid",
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/fibonacci/hex/memory/"+tt.fibParam+"/"+tt.hexParam+"/"+tt.memParam, nil)
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if !tt.expectError {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("Failed to parse JSON response: %v", err)
				}

				data, ok := response["data"].(map[string]interface{})
				if !ok {
					t.Error("Expected 'data' field to be an object")
				}

				if _, ok := data["fibonacci_result"]; !ok {
					t.Error("Expected 'fibonacci_result' in data")
				}
				if _, ok := data["hex_result"]; !ok {
					t.Error("Expected 'hex_result' in data")
				}
				if _, ok := data["memory_result"]; !ok {
					t.Error("Expected 'memory_result' in data")
				}
			}
		})
	}
}

// TestPrimesHexMemory tests the triple combined endpoint with primes
func TestPrimesHexMemory(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		name           string
		primeParam     string
		hexParam       string
		memParam       string
		expectedStatus int
		expectError    bool
	}{
		{
			name:           "Valid triple operation",
			primeParam:     "5",
			hexParam:       "1",
			memParam:       "10",
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "Valid with ranges",
			primeParam:     "3..7",
			hexParam:       "1..3",
			memParam:       "5..15",
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "Invalid prime parameter",
			primeParam:     "invalid",
			hexParam:       "1",
			memParam:       "10",
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:           "Invalid memory parameter",
			primeParam:     "5",
			hexParam:       "1",
			memParam:       "invalid",
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/primes/hex/memory/"+tt.primeParam+"/"+tt.hexParam+"/"+tt.memParam, nil)
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if !tt.expectError {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("Failed to parse JSON response: %v", err)
				}

				data, ok := response["data"].(map[string]interface{})
				if !ok {
					t.Error("Expected 'data' field to be an object")
				}

				if _, ok := data["prime_result"]; !ok {
					t.Error("Expected 'prime_result' in data")
				}
				if _, ok := data["hex_result"]; !ok {
					t.Error("Expected 'hex_result' in data")
				}
				if _, ok := data["memory_result"]; !ok {
					t.Error("Expected 'memory_result' in data")
				}
			}
		})
	}
}

// TestMainFunction tests that main function can be called without panicking
func TestMainFunction(t *testing.T) {
	// We can't easily test the main function directly since it starts a server
	// But we can test that our router setup doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Router setup panicked: %v", r)
		}
	}()

	// Test router creation (similar to what main does)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/", getIndex)
	router.GET("/fibonacci/:f", getFibonacci)
	router.GET("/primes/:p", getPrimes)
	router.GET("/hex/:h", getHexString)
	router.GET("/memory/:m", getMemory)
	router.GET("/fibonacci/hex/:f/:h", getFibonacciHex)
	router.GET("/primes/hex/:p/:h", getPrimesHex)
	router.GET("/fibonacci/hex/memory/:f/:h/:m", fibonacciHexMemory)
	router.GET("/primes/hex/memory/:p/:h/:m", primesHexMemory)

	// Verify router was created successfully
	if router == nil {
		t.Error("Router creation failed")
	}
}