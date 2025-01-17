package utils

import (
	"math/rand"
	"os"
)

// FileExists checks if a file exists
func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

// Get a random number between a start and end value
func GetRandomNumber(start, end int) int {
	return rand.Intn(end) + start
}
