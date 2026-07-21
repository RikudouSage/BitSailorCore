package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
)

func DecryptUserKey(encryptedUserKey string, masterKey []byte) ([]byte, error) {
	encKey, macKey, err := stretchMasterKey(masterKey)
	if err != nil {
		return nil, fmt.Errorf("failed stretching master key: %w", err)
	}

	return decryptEncryptedString(encryptedUserKey, encKey, macKey)
}

func decryptEncryptedString(value string, encKey, macKey []byte) ([]byte, error) {
	parts := strings.SplitN(value, ".", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid encrypted string: %s", value)
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
		return nil, fmt.Errorf("failed decoding iv: %w", err)
	}

	ciphertext, err := base64.StdEncoding.DecodeString(payload[1])
	if err != nil {
		return nil, fmt.Errorf("failed to decode ciphertext: %w", err)
	}

	expectedMAC, err := base64.StdEncoding.DecodeString(payload[2])
	if err != nil {
		return nil, fmt.Errorf("failed to decode mac: %w", err)
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
	if len(iv) != aes.BlockSize {
		return nil, fmt.Errorf("invalid iv length: %d", len(iv))
	}
	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("invalid ciphertext length: %d", len(ciphertext))
	}

	plaintext := make([]byte, len(ciphertext))
	cipher.NewCBCDecrypter(block, iv).CryptBlocks(plaintext, ciphertext)

	return pkcs7Unpad(plaintext, aes.BlockSize)
}

func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	if len(data) == 0 || len(data)%blockSize != 0 {
		return nil, fmt.Errorf("invalid padding size: %d", len(data))
	}

	padLen := int(data[len(data)-1])
	if padLen == 0 || padLen > blockSize || padLen > len(data) {
		return nil, fmt.Errorf("invalid padding")
	}

	for _, b := range data[len(data)-padLen:] {
		if int(b) != padLen {
			return nil, fmt.Errorf("invalid padding")
		}
	}

	return data[:len(data)-padLen], nil
}

func stretchMasterKey(masterKey []byte) (encKey []byte, macKey []byte, err error) {
	encKey, err = hkdfSha256(masterKey, "enc", 32)
	if err != nil {
		return
	}
	macKey, err = hkdfSha256(masterKey, "mac", 32)
	if err != nil {
		return
	}

	return
}

func hkdfSha256(privateKey []byte, info string, length int) ([]byte, error) {
	hashLen := sha256.Size
	if length > 255*hashLen {
		return nil, fmt.Errorf("hkdf expand length too large")
	}

	var result []byte
	var previous []byte

	for counter := byte(1); len(result) < length; counter++ {
		hash := hmac.New(sha256.New, privateKey)
		hash.Write(previous)
		hash.Write([]byte(info))
		hash.Write([]byte{counter})

		previous = hash.Sum(nil)
		result = append(result, previous...)
	}

	return result[:length], nil
}
