package main

import (
	"strconv"
	"strings"
	"testing"
	"time"
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