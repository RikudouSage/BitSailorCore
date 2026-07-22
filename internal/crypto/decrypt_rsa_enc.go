package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
)

func DecryptRSAEncString(encrypted string, privateKey *rsa.PrivateKey) ([]byte, error) {
	parts := strings.SplitN(encrypted, ".", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid encrypted string")
	}

	ciphertext, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed base64 decoding ciphertext: %w", err)
	}

	switch parts[0] {
	case "3":
		return rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, ciphertext, nil)
	case "4":
		return rsa.DecryptOAEP(sha1.New(), rand.Reader, privateKey, ciphertext, nil)
	default:
		return nil, fmt.Errorf("unsupported RSA encrypted string type: %s", parts[0])
	}
}
