package bitwarden

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"go.chrastecky.dev/bitsailor-core/bitwarden/internal/crypto"
	"go.chrastecky.dev/bitsailor-core/bitwarden/internal/dto"
	"go.chrastecky.dev/bitsailor-core/bitwarden/result"
)

func TestGetItemsFiltersDeletedItems(t *testing.T) {
	t.Parallel()

	activeID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	deletedID := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	userKey := dto.Key([]byte("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"))
	activeName, err := crypto.EncryptString("Active item", userKey)
	if err != nil {
		t.Fatalf("EncryptString() returned error: %v", err)
	}
	deletedDate := time.Now()

	vault := &vault{
		vaultData: &result.VaultData{
			Items: []*result.Item{
				{ID: activeID, Type: result.ItemTypeSecureNote, Name: activeName, SecureNote: &result.ItemSecureNote{}},
				{ID: deletedID, Type: result.ItemTypeSecureNote, Name: "not encrypted", DeletedDate: &deletedDate, SecureNote: &result.ItemSecureNote{}},
			},
		},
	}
	session := &result.Session{
		Encryption: &result.EncryptionData{
			UserKey: userKey,
		},
	}

	items, err := vault.GetItems(context.Background(), session)
	if err != nil {
		t.Fatalf("GetItems() returned error: %v", err)
	}

	if len(items) != 1 {
		t.Fatalf("len(items) = %d, want 1", len(items))
	}
	if items[0].ID != activeID {
		t.Fatalf("items[0].ID = %s, want %s", items[0].ID, activeID)
	}
	if items[0].Name != "Active item" {
		t.Fatalf("items[0].Name = %q, want Active item", items[0].Name)
	}
}
