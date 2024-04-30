package random

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// GenerateStringID generates a random string with a length of 8 chars.
//
// This function is useful for generating IDs, or name prefixes/suffixes.
func GenerateStringID() (string, error) {
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("failed to generate random string: %w", err)
	}
	return hex.EncodeToString(b), nil
}
