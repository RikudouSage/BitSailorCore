package crypto

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"strings"

	"go.chrastecky.dev/bitsailor-core/bitwarden/internal/dto"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/hkdf"
	"golang.org/x/crypto/pbkdf2"
)

func DeriveMasterKey(email, password string, keyType KDFType, config *KDFConfig) ([]byte, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	switch keyType {
	case KDFTypeSHA256:
		return deriveSha256(email, password, config.Iterations), nil
	case KDFTypeArgon2ID:
		if config.Memory == nil || config.Parallelism == nil {
			return nil, errors.New("memory and parallelism must be set")
		}

		return deriveArgon2ID(email, password, config.Iterations, *config.Memory, *config.Parallelism), nil
	}

	return nil, fmt.Errorf("unsupported key type: %d", keyType)
}

func DeriveSendKey(seed dto.Key) (dto.Key, error) {
	reader := hkdf.New(
		sha256.New,
		seed,
		[]byte("bitwarden-send"),
		[]byte("send"),
	)

	key := make([]byte, 64)
	if _, err := io.ReadFull(reader, key); err != nil {
		return nil, fmt.Errorf("failed generating key: %w", err)
	}

	return key, nil
}

func deriveSha256(email, password string, iterations int) []byte {
	return pbkdf2.Key(
		[]byte(password),
		[]byte(email),
		iterations,
		32,
		sha256.New,
	)
}

func deriveArgon2ID(email, password string, iterations int, memory int, parallelism int) []byte {
	return argon2.IDKey(
		[]byte(password),
		[]byte(email),
		uint32(iterations),
		uint32(memory)*1024,
		uint8(parallelism),
		32,
	)
}
