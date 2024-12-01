package main

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// fibonacci calculates the nth Fibonacci number.
func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

// getFibonacci handles GET requests to calculate the nth Fibonacci number.
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

// createHexString generates a hex string of n kilobytes.
func createHexString(n int) (string, error) {
	bytes := make([]byte, n*1024)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	hexString := hex.EncodeToString(bytes)
	return hexString[:len(hexString)/2], nil
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

func main() {
	router := gin.Default()
	router.GET("/fibonacci/:f", getFibonacci)
	router.GET("/hex/:h", getHexString)
	router.GET("/fibonacci/hex/:f/:h", getFibonacciHex)

	router.Run("localhost:8080")
}
