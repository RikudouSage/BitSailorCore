package bitwarden

import (
	"testing"

	"go.chrastecky.dev/bitsailor-core/bitwarden/internal/crypto"
)

func TestTokenResponseGetPrivateKeyPrefersAccountKeys(t *testing.T) {
	t.Parallel()

	legacyPrivateKey := "legacy-private-key"
	wrappedPrivateKey := "wrapped-private-key"
	response := &tokenResponse{
		PrivateKey: &legacyPrivateKey,
		AccountKeys: &struct {
			Object                     string `json:"Object"`
			PublicKeyEncryptionKeyPair *struct {
				Object            string  `json:"Object"`
				PublicKey         *string `json:"publicKey"`
				WrappedPrivateKey *string `json:"wrappedPrivateKey"`
			} `json:"publicKeyEncryptionKeyPair"`
		}{
			PublicKeyEncryptionKeyPair: &struct {
				Object            string  `json:"Object"`
				PublicKey         *string `json:"publicKey"`
				WrappedPrivateKey *string `json:"wrappedPrivateKey"`
			}{
				WrappedPrivateKey: &wrappedPrivateKey,
			},
		},
	}

	got := response.GetPrivateKey()
	if got == nil || *got != wrappedPrivateKey {
		t.Fatalf("GetPrivateKey() = %v, want wrapped private key", got)
	}
}

func TestTokenResponseGetPrivateKeyFallsBackToLegacyPrivateKey(t *testing.T) {
	t.Parallel()

	legacyPrivateKey := "legacy-private-key"
	response := &tokenResponse{PrivateKey: &legacyPrivateKey}

	got := response.GetPrivateKey()
	if got == nil || *got != legacyPrivateKey {
		t.Fatalf("GetPrivateKey() = %v, want legacy private key", got)
	}
}

func TestTokenResponseGetUserKeyPrefersMasterKeyWrappedUserKey(t *testing.T) {
	t.Parallel()

	response := &tokenResponse{
		Key: "legacy-key",
		UserDecryptionOptions: &struct {
			Object               string `json:"Object"`
			HasMasterPassword    bool   `json:"HasMasterPassword"`
			MasterPasswordUnlock *struct {
				KDF *struct {
					Iterations  int            `json:"Iterations"`
					KDFType     crypto.KDFType `json:"KdfType"`
					Memory      *int           `json:"Memory"`
					Parallelism *int           `json:"Parallelism"`
				} `json:"Kdf"`
				MasterKeyEncryptedUserKey string `json:"MasterKeyEncryptedUserKey"`
				MasterKeyWrappedUserKey   string `json:"MasterKeyWrappedUserKey"`
				Salt                      string `json:"Salt"`
			} `json:"MasterPasswordUnlock"`
		}{
			MasterPasswordUnlock: &struct {
				KDF *struct {
					Iterations  int            `json:"Iterations"`
					KDFType     crypto.KDFType `json:"KdfType"`
					Memory      *int           `json:"Memory"`
					Parallelism *int           `json:"Parallelism"`
				} `json:"Kdf"`
				MasterKeyEncryptedUserKey string `json:"MasterKeyEncryptedUserKey"`
				MasterKeyWrappedUserKey   string `json:"MasterKeyWrappedUserKey"`
				Salt                      string `json:"Salt"`
			}{
				MasterKeyWrappedUserKey: "wrapped-user-key",
			},
		},
	}

	got := response.GetUserKey()
	if got == nil || *got != "wrapped-user-key" {
		t.Fatalf("GetUserKey() = %v, want wrapped user key", got)
	}
}

func TestTokenResponseGetUserKeyFallsBackToLegacyKey(t *testing.T) {
	t.Parallel()

	response := &tokenResponse{Key: "legacy-key"}

	got := response.GetUserKey()
	if got == nil || *got != "legacy-key" {
		t.Fatalf("GetUserKey() = %v, want legacy key", got)
	}
}

func TestTokenResponseGetUserKeyReturnsNilForMissingKey(t *testing.T) {
	t.Parallel()

	if got := (&tokenResponse{}).GetUserKey(); got != nil {
		t.Fatalf("GetUserKey() = %v, want nil", got)
	}
}
