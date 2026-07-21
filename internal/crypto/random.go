package crypto

import (
	"crypto/rand"
	"fmt"
)

func GenerateRandomBytes(length int) ([]byte, error) {
	buf := make([]byte, length)
	_, err := rand.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("failed generating random bytes: %w", err)
	}

	return buf, nil
}
