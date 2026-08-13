package crypto

import (
	"errors"
	"strings"
	"testing"

	"go.chrastecky.dev/bitsailor-core/bitwarden/internal/dto"
)

func TestEncryptDecryptStringRoundTrip(t *testing.T) {
	t.Parallel()

	key := dto.Key([]byte("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"))
	encrypted, err := EncryptString("secret text", key)
	if err != nil {
		t.Fatalf("EncryptString() returned error: %v", err)
	}
	if !strings.HasPrefix(encrypted, "2.") {
		t.Fatalf("encrypted value = %q, want type 2 prefix", encrypted)
	}

	decrypted, err := DecryptString(encrypted, key)
	if err != nil {
		t.Fatalf("DecryptString() returned error: %v", err)
	}
	if decrypted != "secret text" {
		t.Fatalf("DecryptString() = %q, want secret text", decrypted)
	}
}

func TestDecryptNullableString(t *testing.T) {
	t.Parallel()

	key := dto.Key([]byte("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"))
	if got, err := DecryptNullableString(nil, key); err != nil || got != nil {
		t.Fatalf("DecryptNullableString(nil) = %v, %v; want nil, nil", got, err)
	}

	encrypted, err := EncryptString("secret text", key)
	if err != nil {
		t.Fatalf("EncryptString() returned error: %v", err)
	}
	got, err := DecryptNullableString(&encrypted, key)
	if err != nil {
		t.Fatalf("DecryptNullableString() returned error: %v", err)
	}
	if got == nil || *got != "secret text" {
		t.Fatalf("DecryptNullableString() = %v, want secret text", got)
	}
}

func TestDecryptBytesRejectsMalformedEncryptedStrings(t *testing.T) {
	t.Parallel()

	key := dto.Key([]byte("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"))
	tests := []string{
		"missing-type-separator",
		"1.a|b|c",
		"2.a|b",
		"2.not-base64|b|c",
	}

	for _, encrypted := range tests {
		t.Run(encrypted, func(t *testing.T) {
			t.Parallel()

			_, err := DecryptBytes(encrypted, key)
			if !errors.Is(err, ErrInvalidEncryptedString) {
				t.Fatalf("DecryptBytes() error = %v, want ErrInvalidEncryptedString", err)
			}
		})
	}
}

func TestEncryptDecryptRejectWrongKeyLength(t *testing.T) {
	t.Parallel()

	if _, err := EncryptString("secret text", dto.Key([]byte("short"))); err == nil {
		t.Fatal("EncryptString() error = nil, want key length error")
	}
	if _, err := DecryptString("2.a|b|c", dto.Key([]byte("short"))); err == nil {
		t.Fatal("DecryptString() error = nil, want key length error")
	}
}
