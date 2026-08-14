package bitwarden

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"go.chrastecky.dev/bitsailor-core/bitwarden/internal/crypto"
	"go.chrastecky.dev/bitsailor-core/bitwarden/internal/dto"
	"go.chrastecky.dev/bitsailor-core/bitwarden/internal/types"
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

func TestGetItemsSortsItemsByName(t *testing.T) {
	t.Parallel()

	userKey := dto.Key("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	malformedID := uuid.MustParse("44444444-4444-4444-4444-444444444444")
	names := []string{"Charlie item", "Alpha item", "Bravo item"}
	encryptedNames := make([]string, 0, len(names))
	for _, name := range names {
		encryptedName, err := crypto.EncryptString(name, userKey)
		if err != nil {
			t.Fatalf("EncryptString(%q) returned error: %v", name, err)
		}
		encryptedNames = append(encryptedNames, encryptedName)
	}

	vault := &vault{
		vaultData: &result.VaultData{
			Items: []*result.Item{
				{ID: uuid.MustParse("33333333-3333-3333-3333-333333333333"), Type: result.ItemTypeSecureNote, Name: encryptedNames[0], SecureNote: &result.ItemSecureNote{}},
				{ID: uuid.MustParse("11111111-1111-1111-1111-111111111111"), Type: result.ItemTypeSecureNote, Name: encryptedNames[1], SecureNote: &result.ItemSecureNote{}},
				{ID: malformedID, Type: result.ItemTypeSecureNote, Name: "not encrypted", SecureNote: &result.ItemSecureNote{}},
				{ID: uuid.MustParse("22222222-2222-2222-2222-222222222222"), Type: result.ItemTypeSecureNote, Name: encryptedNames[2], SecureNote: &result.ItemSecureNote{}},
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

	wantNames := []string{"", "Alpha item", "Bravo item", "Charlie item"}
	if len(items) != len(wantNames) {
		t.Fatalf("len(items) = %d, want %d", len(items), len(wantNames))
	}
	for index, wantName := range wantNames {
		if items[index].Name != wantName {
			t.Fatalf("items[%d].Name = %q, want %q", index, items[index].Name, wantName)
		}
	}
	if items[0].ID != malformedID {
		t.Fatalf("items[0].ID = %s, want malformed item %s", items[0].ID, malformedID)
	}
	if !errors.Is(items[0].DecryptionError, crypto.ErrInvalidEncryptedString) {
		t.Fatalf("items[0].DecryptionError = %v, want ErrInvalidEncryptedString", items[0].DecryptionError)
	}
}

func TestGetItemsReturnsInvalidPlaceholderForMalformedItem(t *testing.T) {
	t.Parallel()

	userKey := dto.Key([]byte("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"))
	validName, err := crypto.EncryptString("Valid item", userKey)
	if err != nil {
		t.Fatalf("EncryptString() returned error: %v", err)
	}
	encryptedNotes, err := crypto.EncryptString("Secret note", userKey)
	if err != nil {
		t.Fatalf("EncryptString() returned error: %v", err)
	}
	encryptedFieldName, err := crypto.EncryptString("Secret field name", userKey)
	if err != nil {
		t.Fatalf("EncryptString() returned error: %v", err)
	}
	encryptedFieldValue, err := crypto.EncryptString("Secret field value", userKey)
	if err != nil {
		t.Fatalf("EncryptString() returned error: %v", err)
	}

	validID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	invalidID := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	orgID := uuid.MustParse("33333333-3333-3333-3333-333333333333")
	folderID := uuid.MustParse("44444444-4444-4444-4444-444444444444")
	collectionIDs := []uuid.UUID{
		uuid.MustParse("55555555-5555-5555-5555-555555555555"),
		uuid.MustParse("66666666-6666-6666-6666-666666666666"),
	}
	revisionDate := time.Date(2026, time.August, 8, 10, 30, 0, 0, time.UTC)
	creationDate := time.Date(2026, time.August, 8, 9, 30, 0, 0, time.UTC)
	archivedDate := time.Date(2026, time.August, 8, 11, 30, 0, 0, time.UTC)
	organizationUseTOTP := true
	linkedID := 7
	permissions := &result.ItemPermissions{Delete: true, Restore: true}
	invalidEncryptedKey := "not encrypted"

	vault := &vault{
		vaultData: &result.VaultData{
			Items: []*result.Item{
				{ID: validID, Type: result.ItemTypeSecureNote, Name: validName, SecureNote: &result.ItemSecureNote{}},
				{
					ID:                  invalidID,
					Type:                result.ItemTypeSecureNote,
					Notes:               &encryptedNotes,
					OrganizationUseTOTP: &organizationUseTOTP,
					RevisionDate:        revisionDate,
					Favorite:            true,
					OrganizationID:      orgID,
					Key:                 &invalidEncryptedKey,
					Permissions:         permissions,
					Edit:                true,
					CollectionIDs:       collectionIDs,
					ArchivedDate:        &archivedDate,
					FolderID:            folderID,
					ViewPassword:        true,
					Name:                "not encrypted",
					CreationDate:        creationDate,
					Reprompt:            types.NumBool(true),
					Fields: []*result.Field{
						{Type: result.FieldTypeLinkedID, Name: encryptedFieldName, Value: &encryptedFieldValue, LinkedID: &linkedID},
					},
					SecureNote: &result.ItemSecureNote{Type: 1},
				},
			},
		},
	}
	session := &result.Session{
		Encryption: &result.EncryptionData{
			UserKey:          userKey,
			OrganizationKeys: map[uuid.UUID]dto.Key{orgID: userKey},
		},
	}

	items, err := vault.GetItems(context.Background(), session)
	if err != nil {
		t.Fatalf("GetItems() returned error: %v", err)
	}

	if len(items) != 2 {
		t.Fatalf("len(items) = %d, want 2", len(items))
	}
	var validItem, invalidItem *result.Item
	for _, item := range items {
		switch item.ID {
		case validID:
			validItem = item
		case invalidID:
			invalidItem = item
		}
	}
	if validItem == nil {
		t.Fatalf("items did not contain valid item %s", validID)
	}
	if validItem.Name != "Valid item" {
		t.Fatalf("validItem.Name = %q, want Valid item", validItem.Name)
	}
	if invalidItem == nil {
		t.Fatalf("items did not contain invalid item %s", invalidID)
	}
	if !errors.Is(invalidItem.DecryptionError, crypto.ErrInvalidEncryptedString) {
		t.Fatalf("invalidItem.DecryptionError = %v, want ErrInvalidEncryptedString", invalidItem.DecryptionError)
	}
	if invalidItem.Name != "" {
		t.Fatalf("invalidItem.Name = %q, want empty string", invalidItem.Name)
	}
	if invalidItem.Notes != nil {
		t.Fatalf("invalidItem.Notes = %q, want nil", *invalidItem.Notes)
	}
	if invalidItem.Key != nil {
		t.Fatalf("invalidItem.Key = %q, want nil", *invalidItem.Key)
	}

	if invalidItem.OrganizationUseTOTP != &organizationUseTOTP || invalidItem.Permissions != permissions {
		t.Fatal("invalidItem did not preserve pointer metadata")
	}
	if invalidItem.RevisionDate != revisionDate || invalidItem.CreationDate != creationDate {
		t.Fatal("invalidItem did not preserve item dates")
	}
	if invalidItem.ArchivedDate != &archivedDate {
		t.Fatal("invalidItem did not preserve archived date")
	}
	if !invalidItem.Favorite || !invalidItem.Edit || !invalidItem.ViewPassword || !bool(invalidItem.Reprompt) {
		t.Fatal("invalidItem did not preserve boolean metadata")
	}
	if invalidItem.OrganizationID != orgID || invalidItem.FolderID != folderID {
		t.Fatal("invalidItem did not preserve UUID metadata")
	}
	if len(invalidItem.CollectionIDs) != len(collectionIDs) || invalidItem.CollectionIDs[0] != collectionIDs[0] || invalidItem.CollectionIDs[1] != collectionIDs[1] {
		t.Fatalf("invalidItem.CollectionIDs = %v, want %v", invalidItem.CollectionIDs, collectionIDs)
	}
	if len(invalidItem.Fields) != 1 {
		t.Fatalf("len(invalidItem.Fields) = %d, want 1", len(invalidItem.Fields))
	}
	if invalidItem.Fields[0].Type != result.FieldTypeLinkedID || invalidItem.Fields[0].LinkedID != &linkedID {
		t.Fatal("invalidItem did not preserve field metadata")
	}
	if invalidItem.Fields[0].Name != "" || invalidItem.Fields[0].Value != nil {
		t.Fatal("invalidItem preserved encrypted field payload")
	}
	if invalidItem.SecureNote == nil || invalidItem.SecureNote.Type != 1 {
		t.Fatal("invalidItem did not preserve secure note metadata")
	}
}
