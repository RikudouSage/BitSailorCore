package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"io"
	"strings"
	"testing"
)

func TestFileEncryptWriter(t *testing.T) {
	key := make([]byte, 64)
	for i := range key {
		key[i] = byte(i + 1)
	}

	plaintext := "this plaintext is intentionally longer than one AES block"

	file, err := NewFile(key)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = io.Copy(file, strings.NewReader(plaintext)); err != nil {
		t.Fatal(err)
	}
	if err = file.Finish(); err != nil {
		t.Fatal(err)
	}
	if file.Len() == 0 {
		t.Fatal("expected encrypted file length")
	}

	var encryptedFile bytes.Buffer
	if _, err = io.Copy(&encryptedFile, file); err != nil {
		t.Fatal(err)
	}

	encrypted := encryptedFile.Bytes()
	if encrypted[0] != fileEncryptionType {
		t.Fatalf("expected encryption type %d, got %d", fileEncryptionType, encrypted[0])
	}

	iv := encrypted[1 : 1+aes.BlockSize]
	actualMAC := encrypted[fileHeaderMACOffset:fileHeaderLength]
	ciphertext := encrypted[fileHeaderLength:]

	mac := hmac.New(sha256.New, key[32:])
	mac.Write(iv)
	mac.Write(ciphertext)
	if !hmac.Equal(mac.Sum(nil), actualMAC) {
		t.Fatal("invalid mac")
	}

	block, err := aes.NewCipher(key[:32])
	if err != nil {
		t.Fatal(err)
	}
	decrypted := make([]byte, len(ciphertext))
	cipher.NewCBCDecrypter(block, iv).CryptBlocks(decrypted, ciphertext)

	decrypted, err = pkcs7Unpad(decrypted, aes.BlockSize)
	if err != nil {
		t.Fatal(err)
	}
	if string(decrypted) != plaintext {
		t.Fatalf("expected %q, got %q", plaintext, string(decrypted))
	}
}

func TestFileReadFinishesEncryption(t *testing.T) {
	key := make([]byte, 64)
	for i := range key {
		key[i] = byte(i + 1)
	}

	file, err := NewFile(key)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = io.Copy(file, strings.NewReader("plaintext")); err != nil {
		t.Fatal(err)
	}

	var encrypted bytes.Buffer
	if _, err = io.Copy(&encrypted, file); err != nil {
		t.Fatal(err)
	}
	if encrypted.Len() == 0 {
		t.Fatal("expected encrypted data")
	}
}

func TestFileLenFinishesEncryption(t *testing.T) {
	key := make([]byte, 64)
	for i := range key {
		key[i] = byte(i + 1)
	}

	file, err := NewFile(key)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = io.Copy(file, strings.NewReader("plaintext")); err != nil {
		t.Fatal(err)
	}

	if file.Len() == 0 {
		t.Fatal("expected encrypted file length")
	}

	var encrypted bytes.Buffer
	if _, err = io.Copy(&encrypted, file); err != nil {
		t.Fatal(err)
	}
	if uint64(encrypted.Len()) != file.Len() {
		t.Fatalf("expected read length %d, got %d", file.Len(), encrypted.Len())
	}
}

func TestFileSeek(t *testing.T) {
	key := make([]byte, 64)
	for i := range key {
		key[i] = byte(i + 1)
	}

	file, err := NewFile(key)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = io.Copy(file, strings.NewReader("plaintext")); err != nil {
		t.Fatal(err)
	}

	first, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	if len(first) == 0 {
		t.Fatal("expected encrypted data")
	}

	offset, err := file.Seek(0, io.SeekStart)
	if err != nil {
		t.Fatal(err)
	}
	if offset != 0 {
		t.Fatalf("expected offset 0, got %d", offset)
	}

	second, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(first, second) {
		t.Fatal("expected reread data to match")
	}

	offset, err = file.Seek(-int64(len(second)), io.SeekEnd)
	if err != nil {
		t.Fatal(err)
	}
	if offset != 0 {
		t.Fatalf("expected offset 0, got %d", offset)
	}
}
