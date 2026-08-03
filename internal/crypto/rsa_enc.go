package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"strings"

	"go.chrastecky.dev/bitsailor-core/bitwarden/internal/dto"
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

func EncryptRSAEncBytes(data []byte, key dto.Key) (string, error) {
	pubAny, err := x509.ParsePKIXPublicKey(key)
	if err != nil {
		return "", fmt.Errorf("failed parsing the key as a public key: %w", err)
	}

	pub, ok := pubAny.(*rsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("public key is %T, not *rsa.PublicKey", pubAny)
	}

	ciphertext, err := rsa.EncryptOAEP(sha1.New(), rand.Reader, pub, data, nil)
	if err != nil {
		return "", fmt.Errorf("failed encrypting the data: %w", err)
	}

	return "4." + base64.StdEncoding.EncodeToString(ciphertext), nil
}
