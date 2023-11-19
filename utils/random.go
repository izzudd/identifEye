package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateRandomString(length int) (string, error) {
	// Determine the number of bytes needed
	numBytes := length / 2
	if length%2 != 0 {
		numBytes++
	}

	// Generate random bytes
	randomBytes := make([]byte, numBytes)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// Convert bytes to hex string
	randomString := hex.EncodeToString(randomBytes)[:length]

	return randomString, nil
}
