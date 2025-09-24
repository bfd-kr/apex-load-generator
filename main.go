package main

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// creates a byte slice of size mb and ensures allocation.
func allocateMemory(k int) (err error) {
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
	return nil
}

// getMemory handles GET requests to allocate memory.
func getMemory(c *gin.Context) {
	m := c.Param("m")
	num, err := strconv.Atoi(m)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid number"})
		return
	}
	if err := allocateMemory(num); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "memory allocation failed"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"data": map[string]interface{}{"m": m, "memory": "done"}})
}

// fibonacci calculates the nth Fibonacci number.
// Deprecated: Use generatePrimes() for more predictable CPU load testing.
func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

// generatePrimes generates the first n prime numbers.
func generatePrimes(n int) []int {
	if n <= 0 {
		return []int{}
	}
	primes := []int{2}
	if n == 1 {
		return primes
	}

	for candidate := 3; len(primes) < n; candidate += 2 {
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
		}
	}
	return primes
}

// getFibonacci handles GET requests to calculate the nth Fibonacci number.
// Deprecated: Use getPrimes() for more predictable CPU load testing.
func getFibonacci(c *gin.Context) {
	f := c.Param("f")
	num, err := strconv.Atoi(f)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid number"})
		return
	}
	result := fibonacci(num)
	c.IndentedJSON(http.StatusOK, gin.H{"data": map[string]interface{}{"f": f, "fibonacci": result}})
}

// getPrimes handles GET requests to generate the first n prime numbers.
func getPrimes(c *gin.Context) {
	p := c.Param("p")
	num, err := strconv.Atoi(p)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid number"})
		return
	}
	result := generatePrimes(num)
	c.IndentedJSON(http.StatusOK, gin.H{"data": map[string]interface{}{"p": p, "primes": result}})
}

// createHexString generates a hex string of n kilobytes.
func createHexString(n int) (string, error) {
	hexChars := "0123456789abcdef"
	result := make([]byte, n*1024)
	for i := range result {
		result[i] = hexChars[rand.Intn(16)]
	}
	return string(result), nil
}

// getHexString handles GET requests to generate a hex string of n kilobytes.
func getHexString(c *gin.Context) {
	h := c.Param("h")
	num, err := strconv.Atoi(h)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid number"})
		return
	}
	result, err := createHexString(num)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "could not generate hex string"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"data": map[string]interface{}{"h": h, "hexString": result}})
}

func getFibonacciHex(c *gin.Context) {
	f := c.Param("f")
	h := c.Param("h")
	fNum, err := strconv.Atoi(f)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid number"})
		return
	}
	hNum, err := strconv.Atoi(h)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid number"})
		return
	}
	fResult := fibonacci(fNum)
	hResult, err := createHexString(hNum)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "could not generate hex string"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"data": map[string]interface{}{"f": f, "fibonacci": fResult, "h": h, "hexString": hResult}})
}

// getPrimesHex handles GET requests to generate primes and hex string.
func getPrimesHex(c *gin.Context) {
	p := c.Param("p")
	h := c.Param("h")
	pNum, err := strconv.Atoi(p)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid number"})
		return
	}
	hNum, err := strconv.Atoi(h)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid number"})
		return
	}
	pResult := generatePrimes(pNum)
	hResult, err := createHexString(hNum)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "could not generate hex string"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"data": map[string]interface{}{"p": p, "primes": pResult, "h": h, "hexString": hResult}})
}

// create function fibonacci, hex, memory
func fibonacciHexMemory(c *gin.Context) {
	f := c.Param("f")
	h := c.Param("h")
	m := c.Param("m")
	fNum, err := strconv.Atoi(f)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "f: invalid number"})
		return
	}

	hNum, err := strconv.Atoi(h)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "h: invalid number"})
		return
	}

	mNum, err := strconv.Atoi(m)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "m: invalid number"})
		return
	}

	fResult := fibonacci(fNum)
	hResult, err := createHexString(hNum)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "could not generate hex string"})
		return
	}

	if err := allocateMemory(mNum); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "memory allocation failed"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"data": map[string]interface{}{"f": f, "fibonacci": fResult, "h": h, "hexString": hResult, "m": "done"}})
}

// primesHexMemory handles GET requests to generate primes, hex string, and allocate memory.
func primesHexMemory(c *gin.Context) {
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

	if err := allocateMemory(mNum); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "memory allocation failed"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"data": map[string]interface{}{"p": p, "primes": pResult, "h": h, "hexString": hResult, "m": "done"}})
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
