package crypto

import (
	"crypto/sha256"
	"encoding/base64"

	"go.chrastecky.dev/bitsailor-core/bitwarden/internal/dto"
	"golang.org/x/crypto/pbkdf2"
)

func DeriveSendPassword(password string, seed dto.Key) string {
	hash := pbkdf2.Key(
		[]byte(password),
		seed,
		100_000,
		32,
		sha256.New,
	)

	return base64.StdEncoding.EncodeToString(hash)
}
