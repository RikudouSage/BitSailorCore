package bitwarden

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/google/uuid"
	"go.chrastecky.dev/bitsailor-core/bitwarden/result"
)

func TestVaultDataAccessors(t *testing.T) {
	t.Parallel()

	original := &vault{}
	vaultData := &result.VaultData{}

	clone := original.WithVaultData(vaultData).(*vault)
	if clone == original {
		t.Fatal("WithVaultData() returned original vault, want cloned vault")
	}
	if clone.GetVaultData() != vaultData {
		t.Fatal("GetVaultData() did not return configured vault data")
	}
	if original.GetVaultData() != nil {
		t.Fatal("original vault data was mutated")
	}
}

func TestGetItemMissingVault(t *testing.T) {
	t.Parallel()

	_, err := (&vault{}).GetItem(context.Background(), nil, uuid.New())
	if !errors.Is(err, ErrMissingVault) {
		t.Fatalf("GetItem() error = %v, want ErrMissingVault", err)
	}
}

func TestGetItemNotFound(t *testing.T) {
	t.Parallel()

	_, err := (&vault{vaultData: &result.VaultData{}}).GetItem(context.Background(), nil, uuid.New())
	if !errors.Is(err, ErrItemNotFound) {
		t.Fatalf("GetItem() error = %v, want ErrItemNotFound", err)
	}
}

func TestGetSendMissingVault(t *testing.T) {
	t.Parallel()

	_, err := (&vault{}).GetSend(context.Background(), nil, uuid.New())
	if !errors.Is(err, ErrMissingVault) {
		t.Fatalf("GetSend() error = %v, want ErrMissingVault", err)
	}
}

func TestGetSendNotFound(t *testing.T) {
	t.Parallel()

	_, err := (&vault{vaultData: &result.VaultData{}}).GetSend(context.Background(), nil, uuid.New())
	if !errors.Is(err, ErrItemNotFound) {
		t.Fatalf("GetSend() error = %v, want ErrItemNotFound", err)
	}
}

func TestGetSendsMissingVault(t *testing.T) {
	t.Parallel()

	_, err := (&vault{}).GetSends(context.Background(), nil)
	if !errors.Is(err, ErrMissingVault) {
		t.Fatalf("GetSends() error = %v, want ErrMissingVault", err)
	}
}

func TestSyncRequiresAuthData(t *testing.T) {
	t.Parallel()

	apiURL, err := url.Parse("https://api.example.test")
	if err != nil {
		t.Fatalf("url.Parse() returned error: %v", err)
	}
	vault := &vault{apiURL: apiURL, httpClient: http.DefaultClient, auth: &auth{}}

	if _, err := vault.Sync(context.Background(), nil); err == nil {
		t.Fatal("Sync(nil) error = nil, want error")
	}
	if _, err := vault.Sync(context.Background(), &result.Session{}); err == nil {
		t.Fatal("Sync(session without auth) error = nil, want error")
	}
}
