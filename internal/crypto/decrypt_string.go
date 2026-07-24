package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"

	"go.chrastecky.dev/bitsailor-core/bitwarden/internal/dto"
)

func DecryptNullableString(encrypted *string, userKey []byte) (*string, error) {
	if encrypted == nil {
		return nil, nil
	}

	result, err := DecryptString(*encrypted, userKey)
	if err != nil {
		return nil, err
	}

	return new(result), nil
}

func DecryptString(encrypted string, key []byte) (string, error) {
	decrypted, err := DecryptBytes(encrypted, key)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}

func DecryptBytes(encrypted string, key dto.Key) ([]byte, error) {

	if len(key) != 64 {
		return nil, fmt.Errorf("expected 64-byte user key, got %d", len(key))
	}

	encKey := key[:32]
	macKey := key[32:]

	parts := strings.SplitN(encrypted, ".", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid encrypted string")
	}
	if parts[0] != "2" {
		return nil, fmt.Errorf("unsupported encrypted string type: %s", parts[0])
	}

	payload := strings.Split(parts[1], "|")
	if len(payload) != 3 {
		return nil, fmt.Errorf("invalid encrypted string payload")
	}

	iv, err := base64.StdEncoding.DecodeString(payload[0])
	if err != nil {
		return nil, fmt.Errorf("failed base64 decoding iv: %w", err)
	}

	ciphertext, err := base64.StdEncoding.DecodeString(payload[1])
	if err != nil {
		return nil, fmt.Errorf("failed base64 decoding ciphertext: %w", err)
	}

	expectedMAC, err := base64.StdEncoding.DecodeString(payload[2])
	if err != nil {
		return nil, fmt.Errorf("failed base64 decoding mac: %w", err)
	}

	mac := hmac.New(sha256.New, macKey)
	mac.Write(iv)
	mac.Write(ciphertext)

	if !hmac.Equal(mac.Sum(nil), expectedMAC) {
		return nil, fmt.Errorf("invalid mac")
	}

	block, err := aes.NewCipher(encKey)
	if err != nil {
		return nil, err
	}

	plaintext := make([]byte, len(ciphertext))
	cipher.NewCBCDecrypter(block, iv).CryptBlocks(plaintext, ciphertext)

	result, err := pkcs7Unpad(plaintext, aes.BlockSize)
	if err != nil {
		return nil, fmt.Errorf("failed unpadding plaintext: %w", err)
	}

	return result, nil
}
