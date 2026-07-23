package interop

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"

	"go.chrastecky.dev/bitwarden-client/bitwarden/internal/crypto"
	"go.chrastecky.dev/bitwarden-client/bitwarden/result"
)

func CreateBWCLISessionID(session *result.Session) (BWSession string, protectedValue string, err error) {
	if session == nil {
		return "", "", errors.New("the session is nil")
	}
	if session.Encryption.UserKey == nil {
		return "", "", errors.New("the user key is nil")
	}

	bwSessionRaw, err := crypto.GenerateRandomBytes(64)
	if err != nil {
		return "", "", fmt.Errorf("failed generating random bytes: %w", err)
	}
	bwSessionID := base64.StdEncoding.EncodeToString(bwSessionRaw)

	protectedBytes, err := encryptCLIFileData(session.Encryption.UserKey, bwSessionRaw)
	if err != nil {
		return "", "", err
	}

	protectedValue = base64.StdEncoding.EncodeToString(protectedBytes)
	return bwSessionID, protectedValue, nil
}

func encryptCLIFileData(plaintext []byte, key []byte) ([]byte, error) {
	if len(key) != 64 {
		return nil, fmt.Errorf("expected 64-byte session key, got %d", len(key))
	}

	encKey := key[:32]
	macKey := key[32:]

	iv, err := crypto.GenerateRandomBytes(aes.BlockSize)
	if err != nil {
		return nil, fmt.Errorf("failed generating random bytes for IV: %w", err)
	}

	block, err := aes.NewCipher(encKey)
	if err != nil {
		return nil, err
	}

	padded := crypto.Pkcs7Pad(plaintext, aes.BlockSize)
	ciphertext := make([]byte, len(padded))
	cipher.NewCBCEncrypter(block, iv).CryptBlocks(ciphertext, padded)

	mac := hmac.New(sha256.New, macKey)
	mac.Write(iv)
	mac.Write(ciphertext)
	macBytes := mac.Sum(nil)

	// byte 0 = EncryptionType.AesCbc256_HmacSha256_B64 = 2
	// bytes 1..16 = IV
	// bytes 17..48 = MAC
	// bytes 49.. = ciphertext
	out := make([]byte, 1+len(iv)+len(macBytes)+len(ciphertext))
	out[0] = 2
	copy(out[1:], iv)
	copy(out[1+len(iv):], macBytes)
	copy(out[1+len(iv)+len(macBytes):], ciphertext)

	return out, nil
}
