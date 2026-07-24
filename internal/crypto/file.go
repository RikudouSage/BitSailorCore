package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"hash"
	"io"

	"go.chrastecky.dev/bitsailor-core/bitwarden/internal/dto"
)

const (
	fileEncryptionType  = byte(2)
	fileHeaderLength    = 1 + aes.BlockSize + sha256.Size
	fileHeaderMACOffset = 1 + aes.BlockSize
)

type File struct {
	encrypter  cipher.BlockMode
	mac        hash.Hash
	iv         []byte
	pending    []byte
	ciphertext []byte
	output     []byte
	readAt     int
	finished   bool
	closed     bool
}

func NewFile(key dto.Key) (*File, error) {
	if len(key) != 64 {
		return nil, fmt.Errorf("expected 64-byte file key, got %d", len(key))
	}

	iv, err := GenerateRandomBytes(aes.BlockSize)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key[:32])
	if err != nil {
		return nil, err
	}

	mac := hmac.New(sha256.New, key[32:])
	if _, err = mac.Write(iv); err != nil {
		return nil, err
	}

	return &File{
		encrypter: cipher.NewCBCEncrypter(block, iv),
		mac:       mac,
		iv:        iv,
		pending:   make([]byte, 0, aes.BlockSize),
	}, nil
}

func (receiver *File) Write(data []byte) (int, error) {
	if receiver.closed {
		return 0, fmt.Errorf("encrypted file writer is closed")
	}
	if receiver.finished {
		return 0, fmt.Errorf("encrypted file writer is finished")
	}

	written := len(data)
	receiver.pending = append(receiver.pending, data...)

	encryptLength := len(receiver.pending) - aes.BlockSize
	if encryptLength <= 0 {
		return written, nil
	}
	encryptLength -= encryptLength % aes.BlockSize

	if err := receiver.encryptAndWrite(receiver.pending[:encryptLength]); err != nil {
		return 0, err
	}

	copy(receiver.pending, receiver.pending[encryptLength:])
	receiver.pending = receiver.pending[:len(receiver.pending)-encryptLength]

	return written, nil
}

func (receiver *File) Finish() error {
	if receiver.finished {
		return nil
	}
	if receiver.closed {
		return fmt.Errorf("encrypted file writer is closed")
	}
	receiver.finished = true

	padded := Pkcs7Pad(receiver.pending, aes.BlockSize)
	if err := receiver.encryptAndWrite(padded); err != nil {
		return err
	}

	macSum := receiver.mac.Sum(nil)

	header := make([]byte, fileHeaderLength)
	header[0] = fileEncryptionType
	copy(header[1:], receiver.iv)
	copy(header[fileHeaderMACOffset:], macSum)

	receiver.output = make([]byte, 0, len(header)+len(receiver.ciphertext))
	receiver.output = append(receiver.output, header...)
	receiver.output = append(receiver.output, receiver.ciphertext...)

	return nil
}

func (receiver *File) Read(data []byte) (int, error) {
	if receiver.closed {
		return 0, fmt.Errorf("encrypted file is closed")
	}
	if !receiver.finished {
		if err := receiver.Finish(); err != nil {
			return 0, err
		}
	}
	if receiver.readAt >= len(receiver.output) {
		return 0, io.EOF
	}

	n := copy(data, receiver.output[receiver.readAt:])
	receiver.readAt += n

	return n, nil
}

func (receiver *File) Seek(offset int64, whence int) (int64, error) {
	if receiver.closed {
		return 0, fmt.Errorf("encrypted file is closed")
	}
	if !receiver.finished {
		if err := receiver.Finish(); err != nil {
			return 0, err
		}
	}

	var next int64
	switch whence {
	case io.SeekStart:
		next = offset
	case io.SeekCurrent:
		next = int64(receiver.readAt) + offset
	case io.SeekEnd:
		next = int64(len(receiver.output)) + offset
	default:
		return 0, fmt.Errorf("invalid seek whence: %d", whence)
	}

	if next < 0 || next > int64(len(receiver.output)) {
		return 0, fmt.Errorf("invalid seek offset: %d", next)
	}

	receiver.readAt = int(next)

	return next, nil
}

func (receiver *File) Len() int {
	if !receiver.finished {
		_ = receiver.Finish()
	}

	return len(receiver.output)
}

func (receiver *File) Close() error {
	if receiver.closed {
		return nil
	}
	receiver.closed = true

	return nil
}

func (receiver *File) encryptAndWrite(plaintext []byte) error {
	if len(plaintext) == 0 {
		return nil
	}
	if len(plaintext)%aes.BlockSize != 0 {
		return fmt.Errorf("plaintext length must be a multiple of %d, got %d", aes.BlockSize, len(plaintext))
	}

	ciphertext := make([]byte, len(plaintext))
	receiver.encrypter.CryptBlocks(ciphertext, plaintext)
	if _, err := receiver.mac.Write(ciphertext); err != nil {
		return err
	}
	receiver.ciphertext = append(receiver.ciphertext, ciphertext...)

	return nil
}
