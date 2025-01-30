package testutil

import (
	"crypto/rand"
	"errors"
	"strings"
)

var (
	Sizes = map[int]int{
		1:  2 * 1024,               // 2KB
		2:  4 * 1024,               // 4KB
		3:  8 * 1024,               // 8KB
		4:  16 * 1024,              // 16KB
		5:  32 * 1024,              // 32KB
		6:  64 * 1024,              // 64KB
		7:  128 * 1024,             // 128KB
		8:  256 * 1024,             // 256KB
		9:  512 * 1024,             // 512KB
		10: 1 * 1024 * 1024,        // 1MB
		11: 10 * 1024 * 1024,       // 10MB
		12: 100 * 1024 * 1024,      // 100MB
		13: 1 * 1024 * 1024 * 1024, // 1GB
	}
)

func GenerateMessage(sizeType int) ([]byte, error) {
	targetSize, ok := Sizes[sizeType]
	if !ok {
		return []byte{}, errors.New("invalid size type: must be between 1 and 6")
	}

	// Generate random content
	bytes := make([]byte, targetSize)
	_, err := rand.Read(bytes)
	if err != nil {
		return []byte{}, err
	}

	// Convert random bytes to a string
	for i, b := range bytes {
		// Limit to printable ASCII characters (32 to 126)
		bytes[i] = 32 + (b % 95)
	}

	return bytes, nil
}

func GenerateStringMessage(sizeType int) (string, error) {
	targetSize, ok := Sizes[sizeType]
	if !ok {
		return "", errors.New("invalid size type: must be between 1 and 13")
	}

	// Generate random content
	bytes := make([]byte, targetSize)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	// Convert random bytes to a string
	for i, b := range bytes {
		// Limit to printable ASCII characters (32 to 126)
		bytes[i] = 32 + (b % 95)
	}

	return string(bytes), nil
}

func GenerateRepeatedMessage(sizeType int) ([]byte, error) {
	// Get the target size
	targetSize, ok := Sizes[sizeType]
	if !ok {
		return nil, errors.New("invalid size type: must be between 1 and 13")
	}

	// Define a simple repetitive pattern
	pattern := "ABCDEF1234567890"
	patternLength := len(pattern)

	// Calculate how many times to repeat the pattern
	repeatCount := (targetSize / patternLength) + 1

	// Generate the repeated message
	repeatedMessage := strings.Repeat(pattern, repeatCount)

	// Ensure the message is at least the target size
	if len(repeatedMessage) < targetSize {
		return nil, errors.New("failed to generate repeated message of sufficient size")
	}

	// Return the message trimmed to the exact target size
	return []byte(repeatedMessage[:targetSize]), nil
}
