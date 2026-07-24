package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"go.chrastecky.dev/bitsailor-core/bitwarden/internal/dto"
)

func EncryptString(input string, key dto.Key) (string, error) {
	encrypted, err := EncryptBytes([]byte(input), key)
	if err != nil {
		return "", err
	}

	return encrypted, nil
}

func EncryptBytes(input []byte, key dto.Key) (string, error) {
	if len(key) != 64 {
		return "", fmt.Errorf("expected 64-byte user key, got %d", len(key))
	}

	encKey := key[:32]
	macKey := key[32:]

	iv, err := GenerateRandomBytes(aes.BlockSize)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(encKey)
	if err != nil {
		return "", err
	}

	padded := Pkcs7Pad(input, aes.BlockSize)
	ciphertext := make([]byte, len(padded))
	cipher.NewCBCEncrypter(block, iv).CryptBlocks(ciphertext, padded)

	mac := hmac.New(sha256.New, macKey)
	mac.Write(iv)
	mac.Write(ciphertext)
	macSum := mac.Sum(nil)

	return fmt.Sprintf(
		"2.%s|%s|%s",
		base64.StdEncoding.EncodeToString(iv),
		base64.StdEncoding.EncodeToString(ciphertext),
		base64.StdEncoding.EncodeToString(macSum),
	), nil
}
